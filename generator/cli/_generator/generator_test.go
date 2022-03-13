package generator_test

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"testing"

	"gitlab.com/mnm/bud/internal/gobin"

	"github.com/matryer/is"
	"gitlab.com/mnm/bud/generator/cli/generator"
	"gitlab.com/mnm/bud/internal/testdir"
	"gitlab.com/mnm/bud/pkg/buddy"
	"gitlab.com/mnm/bud/pkg/gen"
)

func mainFile(provider string) func(_ gen.F, file *gen.File) error {
	mainFile := `
package main

import (
	"gitlab.com/mnm/bud/pkg/buddy"
	"app.com/bud/.cli/generator"
	"context"
	"fmt"
	"os"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
		os.Exit(1)
	}
}

func run() error {
	kit, err := buddy.Load(".")
	if err != nil {
		return err
	}
	generator, err := load(kit)
	if err != nil {
		return err
	}
	return generator.Generate(context.TODO())
}

` + provider
	return func(_ gen.F, file *gen.File) error {
		file.Write([]byte(mainFile))
		return nil
	}
}

// run the command
func run(dir string, args ...string) (string, error) {
	kit, err := buddy.Load(dir)
	if err != nil {
		return "", err
	}
	generator := generator.New(kit)
	kit.Generators(map[string]buddy.Generator{
		"bud/.cli/generator/generator.go": gen.FileGenerator(generator),
	})
	provider, err := kit.Wire(&buddy.Function{
		Name:   "load",
		Target: kit.ImportPath(),
		Params: []buddy.Dependency{
			buddy.ToType("gitlab.com/mnm/bud/pkg/buddy", "Kit"),
		},
		Results: []buddy.Dependency{
			buddy.ToType(kit.ImportPath("bud/.cli/generator"), "*Generator"),
			&buddy.ErrorType{},
		},
	})
	if err != nil {
		return "", err
	}
	kit.Generators(map[string]buddy.Generator{
		"bud/main.go": gen.GenerateFile(mainFile(provider.Function())),
	})
	err = kit.Sync("bud", "bud")
	if err != nil {
		return "", err
	}
	tree, err := testdir.Tree(dir)
	if err != nil {
		return "", err
	}
	fmt.Println(tree)
	ctx := context.Background()
	err = gobin.Build(ctx, dir, "bud/main.go", "main")
	if err != nil {
		return "", err
	}
	cmd := exec.Command("./main", args...)
	cmd.Dir = dir
	cmd.Env = []string{
		"HOME=" + os.Getenv("HOME"),
		"PATH=" + os.Getenv("PATH"),
		"GOPATH=" + os.Getenv("GOPATH"),
		"GOCACHE=" + os.Getenv("GOCACHE"),
		"GOMODCACHE=" + testdir.ModCache(dir).Directory(),
		"NO_COLOR=1",
		// TODO: remove once we can write a sum file to the modcache
		"GOPRIVATE=*",
	}
	stdout := new(bytes.Buffer)
	cmd.Stdout = stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return "", err
	}
	return stdout.String(), nil
}

// Tests
func TestEmpty(t *testing.T) {
	t.SkipNow()
	is := is.New(t)
	dir := "_tmp"
	err := os.RemoveAll(dir)
	is.NoErr(err)
	td := testdir.New()
	// td.Files["main.go"] = mainFile
	err = td.Write(dir)
	is.NoErr(err)
	stdout, err := run(dir)
	is.NoErr(err)
	fmt.Println(stdout)
}
