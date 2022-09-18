package transform_test

import (
	"context"
	"io/fs"
	"os"
	"testing"

	transform "github.com/livebud/bud/framework/transform2"
	"github.com/livebud/bud/package/di"
	"github.com/livebud/bud/package/log"

	"github.com/livebud/bud/package/parser"

	"github.com/livebud/bud/framework"

	"github.com/livebud/bud/internal/is"
	"github.com/livebud/bud/internal/testdir"
	"github.com/livebud/bud/package/budfs"
	"github.com/livebud/bud/package/gomod"
	"github.com/livebud/bud/package/log/testlog"
)

func loadBFS(log log.Interface, module *gomod.Module) *budfs.FileSystem {
	bfs := budfs.New(module, log)
	flag := &framework.Flag{
		Embed:  false,
		Minify: false,
		Hot:    true,
		Stdin:  os.Stdin,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
		Env:    os.Environ(),
	}
	parser := parser.New(bfs, module)
	injector := di.New(bfs, log, module, parser)
	bfs.FileServer("bud/transform", transform.New(flag, injector, log, module, parser))
	return bfs
}

func TestEmpty(t *testing.T) {
	t.Skip()
}

func TestSvelteToSvelte(t *testing.T) {
	is := is.New(t)
	log := testlog.New()
	ctx := context.Background()
	dir := t.TempDir()
	td := testdir.New(dir)
	td.Files["view/index.svelte"] = `<h1>hello</h1>`
	td.Files["transform/doubler/transform.go"] = `
		package doubler
		import "github.com/livebud/bud/framework/transform2/transformrt"
		type Transform struct {}
		func (t *Transform) SvelteToSvelte(file *transformrt.File) error {
			file.Data = append(file.Data, file.Data...)
			return nil
		}
	`
	is.NoErr(td.Write(ctx))
	module, err := gomod.Find(dir)
	is.NoErr(err)
	bfs := loadBFS(log, module)
	defer bfs.Close()
	code, err := fs.ReadFile(bfs, "bud/transform/view/index.svelte")
	is.NoErr(err)
	is.Equal(string(code), `<h1>hello</h1><h1>hello</h1>`)
}

func TestSvelteToSvelteToJSX(t *testing.T) {
	is := is.New(t)
	log := testlog.New()
	ctx := context.Background()
	dir := t.TempDir()
	td := testdir.New(dir)
	td.Files["view/index.svelte"] = `<h1>hello</h1>`
	td.Files["transform/doubler/transform.go"] = `
		package doubler
		import "github.com/livebud/bud/framework/transform2/transformrt"
		type Transform struct {}
		func (t *Transform) SvelteToSvelte(file *transformrt.File) error {
			file.Data = append(file.Data, file.Data...)
			return nil
		}
	`
	td.Files["transform/jsx/transform.go"] = `
		package jsx
		import "github.com/livebud/bud/framework/transform2/transformrt"
		type Transform struct {}
		func (t *Transform) SvelteToJsx(file *transformrt.File) error {
			file.Data = []byte("export default function() { return <>" + string(file.Data) + "</> }")
			return nil
		}
	`
	is.NoErr(td.Write(ctx))
	module, err := gomod.Find(dir)
	is.NoErr(err)
	bfs := loadBFS(log, module)
	defer bfs.Close()
	code, err := fs.ReadFile(bfs, "bud/transform/view/index.svelte..jsx")
	is.NoErr(err)
	is.Equal(string(code), `export default function() { return <><h1>hello</h1><h1>hello</h1></> }`)
}

func TestFaviconToFavicon(t *testing.T) {
	is := is.New(t)
	log := testlog.New()
	ctx := context.Background()
	dir := t.TempDir()
	td := testdir.New(dir)
	td.BFiles["public/favicon.ico"] = []byte{0x01, 0x02, 0x03}
	td.Files["transform/doubler/transform.go"] = `
		package doubler
		import "github.com/livebud/bud/framework/transform2/transformrt"
		type Transform struct {}
		func (t *Transform) SvelteToSvelte(file *transformrt.File) error {
			file.Data = append(file.Data, file.Data...)
			return nil
		}
		func (t *Transform) IcoToIco(file *transformrt.File) error {
			file.Data = append(file.Data, file.Data...)
			return nil
		}
	`
	is.NoErr(td.Write(ctx))
	module, err := gomod.Find(dir)
	is.NoErr(err)
	bfs := loadBFS(log, module)
	defer bfs.Close()
	code, err := fs.ReadFile(bfs, "bud/transform/public/favicon.ico")
	is.NoErr(err)
	is.Equal(code, []byte{0x01, 0x02, 0x03, 0x01, 0x02, 0x03})
}

func TestMdToSsrJsAndDomJs(t *testing.T) {
	is := is.New(t)
	log := testlog.New()
	ctx := context.Background()
	dir := t.TempDir()
	td := testdir.New(dir)
	td.Files["view/index.md"] = `# hello`
	td.Files["transform/svelte/transform.go"] = `
		package svelte
		import "github.com/livebud/bud/framework/transform2/transformrt"
		type Transform struct {}
		func (t *Transform) SvelteToSsrJs(file *transformrt.File) error {
			file.Data = []byte("module.exports = function() { return " + string(file.Data) + " }")
			return nil
		}
		func (t *Transform) SvelteToDomJs(file *transformrt.File) error {
			file.Data = []byte("export default function() { return " + string(file.Data) + " }")
			return nil
		}
	`
	td.Files["transform/tailwind/transform.go"] = `
		package tailwind
		import "github.com/livebud/bud/framework/transform2/transformrt"
		import "bytes"
		type Transform struct {}
		func (t *Transform) SvelteToSvelte(file *transformrt.File) error {
			file.Data = bytes.Replace(file.Data, []byte("class='bg-red-100'"), []byte("style='color: red'"), -1)
			return nil
		}
	`
	td.Files["transform/markdoc/transform.go"] = `
		package markdoc
		import "github.com/livebud/bud/framework/transform2/transformrt"
		import "bytes"
		type Transform struct {}
		func (t *Transform) MdToSvelte(file *transformrt.File) error {
			file.Data = bytes.TrimPrefix(file.Data, []byte("# "))
			file.Data = []byte("<h1 class='bg-red-100'>" + string(file.Data) + "</h1>")
			return nil
		}
	`
	is.NoErr(td.Write(ctx))
	module, err := gomod.Find(dir)
	is.NoErr(err)
	bfs := loadBFS(log, module)
	defer bfs.Close()
	code, err := fs.ReadFile(bfs, "bud/transform/view/index.md..ssr.js")
	is.NoErr(err)
	is.Equal(string(code), `module.exports = function() { return <h1 style='color: red'>hello</h1> }`)
	code, err = fs.ReadFile(bfs, "bud/transform/view/index.md..dom.js")
	is.NoErr(err)
	is.Equal(string(code), `export default function() { return <h1 style='color: red'>hello</h1> }`)
}
