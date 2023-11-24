package main

//go:generate go run . generate asset/.embed

import (
	"context"
	"os"

	"github.com/livebud/bud/_example/jack/asset"
	"github.com/livebud/bud/_example/jack/env"
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
// $ direnv exec ./_example/jack go run -C _example/jack -mod=mod .
//
// Generate
// $ direnv exec ./_example/jack go generate -C _example/jack -mod=mod .
//
// Build
// $ direnv exec ./_example/jack go build -C _example/jack -mod=mod -tags=embed -o main .
// $ direnv exec ./_example/jack ./_example/jack/main

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
	fsys := asset.Load(env, log, module)
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
