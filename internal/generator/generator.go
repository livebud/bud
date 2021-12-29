package generator

import (
	"context"
	"errors"

	"gitlab.com/mnm/bud/fsync"

	"gitlab.com/mnm/bud/gen"
	"gitlab.com/mnm/bud/go/mod"
	"gitlab.com/mnm/bud/internal/di"
	"gitlab.com/mnm/bud/internal/generator/action"
	"gitlab.com/mnm/bud/internal/generator/command"
	"gitlab.com/mnm/bud/internal/generator/generate"
	"gitlab.com/mnm/bud/internal/generator/gomod"
	"gitlab.com/mnm/bud/internal/generator/maingo"
	"gitlab.com/mnm/bud/internal/generator/program"
	"gitlab.com/mnm/bud/internal/generator/public"
	"gitlab.com/mnm/bud/internal/generator/router"
	"gitlab.com/mnm/bud/internal/generator/transform"
	"gitlab.com/mnm/bud/internal/generator/view"
	"gitlab.com/mnm/bud/internal/generator/web"
	"gitlab.com/mnm/bud/internal/modcache"
	"gitlab.com/mnm/bud/internal/parser"
	"gitlab.com/mnm/bud/plugin"
	"gitlab.com/mnm/bud/vfs"
)

type option struct {
	Embed    bool
	Hot      bool
	Minify   bool
	Replaces []*gomod.Replace
	Cache    *modcache.Cache
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
		option.Replaces = append(option.Replaces, &gomod.Replace{Old: from, New: to})
	}
}

func WithCache(mc *modcache.Cache) func(*option) {
	return func(option *option) {
		option.Cache = mc
	}
}

func toType(importPath, dataType string) *di.Type {
	return &di.Type{Import: importPath, Type: dataType}
}

func virtualModule(mf *mod.Finder) (*mod.Module, error) {
	return mf.Parse("go.mod", []byte(`module app.com`))
}

func Load(appFS vfs.ReadWritable, options ...Option) (*Generator, error) {
	option := &option{
		Embed:  false,
		Hot:    true,
		Minify: false,
		Cache:  modcache.Default(),
	}
	for _, fn := range options {
		fn(option)
	}
	genFS := gen.New(appFS)
	modFinder := mod.New(mod.WithFS(vfs.SingleFlight(vfs.GitIgnore(genFS))), mod.WithCache(option.Cache))
	module, err := modFinder.Find(".")
	if err != nil {
		if !errors.Is(err, mod.ErrCantInfer) {
			return nil, err
		}
		module, err = virtualModule(modFinder)
		if err != nil {
			return nil, err
		}
	}
	parser := parser.New(module)
	injector := di.New(module, parser, di.Map{
		toType("gitlab.com/mnm/bud/gen", "FS"):        toType("gitlab.com/mnm/bud/gen", "*FileSystem"),
		toType("gitlab.com/mnm/bud/js", "VM"):         toType("gitlab.com/mnm/bud/js/v8client", "*Client"),
		toType("gitlab.com/mnm/bud/view", "Renderer"): toType("gitlab.com/mnm/bud/view", "*Server"),
	})
	genFS.Add(map[string]gen.Generator{
		"go.mod": gen.FileGenerator(&gomod.Generator{
			FS: appFS,
			Go: &gomod.Go{
				Version: "1.17",
			},
			ModFinder: modFinder,
			Requires: []*gomod.Require{
				{
					Path:    `gitlab.com/mnm/bud`,
					Version: `v0.0.0-20211017185247-da18ff96a31f`,
				},
			},
			Replaces: option.Replaces,
		}),
		"bud/plugin": gen.DirGenerator(&plugin.Generator{
			Module: module,
		}),
		"bud/generate/main.go": gen.FileGenerator(&generate.Generator{
			Module: module,
			Embed:  option.Embed,
			Hot:    option.Hot,
			Minify: option.Minify,
		}),
		// TODO: separate the following from the generators to give the generators
		// a chance to add files that are picked up by these compiler plugins.
		"bud/command/command.go": gen.FileGenerator(&command.Generator{
			Module: module,
			Parser: parser,
		}),
		"bud/action/action.go": gen.FileGenerator(&action.Generator{
			Injector: injector,
			Module:   module,
			Parser:   parser,
		}),
		"bud/transform/transform.go": gen.FileGenerator(&transform.Generator{
			Module: module,
		}),
		"bud/view/view.go": gen.FileGenerator(&view.Generator{
			Module: module,
		}),
		"bud/public/public.go": gen.FileGenerator(&public.Generator{
			Module: module,
			Embed:  option.Embed,
			Minify: option.Minify,
		}),
		"bud/router/router.go": gen.FileGenerator(&router.Generator{
			Module: module,
		}),
		"bud/web/web.go": gen.FileGenerator(&web.Generator{
			Module: module,
		}),
		"bud/program/program.go": gen.FileGenerator(&program.Generator{
			Module:   module,
			Injector: injector,
		}),
		"bud/main.go": gen.FileGenerator(&maingo.Generator{
			Module: module,
		}),
	})
	return &Generator{appFS, genFS, module}, nil
}

type Generator struct {
	appFS  vfs.ReadWritable
	genFS  *gen.FileSystem
	module *mod.Module
}

func (g *Generator) Module() *mod.Module {
	return g.module
}

func (g *Generator) Add(generators map[string]gen.Generator) {
	g.genFS.Add(generators)
}

func (g *Generator) Generate(ctx context.Context) error {
	if err := fsync.Dir(g.module, ".", vfs.GitIgnoreRW(g.appFS), "."); err != nil {
		return err
	}
	return nil
}
