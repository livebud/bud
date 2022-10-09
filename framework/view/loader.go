package view

import (
	"context"
	"io/fs"
	"path"

	"github.com/livebud/bud/framework"
	"github.com/livebud/bud/framework/view/dom"
	"github.com/livebud/bud/framework/view/ssr"

	"github.com/livebud/bud/framework/transform/transformrt"
	"github.com/livebud/bud/internal/bail"
	"github.com/livebud/bud/internal/embed"
	"github.com/livebud/bud/internal/entrypoint"
	"github.com/livebud/bud/internal/imports"
	"github.com/livebud/bud/package/budfs"
	"github.com/livebud/bud/package/gomod"
)

func Load(
	ctx context.Context,
	fsys budfs.FS,
	module *gomod.Module,
	transform *transformrt.Map,
	flag *framework.Flag,
) (*State, error) {
	return (&loader{
		fsys:      fsys,
		module:    module,
		transform: transform,
		flag:      flag,
		imports:   imports.New(),
	}).Load(ctx)
}

type loader struct {
	fsys      budfs.FS
	module    *gomod.Module
	transform *transformrt.Map
	flag      *framework.Flag

	bail.Struct
	imports *imports.Set
}

func (l *loader) Load(ctx context.Context) (state *State, err error) {
	defer l.Recover2(&err, "view: unable to load")
	state = &State{}
	views, err := entrypoint.List(l.fsys, "view")
	if err != nil {
		return nil, err
	} else if len(views) == 0 {
		return nil, fs.ErrNotExist
	}
	// Load the embeds
	if l.flag.Embed {
		// Add SSR
		ssrCompiler := ssr.New(l.module, l.transform.SSR)
		ssrCode, err := ssrCompiler.Compile(ctx, l.fsys)
		if err != nil {
			return nil, err
		}
		state.Embeds = append(state.Embeds, &embed.File{
			Path: "bud/view/_ssr.js",
			Data: ssrCode,
		})
		// Bundle client-side files
		domCompiler := dom.New(l.module, l.transform.DOM)
		files, err := domCompiler.Compile(ctx, l.fsys)
		if err != nil {
			return nil, err
		}
		for _, file := range files {
			// Lowercase the path for the router which requires lowercase routes
			filePath := path.Join("bud/view", file.Path)
			state.Embeds = append(state.Embeds, &embed.File{
				Path: filePath,
				Data: file.Contents,
			})
			state.Routes = append(state.Routes, "/"+filePath)
		}
	} else {
		// Load the routes as references
		for _, view := range views {
			// Add the entrypoint
			state.Routes = append(state.Routes, "/"+view.Client)
			// Add the dynamic import
			state.Routes = append(state.Routes, "/bud/"+string(view.Page))
		}
		// Add node modules if we're not bundling
		state.Routes = append(state.Routes, "/bud/node_modules/:module*")
	}
	// Add the imports
	l.imports.AddStd("io/fs", "net/http")
	l.imports.AddNamed("router", "github.com/livebud/bud/package/router")
	l.imports.AddNamed("virtual", "github.com/livebud/bud/package/virtual")
	l.imports.AddNamed("viewrt", "github.com/livebud/bud/framework/view/viewrt")
	state.Imports = l.imports.List()
	return state, nil
}
