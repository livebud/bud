package gen_test

import (
	"errors"
	"io/fs"
	"testing"
	"time"

	"github.com/matryer/is"
	"gitlab.com/mnm/bud/gen"
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
