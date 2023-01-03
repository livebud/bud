package transform_test

// func loadBFS(log log.Log, module *gomod.Module) (*bfs.FS, error) {
// 	flag := &framework.Flag{
// 		Embed:  false,
// 		Minify: false,
// 		Hot:    true,
// 	}
// 	return bfs.Load(flag, log, module)
// }

// func TestEmpty(t *testing.T) {
// 	is := is.New(t)
// 	ctx := context.Background()
// 	dir := t.TempDir()
// 	td := testdir.New(dir)
// 	is.NoErr(td.Write(ctx))
// 	cli := testcli.New(dir)
// 	_, err := cli.Run(ctx, "build", "--embed=false")
// 	is.NoErr(err)
// 	is.NoErr(td.NotExists("bud/internal/generator/transform"))
// }

// func TestSvelteToSvelte(t *testing.T) {
// 	is := is.New(t)
// 	log := testlog.New()
// 	ctx := context.Background()
// 	dir := t.TempDir()
// 	td := testdir.New(dir)
// 	td.Files["view/index.svelte"] = `<h1>hello</h1>`
// 	td.Files["transform/doubler/transform.go"] = `
// 		package doubler
// 		import "github.com/livebud/bud/framework/transform2/transformrt"
// 		type Transform struct {}
// 		func (t *Transform) SvelteToSvelte(file *transformrt.File) error {
// 			file.Data = append(file.Data, file.Data...)
// 			return nil
// 		}
// 	`
// 	is.NoErr(td.Write(ctx))
// 	module, err := gomod.Find(dir)
// 	is.NoErr(err)
// 	bfs, err := loadBFS(log, module)
// 	is.NoErr(err)
// 	is.NoErr(bfs.Sync())
// 	defer bfs.Close()
// 	is.NoErr(td.Exists("bud/internal/generator/transform/transform.go"))
// 	is.NoErr(td.NotExists("bud/service/transform"))
// 	code, err := fs.ReadFile(bfs, "bud/service/transform/svelte/view/index.svelte")
// 	is.NoErr(err)
// 	is.Equal(string(code), `<h1>hello</h1><h1>hello</h1>`)
// }

// func TestSvelteToSvelteToJSX(t *testing.T) {
// 	is := is.New(t)
// 	log := testlog.New()
// 	ctx := context.Background()
// 	dir := t.TempDir()
// 	td := testdir.New(dir)
// 	td.Files["view/index.svelte"] = `<h1>hello</h1>`
// 	td.Files["transform/doubler/transform.go"] = `
// 		package doubler
// 		import "github.com/livebud/bud/framework/transform2/transformrt"
// 		type Transform struct {}
// 		func (t *Transform) SvelteToSvelte(file *transformrt.File) error {
// 			file.Data = append(file.Data, file.Data...)
// 			return nil
// 		}
// 	`
// 	td.Files["transform/jsx/transform.go"] = `
// 		package jsx
// 		import "github.com/livebud/bud/framework/transform2/transformrt"
// 		type Transform struct {}
// 		func (t *Transform) SvelteToJsx(file *transformrt.File) error {
// 			file.Data = []byte("export default function() { return <>" + string(file.Data) + "</> }")
// 			return nil
// 		}
// 	`
// 	is.NoErr(td.Write(ctx))
// 	module, err := gomod.Find(dir)
// 	is.NoErr(err)
// 	bfs, err := loadBFS(log, module)
// 	is.NoErr(err)
// 	is.NoErr(bfs.Sync())
// 	defer bfs.Close()
// 	code, err := fs.ReadFile(bfs, "bud/service/transform/jsx/view/index.svelte")
// 	is.NoErr(err)
// 	is.Equal(string(code), `export default function() { return <><h1>hello</h1><h1>hello</h1></> }`)
// }

// func TestFaviconToFavicon(t *testing.T) {
// 	is := is.New(t)
// 	log := testlog.New()
// 	ctx := context.Background()
// 	dir := t.TempDir()
// 	td := testdir.New(dir)
// 	td.BFiles["public/favicon.ico"] = []byte{0x01, 0x02, 0x03}
// 	td.Files["transform/doubler/transform.go"] = `
// 		package doubler
// 		import "github.com/livebud/bud/framework/transform2/transformrt"
// 		type Transform struct {}
// 		func (t *Transform) SvelteToSvelte(file *transformrt.File) error {
// 			file.Data = append(file.Data, file.Data...)
// 			return nil
// 		}
// 		func (t *Transform) IcoToIco(file *transformrt.File) error {
// 			file.Data = append(file.Data, file.Data...)
// 			return nil
// 		}
// 	`
// 	is.NoErr(td.Write(ctx))
// 	module, err := gomod.Find(dir)
// 	is.NoErr(err)
// 	bfs, err := loadBFS(log, module)
// 	is.NoErr(err)
// 	is.NoErr(bfs.Sync())
// 	defer bfs.Close()
// 	code, err := fs.ReadFile(bfs, "bud/service/transform/ico/public/favicon.ico")
// 	is.NoErr(err)
// 	is.Equal(code, []byte{0x01, 0x02, 0x03, 0x01, 0x02, 0x03})
// }

// func TestMdToSsrJsAndDomJs(t *testing.T) {
// 	is := is.New(t)
// 	log := testlog.New()
// 	ctx := context.Background()
// 	dir := t.TempDir()
// 	td := testdir.New(dir)
// 	td.Files["view/index.md"] = `# hello`
// 	td.Files["transform/svelte/transform.go"] = `
// 		package svelte
// 		import "github.com/livebud/bud/framework/transform2/transformrt"
// 		type Transform struct {}
// 		func (t *Transform) SvelteToSsrJs(file *transformrt.File) error {
// 			file.Data = []byte("module.exports = function() { return " + string(file.Data) + " }")
// 			return nil
// 		}
// 		func (t *Transform) SvelteToDomJs(file *transformrt.File) error {
// 			file.Data = []byte("export default function() { return " + string(file.Data) + " }")
// 			return nil
// 		}
// 	`
// 	td.Files["transform/tailwind/transform.go"] = `
// 		package tailwind
// 		import "github.com/livebud/bud/framework/transform2/transformrt"
// 		import "bytes"
// 		type Transform struct {}
// 		func (t *Transform) SvelteToSvelte(file *transformrt.File) error {
// 			file.Data = bytes.Replace(file.Data, []byte("class='bg-red-100'"), []byte("style='color: red'"), -1)
// 			return nil
// 		}
// 	`
// 	td.Files["transform/markdoc/transform.go"] = `
// 		package markdoc
// 		import "github.com/livebud/bud/framework/transform2/transformrt"
// 		import "bytes"
// 		type Markdoc struct {}
// 		func (m *Markdoc) Compile(data []byte) ([]byte, error) {
// 			data = bytes.TrimPrefix(data, []byte("# "))
// 			data = []byte("<h1 class='bg-red-100'>" + string(data) + "</h1>")
// 			return data, nil
// 		}
// 		type Transform struct {
// 			Markdoc *Markdoc
// 		}
// 		func (t *Transform) MdToSvelte(file *transformrt.File) (err error) {
// 			file.Data, err = t.Markdoc.Compile(file.Data)
// 			if err != nil {
// 				return err
// 			}
// 			return nil
// 		}
// 	`
// 	is.NoErr(td.Write(ctx))
// 	module, err := gomod.Find(dir)
// 	is.NoErr(err)
// 	bfs, err := loadBFS(log, module)
// 	is.NoErr(err)
// 	is.NoErr(bfs.Sync())
// 	defer bfs.Close()
// 	code, err := fs.ReadFile(bfs, "bud/service/transform/ssr.js/view/index.md")
// 	is.NoErr(err)
// 	is.Equal(string(code), `module.exports = function() { return <h1 style='color: red'>hello</h1> }`)
// 	code, err = fs.ReadFile(bfs, "bud/service/transform/dom.js/view/index.md")
// 	is.NoErr(err)
// 	is.Equal(string(code), `export default function() { return <h1 style='color: red'>hello</h1> }`)
// }
