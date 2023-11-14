package generator

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/livebud/bud/pkg/gen"
	"github.com/livebud/bud/pkg/gen/gencache"
	"github.com/livebud/bud/pkg/logs"
	"github.com/livebud/bud/pkg/mod"
	"github.com/livebud/bud/pkg/u"
	"github.com/livebud/bud/pkg/view/css"
	"github.com/livebud/bud/pkg/view/preact"
	"github.com/livebud/bud/pkg/virt"
)

func New(log logs.Log) fs.FS {
	module := u.Must(mod.Find())
	preact := preact.New(module, preact.WithEnv(map[string]any{
		"API_URL":            os.Getenv("API_URL"),
		"SLACK_CLIENT_ID":    os.Getenv("SLACK_CLIENT_ID"),
		"SLACK_REDIRECT_URL": os.Getenv("SLACK_REDIRECT_URL"),
		"SLACK_SCOPE":        os.Getenv("SLACK_SCOPE"),
		"SLACK_USER_SCOPE":   os.Getenv("SLACK_USER_SCOPE"),
		"STRIPE_CLIENT_KEY":  os.Getenv("STRIPE_CLIENT_KEY"),
	}))
	css := css.New(module)
	gfs := gen.New(gencache.Discard(), virt.Map{}, log)
	gfs.GenerateFile("view/layout.tsx", func(fsys gen.FS, file *gen.File) error {
		outfile, err := preact.CompileSSR("./" + file.Path())
		if err != nil {
			return err
		}
		file.Data = outfile.Contents
		return nil
	})
	gfs.GenerateFile("view/index.tsx", func(fsys gen.FS, file *gen.File) error {
		outfile, err := preact.CompileSSR("./" + file.Path())
		if err != nil {
			return err
		}
		file.Data = outfile.Contents
		return nil
	})
	gfs.GenerateFile("view/faq.tsx", func(fsys gen.FS, file *gen.File) error {
		outfile, err := preact.CompileSSR("./" + file.Path())
		if err != nil {
			return err
		}
		file.Data = outfile.Contents
		return nil
	})
	gfs.GenerateFile("view/index.tsx.js", func(fsys gen.FS, file *gen.File) error {
		outfile, err := preact.CompileDOM("./" + strings.TrimSuffix(file.Path(), ".js"))
		if err != nil {
			return err
		}
		file.Data = outfile.Contents
		return nil
	})
	gfs.GenerateFile("view/faq.tsx.js", func(fsys gen.FS, file *gen.File) error {
		outfile, err := preact.CompileDOM("./" + strings.TrimSuffix(file.Path(), ".js"))
		if err != nil {
			return err
		}
		file.Data = outfile.Contents
		return nil
	})
	gfs.GenerateFile("view/layout.css", func(fsys gen.FS, file *gen.File) error {
		outfile, err := css.Compile(file.Path())
		if err != nil {
			return err
		}
		file.Data = outfile.Contents
		return nil
	})
	gfs.GenerateFile("view/index.css", func(fsys gen.FS, file *gen.File) error {
		outfile, err := css.Compile(file.Path())
		if err != nil {
			return err
		}
		file.Data = outfile.Contents
		return nil
	})
	gfs.GenerateFile("view/faq.css", func(fsys gen.FS, file *gen.File) error {
		outfile, err := css.Compile(file.Path())
		if err != nil {
			return err
		}
		file.Data = outfile.Contents
		return nil
	})
	gfs.ServeFile("public", func(fsys gen.FS, file *gen.File) error {
		code, err := fs.ReadFile(module, file.Target())
		if err != nil {
			return err
		}
		file.Data = code
		return nil
	})
	gfs.GenerateDir("public", func(fsys gen.FS, dir *gen.Dir) error {
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
			dir.GenerateFile(rel, func(fsys gen.FS, file *gen.File) error {
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
