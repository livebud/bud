package main

import (
	"context"
	"io/fs"
	"net/http"
	"os"

	"github.com/livebud/bud/framework"
	"github.com/livebud/bud/internal/current"
	"github.com/livebud/bud/package/es"
	"github.com/livebud/bud/package/gomod"
	"github.com/livebud/bud/package/log"
	"github.com/livebud/bud/package/log/console"
	"github.com/livebud/bud/package/router"
	"github.com/livebud/bud/package/transpiler"
	"github.com/livebud/bud/package/viewer"
	"github.com/livebud/bud/package/viewer/svelte"
	"github.com/livebud/bud/package/virtual"
	"github.com/livebud/js"
	"github.com/livebud/js/goja"
)

func main() {
	if err := run(); err != nil {
		console.Error(err.Error())
		os.Exit(1)
	}
}

func run() error {
	// return serve()
	fsys := virtual.Tree{}
	if err := bundle(fsys); err != nil {
		return err
	}
	return static(fsys)
}

func serve() error {
	dir, err := current.Directory()
	if err != nil {
		return err
	}
	log := log.New(console.New(os.Stderr))
	module := gomod.New(dir)
	pages, err := viewer.Find(module)
	if err != nil {
		return err
	}
	svelte, err := loadViewer(log, module, pages)
	if err != nil {
		return err
	}
	router := router.New()
	if err := svelte.Mount(router); err != nil {
		return err
	}
	for _, page := range pages {
		router.Get(page.Route, svelte.Handler(page))
	}

	log.Info("listening on http://localhost:3000")
	return http.ListenAndServe(":3000", router)
}

func bundle(fsys virtual.Tree) error {
	dir, err := current.Directory()
	if err != nil {
		return err
	}
	log := log.New(console.New(os.Stderr))
	module := gomod.New(dir)
	pages, err := viewer.Find(module)
	if err != nil {
		return err
	}
	svelte, err := loadViewer(log, module, pages)
	if err != nil {
		return err
	}
	ctx := context.Background()
	return svelte.Bundle(ctx, fsys)
}

func static(fsys fs.FS) error {
	dir, err := current.Directory()
	if err != nil {
		return err
	}
	log := log.New(console.New(os.Stderr))
	module := gomod.New(dir)
	pages, err := viewer.Find(module)
	if err != nil {
		return err
	}
	svelte, err := loadStatic(fsys, log, pages)
	if err != nil {
		return err
	}
	router := router.New()
	if err := svelte.Mount(router); err != nil {
		return err
	}
	for _, page := range pages {
		router.Get(page.Route, svelte.Handler(page))
	}
	log.Info("listening on http://localhost:3000")
	return http.ListenAndServe(":3000", router)
}

func loadViewer(log log.Log, module *gomod.Module, pages map[string]*viewer.Page) (*svelte.Viewer, error) {
	js := goja.New(&js.Console{
		Log:   os.Stdout,
		Error: os.Stderr,
	})
	flag := &framework.Flag{}
	esb := es.New(flag, log)
	svelteCompiler, err := svelte.Load(flag, js)
	if err != nil {
		return nil, err
	}
	tr := transpiler.New()
	tr.Add(".svelte", ".ssr.js", func(ctx context.Context, file *transpiler.File) error {
		ssr, err := svelteCompiler.SSR(ctx, file.Path(), file.Data)
		if err != nil {
			return err
		}
		file.Data = []byte(ssr.JS)
		return nil
	})
	tr.Add(".svelte", ".dom.js", func(ctx context.Context, file *transpiler.File) error {
		dom, err := svelteCompiler.DOM(ctx, file.Path(), file.Data)
		if err != nil {
			return err
		}
		file.Data = []byte(dom.JS)
		return nil
	})
	viewer := svelte.New(esb, flag, js, log, module, pages, tr)
	return viewer, nil
}

func loadStatic(fsys fs.FS, log log.Log, pages map[string]*viewer.Page) (*svelte.StaticViewer, error) {
	js := goja.New(&js.Console{
		Log:   os.Stdout,
		Error: os.Stderr,
	})
	flag := &framework.Flag{}
	svelteCompiler, err := svelte.Load(flag, js)
	if err != nil {
		return nil, err
	}
	tr := transpiler.New()
	tr.Add(".svelte", ".ssr.js", func(ctx context.Context, file *transpiler.File) error {
		ssr, err := svelteCompiler.SSR(ctx, file.Path(), file.Data)
		if err != nil {
			return err
		}
		file.Data = []byte(ssr.JS)
		return nil
	})
	tr.Add(".svelte", ".dom.js", func(ctx context.Context, file *transpiler.File) error {
		dom, err := svelteCompiler.DOM(ctx, file.Path(), file.Data)
		if err != nil {
			return err
		}
		file.Data = []byte(dom.JS)
		return nil
	})
	viewer := svelte.Static(fsys, js, log, pages)
	return viewer, nil
}
