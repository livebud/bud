package generator_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/matryer/is"
	"gitlab.com/mnm/bud/generator/cli/generator"
	"gitlab.com/mnm/bud/internal/testdir"
	"gitlab.com/mnm/bud/pkg/buddy"
	"gitlab.com/mnm/bud/pkg/gen"
)

const mainFile = `
package main

import "app.com/bud/.cli/generator"

func main() {
	println(generator.New)
}
`

func TestHelp(t *testing.T) {
	is := is.New(t)
	td := testdir.New()
	td.Files["main.go"] = mainFile
	dir := "_tmp"
	is.NoErr(os.RemoveAll(dir))
	err := td.Write(dir)
	is.NoErr(err)
	kit, err := buddy.Load(dir)
	is.NoErr(err)
	generator := generator.New(kit)
	kit.Generator("bud/.cli/generator/generator.go", gen.FileGenerator(generator))
	err = kit.Sync("bud", "bud")
	is.NoErr(err)
	tree, err := td.Tree(dir)
	is.NoErr(err)
	fmt.Println(tree)
	// err = driver.Expand(ctx, &bud.Expand{})
	// is.NoErr(err)
	// cmd := exec.Command("./bud/cli", "-h")
	// cmd.Dir = dir
	// stdout := new(bytes.Buffer)
	// cmd.Stdout = stdout
	// cmd.Stderr = os.Stderr
	// err = cmd.Run()
	// is.NoErr(err)
	// fmt.Println(stdout.String())
}
