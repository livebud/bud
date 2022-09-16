package generator_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/lithammer/dedent"
	"github.com/livebud/bud/internal/cli/testcli"
	"github.com/livebud/bud/internal/is"
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
			"github.com/livebud/bud/package/budfs"
		)
		type Generator struct {}
		func (g *Generator) GenerateDir(fsys budfs.FS, dir *budfs.Dir) error {
			dir.GenerateFile("tailwind.css", func(fsys budfs.FS, file *budfs.File) error {
				file.Data = []byte("/** tailwind **/")
				return nil
			})
			dir.GenerateFile("preflight.css", func(fsys budfs.FS, file *budfs.File) error {
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
			"github.com/livebud/bud/package/budfs"
			"io/fs"
		)
		type Generator struct {
			Markdoc *markdoc.Compiler
		}
		func (g *Generator) GenerateDir(fsys budfs.FS, dir *budfs.Dir) error {
			dir.GenerateFile("view/index.md", g.compile)
			dir.GenerateFile("view/about/index.md", g.compile)
			return nil
		}
		func (g *Generator) compile(fsys budfs.FS, file *budfs.File) error {
			data, err := fs.ReadFile(fsys, file.Path())
			if err != nil {
				return err
			}
			result, err := g.Markdoc.Compile(fsys.Context(), string(data))
			if err != nil {
				return err
			}
			file.Data = []byte(result)
			return nil
		}
	`
	td.Files["generator/web/viewer/viewer.go"] = `
		package viewer
		import (
			"github.com/livebud/bud/package/budfs"
		)
		type Generator struct {
		}
		func (g *Generator) GenerateDir(fsys budfs.FS, dir *budfs.Dir) error {
			dir.GenerateFile("viewer.go", func(fsys budfs.FS, file *budfs.File) error {
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
	is.NoErr(td.Exists("bud/tmp/generate/main.go"))
	is.NoErr(td.Exists("bud/internal/generator/tailwind/tailwind.css"))
	is.NoErr(td.Exists("bud/internal/generator/tailwind/preflight.css"))
	is.NoErr(td.Exists("bud/internal/generator/markdoc/view/index.md"))
	is.NoErr(td.Exists("bud/internal/generator/markdoc/view/about/index.md"))
	is.NoErr(td.Exists("bud/internal/generator/web/viewer/viewer.go"))
	data, err := os.ReadFile(td.Path("bud/internal/generator/tailwind/tailwind.css"))
	is.NoErr(err)
	is.Equal(string(data), "/** tailwind **/")
	data, err = os.ReadFile(td.Path("bud/internal/generator/tailwind/preflight.css"))
	is.NoErr(err)
	is.Equal(string(data), "/** preflight **/")
	data, err = os.ReadFile(td.Path("bud/internal/generator/markdoc/view/index.md"))
	is.NoErr(err)
	is.Equal(string(data), "# Index # Index")
	data, err = os.ReadFile(td.Path("bud/internal/generator/markdoc/view/about/index.md"))
	is.NoErr(err)
	is.Equal(string(data), "# About # About")
	data, err = os.ReadFile(td.Path("bud/internal/generator/web/viewer/viewer.go"))
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
	is.True(err != nil)
	is.In(err.Error(), `generator: no Generator struct in "app.com/generator/web/transform"`)
}

func TestMissingGenerateDirMethod(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()
	dir := t.TempDir()
	td := testdir.New(dir)
	td.Files["generator/web/transform/transform.go"] = `
		package transform
		type Generator struct {}
		func (g *Generator) Generate() error { return nil }
	`
	is.NoErr(td.Write(ctx))
	cli := testcli.New(dir)
	_, err := cli.Run(ctx, "build", "--embed=false")
	is.True(err != nil)
	is.In(err.Error(), `generator: no (*Generator).GenerateDir(...) method in "app.com/generator/web/transform"`)
}

func TestSyntaxError(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()
	dir := t.TempDir()
	td := testdir.New(dir)
	td.Files["generator/tailwind/tailwind.go"] = `
		package tailwind
		import (
			"github.com/livebud/bud/package/budfs"
		)
		type Generator struct {}
		func (g *Generator) GenerateDir(fsys budfs.FS, dir *budfs.Dir) error {
			return "not an error"
		}
	`
	is.NoErr(td.Write(ctx))
	cli := testcli.New(dir)
	res, err := cli.Run(ctx, "build", "--embed=false")
	is.True(err != nil)
	is.In(err.Error(), `exit status 2`)
	is.In(res.Stderr(), `string does not implement error (missing Error method)`)
}

func TestUpdateGenerator(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()
	dir := t.TempDir()
	td := testdir.New(dir)
	td.Files["generator/tailwind/tailwind.go"] = `
		package tailwind
		import (
			"github.com/livebud/bud/package/budfs"
		)
		type Generator struct {}
		func (g *Generator) GenerateDir(fsys budfs.FS, dir *budfs.Dir) error {
			dir.GenerateFile("tailwind.css", func(fsys budfs.FS, file *budfs.File) error {
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
	data, err := os.ReadFile(td.Path("bud/internal/generator/tailwind/tailwind.css"))
	is.NoErr(err)
	is.Equal(string(data), "/** tailwind **/")
	// Update generator
	generatorFile := filepath.Join(dir, "generator", "tailwind", "tailwind.go")
	is.NoErr(os.WriteFile(generatorFile, []byte(dedent.Dedent(`
		package tailwind
		import (
			"github.com/livebud/bud/package/budfs"
		)
		type Generator struct {
		}
		func (g *Generator) GenerateDir(fsys budfs.FS, dir *budfs.Dir) error {
			dir.GenerateFile("tailwind.css", func(fsys budfs.FS, file *budfs.File) error {
				file.Data = []byte("/** tailwind2 **/")
				return nil
			})
			dir.GenerateFile("preflight.css", func(fsys budfs.FS, file *budfs.File) error {
				file.Data = []byte("/** preflight **/")
				return nil
			})
			return nil
		}
	`)), 0644))
	// Wait for the app to be ready again
	ctx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()
	is.NoErr(app.Ready(ctx))
	// Check for preflight
	data, err = os.ReadFile(td.Path("bud/internal/generator/tailwind/preflight.css"))
	is.NoErr(err)
	is.Equal(string(data), "/** preflight **/")
	// Check that tailwind has been updated
	data, err = os.ReadFile(td.Path("bud/internal/generator/tailwind/tailwind.css"))
	is.NoErr(err)
	is.Equal(string(data), "/** tailwind2 **/")
}

func TestRemoveGenerator(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()
	dir := t.TempDir()
	td := testdir.New(dir)
	td.Files["generator/tailwind/tailwind.go"] = `
		package tailwind
		import (
			"github.com/livebud/bud/package/budfs"
		)
		type Generator struct {}
		func (g *Generator) GenerateDir(fsys budfs.FS, dir *budfs.Dir) error {
			dir.GenerateFile("tailwind.css", func(fsys budfs.FS, file *budfs.File) error {
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
	data, err := os.ReadFile(td.Path("bud/internal/generator/tailwind/tailwind.css"))
	is.NoErr(err)
	is.Equal(string(data), "/** tailwind **/")
	// Update generator
	generatorFile := filepath.Join(dir, "generator", "tailwind", "tailwind.go")
	is.NoErr(os.WriteFile(generatorFile, []byte(dedent.Dedent(`
		package tailwind
		import (
			"github.com/livebud/bud/package/budfs"
		)
		type Generator struct {
		}
		func (g *Generator) GenerateDir(fsys budfs.FS, dir *budfs.Dir) error {
			return nil
		}
	`)), 0644))
	// Wait for the app to be ready again
	ctx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()
	is.NoErr(app.Ready(ctx))
	// Check that tailwind has been updated
	data, err = os.ReadFile(td.Path("bud/internal/generator/tailwind/tailwind.css"))
	is.True(os.IsNotExist(err))
	is.Equal(data, nil)
}
