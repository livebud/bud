package transpiler_test

import (
	"context"
	"os"
	"testing"

	"github.com/livebud/bud/internal/is"
	"github.com/livebud/bud/internal/testcli"
	"github.com/livebud/bud/internal/testdir"
	"github.com/livebud/bud/internal/versions"
)

func TestSvelteToSvelteTranspiler(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()
	dir := t.TempDir()
	td := testdir.New(dir)
	td.NodeModules["svelte"] = versions.Svelte
	td.NodeModules["livebud"] = "*"
	td.Files["view/index.svelte"] = `<h1>hello</h1>`
	td.Files["transpiler/doubler/doubler.go"] = `
		package doubler
		import "github.com/livebud/bud/runtime/transpiler"
		type Transpiler struct {}
		func (t *Transpiler) SvelteToSvelte(file *transpiler.File) error {
			file.Data = append(file.Data, file.Data...)
			return nil
		}
	`
	// Since transpiler is used internally within generators, we create a custom
	// test generator to access the transpiler.
	td.Files["generator/svelte/svelte.go"] = `
		package svelte
		import (
			"github.com/livebud/bud/package/genfs"
			"github.com/livebud/bud/runtime/transpiler"
		)
		type Generator struct {}
		func (g *Generator) Generate(fsys genfs.FS, dir *genfs.Dir) error {
			dir.GenerateFile("svelte.txt", func(fsys genfs.FS, file *genfs.File) error {
				data, err := transpiler.TranspileFile(fsys, "view/index.svelte", ".svelte")
				if err != nil {
					return err
				}
				file.Data = data
				return nil
			})
			return nil
		}
	`
	is.NoErr(td.Write(ctx))
	cli := testcli.New(dir)
	_, err := cli.Run(ctx, "build", "--embed=false")
	is.NoErr(err)
	is.NoErr(td.Exists("bud/internal/svelte/svelte.txt"))
	data, err := os.ReadFile(td.Path("bud/internal/svelte/svelte.txt"))
	is.NoErr(err)
	is.Equal(string(data), `<h1>hello</h1><h1>hello</h1>`)
}
