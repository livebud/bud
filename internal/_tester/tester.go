package tester

import (
	"bytes"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"gitlab.com/mnm/bud/fsync"
	"gitlab.com/mnm/bud/gen"
	"gitlab.com/mnm/bud/go/mod"
	"gitlab.com/mnm/bud/internal/di"
	"gitlab.com/mnm/bud/internal/modcache"
	"gitlab.com/mnm/bud/internal/parser"

	"github.com/lithammer/dedent"
	"github.com/matryer/is"
	"gitlab.com/mnm/bud/vfs"
)

func New(t testing.TB) *Tester {
	t.Helper()
	is := is.New(t)
	appDir := filepath.Join("_tmp", t.Name())
	err := os.MkdirAll(appDir, 0755)
	is.NoErr(err)
	t.Cleanup(cleanup(t, "_tmp", appDir))
	env := env{"PWD": appDir}
	return &Tester{
		TB:       t,
		name:     t.Name(),
		env:      env,
		appDir:   appDir,
		modCache: modcache.Default(),
	}
}

// Cleanup individual files and root if no files left
func cleanup(t testing.TB, root, dir string) func() {
	t.Helper()
	is := is.New(t)
	return func() {
		if t.Failed() {
			return
		}
		is.NoErr(os.RemoveAll(dir))
		fis, err := os.ReadDir(root)
		if err != nil {
			return
		}
		if len(fis) > 0 {
			return
		}
		is.NoErr(os.RemoveAll(root))
	}
}

type env map[string]string

func (env env) List() (list []string) {
	for key, value := range env {
		list = append(list, key+"="+value)
	}
	return list
}

type Tester struct {
	testing.TB
	name     string
	env      env
	appDir   string
	modCache *modcache.Cache

	// Filled in dynamically
	module   *mod.Module
	appFS    vfs.ReadWritable
	genFS    *gen.FileSystem
	parser   *parser.Parser
	injector *di.Injector
}

func (t *Tester) Env(key, value string) *Tester {
	t.Helper()
	t.env[key] = value
	return t
}

// func (t *Tester) GoModCache(dir string) *Tester {
// 	t.Helper()
// 	t.cacheDir = dir
// 	t.env = append(t.env, "GOMODCACHE="+dir)
// 	return t
// }

func (t *Tester) GoPrivate(glob string) *Tester {
	t.Helper()
	t.env["GOPRIVATE"] = glob
	return t
}

func (t *Tester) NoColor() *Tester {
	t.Helper()
	t.env["NO_COLOR"] = "1"
	return t
}

func redent(s string) string {
	return strings.TrimSpace(dedent.Dedent(s)) + "\n"
}

func replaceBud(t testing.TB, code string) string {
	is := is.New(t)
	wd, err := os.Getwd()
	is.NoErr(err)
	budModule, err := mod.New().Find(wd)
	is.NoErr(err)
	module, err := mod.New().Parse("go.mod", []byte(code))
	is.NoErr(err)
	err = module.File().Replace("gitlab.com/mnm/bud", budModule.Directory())
	is.NoErr(err)
	return string(module.File().Format())
}

func (t *Tester) Files(files map[string]string) {
	t.Helper()
	is := is.New(t)
	for path, file := range files {
		if path == "go.mod" {
			files[path] = replaceBud(t, file)
			continue
		}
		files[path] = redent(file)
	}
	err := vfs.Write(t.appDir, vfs.Map(files))
	is.NoErr(err)
}

func (t *Tester) Modules(modules map[string]map[string]string) *Tester {
	t.Helper()
	is := is.New(t)
	cacheDir := t.TempDir()
	t.modCache = modcache.New(cacheDir)
	err := t.modCache.Write(modules)
	is.NoErr(err)
	return t
}

func (t *Tester) Exists(path string) bool {
	is := is.New(t)
	path = filepath.Join(t.appDir, path)
	if _, err := os.Stat(path); err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return false
		}
		is.NoErr(err)
	}
	return true
}

func (t *Tester) AppFS() vfs.ReadWritable {
	if t.appFS != nil {
		return t.appFS
	}
	t.appFS = vfs.OS(t.appDir)
	return t.appFS
}

func (t *Tester) GenFS() *gen.FileSystem {
	if t.genFS != nil {
		return t.genFS
	}
	t.genFS = gen.New(t.AppFS())
	return t.genFS
}

func (t *Tester) Sync() error {
	return fsync.Dir(t.GenFS(), ".", t.AppFS(), ".")
}

func (t *Tester) Module() *mod.Module {
	t.Helper()
	is := is.New(t)
	if t.module != nil {
		return t.module
	}
	modFinder := mod.New(mod.WithCache(t.modCache), mod.WithFS(t.GenFS()))
	module, err := modFinder.Find(".")
	is.NoErr(err)
	t.module = module
	return t.module
}

func (t *Tester) Parser() *parser.Parser {
	if t.parser != nil {
		return t.parser
	}
	return parser.New(t.Module())
}

func (t *Tester) Injector(tm di.Map) *di.Injector {
	if t.injector != nil {
		return t.injector
	}
	return di.New(t.Module(), t.Parser(), tm)
}

// TODO: consolidate with Start
func (t *Tester) Run(command string, args ...string) *CommandResult {
	t.Helper()
	cmd := exec.Command(command, args...)
	cmd.Env = t.env.List()
	cmd.Dir = t.appDir
	stdout := new(bytes.Buffer)
	cmd.Stdout = stdout
	stderr := new(bytes.Buffer)
	cmd.Stderr = stderr
	err := cmd.Run()
	return &CommandResult{
		t:      t,
		err:    err,
		Stdout: &Stdio{stdout, t},
		Stderr: &Stdio{stderr, t},
	}
}

func (t *Tester) GoRun(mainPath string, args ...string) *CommandResult {

	args = append([]string{"run", "-mod=mod", mainPath}, args...)
	return t.Run("go", args...)
}

type Stdio struct {
	*bytes.Buffer
	t *Tester
}

func (s *Stdio) Equal(expect string) {
	is := is.New(s.t)
	is.Equal(s.String(), expect)
}

func (s *Stdio) Match(pattern string) {
	s.t.Match(s.String(), pattern)
}

type CommandResult struct {
	t      *Tester
	err    error
	Stdout *Stdio
	Stderr *Stdio
}

func (c *CommandResult) NoErr() *CommandResult {
	is := is.New(c.t)
	if c.err != nil {
		fmt.Fprint(os.Stderr, c.Stderr.String())
		is.NoErr(c.err)
	}
	return c
}

func (c *CommandResult) Error() string {
	if c.err == nil {
		return ""
	}
	return c.err.Error()
}

func (t *Tester) Start(command string, args ...string) (*Command, error) {
	t.Helper()
	is := is.New(t)
	cmd := exec.Command(command, args...)
	cmd.Env = t.env.List()
	cmd.Dir = t.appDir
	stdout := new(bytes.Buffer)
	stderr := new(bytes.Buffer)
	cmd.Stdout = stdout
	cmd.Stderr = stderr
	err := cmd.Start()
	is.NoErr(err)
	return &Command{cmd, stdout, stderr}, err
}

type Command struct {
	cmd    *exec.Cmd
	stdout *bytes.Buffer
	stderr *bytes.Buffer
}

func (c *Command) Stdout() string {
	return c.stdout.String()
}

func (c *Command) Stderr() string {
	return c.stderr.String()
}

func (c *Command) Wait() error {
	return c.cmd.Wait()
}

func (c *Command) Close() error {
	p := c.cmd.Process
	if p != nil {
		if err := p.Signal(os.Interrupt); err != nil {
			p.Kill()
		}
	}
	return c.cmd.Wait()
}

func (t *Tester) Match(actual, pattern string) {
	t.Helper()
}

func (t *Tester) WaitFile(path string, deadline time.Duration) {

}

func (t *Tester) WaitPort(port int, deadline time.Duration) {

}
