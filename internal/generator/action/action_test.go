package action_test

import (
	"bytes"
	"context"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"

	"github.com/matryer/is"
	"gitlab.com/mnm/bud/fsync"
	"gitlab.com/mnm/bud/gen"
	"gitlab.com/mnm/bud/internal/di"
	"gitlab.com/mnm/bud/internal/generator/action"
	"gitlab.com/mnm/bud/internal/generator/command"
	"gitlab.com/mnm/bud/internal/generator/maingo"
	"gitlab.com/mnm/bud/internal/generator/program"
	"gitlab.com/mnm/bud/internal/generator/web"
	"gitlab.com/mnm/bud/internal/modcache"
	"gitlab.com/mnm/bud/internal/modtest"
	"gitlab.com/mnm/bud/internal/parser"
	"gitlab.com/mnm/bud/internal/test"
	"gitlab.com/mnm/bud/socket"
	"gitlab.com/mnm/bud/vfs"
)

// func redent(s string) string {
// 	return strings.TrimSpace(dedent.Dedent(s)) + "\n"
// }

// func isEqual(t testing.TB, actual, expect string) {
// 	diff.TestString(t, redent(expect), redent(actual))
// }

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
		"bud/action/action.go": gen.FileGenerator(&action.Generator{
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
	gitlab.com/mnm/bud v0.0.0
	gitlab.com/mnm/bud-tailwind v0.0.0-20211228175933-3ca601f1a518
)
`

const hnClient = `
package hn

import "context"

func New() *Client {
	return &Client{"https://news.ycombinator.com/"}
}

type Client struct {
	baseURL string
}

func (c *Client) FrontPage(ctx context.Context) (string, error) {
	return "https://news.ycombinator.com/", nil
}

func (c *Client) Find(ctx context.Context, id string) (string, error) {
	return "https://news.ycombinator.com/item?id=" + id, nil
}
`

func TestNoController(t *testing.T) {
	t.SkipNow()
}

func TestNoMethods(t *testing.T) {
	t.SkipNow()
}

func TestSimple(t *testing.T) {
	t.SkipNow()
	is := is.New(t)
	generator := test.Generator(t)
	generator.Files["go.mod"] = goMod
	generator.Files["internal/hn/client.go"] = hnClient
	generator.Files["action/action.go"] = `
		package action

		import (
			"context"

			"app.com/internal/hn"
		)

		type Controller struct {
			HN *hn.Client
		}

		func (c *Controller) Index(ctx context.Context) (string, error) {
			return c.HN.FrontPage(ctx)
		}

		func (c *Controller) Show(ctx context.Context, id string) (string, error) {
			return c.HN.Find(ctx, id)
		}
	`
	app, err := generator.Generate()
	is.NoErr(err)
	app.Exists("bud/action/action.go")
	app.Exists("bud/main.go")
	socketPath := filepath.Join(t.TempDir(), "tmp.sock")
	listener, err := socket.Listen(socketPath)
	is.NoErr(err)
	defer listener.Close()
	files, env, err := socket.Files(listener)
	is.NoErr(err)
	app.ExtraFiles(files...)
	app.Env(env.Key(), env.Value())
	cmd, err := app.Start()
	is.NoErr(err)
	defer cmd.Close()
	transport, err := socket.Transport(socketPath)
	is.NoErr(err)
	client := http.Client{
		Timeout:   time.Second,
		Transport: transport,
	}
	res, err := client.Get("http://host/")
	is.NoErr(err)
	is.Equal(res.StatusCode, 200)
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	is.NoErr(err)
	is.Equal(string(body), `https://news.ycombinator.com/`)
}
