package generator

import (
	"context"
	"errors"

	"gitlab.com/mnm/bud/internal/fscache"
	"gitlab.com/mnm/bud/internal/gitignore"
	"gitlab.com/mnm/bud/pkg/budfs"
	"gitlab.com/mnm/bud/pkg/gen"

	"gitlab.com/mnm/bud/internal/fsync"

	"gitlab.com/mnm/bud/generator/public"
	"gitlab.com/mnm/bud/internal/di"
	"gitlab.com/mnm/bud/internal/generator/action"
	"gitlab.com/mnm/bud/internal/generator/command"
	"gitlab.com/mnm/bud/internal/generator/maingo"
	"gitlab.com/mnm/bud/internal/generator/program"
	"gitlab.com/mnm/bud/internal/generator/transform"
	"gitlab.com/mnm/bud/internal/generator/view"
	"gitlab.com/mnm/bud/internal/generator/web"
	"gitlab.com/mnm/bud/internal/modcache"
	"gitlab.com/mnm/bud/internal/parser"
	"gitlab.com/mnm/bud/pkg/gomod"
	"gitlab.com/mnm/bud/pkg/vfs"
)

type option struct {
	Embed    bool
	Hot      bool
	Minify   bool
	Replaces []*gomod.Replace
	ModCache *modcache.Cache
	FSCache  *fscache.Cache
}

type Option func(*option)

func WithEmbed(embed bool) func(*option) {
	return func(option *option) {
		option.Embed = embed
	}
}

func WithHot(hot bool) func(*option) {
	return func(option *option) {
		option.Hot = hot
	}
}

func WithMinify(minify bool) func(*option) {
	return func(option *option) {
		option.Minify = minify
	}
}

func WithReplace(from, to string) func(*option) {
	return func(option *option) {
		option.Replaces = append(option.Replaces, &gomod.Replace{
			Old: gomod.Version{Path: from, Version: ""},
			New: gomod.Version{Path: to, Version: ""},
		})
	}
}

func WithModCache(mc *modcache.Cache) func(*option) {
	return func(option *option) {
		option.ModCache = mc
	}
}

func WithFSCache(fc *fscache.Cache) func(*option) {
	return func(option *option) {
		option.FSCache = fc
	}
}

func toType(importPath, dataType string) *di.Type {
	return &di.Type{Import: importPath, Type: dataType}
}

func Load(dir string, options ...Option) (*Generator, error) {
	appFS := vfs.OS(dir)
	option := &option{
		Embed:    false,
		Hot:      true,
		Minify:   false,
		ModCache: modcache.Default(),
		FSCache:  nil,
	}
	for _, fn := range options {
		fn(option)
	}
	// Find go.mod
	module, err := gomod.Find(dir, gomod.WithModCache(option.ModCache), gomod.WithFSCache(option.FSCache))
	if err != nil {
		if !errors.Is(err, gomod.ErrFileNotFound) {
			return nil, err
		}
		module, err = gomod.Parse(dir, []byte(`module app.com`), gomod.WithModCache(option.ModCache), gomod.WithFSCache(option.FSCache))
		if err != nil {
			return nil, err
		}
	}
	// Load the bud filesystem
	bfs, err := budfs.Load(module, budfs.WithFSCache(option.FSCache))
	if err != nil {
		return nil, err
	}
	parser := parser.New(bfs, module)
	injector := di.New(bfs, module, parser, di.Map{
		toType("gitlab.com/mnm/bud/bfs", "FS"):        toType("gitlab.com/mnm/bud/pkg/gen", "*FileSystem"),
		toType("gitlab.com/mnm/bud/pkg/gen", "FS"):    toType("gitlab.com/mnm/bud/pkg/gen", "*FileSystem"),
		toType("gitlab.com/mnm/bud/pkg/js", "VM"):     toType("gitlab.com/mnm/bud/pkg/js/v8client", "*Client"),
		toType("gitlab.com/mnm/bud/runtime/view", "Renderer"): toType("gitlab.com/mnm/bud/runtime/view", "*Server"),
	})

	// go.mod generator
	// bfs.Entry("go.mod", gen.FileGenerator(&gomod.Generator{
	// 	Module: module,
	// }))

	// generate generator
	// bfs.Entry("bud/generate/main.go", gen.FileGenerator(&generate.Generator{
	// 	BFS:    bfs,
	// 	Module: module,
	// 	Embed:  option.Embed,
	// 	Hot:    option.Hot,
	// 	Minify: option.Minify,
	// }))

	// TODO: separate the following from the generators to give the generators
	// a chance to add files that are picked up by these compiler plugins.
	bfs.Entry("bud/command/command.go", gen.FileGenerator(&command.Generator{
		BFS:    bfs,
		Module: module,
		Parser: parser,
	}))

	// action generator
	bfs.Entry("bud/action/action.go", gen.FileGenerator(&action.Generator{
		BFS:      bfs,
		Injector: injector,
		Module:   module,
		Parser:   parser,
	}))

	// transform generator
	bfs.Entry("bud/transform/transform.go", gen.FileGenerator(&transform.Generator{
		BFS:    bfs,
		Module: module,
	}))

	bfs.Entry("bud/view/view.go", gen.FileGenerator(&view.Generator{
		BFS:    bfs,
		Module: module,
	}))

	bfs.Entry("bud/public/public.go", gen.FileGenerator(&public.Generator{
		BFS:    bfs,
		Module: module,
		Embed:  option.Embed,
		Minify: option.Minify,
	}))

	bfs.Entry("bud/web/web.go", gen.FileGenerator(&web.Generator{
		BFS:    bfs,
		Module: module,
		Parser: parser,
	}))

	bfs.Entry("bud/program/program.go", gen.FileGenerator(&program.Generator{
		BFS:      bfs,
		Module:   module,
		Injector: injector,
	}))

	bfs.Entry("bud/main.go", gen.FileGenerator(&maingo.Generator{
		BFS:    bfs,
		Module: module,
	}))

	return &Generator{appFS, bfs, module}, nil
}

type Generator struct {
	appFS  vfs.ReadWritable
	bfs    budfs.FS
	module *gomod.Module
}

func (g *Generator) Module() *gomod.Module {
	return g.module
}

func (g *Generator) Generate(ctx context.Context) error {
	skipOption := fsync.WithSkip(
		gitignore.New(g.appFS),
		// Don't delete files that were pre-generated.
		func(name string, isDir bool) bool {
			return isDir && (name == "bud/generate" || name == "bud/generator")
		},
	)
	// Sync bud
	if err := fsync.Dir(vfs.SingleFlight(g.bfs), "bud", g.appFS, "bud", skipOption); err != nil {
		return err
	}
	return nil
}
