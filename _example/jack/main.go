package main

//go:generate go run . generate

import (
	"context"
	"embed"
	"io/fs"
	"os"

	"github.com/livebud/bud/_example/jack/env"
	"github.com/livebud/bud/_example/jack/generator"
	"github.com/livebud/bud/_example/jack/router"
	"github.com/livebud/bud/pkg/cli"
	"github.com/livebud/bud/pkg/cmds"
	"github.com/livebud/bud/pkg/ldflag"
	"github.com/livebud/bud/pkg/logs"
	"github.com/livebud/bud/pkg/mod"
	"github.com/livebud/bud/pkg/sse"
	"github.com/livebud/bud/pkg/u"
)

// Dev
// - direnv exec ./_example/jack go run -C _example/jack .
//
// Prod
// - go generate && go build -ldflags "-X github.com/livebud/bud/pkg/ldflag.embed=.bud" -o /tmp/main main.go
// - (cd /tmp && ./main)
//
// TODO: consider using overlays instead

// var generate = flag.Bool("generate", false, "generate")

// TODO: support embedding files without needing .bud/ with a file to be present
//
//go:embed .bud/**
var embeddedFS embed.FS

func loadGenerator(env *env.Env, log logs.Log, mod *mod.Module) fs.FS {
	if ldflag.Embed() != "" {
		log.Info("Using embedded assets")
		return u.Must(fs.Sub(embeddedFS, ".bud"))
	}
	return generator.New(env, log, mod)
}

func loadModule(log logs.Log, dir string) *mod.Module {
	if ldflag.Embed() != "" {
		log.Info("Using embedded module")
		return mod.New(".")
	}
	return mod.MustFind(dir)
}

func main() {
	env := u.Must(env.Load())
	log := logs.Default()
	module := loadModule(log, ".")
	fsys := loadGenerator(env, log, module)
	sse := sse.New(log)
	router := router.New(fsys, sse)

	serve := cmds.Serve(log, router)
	watch := cmds.Watch(log, module, sse)
	generate := cmds.Generate(log, fsys, module)
	cmd := cli.New("jack", "jack CLI")
	cmd.Command("watch", "watch for file changes").Add(watch)
	cmd.Command("generate", "generate the app").Add(generate)
	cmd.Command("serve", "serve the app").Add(serve)
	cmd.Add(cli.All(serve, watch))

	ctx := context.Background()
	if err := cmd.Parse(ctx, os.Args[1:]...); err != nil {
		log.Error(err)
		os.Exit(1)
	}
}
