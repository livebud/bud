package transpiler_test

import (
	"context"
	"os"
	"path"
	"testing"

	"github.com/livebud/bud/internal/is"
	"github.com/livebud/bud/internal/testcli"
	"github.com/livebud/bud/internal/testdir"
	"github.com/livebud/bud/internal/versions"
)

type Transpile struct {
	FromPath string // e.g. view/index.svelte
	ToExt    string // e.g. svelte, jsx, etc. (no dot)
}

// Since transpiler is used internally within generators, we create a custom
// test generator to access the transpiler.
func addFiles(td *testdir.Dir, transpiles []Transpile) {
	for _, t := range transpiles {
		toPath := path.Join("generator", t.ToExt, t.ToExt+".go")
		td.Files[toPath] = `
			package svelte
			import (
				"github.com/livebud/bud/package/genfs"
				"github.com/livebud/bud/runtime/transpiler"
			)
			type Generator struct {}
			func (g *Generator) Generate(fsys genfs.FS, dir *genfs.Dir) error {
				dir.GenerateFile(` + "`" + t.FromPath + "`" + `, func(fsys genfs.FS, file *genfs.File) error {
					data, err := transpiler.TranspileFile(fsys, ` + "`" + t.FromPath + "`" + `, ".` + t.ToExt + `")
					if err != nil {
						return err
					}
					file.Data = data
					return nil
				})
				return nil
			}
		`
	}
}

func TestEmpty(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()
	dir := t.TempDir()
	td := testdir.New(dir)
	is.NoErr(td.Write(ctx))
	cli := testcli.New(dir)
	_, err := cli.Run(ctx, "build", "--embed=false")
	is.NoErr(err)
}

func TestSvelteToSvelte(t *testing.T) {
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
	addFiles(td, []Transpile{
		{"view/index.svelte", "svelte"},
	})
	is.NoErr(td.Write(ctx))
	cli := testcli.New(dir)
	_, err := cli.Run(ctx, "build", "--embed=false")
	is.NoErr(err)
	is.NoErr(td.Exists("bud/internal/svelte/view/index.svelte"))
	data, err := os.ReadFile(td.Path("bud/internal/svelte/view/index.svelte"))
	is.NoErr(err)
	is.Equal(string(data), `<h1>hello</h1><h1>hello</h1>`)
}

func TestSvelteToSvelteToJSX(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()
	dir := t.TempDir()
	td := testdir.New(dir)
	td.NodeModules["svelte"] = versions.Svelte
	td.NodeModules["livebud"] = "*"
	td.Files["view/index.svelte"] = `<h1>hello</h1>`
	td.Files["transpiler/doubler/transpiler.go"] = `
		package doubler
		import "github.com/livebud/bud/runtime/transpiler"
		type Transpiler struct {}
		func (t *Transpiler) SvelteToSvelte(file *transpiler.File) error {
			file.Data = append(file.Data, file.Data...)
			return nil
		}
	`
	td.Files["transpiler/jsx/transpiler.go"] = `
		package jsx
		import "github.com/livebud/bud/runtime/transpiler"
		type Transpiler struct {}
		func (t *Transpiler) SvelteToJsx(file *transpiler.File) error {
			file.Data = []byte("export default function() { return <>" + string(file.Data) + "</> }")
			return nil
		}
	`
	addFiles(td, []Transpile{
		{"view/index.svelte", "jsx"},
	})
	is.NoErr(td.Write(ctx))
	cli := testcli.New(dir)
	_, err := cli.Run(ctx, "build", "--embed=false")
	is.NoErr(err)
	is.NoErr(td.Exists("bud/internal/jsx/view/index.svelte"))
	data, err := os.ReadFile(td.Path("bud/internal/jsx/view/index.svelte"))
	is.NoErr(err)
	is.Equal(string(data), `export default function() { return <><h1>hello</h1><h1>hello</h1></> }`)
}

func TestFaviconToFavicon(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()
	dir := t.TempDir()
	td := testdir.New(dir)
	td.BFiles["public/favicon.ico"] = []byte{0x01, 0x02, 0x03}
	td.Files["transpiler/doubler/transpiler.go"] = `
		package doubler
		import "github.com/livebud/bud/runtime/transpiler"
		type Transpiler struct {}
		func (t *Transpiler) SvelteToSvelte(file *transpiler.File) error {
			file.Data = append(file.Data, file.Data...)
			return nil
		}
		func (t *Transpiler) IcoToIco(file *transpiler.File) error {
			file.Data = append(file.Data, file.Data...)
			return nil
		}
	`
	addFiles(td, []Transpile{
		{"public/favicon.ico", "ico"},
	})
	is.NoErr(td.Write(ctx))
	cli := testcli.New(dir)
	_, err := cli.Run(ctx, "build", "--embed=false")
	is.NoErr(err)
	is.NoErr(td.Exists("bud/internal/ico/public/favicon.ico"))
	data, err := os.ReadFile(td.Path("bud/internal/ico/public/favicon.ico"))
	is.NoErr(err)
	is.Equal(data, []byte{0x01, 0x02, 0x03, 0x01, 0x02, 0x03})
}

func TestMdToSsrJsAndDomJs(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()
	dir := t.TempDir()
	td := testdir.New(dir)
	td.NodeModules["svelte"] = versions.Svelte
	td.NodeModules["livebud"] = "*"
	td.Files["view/index.md"] = `# hello`
	td.Files["transpiler/svelte/transpiler.go"] = `
		package svelte
		import "github.com/livebud/bud/runtime/transpiler"
		type Transpiler struct {}
		func (t *Transpiler) SvelteToSsrJs(file *transpiler.File) error {
			file.Data = []byte("module.exports = function() { return " + string(file.Data) + " }")
			return nil
		}
		func (t *Transpiler) SvelteToDomJs(file *transpiler.File) error {
			file.Data = []byte("export default function() { return " + string(file.Data) + " }")
			return nil
		}
	`
	td.Files["transpiler/tailwind/transpiler.go"] = `
		package tailwind
		import "github.com/livebud/bud/runtime/transpiler"
		import "bytes"
		type Transpiler struct {}
		func (t *Transpiler) SvelteToSvelte(file *transpiler.File) error {
			file.Data = bytes.Replace(file.Data, []byte("class='bg-red-100'"), []byte("style='color: red'"), -1)
			return nil
		}
	`
	td.Files["transpiler/markdoc/transpiler.go"] = `
		package markdoc
		import "github.com/livebud/bud/runtime/transpiler"
		import "bytes"
		type Markdoc struct {}
		func (m *Markdoc) Compile(data []byte) ([]byte, error) {
			data = bytes.TrimPrefix(data, []byte("# "))
			data = []byte("<h1 class='bg-red-100'>" + string(data) + "</h1>")
			return data, nil
		}
		type Transpiler struct {
			Markdoc *Markdoc
		}
		func (t *Transpiler) MdToSvelte(file *transpiler.File) (err error) {
			file.Data, err = t.Markdoc.Compile(file.Data)
			if err != nil {
				return err
			}
			return nil
		}
	`
	addFiles(td, []Transpile{
		{"view/index.md", "ssr.js"},
		{"view/index.md", "dom.js"},
	})
	is.NoErr(td.Write(ctx))
	cli := testcli.New(dir)
	_, err := cli.Run(ctx, "build", "--embed=false")
	is.NoErr(err)

	// SSR
	is.NoErr(td.Exists("bud/internal/ssr.js/view/index.md"))
	data, err := os.ReadFile(td.Path("bud/internal/ssr.js/view/index.md"))
	is.NoErr(err)
	is.Equal(string(data), `module.exports = function() { return <h1 style='color: red'>hello</h1> }`)

	// DOM
	is.NoErr(td.Exists("bud/internal/dom.js/view/index.md"))
	data, err = os.ReadFile(td.Path("bud/internal/dom.js/view/index.md"))
	is.NoErr(err)
	is.Equal(string(data), `export default function() { return <h1 style='color: red'>hello</h1> }`)
}
