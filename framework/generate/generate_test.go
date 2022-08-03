package generate_test

import (
	"context"
	"os"
	"testing"

	"github.com/livebud/bud/internal/cli/testcli"
	"github.com/livebud/bud/internal/is"
	"github.com/livebud/bud/internal/testdir"
)

func TestTailwindGenerator(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()
	dir := t.TempDir()
	td := testdir.New(dir)
	td.Files["generator/tailwind/tailwind.go"] = `
		package tailwind
		type Generator struct {
		}
		func (g *Generator) Register(genfs *overlay.FileSystem, flag *framework.Flag) {
			fmt.Println("registering tailwind")
			genfs.FileGenerator("tailwind.css", GenerateFile(ctx context.Context, fsys overlay.F, file *overlay.File) error {
				fmt.Println("running generator!")
				file.Data = []byte("/** tailwind **/")
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
	is.NoErr(td.Exists("bud/internal/generator/tailwind/tailwind.css"))
	data, err := os.ReadFile(td.Path("bud/internal/generator/tailwind/tailwind.css"))
	is.NoErr(err)
	is.Equal(string(data), "/** tailwind **/")
}
