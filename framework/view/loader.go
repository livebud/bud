package view

import (
	"context"
	"errors"
	"io/fs"
	"path"
	"strings"

	"github.com/livebud/bud/package/vfs"

	"github.com/livebud/bud/framework"
	"github.com/livebud/bud/runtime/view/dom"
	"github.com/livebud/bud/runtime/view/ssr"

	"github.com/livebud/bud/internal/bail"
	"github.com/livebud/bud/internal/embed"
	"github.com/livebud/bud/internal/embedded"
	"github.com/livebud/bud/internal/entrypoint"
	"github.com/livebud/bud/internal/imports"
	"github.com/livebud/bud/package/gomod"
	"github.com/livebud/bud/runtime/transform"
)

func Load(
	ctx context.Context,
	fsys fs.FS,
	module *gomod.Module,
	transform *transform.Map,
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
	fsys      fs.FS
	module    *gomod.Module
	transform *transform.Map
	flag      *framework.Flag

	bail.Struct
	imports *imports.Set
}

func (l *loader) Load(ctx context.Context) (state *State, err error) {
	defer l.Recover2(&err, "view: unable to load")
	state = &State{
		Flag: l.flag,
	}
	views, err := entrypoint.List(l.fsys, "view")
	if err != nil {
		return nil, err
	} else if len(views) == 0 {
		return nil, fs.ErrNotExist
	}
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
		// Add DOM
		domCompiler := dom.New(l.module, l.transform.DOM)
		files, err := domCompiler.Compile(ctx, l.fsys)
		if err != nil {
			return nil, err
		}
		for _, file := range files {
			// TODO: decide if we should be doing strings.ToLower here. It's needed
			// because the router always lower-cases things, but the generated chunks
			// contain are upper values
			state.Embeds = append(state.Embeds, &embed.File{
				Path: path.Join("bud/view", strings.ToLower(file.Path)) + ".js",
				Data: file.Contents,
			})
		}
		// Add default layout.css
		if err := vfs.Exist(l.fsys, "view/layout.css"); err != nil {
			if !errors.Is(err, fs.ErrNotExist) {
				return nil, err
			}
			state.Embeds = append(state.Embeds, &embed.File{
				Path: "bud/view/layout.css",
				Data: embedded.Layout(),
			})
		}
	}
	// fmt.Println(l.Flag.Embed, l.Transform.SSR, views)
	if l.flag.Embed {
		l.imports.AddNamed("overlay", "github.com/livebud/bud/package/overlay")
		l.imports.AddNamed("mod", "github.com/livebud/bud/package/gomod")
		l.imports.AddNamed("js", "github.com/livebud/bud/package/js")
	} else {
		l.imports.AddNamed("budproxy", "github.com/livebud/bud/package/budproxy")
	}
	l.imports.AddNamed("view", "github.com/livebud/bud/runtime/view")
	state.Imports = l.imports.List()
	return state, nil
}
