package web_test

import (
	"bytes"
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/lithammer/dedent"
	"github.com/matryer/is"
	"github.com/matthewmueller/diff"
	"gitlab.com/mnm/bud/fsync"
	"gitlab.com/mnm/bud/gen"
	"gitlab.com/mnm/bud/internal/di"
	"gitlab.com/mnm/bud/internal/generator/command"
	"gitlab.com/mnm/bud/internal/generator/maingo"
	"gitlab.com/mnm/bud/internal/generator/program"
	"gitlab.com/mnm/bud/internal/generator/web"
	"gitlab.com/mnm/bud/internal/modcache"
	"gitlab.com/mnm/bud/internal/modtest"
	"gitlab.com/mnm/bud/internal/parser"
	"gitlab.com/mnm/bud/vfs"
)

func redent(s string) string {
	return strings.TrimSpace(dedent.Dedent(s)) + "\n"
}

func goRun(cacheDir, appDir string, args ...string) (string, error) {
	ctx := context.Background()
	mainPath := filepath.Join("bud", "main.go")
	args = append([]string{"run", "-mod=mod", mainPath}, args...)
	cmd := exec.CommandContext(ctx, "go", args...)
	cmd.Env = append(os.Environ(), "GOMODCACHE="+cacheDir, "GOPRIVATE=*", "NO_COLOR=1")
	stdout := new(bytes.Buffer)
	cmd.Stdout = stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	cmd.Dir = appDir
	err := cmd.Run()
	if err != nil {
		return "", err
	}
	return stdout.String(), nil
}

// Generate commands
func generate(t testing.TB, m modtest.Module) func(args ...string) string {
	is := is.New(t)
	m.AppDir = t.TempDir()
	m.CacheDir = modcache.Default().Directory()
	appfs := vfs.OS(m.AppDir)
	genfs := gen.New(appfs)
	m.FS = genfs
	module := modtest.Make(t, m)
	parser := parser.New(module)
	injector := di.New(module, parser, di.Map{})
	genfs.Add(map[string]gen.Generator{
		"bud/main.go": gen.FileGenerator(&maingo.Generator{
			Module: module,
		}),
		"bud/program/program.go": gen.FileGenerator(&program.Generator{
			Module:   module,
			Injector: injector,
		}),
		"bud/command/command.go": gen.FileGenerator(&command.Generator{
			Module: module,
			Parser: parser,
		}),
		"bud/web/web.go": gen.FileGenerator(&web.Generator{
			Module: module,
		}),
	})
	err := fsync.Dir(genfs, "bud", appfs, "bud")
	is.NoErr(err)
	return func(args ...string) string {
		stdout, err := goRun(m.CacheDir, m.AppDir, args...)
		if err != nil {
			return err.Error()
		}
		return stdout
	}
}

const goMod = `
module app.com

require (
  github.com/hexops/valast v1.4.1
	gitlab.com/mnm/bud v0.0.0
)
`

func isEqual(t testing.TB, actual, expect string) {
	diff.TestString(t, redent(expect), redent(actual))
}

func TestEmpty(t *testing.T) {
	t.SkipNow()
	run := generate(t, modtest.Module{
		Files: map[string]string{
			"go.mod": goMod,
		},
	})
	isEqual(t, run("-h"), `
		Usage:
		  app
	`)
}
