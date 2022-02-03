package main

import (
	"context"
	"encoding/gob"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"gitlab.com/mnm/bud/pkg/socket"

	"gitlab.com/mnm/bud/internal/dsync"
	"gitlab.com/mnm/bud/pkg/vfs"

	"gitlab.com/mnm/bud/pkg/budfs"
	"gitlab.com/mnm/bud/pkg/gen"

	"github.com/mattn/go-isatty"
	"gitlab.com/mnm/bud/generator/generate"
	generatorGenerator "gitlab.com/mnm/bud/generator/generator"
	"gitlab.com/mnm/bud/internal/gitignore"
	"gitlab.com/mnm/bud/internal/gobin"
	"gitlab.com/mnm/bud/pkg/di"
	"gitlab.com/mnm/bud/pkg/generator"
	v8 "gitlab.com/mnm/bud/pkg/js/v8"
	"gitlab.com/mnm/bud/pkg/parser"

	"gitlab.com/mnm/bud/pkg/gomod"

	"gitlab.com/mnm/bud/pkg/commander"

	"gitlab.com/mnm/bud/pkg/log/console"
)

func main() {
	if err := do(); err != nil {
		if !isExitStatus(err) {
			console.Error(err.Error())
		}
		os.Exit(1)
	}
}

func do() error {
	// $ bud
	bud := new(bud)
	cli := commander.New("bud")
	cli.Flag("chdir", "Change the working directory").Short('C').String(&bud.Chdir).Default(".")
	cli.Args("command", "custom command").Strings(&bud.Args)
	cli.Run(bud.Run)

	{ // $ bud run
		cmd := &runCommand{bud: bud}
		cli := cli.Command("run2", "run the development server")
		cli.Flag("embed", "embed the assets").Bool(&bud.Embed).Default(false)
		cli.Flag("hot", "hot reload the frontend").Bool(&bud.Hot).Default(true)
		cli.Flag("minify", "minify the assets").Bool(&bud.Minify).Default(false)
		cli.Flag("port", "port").Int(&cmd.Port).Default(3000)
		cli.Run(cmd.Run)
	}

	{ // $ bud run
		cmd := &runCommand2{bud: bud}
		cli := cli.Command("run", "run the development server")
		cli.Flag("embed", "embed the assets").Bool(&bud.Embed).Default(false)
		cli.Flag("hot", "hot reload the frontend").Bool(&bud.Hot).Default(true)
		cli.Flag("minify", "minify the assets").Bool(&bud.Minify).Default(false)
		cli.Flag("port", "port").String(&cmd.Port).Default("3000")
		cli.Run(cmd.Run)
	}

	{ // $ bud build
		cmd := &buildCommand{bud: bud}
		cli := cli.Command("build", "build the production server")
		cli.Flag("embed", "embed the assets").Bool(&bud.Embed).Default(true)
		cli.Flag("hot", "hot reload the frontend").Bool(&bud.Hot).Default(false)
		cli.Flag("minify", "minify the assets").Bool(&bud.Minify).Default(true)
		cli.Run(cmd.Run)
	}

	{ // $ bud tool
		cli := cli.Command("tool", "extra tools")

		{ // $ bud tool di
			cmd := &diCommand{bud: bud}
			cli := cli.Command("di", "dependency injection generator")
			cli.Flag("dependency", "generate dependency provider").Short('d').Strings(&cmd.Dependencies)
			cli.Flag("external", "mark dependency as external").Short('e').Strings(&cmd.Externals).Optional()
			cli.Flag("map", "map interface types to concrete types").Short('m').StringMap(&cmd.Map).Optional()
			cli.Flag("target", "target import path").Short('t').String(&cmd.Target)
			cli.Flag("hoist", "hoist dependencies that depend on externals").Bool(&cmd.Hoist).Default(false)
			cli.Flag("verbose", "verbose logging").Short('v').Bool(&cmd.Verbose).Default(false)
			cli.Run(cmd.Run)
		}

		{ // $ bud tool v8
			cmd := &v8Command{bud: bud}
			cli := cli.Command("v8", "Execute Javascript with V8 from stdin")
			cli.Run(cmd.Run)

			{ // $ bud tool v8 client
				cmd := &v8ClientCommand{bud: bud}
				cli := cli.Command("client", "V8 client used during development")
				cli.Run(cmd.Run)
			}
		}
	}

	return cli.Parse(os.Args[1:])
}

type bud struct {
	Chdir  string
	Embed  bool
	Hot    bool
	Minify bool
	Args   []string
}

func (c *bud) Build(ctx context.Context, dir string) (string, error) {
	generator, err := generator.Load(dir)
	if err != nil {
		return "", err
	}
	if err := generator.Generate(ctx); err != nil {
		return "", err
	}
	mainPath := filepath.Join(dir, "bud", "main.go")
	// Check to see if we generated a main.go
	if _, err := os.Stat(mainPath); err != nil {
		return "", err
	}
	cacheDir, err := os.UserCacheDir()
	if err != nil {
		return "", err
	}
	// Building over an existing binary is faster for some reason, so we'll use
	// the cache directory for a consistent place to output builds
	binPath := filepath.Join(cacheDir, filepath.ToSlash(generator.Module().Import()), "bud", "main")
	if err := gobin.Build(ctx, dir, mainPath, binPath); err != nil {
		return "", err
	}
	return binPath, nil
}

// Run a custom command
func (c *bud) Run(ctx context.Context) error {
	// Find the project directory
	dir, err := gomod.Absolute(c.Chdir)
	if err != nil {
		return err
	}
	// Generate the code
	binPath, err := c.Build(ctx, dir)
	if err != nil {
		if !errors.Is(err, fs.ErrNotExist) {
			return err
		}
		return fmt.Errorf("unknown command %q", c.Args)
	}
	// Run the built binary
	cmd := exec.Command(binPath, c.Args...)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	err = cmd.Run()
	if err != nil {
		return err
	}
	return nil
}

type runCommand2 struct {
	bud  *bud
	Port string
	Args []string
}

func (c *runCommand2) Run(ctx context.Context) error {
	listener, err := socket.Load(c.Port)
	if err != nil {
		return err
	}
	console.Info("Listening on http://%s", listener.Addr().String())
	module, err := gomod.Find(c.bud.Chdir)
	if err != nil {
		return err
	}
	bfs, err := budfs.Load(module)
	if err != nil {
		return err
	}
	bfs.Entry("bud/generator/generator.go", gen.FileGenerator(&generatorGenerator.Generator{
		BFS:    bfs,
		Module: module,
		Embed:  c.bud.Embed,
		Hot:    c.bud.Hot,
		Minify: c.bud.Minify,
	}))
	parser := parser.New(bfs, module)
	injector := di.New(bfs, module, parser)
	bfs.Entry("bud/generate/main.go", gen.FileGenerator(&generate.Generator{
		BFS:      bfs,
		Injector: injector,
		Module:   module,
		Embed:    c.bud.Embed,
		Hot:      c.bud.Hot,
		Minify:   c.bud.Minify,
	}))
	appFS := vfs.OS(module.Directory())
	skipOption := dsync.WithSkip(
		gitignore.FromFS(appFS),
		// Keep bud/main around to improve build caching
		func(name string, isDir bool) bool {
			return !isDir && name == "bud/main"
		},
	)
	if err := dsync.Dir(vfs.SingleFlight(bfs), "bud", appFS, "bud", skipOption); err != nil {
		return err
	}
	mainPath := filepath.Join(module.Directory(), "bud", "generate", "main.go")
	// Check to see if we generated a main.go
	if _, err := os.Stat(mainPath); err != nil {
		return err
	}
	// Building over an existing binary is faster for some reason, so we'll use
	// the cache directory for a consistent place to output builds
	binPath := filepath.Join(module.Directory(), "bud", "generate", "main")
	if err := gobin.Build(ctx, module.Directory(), mainPath, binPath); err != nil {
		return err
	}
	// Pass the socket through
	files, env, err := socket.Files(listener)
	if err != nil {
		return err
	}
	// Run the generator
	cmd := exec.Command(binPath, c.Args...)
	cmd.Env = append(os.Environ(), string(env))
	cmd.ExtraFiles = files
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Dir = module.Directory()
	err = cmd.Run()
	if err != nil {
		return err
	}
	return nil
}

type runCommand struct {
	bud  *bud
	Port int
	Args []string
}

func (c *runCommand) Run(ctx context.Context) error {
	// Find the project directory
	dir, err := gomod.Absolute(c.bud.Chdir)
	if err != nil {
		return err
	}
	// Generate the code
	binPath, err := c.bud.Build(ctx, dir)
	if err != nil {
		if !errors.Is(err, fs.ErrNotExist) {
			return err
		}
		// TODO: improve the welcome server
		address := fmt.Sprintf(":%d", c.Port)
		console.Info("Listening on http://localhost%s", address)
		return http.ListenAndServe(address, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("Welcome Server!\n"))
		}))
	}
	// Run the app
	cmd := exec.Command(binPath, c.Args...)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	err = cmd.Run()
	if err != nil {
		return err
	}
	return nil
}

type buildCommand struct {
	bud *bud
}

func (c *buildCommand) Run(ctx context.Context) error {
	// Find the project directory
	dir, err := gomod.Absolute(c.bud.Chdir)
	if err != nil {
		return err
	}
	// Build the code
	if _, err := c.bud.Build(ctx, dir); err != nil {
		return err
	}
	return nil
}

type diCommand struct {
	bud          *bud
	Target       string
	Map          map[string]string
	Dependencies []string
	Externals    []string
	Hoist        bool
	Verbose      bool
}

func (c *diCommand) Run(ctx context.Context) error {
	module, err := gomod.Find(c.bud.Chdir)
	if err != nil {
		return err
	}
	// TODO: should budfs be empty or fully-loaded with generators?
	bfs, err := budfs.Load(module)
	if err != nil {
		return err
	}
	parser := parser.New(bfs, module)
	fn := &di.Function{
		Hoist: c.Hoist,
	}
	fn.Target, err = c.toImportPath(module, c.Target)
	if err != nil {
		return err
	}
	// Add the type mapping
	for from, to := range c.Map {
		fromDep, err := c.toDependency(module, from)
		if err != nil {
			return err
		}
		toDep, err := c.toDependency(module, to)
		if err != nil {
			return err
		}
		fn.Aliases[fromDep] = toDep
	}
	// Add the dependencies
	for _, dependency := range c.Dependencies {
		dep, err := c.toDependency(module, dependency)
		if err != nil {
			return err
		}
		fn.Results = append(fn.Results, dep)
	}
	// Add the externals
	for _, external := range c.Externals {
		ext, err := c.toDependency(module, external)
		if err != nil {
			return err
		}
		fn.Params = append(fn.Params, ext)
	}
	injector := di.New(bfs, module, parser)
	node, err := injector.Load(fn)
	if err != nil {
		return err
	}
	if c.Verbose {
		fmt.Println(node.Print())
	}
	provider := node.Generate("Load", fn.Target)
	fmt.Fprintln(os.Stdout, provider.File())
	return nil
}

// This should handle both stdlib (e.g. "net/http"), directories (e.g. "web"),
// and dependencies
func (c *diCommand) toImportPath(module *gomod.Module, importPath string) (string, error) {
	importPath = strings.Trim(importPath, "\"")
	maybeDir := module.Directory(importPath)
	if _, err := os.Stat(maybeDir); err == nil {
		importPath, err = module.ResolveImport(maybeDir)
		if err != nil {
			return "", fmt.Errorf("di: unable to resolve import %s because %+s", importPath, err)
		}
	}
	return importPath, nil
}

func (c *diCommand) toDependency(module *gomod.Module, dependency string) (di.Dependency, error) {
	i := strings.LastIndex(dependency, ".")
	if i < 0 {
		return nil, fmt.Errorf("di: external must have form '<import>.<type>'. got %q ", dependency)
	}
	importPath, err := c.toImportPath(module, dependency[0:i])
	if err != nil {
		return nil, err
	}
	dataType := dependency[i+1:]
	// Create the dependency
	return &di.Type{
		Import: importPath,
		Type:   dataType,
	}, nil
}

type v8Command struct {
	bud *bud
}

func (c *v8Command) Run(ctx context.Context) error {
	script, err := c.getScript()
	if err != nil {
		return err
	}
	vm := v8.New()
	result, err := vm.Eval("script.js", script)
	if err != nil {
		return err
	}
	fmt.Fprintln(os.Stdout, result)
	return nil
}

func (c *v8Command) getScript() (string, error) {
	code, err := ioutil.ReadAll(stdin())
	if err != nil {
		return "", err
	}
	script := string(code)
	if script == "" {
		return "", errors.New("missing script to evaluate")
	}
	return script, nil
}

type v8ClientCommand struct {
	bud *bud
}

func (c *v8ClientCommand) Run(ctx context.Context) error {
	vm := v8.New()
	dec := gob.NewDecoder(os.Stdin)
	enc := gob.NewEncoder(os.Stdout)
	for {
		var expr string
		if err := dec.Decode(&expr); err != nil {
			if errors.Is(err, io.EOF) {
				return nil
			}
			return err
		}
		result, err := vm.Eval("<stdin>", string(expr))
		if err != nil {
			return err
		}
		if err := enc.Encode(result); err != nil {
			return err
		}
	}
}

// input from stdin or empty object by default.
func stdin() io.Reader {
	if isatty.IsTerminal(os.Stdin.Fd()) {
		return strings.NewReader("")
	}
	return os.Stdin
}

func toType(importPath, dataType string) *di.Type {
	return &di.Type{Import: importPath, Type: dataType}
}

func isExitStatus(err error) bool {
	return strings.Contains(err.Error(), "exit status ")
}
