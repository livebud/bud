package testgen_test

import (
	"context"
	"fmt"
	"io/fs"
	"testing"

	"github.com/livebud/bud/framework/generator"

	"github.com/livebud/bud/internal/dag"
	"github.com/livebud/bud/internal/is"

	"github.com/livebud/bud/package/genfs"
	"github.com/livebud/bud/package/gomod"
	"github.com/livebud/bud/package/log/testlog"
	"github.com/livebud/bud/package/parser"
	"github.com/livebud/bud/package/testdir"
	"github.com/livebud/bud/package/testgen"
)

const mainGen = `
	package main
	import (
		"fmt"
		"app.com/generator"
		"github.com/livebud/bud/package/genfs"
		"github.com/livebud/bud/package/log"
		"github.com/livebud/bud/package/gomod"
	)
	func main() {
		module := gomod.MustFind(".")
		fsys := genfs.New(nil, module, log.Discard)
		_, _ = fsys, generator.New
		fmt.Printf("Hello, World!")
	}
`

func TestGenerator(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()
	td, err := testdir.Load()
	is.NoErr(err)
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
	is.NoErr(td.Write(ctx))
	log := testlog.New()

	fsys := genfs.New(dag.Discard, td, log)
	module, err := gomod.Find(td.Directory())
	is.NoErr(err)
	parser := parser.New(td, module)
	generator := generator.New(log, module, parser)
	generator.Glob = "{generator/**/*.go,bud/internal/generator/*/*.go}"
	fsys.FileGenerator("generator/generator.go", generator)
	fsys.FileGenerator("main.go", &genfs.Embed{Data: []byte(mainGen)})
	code, err := fs.ReadFile(fsys, "generator/generator.go")
	is.NoErr(err)
	fmt.Println(string(code))
	gen := testgen.New(fsys, module)
	stdout, stderr, err := gen.Run(ctx, "main.go")
	is.NoErr(err)
	// TODO: Finish up the test.
	is.Equal(stdout.String(), "Hello, World!")
	is.Equal(stderr.String(), "")
}
