package generator_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/lithammer/dedent"
	"github.com/livebud/bud/internal/is"
	"github.com/livebud/bud/internal/testcli"
	"github.com/livebud/bud/internal/testdir"
)

func TestGenerators(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()
	dir := t.TempDir()
	td := testdir.New(dir)
	td.Files["generator/tailwind/tailwind.go"] = `
		package tailwind
		import (
			"github.com/livebud/bud/package/genfs"
		)
		type Generator struct {}
		func (g *Generator) Generate(fsys genfs.FS, dir *genfs.Dir) error {
			dir.GenerateFile("tailwind.css", func(fsys genfs.FS, file *genfs.File) error {
				file.Data = []byte("/** tailwind **/")
				return nil
			})
			dir.GenerateFile("preflight.css", func(fsys genfs.FS, file *genfs.File) error {
				file.Data = []byte("/** preflight **/")
				return nil
			})
			return nil
		}
	`
	td.Files["internal/markdoc/markdoc.go"] = `
		package markdoc
		import "context"
		type Compiler struct {}
		func (c *Compiler) Compile(ctx context.Context, input string) (string, error) {
			return input + " " + input, nil
		}
	`
	td.Files["view/index.md"] = `# Index`
	td.Files["view/about/index.md"] = `# About`
	td.Files["generator/markdoc/markdoc.go"] = `
		package markdoc
		import (
			"app.com/internal/markdoc"
			"github.com/livebud/bud/package/genfs"
			"context"
			"io/fs"
			"path"
		)
		type Generator struct {
			Markdoc *markdoc.Compiler
		}
		func (g *Generator) Generate(fsys genfs.FS, dir *genfs.Dir) error {
			dir.ServeFile(".", g.compile)
			dir.GenerateFile("markdoc.go", g.generate)
			return nil
		}
		func (g *Generator) generate(fsys genfs.FS, file *genfs.File) error {
			mds, err := fs.Glob(fsys, "view/**.md")
			if err != nil {
				return err
			}
			out := "package markdoc "
			for _, md := range mds {
				data, err := fs.ReadFile(fsys, path.Join("bud/internal/markdoc", md))
				if err != nil {
					return err
				}
				out += string(data) + " "
			}
			file.Data = []byte(out)
			return nil
		}
		func (g *Generator) compile(fsys genfs.FS, file *genfs.File) error {
			data, err := fs.ReadFile(fsys, file.Relative())
			if err != nil {
				return err
			}
			result, err := g.Markdoc.Compile(context.Background(), string(data))
			if err != nil {
				return err
			}
			file.Data = []byte(result)
			return nil
		}
	`
	td.Files["generator/frontend/viewer/viewer.go"] = `
		package viewer
		import (
			"github.com/livebud/bud/package/genfs"
		)
		type Generator struct {
		}
		func (g *Generator) Generate(fsys genfs.FS, dir *genfs.Dir) error {
			dir.GenerateFile("viewer.go", func(fsys genfs.FS, file *genfs.File) error {
				file.Data = []byte("package viewer")
				return nil
			})
			return nil
		}
	`
	is.NoErr(td.Write(ctx))
	cli := testcli.New(dir)
	_, err := cli.Run(ctx, "build", "--embed=false")
	is.NoErr(err)
	is.NoErr(td.Exists("bud/internal/tailwind/tailwind.css"))
	is.NoErr(td.Exists("bud/internal/tailwind/preflight.css"))
	is.NoErr(td.Exists("bud/internal/markdoc/markdoc.go"))
	is.NoErr(td.Exists("bud/internal/frontend/viewer/viewer.go"))
	data, err := os.ReadFile(td.Path("bud/internal/tailwind/tailwind.css"))
	is.NoErr(err)
	is.Equal(string(data), "/** tailwind **/")
	data, err = os.ReadFile(td.Path("bud/internal/tailwind/preflight.css"))
	is.NoErr(err)
	is.Equal(string(data), "/** preflight **/")
	data, err = os.ReadFile(td.Path("bud/internal/markdoc/markdoc.go"))
	is.NoErr(err)
	is.Equal(string(data), "package markdoc # About # About # Index # Index ")
	data, err = os.ReadFile(td.Path("bud/internal/frontend/viewer/viewer.go"))
	is.NoErr(err)
	is.Equal(string(data), "package viewer")
}

func TestMissingGenerator(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()
	dir := t.TempDir()
	td := testdir.New(dir)
	td.Files["generator/web/transform/transform.go"] = `
		package transform
	`
	is.NoErr(td.Write(ctx))
	cli := testcli.New(dir)
	_, err := cli.Run(ctx, "build", "--embed=false")
	is.NoErr(err)
	is.NoErr(td.NotExists("bud/command/generate/main.go"))
	is.NoErr(td.NotExists("bud/command/generate/main"))
}

func TestMissingMatchingMethod(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()
	dir := t.TempDir()
	td := testdir.New(dir)
	td.Files["generator/web/transform/transform.go"] = `
		package transform
		type Generator struct {}
		func (g *Generator) generate() error { return nil }
	`
	is.NoErr(td.Write(ctx))
	cli := testcli.New(dir)
	_, err := cli.Run(ctx, "build", "--embed=false")
	is.NoErr(err)
	is.NoErr(td.NotExists("bud/command/generate/main"))
}

func TestSyntaxError(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()
	dir := t.TempDir()
	td := testdir.New(dir)
	td.Files["generator/tailwind/tailwind.go"] = `
		package tailwind
		import (
			"github.com/livebud/bud/package/genfs"
		)
		type Generator struct {}
		func (g *Generator) Generate(fsys genfs.FS, dir *genfs.Dir) error {
			"ok"
		}
	`
	is.NoErr(td.Write(ctx))
	cli := testcli.New(dir)
	res, err := cli.Run(ctx, "build", "--embed=false")
	is.True(err != nil)
	is.In(err.Error(), `exit status 2`)
	is.In(res.Stderr(), `"ok"`)
	is.In(res.Stderr(), `not used`)
}

func TestUpdateGenerator(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()
	dir := t.TempDir()
	td := testdir.New(dir)
	td.Files["generator/tailwind/tailwind.go"] = `
		package tailwind
		import (
			"github.com/livebud/bud/package/genfs"
		)
		type Generator struct {}
		func (g *Generator) Generate(fsys genfs.FS, dir *genfs.Dir) error {
			dir.GenerateFile("tailwind.css", func(fsys genfs.FS, file *genfs.File) error {
				file.Data = []byte("/** tailwind **/")
				return nil
			})
			return nil
		}
	`
	is.NoErr(td.Write(ctx))
	cli := testcli.New(dir)
	app, err := cli.Start(ctx, "run")
	is.NoErr(err)
	defer app.Close()
	// Check for tailwind
	data, err := os.ReadFile(td.Path("bud/internal/tailwind/tailwind.css"))
	is.NoErr(err)
	is.Equal(string(data), "/** tailwind **/")
	// Update generator
	generatorFile := filepath.Join(dir, "generator", "tailwind", "tailwind.go")
	is.NoErr(os.WriteFile(generatorFile, []byte(dedent.Dedent(`
		package tailwind
		import (
			"github.com/livebud/bud/package/genfs"
		)
		type Generator struct {
		}
		func (g *Generator) Generate(fsys genfs.FS, dir *genfs.Dir) error {
			dir.GenerateFile("tailwind.css", func(fsys genfs.FS, file *genfs.File) error {
				file.Data = []byte("/** tailwind2 **/")
				return nil
			})
			dir.GenerateFile("preflight.css", func(fsys genfs.FS, file *genfs.File) error {
				file.Data = []byte("/** preflight **/")
				return nil
			})
			return nil
		}
	`)), 0644))
	// Wait for the app to be ready again
	readyCtx, cancel := context.WithTimeout(ctx, 15*time.Second)
	is.NoErr(app.Ready(readyCtx))
	cancel()
	// Check for preflight
	data, err = os.ReadFile(td.Path("bud/internal/tailwind/preflight.css"))
	is.NoErr(err)
	is.Equal(string(data), "/** preflight **/")
	// Check that tailwind has been updated
	data, err = os.ReadFile(td.Path("bud/internal/tailwind/tailwind.css"))
	is.NoErr(err)
	is.Equal(string(data), "/** tailwind2 **/")
}

func TestDeleteGenerator(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()
	dir := t.TempDir()
	td := testdir.New(dir)
	td.Files["generator/tailwind/tailwind.go"] = `
		package tailwind
		import (
			"github.com/livebud/bud/package/genfs"
		)
		type Generator struct {}
		func (g *Generator) Generate(fsys genfs.FS, dir *genfs.Dir) error {
			dir.GenerateFile("tailwind.css", func(fsys genfs.FS, file *genfs.File) error {
				file.Data = []byte("/** tailwind **/")
				return nil
			})
			return nil
		}
	`
	is.NoErr(td.Write(ctx))
	cli := testcli.New(dir)
	app, err := cli.Start(ctx, "run")
	is.NoErr(err)
	defer app.Close()
	// Check for tailwind
	data, err := os.ReadFile(td.Path("bud/internal/tailwind/tailwind.css"))
	is.NoErr(err)
	is.Equal(string(data), "/** tailwind **/")
	// Update generator
	generatorFile := filepath.Join(dir, "generator", "tailwind", "tailwind.go")
	is.NoErr(os.WriteFile(generatorFile, []byte(dedent.Dedent(`
		package tailwind
		import (
			"github.com/livebud/bud/package/genfs"
		)
		type Generator struct {}
		func (g *Generator) Generate(fsys genfs.FS, dir *genfs.Dir) error {
			return nil
		}
	`)), 0644))
	// Wait for the app to be ready again
	readyCtx, cancel := context.WithTimeout(ctx, 15*time.Second)
	is.NoErr(app.Ready(readyCtx))
	cancel()
	// Check that tailwind has been updated
	data, err = os.ReadFile(td.Path("bud/internal/tailwind/tailwind.css"))
	is.True(os.IsNotExist(err))
	is.Equal(data, nil)
}

func TestChangeGenerator(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()
	dir := t.TempDir()
	td := testdir.New(dir)
	td.Files["generator/tailwind/tailwind.go"] = `
		package tailwind
		import (
			"github.com/livebud/bud/package/genfs"
		)
		type Generator struct {}
		func (g *Generator) Generate(fsys genfs.FS, dir *genfs.Dir) error {
			dir.GenerateFile("tailwind.css", func(fsys genfs.FS, file *genfs.File) error {
				file.Data = []byte("/** tailwind **/")
				return nil
			})
			return nil
		}
	`
	is.NoErr(td.Write(ctx))
	cli := testcli.New(dir)
	app, err := cli.Start(ctx, "run")
	is.NoErr(err)
	defer app.Close()
	// Check for tailwind
	data, err := os.ReadFile(td.Path("bud/internal/tailwind/tailwind.css"))
	is.NoErr(err)
	is.Equal(string(data), "/** tailwind **/")
	// Update generator
	generatorFile := filepath.Join(dir, "generator", "tailwind", "tailwind.go")
	is.NoErr(os.WriteFile(generatorFile, []byte(dedent.Dedent(`
		package tailwind
		import (
			"github.com/livebud/bud/package/genfs"
		)
		type Generator struct {}
		func (g *Generator) Generate(fsys genfs.FS, dir *genfs.Dir) error {
			dir.GenerateFile("preflight.css", func(fsys genfs.FS, file *genfs.File) error {
				file.Data = []byte("/** preflight **/")
				return nil
			})
			return nil
		}
	`)), 0644))
	// Wait for the app to be ready again
	readyCtx, cancel := context.WithTimeout(ctx, 15*time.Second)
	is.NoErr(app.Ready(readyCtx))
	cancel()
	// Check that tailwind has been updated
	data, err = os.ReadFile(td.Path("bud/internal/tailwind/tailwind.css"))
	is.True(os.IsNotExist(err))
	is.Equal(data, nil)
	// Check for preflight
	data, err = os.ReadFile(td.Path("bud/internal/tailwind/preflight.css"))
	is.NoErr(err)
	is.Equal(string(data), "/** preflight **/")
}

func TestPkgGenerator(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()
	dir := t.TempDir()
	td := testdir.New(dir)
	td.Files["generator/tailwind/tailwind.go"] = `
		package tailwind
		import (
			"github.com/livebud/bud/package/genfs"
		)
		type Generator struct {}
		func (g *Generator) GeneratePkg(fsys genfs.FS, dir *genfs.Dir) error {
			dir.GenerateFile("tailwind.css", func(fsys genfs.FS, file *genfs.File) error {
				file.Data = []byte("/** tailwind **/")
				return nil
			})
			return nil
		}
	`
	is.NoErr(td.Write(ctx))
	cli := testcli.New(dir)
	res, err := cli.Run(ctx, "build", "--embed=false")
	is.NoErr(err)
	is.Equal(res.Stderr(), "")
	is.Equal(res.Stdout(), "")
	is.NoErr(td.Exists("bud/pkg/tailwind/tailwind.css"))
	data, err := os.ReadFile(td.Path("bud/pkg/tailwind/tailwind.css"))
	is.NoErr(err)
	is.Equal(string(data), "/** tailwind **/")
}

func TestCmdGenerator(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()
	dir := t.TempDir()
	td := testdir.New(dir)
	td.Files["generator/tailwind/tailwind.go"] = `
		package tailwind
		import (
			"github.com/livebud/bud/package/genfs"
		)
		type Generator struct {}
		func (g *Generator) GenerateCmd(fsys genfs.FS, dir *genfs.Dir) error {
			dir.GenerateFile("main.go", func(fsys genfs.FS, file *genfs.File) error {
				file.Data = []byte("package main\nfunc main() {}")
				return nil
			})
			return nil
		}
	`
	is.NoErr(td.Write(ctx))
	cli := testcli.New(dir)
	res, err := cli.Run(ctx, "build", "--embed=false")
	is.NoErr(err)
	is.Equal(res.Stderr(), "")
	is.Equal(res.Stdout(), "")
	is.NoErr(td.Exists("bud/cmd/tailwind/main.go"))
	data, err := os.ReadFile(td.Path("bud/cmd/tailwind/main.go"))
	is.NoErr(err)
	is.Equal(string(data), "package main\nfunc main() {}")
}

func TestGeneratorServer(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()
	dir := t.TempDir()
	td := testdir.New(dir)
	td.Files["generator/tailwind/tailwind.go"] = `
		package tailwind
		import (
			"github.com/livebud/bud/package/genfs"
		)
		type Generator struct {}
		func (g *Generator) Serve(fsys genfs.FS, file *genfs.File) error {
			file.Data = []byte(file.Relative())
			return nil
		}
	`
	td.Files["generator/view/view.go"] = `
		package tailwind
		import (
			"github.com/livebud/bud/package/genfs"
			"io/fs"
		)
		type Generator struct {}
		func (g *Generator) Generate(fsys genfs.FS, dir *genfs.Dir) error {
			dir.GenerateFile("view.go", func(fsys genfs.FS, file *genfs.File) error {
				code, err := fs.ReadFile(fsys, "bud/internal/tailwind/preflight.css")
				if err != nil {
					return err
				}
				file.Data = code
				return nil
			})
			return nil
		}
	`
	is.NoErr(td.Write(ctx))
	cli := testcli.New(dir)
	res, err := cli.Run(ctx, "build", "--embed=false")
	is.NoErr(err)
	is.Equal(res.Stderr(), "")
	is.Equal(res.Stdout(), "")
	is.NoErr(td.Exists("bud/internal/view/view.go"))
	data, err := os.ReadFile(td.Path("bud/internal/view/view.go"))
	is.NoErr(err)
	is.Equal(string(data), "preflight.css")
}

// This test is trippy, but amazing. We're testing that we can create a custom
// generator in $APP/generator/web/health that writes a health.go file to
// $APP/bud/internal/web/health.go that contains a Handler that the built-in
// web generator know how to hook up to the web server to expose the /health
// endpoint.
func TestCustomWebGenerator(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()
	dir := t.TempDir()
	td := testdir.New(dir)
	const healthHandler = `
		package health
		import "github.com/livebud/bud/package/router"
		import "net/http"
		type Handler struct {}
		func (h *Handler) Register(r *router.Router) {
			r.Get("/health", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Write([]byte("all good"))
			}))
		}
	`
	td.Files["generator/web/health/health.go"] = `
		package health
		import "github.com/livebud/bud/package/genfs"
		type Generator struct {}
		func (g *Generator) Generate(fsys genfs.FS, dir *genfs.Dir) error {
			dir.GenerateFile("health.go", func(fsys genfs.FS, file *genfs.File) error {
				file.Data = []byte(` + "`" + healthHandler + "`" + `)
				return nil
			})
			return nil
		}
	`
	is.NoErr(td.Write(ctx))
	cli := testcli.New(dir)
	app, err := cli.Start(ctx, "run")
	is.NoErr(err)
	defer app.Close()
	res, err := app.Get("/health")
	is.NoErr(err)
	is.Equal(res.Status(), 200)
	is.Equal(res.Body().String(), "all good")
	is.NoErr(app.Close())
}
