package plugin_test

import (
	"context"
	"io/fs"
	"os"
	"strings"
	"testing"

	"github.com/matryer/is"
	"gitlab.com/mnm/bud/gen"
	"gitlab.com/mnm/bud/go/mod"
	"gitlab.com/mnm/bud/internal/gobin"
	"gitlab.com/mnm/bud/internal/modcache"
	"gitlab.com/mnm/bud/plugin"
	"gitlab.com/mnm/bud/vfs"
)

func TestPlugin(t *testing.T) {
	is := is.New(t)
	dir := "_tmp"
	is.NoErr(os.RemoveAll(dir))
	is.NoErr(os.MkdirAll(dir, 0755))
	defer func() {
		if !t.Failed() {
			is.NoErr(os.RemoveAll(dir))
		}
	}()
	err := vfs.WriteTo(dir, vfs.Map{
		"go.mod": `module test.mod`,
	})
	is.NoErr(err)
	ctx := context.Background()
	err = gobin.Get(ctx, dir, "gitlab.com/mnm/testdata/bud-tailwind", "gitlab.com/mnm/testdata/bud-markdown")
	is.NoErr(err)
	modfile, err := mod.FindIn(modcache.Default(), dir)
	is.NoErr(err)
	dirfs := os.DirFS(dir)
	bf := gen.New(dirfs)
	bf.Add(map[string]gen.Generator{
		"bud/plugin": gen.DirGenerator(&plugin.Generator{Modfile: modfile}),
	})
	fis, err := fs.ReadDir(bf, "bud/plugin")
	is.NoErr(err)
	is.Equal(len(fis), 2) // expected 2 plugins
	is.Equal(fis[0].Name(), "bud-markdown")
	is.Equal(fis[0].IsDir(), true)
	is.Equal(fis[1].Name(), "bud-tailwind")
	is.Equal(fis[1].IsDir(), true)

	// Markdown
	fis, err = fs.ReadDir(bf, "bud/plugin/bud-markdown")
	is.NoErr(err)
	is.Equal(len(fis), 1) // expected 1 markdown dir
	is.Equal(fis[0].Name(), "transform")
	is.Equal(fis[0].IsDir(), true)
	fis, err = fs.ReadDir(bf, "bud/plugin/bud-markdown/transform")
	is.NoErr(err)
	is.Equal(len(fis), 1) // expected 1 markdown dir
	is.Equal(fis[0].Name(), "markdown")
	is.Equal(fis[0].IsDir(), true)
	fis, err = fs.ReadDir(bf, "bud/plugin/bud-markdown/transform/markdown")
	is.NoErr(err)
	is.Equal(len(fis), 1) // expected 1 markdown file
	is.Equal(fis[0].Name(), "markdown.go")
	is.Equal(fis[0].IsDir(), false)
	data, err := fs.ReadFile(bf, "bud/plugin/bud-markdown/transform/markdown/markdown.go")
	is.NoErr(err)
	is.True(strings.Contains(string(data), `MDToSvelte`))

	// Tailwind
	fis, err = fs.ReadDir(bf, "bud/plugin/bud-tailwind")
	is.NoErr(err)
	is.Equal(len(fis), 3) // expected 2 tailwind dirs
	is.Equal(fis[0].Name(), "internal")
	is.Equal(fis[0].IsDir(), true)
	is.Equal(fis[1].Name(), "public")
	is.Equal(fis[1].IsDir(), true)
	is.Equal(fis[2].Name(), "transform")
	is.Equal(fis[2].IsDir(), true)
	// Public
	fis, err = fs.ReadDir(bf, "bud/plugin/bud-tailwind/public")
	is.NoErr(err)
	is.Equal(len(fis), 1) // expected 1 markdown file
	is.Equal(fis[0].Name(), "tailwind")
	is.Equal(fis[0].IsDir(), true)
	fis, err = fs.ReadDir(bf, "bud/plugin/bud-tailwind/public/tailwind")
	is.NoErr(err)
	is.Equal(len(fis), 1) // expected 1 markdown file
	is.Equal(fis[0].Name(), "preflight.css")
	is.Equal(fis[0].IsDir(), false)
	data, err = fs.ReadFile(bf, "bud/plugin/bud-tailwind/public/tailwind/preflight.css")
	is.NoErr(err)
	is.True(strings.Contains(string(data), `/* Preflight.css */`))
	// Transform
	fis, err = fs.ReadDir(bf, "bud/plugin/bud-tailwind/transform")
	is.NoErr(err)
	is.Equal(len(fis), 1) // expected 1 markdown file
	is.Equal(fis[0].Name(), "tailwind")
	is.Equal(fis[0].IsDir(), true)
	fis, err = fs.ReadDir(bf, "bud/plugin/bud-tailwind/transform/tailwind")
	is.NoErr(err)
	is.Equal(len(fis), 1) // expected 1 markdown file
	is.Equal(fis[0].Name(), "tailwind.go")
	is.Equal(fis[0].IsDir(), false)
	data, err = fs.ReadFile(bf, "bud/plugin/bud-tailwind/transform/tailwind/tailwind.go")
	is.NoErr(err)
	is.True(strings.Contains(string(data), `SvelteToSvelte`))
}
