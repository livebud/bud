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
	"github.com/livebud/bud/package/gomod"
)

func Load(
	ctx context.Context,
	fsys fs.FS,
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
	fsys      fs.FS
	module    *gomod.Module
	transform *transformrt.Map
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
			state.Embeds = append(state.Embeds, &embed.File{
				Path: path.Join("bud/view", file.Path),
				Data: file.Contents,
			})
		}
	}
	// fmt.Println(l.Flag.Embed, l.Transform.SSR, views)
	if l.flag.Embed {
		l.imports.AddNamed("overlay", "github.com/livebud/bud/package/overlay")
		l.imports.AddNamed("mod", "github.com/livebud/bud/package/gomod")
		l.imports.AddNamed("js", "github.com/livebud/bud/package/js")
	} else {
		l.imports.AddNamed("budclient", "github.com/livebud/bud/package/budclient")
	}
	l.imports.AddNamed("viewrt", "github.com/livebud/bud/framework/view/viewrt")
	state.Imports = l.imports.List()
	return state, nil
}
