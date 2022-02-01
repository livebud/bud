package generator

import (
	"context"
	"io/fs"

	"gitlab.com/mnm/bud/2/genfs"

	"gitlab.com/mnm/bud/2/pluginfs"

	"gitlab.com/mnm/bud/2/cachefs"
	"gitlab.com/mnm/bud/2/singleflight"

	"gitlab.com/mnm/bud/fsync"

	"gitlab.com/mnm/bud/2/mod"
	"gitlab.com/mnm/bud/internal/di"
	"gitlab.com/mnm/bud/internal/generator/gomod"
	"gitlab.com/mnm/bud/pkg/modcache"
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

func Load(dir string, options ...Option) (*Generator, error) {
	option := &option{
		Embed:  false,
		Hot:    true,
		Minify: false,
		Cache:  modcache.Default(),
	}
	for _, fn := range options {
		fn(option)
	}
	module, err := mod.Find(dir, mod.WithCache(option.Cache))
	if err != nil {
		return nil, err
	}
	store := cachefs.Cache()
	loader := singleflight.New()
	cache1 := cachefs.New(module, loader, store)
	plugin, err := pluginfs.Load(cache1, module)
	if err != nil {
		return nil, err
	}
	cache2 := cachefs.New(plugin, loader, store)
	genfs := genfs.New(cache2)
	cache3 := cachefs.New(genfs, loader, store)

	// genFS := gen.New(fsys)
	// modFinder := mod.New(mod.WithFS(vfs.SingleFlight(vfs.GitIgnore(genFS))), mod.WithCache(option.Cache))
	// module, err := modFinder.Find(".")
	// if err != nil {
	// 	if !errors.Is(err, mod.ErrCantInfer) {
	// 		return nil, err
	// 	}
	// 	module, err = virtualModule(modFinder)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// }
	// parser := parser.New(module)
	// injector := di.New(module, parser, di.Map{
	// 	toType("gitlab.com/mnm/bud/pkg/gen", "FS"):        toType("gitlab.com/mnm/bud/pkg/gen", "*FileSystem"),
	// 	toType("gitlab.com/mnm/bud/pkg/js", "VM"):         toType("gitlab.com/mnm/bud/pkg/js/v8client", "*Client"),
	// 	toType("gitlab.com/mnm/bud/runtime/view", "Renderer"): toType("gitlab.com/mnm/bud/runtime/view", "*Server"),
	// })
	// _ = injector
	// genFS.Add(map[string]gen.Generator{
	// 	"go.mod": gen.FileGenerator(&gomod.Generator{
	// 		FS: appFS,
	// 		Go: &gomod.Go{
	// 			Version: "1.17",
	// 		},
	// 		ModFinder: modFinder,
	// 		Requires: []*gomod.Require{
	// 			{
	// 				Path:    `gitlab.com/mnm/bud`,
	// 				Version: `v0.0.0-20211017185247-da18ff96a31f`,
	// 			},
	// 		},
	// 		Replaces: option.Replaces,
	// 	}),
	// 	"bud/plugin": gen.DirGenerator(&plugin.Generator{
	// 		Module: module,
	// 	}),
	// 	"bud/generate/main.go": gen.FileGenerator(&generate.Generator{
	// 		Module: module,
	// 		Embed:  option.Embed,
	// 		Hot:    option.Hot,
	// 		Minify: option.Minify,
	// 	}),
	// 	// TODO: separate the following from the generators to give the generators
	// 	// a chance to add files that are picked up by these compiler plugins.
	// 	"bud/command/command.go": gen.FileGenerator(&command.Generator{
	// 		Module: module,
	// 		Parser: parser,
	// 	}),
	// 	"bud/action/action.go": gen.FileGenerator(&action.Generator{
	// 		Injector: injector,
	// 		Module:   module,
	// 		Parser:   parser,
	// 	}),
	// 	"bud/transform/transform.go": gen.FileGenerator(&transform.Generator{
	// 		Module: module,
	// 	}),
	// 	"bud/view/view.go": gen.FileGenerator(&view.Generator{
	// 		Module: module,
	// 	}),
	// 	"bud/public/public.go": gen.FileGenerator(&public.Generator{
	// 		Module: module,
	// 		Embed:  option.Embed,
	// 		Minify: option.Minify,
	// 	}),
	// 	"bud/web/web.go": gen.FileGenerator(&web.Generator{
	// 		Module: module,
	// 		Parser: parser,
	// 	}),
	// 	"bud/program/program.go": gen.FileGenerator(&program.Generator{
	// 		Module:   module,
	// 		Injector: injector,
	// 	}),
	// 	"bud/main.go": gen.FileGenerator(&maingo.Generator{
	// 		Module: module,
	// 	}),
	// })
	return &Generator{fsys}, nil
}

type Generator struct {
	fsys fs.FS
}

func (g *Generator) Generate(ctx context.Context, fsys vfs.ReadWritable) error {
	if err := fsync.Dir(g.fsys, ".", fsys, "."); err != nil {
		return err
	}
	return nil
}
