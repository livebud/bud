package generator_test

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"

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
			"context"
			"github.com/livebud/bud/package/overlay"
		)
		type Generator struct {
		}
		func (g *Generator) GenerateDir(ctx context.Context, fsys overlay.F, dir *overlay.Dir) error {
			dir.GenerateFile("tailwind.css", func(ctx context.Context, fsys overlay.F, file *overlay.File) error {
				file.Data = []byte("/** tailwind **/")
				return nil
			})
			dir.GenerateFile("preflight.css", func(ctx context.Context, fsys overlay.F, file *overlay.File) error {
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
			"context"
			"github.com/livebud/bud/package/overlay"
			"io/fs"
			"strings"
		)
		type Generator struct {
			Markdoc *markdoc.Compiler
		}
		func (g *Generator) GenerateDir(ctx context.Context, fsys overlay.F, dir *overlay.Dir) error {
			dir.GenerateFile("view/index.md", g.compile)
			dir.GenerateFile("view/about/index.md", g.compile)
			return nil
		}
		func (g *Generator) compile(ctx context.Context, fsys overlay.F, file *overlay.File) error {
			data, err := fs.ReadFile(fsys, strings.TrimPrefix(file.Path(), "bud/internal/generator/markdoc/"))
			if err != nil {
				return err
			}
			result, err := g.Markdoc.Compile(ctx, string(data))
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
			"context"
			"github.com/livebud/bud/package/overlay"
		)
		type Generator struct {
		}
		func (g *Generator) GenerateDir(ctx context.Context, fsys overlay.F, dir *overlay.Dir) error {
			dir.GenerateFile("viewer.go", func(ctx context.Context, fsys overlay.F, file *overlay.File) error {
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
	is.NoErr(td.Exists("bud/internal/generate/main.go"))
	is.NoErr(td.Exists("bud/internal/generate/generator/generator.go"))
	is.NoErr(td.Exists("bud/internal/generator/tailwind/tailwind.css"))
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
		package transform
		import (
			"context"
			"github.com/livebud/bud/package/overlay"
		)
		type Generator struct {}
		func (g *Generator) GenerateDir(ctx context.Context, fsys overlay.F, dir *overlay.Dir) error {
			return "oh noz"
		}
	`
	is.NoErr(td.Write(ctx))
	cli := testcli.New(dir)
	res, err := cli.Run(ctx, "build", "--embed=false")
	is.True(err != nil)
	is.In(err.Error(), `exit status 2`)
	is.Equal(res.Stderr(), ``)
	is.Equal(res.Stdout(), ``)
}

func TestChange(t *testing.T) {
	t.SkipNow()
	is := is.New(t)
	ctx := context.Background()
	dir := t.TempDir()
	td := testdir.New(dir)
	td.Files["generator/tailwind/tailwind.go"] = `
		package tailwind
		import (
			"context"
			"github.com/livebud/bud/package/overlay"
		)
		type Generator struct {
		}
		func (g *Generator) GenerateDir(ctx context.Context, fsys overlay.F, dir *overlay.Dir) error {
			dir.GenerateFile("tailwind.css", func(ctx context.Context, fsys overlay.F, file *overlay.File) error {
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
			"context"
			"github.com/livebud/bud/package/overlay"
		)
		type Generator struct {
		}
		func (g *Generator) GenerateDir(ctx context.Context, fsys overlay.F, dir *overlay.Dir) error {
			dir.GenerateFile("preflight.css", func(ctx context.Context, fsys overlay.F, file *overlay.File) error {
				file.Data = []byte("/** preflight **/")
				return nil
			})
			return nil
		}
	`)), 0644))
	// Wait for the app to be ready again
	fmt.Println("waiting for app to be ready again")
	is.NoErr(app.Ready(ctx))
	fmt.Println("waited for app to be ready again")
	// Check for preflight
	data, err = os.ReadFile(td.Path("bud/internal/generator/tailwind/preflight.css"))
	is.NoErr(err)
	is.Equal(string(data), "/** preflight **/")
	// Tailwind.css should have been removed
	is.NoErr(td.NotExists("bud/internal/generator/tailwind/tailwind.css"))
}
