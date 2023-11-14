package main

//go:generate go run . -generate

import (
	"context"
	"embed"
	"encoding/json"
	"flag"
	"io/fs"
	"net/http"
	"os"

	"github.com/livebud/bud/_example/jack/controller"
	"github.com/livebud/bud/_example/jack/generator"
	"github.com/livebud/bud/pkg/ldflag"
	"github.com/livebud/bud/pkg/logs"
	"github.com/livebud/bud/pkg/mux"
	"github.com/livebud/bud/pkg/sse"
	"github.com/livebud/bud/pkg/ssr"
	"github.com/livebud/bud/pkg/u"
	"github.com/livebud/bud/pkg/virt"
	"github.com/livebud/bud/pkg/watcher"
	"golang.org/x/sync/errgroup"
)

// Dev
// - direnv exec ./_example/jack go run -C _example/jack .
//
// Prod
// - go generate && go build -ldflags "-X github.com/livebud/bud/pkg/ldflag.embed=.bud" -o /tmp/main main.go
// - /tmp/main

var generate = flag.Bool("generate", false, "generate")

// TODO: support embedding files without needing .bud/ with a file to be present
//
//go:embed .bud/**
var embeddedFS embed.FS

func loadFS(log logs.Log) fs.FS {
	if ldflag.Embed() != "" {
		log.Info("Using embedded assets")
		return u.Must(fs.Sub(embeddedFS, ".bud"))
	}
	return generator.New(log)
}

func main() {
	flag.Parse()
	log := logs.Default()
	router := mux.New()
	se := sse.New(log)
	ctx := context.Background()

	fsys := loadFS(log)
	if *generate {
		if err := virt.Sync(log, fsys, virt.OS(".bud")); err != nil {
			log.Error(err)
			os.Exit(1)
		}
		return
	}
	// TODO: disable live reload in production
	ssr := ssr.New(fsys, "/live.js")
	root := controller.New(ssr)
	router.Get("/live.js", se)
	router.Add(root)
	// router.Get("/view/{path*}", http.FileServer(http.FS(gfs)))
	router.Get("/view/{path*}", http.FileServer(http.FS(fsys)))
	// publicFS, err := fs.Sub(gfs, "public")
	publicFS, err := fs.Sub(fsys, "public")
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}
	router.Get("/{path*}", http.FileServer(http.FS(publicFS)))
	eg := new(errgroup.Group)
	eg.Go(func() error {
		log.Infof("Listening on http://localhost%s", ":8080")
		return http.ListenAndServe(":8080", router)
	})
	eg.Go(func() error {
		return watcher.Watch(ctx, ".", func(events *watcher.Events) error {
			eventData, err := json.Marshal(events)
			if err != nil {
				log.Error(err)
			}
			if err := se.Publish(ctx, &sse.Event{Data: []byte(eventData)}); err != nil {
				log.Error(err)
			}
			return nil
		})
	})
	if err := eg.Wait(); err != nil {
		log.Error(err)
		os.Exit(1)
	}
}
