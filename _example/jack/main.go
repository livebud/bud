package main

//go:generate go run . -generate

import (
	"context"
	"embed"
	"encoding/json"
	"flag"
	"io"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/livebud/bud/pkg/genfs"
	"github.com/livebud/bud/pkg/genfs/genfscache"
	"github.com/livebud/bud/pkg/logs"
	"github.com/livebud/bud/pkg/mod"
	"github.com/livebud/bud/pkg/mux"
	"github.com/livebud/bud/pkg/request"
	"github.com/livebud/bud/pkg/slots"
	"github.com/livebud/bud/pkg/sse"
	"github.com/livebud/bud/pkg/ssr"
	"github.com/livebud/bud/pkg/u"
	"github.com/livebud/bud/pkg/view"
	"github.com/livebud/bud/pkg/view/css"
	"github.com/livebud/bud/pkg/view/preact"
	"github.com/livebud/bud/pkg/virt"
	"github.com/livebud/bud/pkg/watcher"
	"golang.org/x/sync/errgroup"
)

// Dev
// - direnv exec ./_example/jack go run -C _example/jack .
//
// Prod
// - go generate && go build -ldflags "-X main.embedded=true" -o /tmp/main main.go
// - /tmp/main

type VNode struct {
	Name     string            `json:"name,omitempty"`
	Attrs    map[string]string `json:"attrs,omitempty"`
	Children any               `json:"children,omitempty"`
	Value    string            `json:"value,omitempty"`
}

var generate = flag.Bool("generate", false, "generate")

var embedded string

// TODO: support embedding files without needing .bud/ with a file to be present
//
//go:embed .bud/**
var embeddedFS embed.FS

func loadModule() *mod.Module {
	if embedded == "true" {
		return mod.New(".bud")
	}
	return u.Must(mod.Find())
}

func loadFS(log logs.Log, module *mod.Module, preact *preact.Viewer, css *css.Viewer) fs.FS {
	if embedded == "true" {
		log.Info("Using embedded assets")
		return u.Must(fs.Sub(embeddedFS, ".bud"))
	}
	gfs := genfs.New(genfscache.Discard(), virt.Map{}, log)
	gfs.GenerateFile("view/layout.tsx", func(fsys genfs.FS, file *genfs.File) error {
		outfile, err := preact.CompileSSR("./" + file.Path())
		if err != nil {
			return err
		}
		file.Data = outfile.Contents
		return nil
	})
	gfs.GenerateFile("view/index.tsx", func(fsys genfs.FS, file *genfs.File) error {
		outfile, err := preact.CompileSSR("./" + file.Path())
		if err != nil {
			return err
		}
		file.Data = outfile.Contents
		return nil
	})
	gfs.GenerateFile("view/faq.tsx", func(fsys genfs.FS, file *genfs.File) error {
		outfile, err := preact.CompileSSR("./" + file.Path())
		if err != nil {
			return err
		}
		file.Data = outfile.Contents
		return nil
	})
	gfs.GenerateFile("view/index.tsx.js", func(fsys genfs.FS, file *genfs.File) error {
		outfile, err := preact.CompileDOM("./" + strings.TrimSuffix(file.Path(), ".js"))
		if err != nil {
			return err
		}
		file.Data = outfile.Contents
		return nil
	})
	gfs.GenerateFile("view/faq.tsx.js", func(fsys genfs.FS, file *genfs.File) error {
		outfile, err := preact.CompileDOM("./" + strings.TrimSuffix(file.Path(), ".js"))
		if err != nil {
			return err
		}
		file.Data = outfile.Contents
		return nil
	})
	gfs.GenerateFile("view/layout.css", func(fsys genfs.FS, file *genfs.File) error {
		outfile, err := css.Compile(file.Path())
		if err != nil {
			return err
		}
		file.Data = outfile.Contents
		return nil
	})
	gfs.GenerateFile("view/index.css", func(fsys genfs.FS, file *genfs.File) error {
		outfile, err := css.Compile(file.Path())
		if err != nil {
			return err
		}
		file.Data = outfile.Contents
		return nil
	})
	gfs.GenerateFile("view/faq.css", func(fsys genfs.FS, file *genfs.File) error {
		outfile, err := css.Compile(file.Path())
		if err != nil {
			return err
		}
		file.Data = outfile.Contents
		return nil
	})
	gfs.ServeFile("public", func(fsys genfs.FS, file *genfs.File) error {
		code, err := fs.ReadFile(module, file.Target())
		if err != nil {
			return err
		}
		file.Data = code
		return nil
	})
	gfs.GenerateDir("public", func(fsys genfs.FS, dir *genfs.Dir) error {
		return fs.WalkDir(module, dir.Path(), func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if d.IsDir() {
				return nil
			}
			rel, err := filepath.Rel(dir.Path(), path)
			if err != nil {
				return err
			}
			dir.GenerateFile(rel, func(fsys genfs.FS, file *genfs.File) error {
				data, err := fs.ReadFile(module, path)
				if err != nil {
					return err
				}
				file.Data = data
				return nil
			})
			return nil
		})
	})
	return gfs
}

func main() {
	flag.Parse()
	log := logs.Default()
	module := loadModule()
	router := mux.New()
	se := sse.New(log)
	ctx := context.Background()
	preact := preact.New(module, preact.WithEnv(map[string]any{
		"API_URL":            os.Getenv("API_URL"),
		"SLACK_CLIENT_ID":    os.Getenv("SLACK_CLIENT_ID"),
		"SLACK_REDIRECT_URL": os.Getenv("SLACK_REDIRECT_URL"),
		"SLACK_SCOPE":        os.Getenv("SLACK_SCOPE"),
		"SLACK_USER_SCOPE":   os.Getenv("SLACK_USER_SCOPE"),
		"STRIPE_CLIENT_KEY":  os.Getenv("STRIPE_CLIENT_KEY"),
	}))
	css := css.New(module)
	fsys := loadFS(log, module, preact, css)
	if *generate {
		if err := virt.Sync(log, fsys, module.Sub(".bud")); err != nil {
			log.Error(err)
			os.Exit(1)
		}
		return
	}
	// TODO: disable live reload in production
	ssr := ssr.New(fsys, "/live.js")
	router.Get("/live.js", se)
	router.Layout("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		slot, err := slots.FromContext(r.Context())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		page, err := io.ReadAll(slot)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		heads, err := io.ReadAll(slot.Slot("heads"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if err := ssr.Render(w, "view/layout.tsx", &view.Data{
			Props: map[string]interface{}{
				"page":  string(page),
				"heads": json.RawMessage(heads),
			},
		}); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}))
	router.Get("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var props struct {
			Success *bool `json:"success"`
		}
		if err := request.Unmarshal(r, &props); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		slot, err := slots.FromContext(r.Context())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if err := ssr.Render(slot, "view/index.tsx", &view.Data{Props: props, Slots: slot}); err != nil {
			// TODO: support live reload on error pages
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}))
	router.Get("/faq", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var props struct {
			Success *bool `json:"success"`
		}
		if err := request.Unmarshal(r, &props); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		slot, err := slots.FromContext(r.Context())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if err := ssr.Render(slot, "view/faq.tsx", &view.Data{Props: props, Slots: slot}); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}))
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
