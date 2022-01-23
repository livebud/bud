package gen_test

import (
	"errors"
	"io/fs"
	"testing"
	"time"

	"github.com/matryer/is"
	"gitlab.com/mnm/bud/gen"
	"gitlab.com/mnm/bud/vfs"
)

func TestOpen(t *testing.T) {
	is := is.New(t)
	now := time.Now()
	efs := gen.EFS{
		"duo/view/index.svelte": &gen.Embed{
			Data:    []byte(`<h1>index</h1>`),
			Mode:    fs.FileMode(0644),
			ModTime: now,
		},
		"duo/view/about/about.svelte": &gen.Embed{
			Data:    []byte(`<h1>about</h1>`),
			Mode:    fs.FileMode(0644),
			ModTime: now,
		},
		"duo/public/favicon.ico": &gen.Embed{
			Data:    []byte(`favicon.ico`),
			Mode:    fs.FileMode(0644),
			ModTime: now,
		},
	}
	// duo/view/index.svelte
	code, err := fs.ReadFile(efs, "duo/view/index.svelte")
	is.NoErr(err)
	is.Equal(string(code), `<h1>index</h1>`)
	stat, err := fs.Stat(efs, "duo/view/index.svelte")
	is.NoErr(err)
	is.Equal(stat.ModTime(), now)
	is.Equal(stat.Mode(), fs.FileMode(0644))
	is.Equal(stat.IsDir(), false)

	// duo/view/about/about.svelte
	code, err = fs.ReadFile(efs, "duo/view/about/about.svelte")
	is.NoErr(err)
	is.Equal(string(code), `<h1>about</h1>`)
	stat, err = fs.Stat(efs, "duo/view/about/about.svelte")
	is.NoErr(err)
	is.Equal(stat.ModTime(), now)
	is.Equal(stat.Mode(), fs.FileMode(0644))
	is.Equal(stat.IsDir(), false)

	// duo/public/favicon.ico
	code, err = fs.ReadFile(efs, "duo/public/favicon.ico")
	is.NoErr(err)
	is.Equal(string(code), `favicon.ico`)
	stat, err = fs.Stat(efs, "duo/public/favicon.ico")
	is.NoErr(err)
	is.Equal(stat.ModTime(), now)
	is.Equal(stat.Mode(), fs.FileMode(0644))
	is.Equal(stat.IsDir(), false)

	// duo/public
	des, err := fs.ReadDir(efs, "duo/public")
	is.Equal(errors.Is(err, fs.ErrNotExist), true)
	is.Equal(des, nil)
}

func TestAdd(t *testing.T) {
	is := is.New(t)
	now := time.Now()
	efs := gen.EFS{
		"duo/view/index.svelte": &gen.Embed{
			Data:    []byte(`<h1>index</h1>`),
			Mode:    fs.FileMode(0644),
			ModTime: now,
		},
	}
	efs.Add(map[string]gen.Generator{
		"duo/view/about/about.svelte": &gen.Embed{
			Data:    []byte(`<h1>about</h1>`),
			Mode:    fs.FileMode(0644),
			ModTime: now,
		},
	})
	efs.Add(map[string]gen.Generator{
		"duo/public/favicon.ico": &gen.Embed{
			Data:    []byte(`favicon.ico`),
			Mode:    fs.FileMode(0644),
			ModTime: now,
		},
	})

	// duo/view/index.svelte
	code, err := fs.ReadFile(efs, "duo/view/index.svelte")
	is.NoErr(err)
	is.Equal(string(code), `<h1>index</h1>`)
	stat, err := fs.Stat(efs, "duo/view/index.svelte")
	is.NoErr(err)
	is.Equal(stat.ModTime(), now)
	is.Equal(stat.Mode(), fs.FileMode(0644))
	is.Equal(stat.IsDir(), false)

	// duo/view/about/about.svelte
	code, err = fs.ReadFile(efs, "duo/view/about/about.svelte")
	is.NoErr(err)
	is.Equal(string(code), `<h1>about</h1>`)
	stat, err = fs.Stat(efs, "duo/view/about/about.svelte")
	is.NoErr(err)
	is.Equal(stat.ModTime(), now)
	is.Equal(stat.Mode(), fs.FileMode(0644))
	is.Equal(stat.IsDir(), false)

	// duo/public/favicon.ico
	code, err = fs.ReadFile(efs, "duo/public/favicon.ico")
	is.NoErr(err)
	is.Equal(string(code), `favicon.ico`)
	stat, err = fs.Stat(efs, "duo/public/favicon.ico")
	is.NoErr(err)
	is.Equal(stat.ModTime(), now)
	is.Equal(stat.Mode(), fs.FileMode(0644))
	is.Equal(stat.IsDir(), false)

	// duo/public
	des, err := fs.ReadDir(efs, "duo/public")
	is.Equal(errors.Is(err, fs.ErrNotExist), true)
	is.Equal(des, nil)
}

func TestGoModGoMod(t *testing.T) {
	is := is.New(t)
	fsys := vfs.Map{
		"app.go": []byte("package app\nimport \"mod.test/module\"\nvar a = module.Answer"),
	}
	genfs := gen.New(fsys)
	genfs.Add(map[string]gen.Generator{
		"go.mod": gen.GenerateFile(func(f gen.F, file *gen.File) error {
			file.Write([]byte("module app.com\nrequire mod.test/module v1.2.4"))
			return nil
		}),
	})
	stat, err := fs.Stat(genfs, "go.mod/go.mod")
	is.True(err != nil)
	is.True(errors.Is(err, fs.ErrNotExist))
	is.Equal(stat, nil)
	stat, err = fs.Stat(genfs, "go.mod")
	is.NoErr(err)
	is.Equal(stat.Name(), "go.mod")
}

// TODO: support passing embeds into gen.Generator
func TestGoModGoModEmbed(t *testing.T) {
	t.SkipNow()
	is := is.New(t)
	fsys := vfs.Map{
		"app.go": []byte("package app\nimport \"mod.test/module\"\nvar a = module.Answer"),
	}
	genfs := gen.New(fsys)
	genfs.Add(map[string]gen.Generator{
		"go.mod": &gen.Embed{Data: []byte("module app.com\nrequire mod.test/module v1.2.4")},
	})
	stat, err := fs.Stat(genfs, "go.mod/go.mod")
	is.True(err != nil)
	is.True(errors.Is(err, fs.ErrNotExist))
	is.Equal(stat, nil)
	stat, err = fs.Stat(genfs, "go.mod")
	is.NoErr(err)
	is.Equal(stat.Name(), "go.mod")
}
