package view

import (
	"context"
	"io/fs"
	"path"
	"strings"

	"github.com/livebud/bud/runtime/command"
	"github.com/livebud/bud/runtime/view/dom"
	"github.com/livebud/bud/runtime/view/ssr"

	"github.com/livebud/bud/internal/embed"
	"github.com/livebud/bud/internal/entrypoint"
	"github.com/livebud/bud/internal/imports"
	"github.com/livebud/bud/package/gomod"
	"github.com/livebud/bud/runtime/transform"
)

type parser struct {
	Flag      *command.Flag
	Module    *gomod.Module
	Imports   *imports.Set
	Transform *transform.Map
}

type State struct {
	Imports []*imports.Import
	Flag    *command.Flag
	Embeds  []*embed.File
}

func (p *parser) Parse(fsys fs.FS, ctx context.Context) (*State, error) {
	state := &State{
		Flag: p.Flag,
	}
	views, err := entrypoint.List(fsys, "view")
	if err != nil {
		return nil, err
	} else if len(views) == 0 {
		return nil, fs.ErrNotExist
	}
	if p.Flag.Embed {
		// Add SSR
		ssrCompiler := ssr.New(p.Module, p.Transform.SSR)
		ssrCode, err := ssrCompiler.Compile(ctx, fsys)
		if err != nil {
			return nil, err
		}
		state.Embeds = append(state.Embeds, &embed.File{
			Path: "bud/view/_ssr.js",
			Data: ssrCode,
		})
		// Add DOM
		domCompiler := dom.New(p.Module, p.Transform.DOM)
		files, err := domCompiler.Compile(ctx, fsys)
		if err != nil {
			return nil, err
		}
		for _, file := range files {
			// TODO: decide if we should be doing strings.ToLower here. It's needed
			// because the router always lower-cases things, but the generated chunks
			// contain are upper values
			state.Embeds = append(state.Embeds, &embed.File{
				Path: path.Join("bud/view", strings.ToLower(file.Path)),
				Data: file.Contents,
			})
		}
	}
	// fmt.Println(p.Flag.Embed, p.Transform.SSR, views)
	p.Imports.AddNamed("transform", p.Module.Import("bud/.cli/transform"))
	p.Imports.AddNamed("overlay", "github.com/livebud/bud/package/overlay")
	p.Imports.AddNamed("mod", "github.com/livebud/bud/package/gomod")
	p.Imports.AddNamed("js", "github.com/livebud/bud/package/js")
	p.Imports.AddNamed("view", "github.com/livebud/bud/runtime/view")
	state.Imports = p.Imports.List()
	return state, nil
}
