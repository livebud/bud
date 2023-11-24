//go:build !embed

package asset

import (
	"io/fs"
	"path/filepath"
	"strings"

	"github.com/livebud/bud/_example/jack/env"
	"github.com/livebud/bud/pkg/gen"
	"github.com/livebud/bud/pkg/gen/gencache"
	"github.com/livebud/bud/pkg/logs"
	"github.com/livebud/bud/pkg/mod"
	"github.com/livebud/bud/pkg/preact"
	"github.com/livebud/bud/pkg/view/css"
	"github.com/livebud/bud/pkg/virt"
)

func Load(env *env.Env, log logs.Log, module *mod.Module) fs.FS {
	preact := preact.New(module, preact.WithEnv(env))
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
