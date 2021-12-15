package command_test

import (
	"bytes"
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"gitlab.com/mnm/bud/fsync"

	"gitlab.com/mnm/bud/internal/di"
	"gitlab.com/mnm/bud/internal/parser"

	"gitlab.com/mnm/bud/gen"

	"gitlab.com/mnm/bud/go/mod"
	"gitlab.com/mnm/bud/internal/generator/command"
	"gitlab.com/mnm/bud/internal/modcache"

	"github.com/lithammer/dedent"
	"github.com/matryer/is"
	"github.com/matthewmueller/diff"
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
	cmd.Env = append(os.Environ(), "GOMODCACHE="+cacheDir, "GOPRIVATE=*")
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

var webServer = gen.GenerateFile(func(f gen.F, file *gen.File) error {
	file.Write([]byte(redent(`
			package web

			import (
				"context"
			)

			type Server struct {
			}

			func (s *Server) ListenAndServe(ctx context.Context, address string) error {
				return nil
			}
		`)))
	return nil
})

var mainGo = gen.GenerateFile(func(f gen.F, file *gen.File) error {
	file.Write([]byte(redent(`
			package main

			import (
				command "app.com/bud/command"
				os "os"
			)

			func main() {
				os.Exit(command.Parse(os.Args[1:]...))
			}
		`)))
	return nil
})

type Test struct {
	Skip   bool
	Files  map[string]string
	Args   []string
	Expect string
}

func runTest(t *testing.T, test Test) {
	t.Helper()
	if test.Skip {
		t.SkipNow()
	}
	is := is.New(t)
	wd, err := os.Getwd()
	is.NoErr(err)
	appDir := t.TempDir()
	modCache := modcache.Default()
	modFinder := mod.New(mod.WithCache(modCache))
	budModule, err := modFinder.Find(wd)
	is.NoErr(err)
	// Write application files
	if test.Files != nil {
		vmap := vfs.Map{}
		for path, code := range test.Files {
			switch path {
			case "go.mod":
				module, err := modFinder.Parse(path, []byte(code))
				is.NoErr(err)
				err = module.File().Replace("gitlab.com/mnm/bud", budModule.Directory())
				is.NoErr(err)
				vmap[path] = string(module.File().Format())
			default:
				vmap[path] = redent(code)
			}
		}
		err := vfs.Write(appDir, vmap)
		is.NoErr(err)
	}
	// Setup genFS
	appFS := vfs.OS(appDir)
	genFS := gen.New(appFS)
	modFinder = mod.New(mod.WithCache(modCache), mod.WithFS(genFS))
	module, err := modFinder.Find(".")
	is.NoErr(err)
	parser := parser.New(module)
	injector := di.New(module, parser, di.Map{})
	genFS.Add(map[string]gen.Generator{
		"bud/main.go":    mainGo,
		"bud/web/web.go": webServer,
		"bud/command/command.go": gen.FileGenerator(&command.Generator{
			Module:   module,
			Injector: injector,
		}),
	})
	err = fsync.Dir(genFS, "bud", appFS, "bud")
	is.NoErr(err)
	stdout, err := goRun(modCache.Directory(), appDir, test.Args...)
	is.NoErr(err)
	diff.TestString(t, test.Expect, stdout)
}

const goMod = `
module app.com

require (
  github.com/hexops/valast v1.4.1
	gitlab.com/mnm/bud v0.0.0
)
`

func TestRoot(t *testing.T) {
	runTest(t, Test{
		Files: map[string]string{
			"go.mod": goMod,
		},
		Args:   []string{"-h"},
		Expect: ``,
	})
}
