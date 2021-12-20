package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/mattn/go-isatty"
	"gitlab.com/mnm/bud/internal/di"
	"gitlab.com/mnm/bud/internal/gobin"
	"gitlab.com/mnm/bud/internal/parser"
	v8 "gitlab.com/mnm/bud/js/v8"

	"gitlab.com/mnm/bud/internal/generator/command"
	"gitlab.com/mnm/bud/internal/generator/controller"
	"gitlab.com/mnm/bud/internal/generator/generator"
	"gitlab.com/mnm/bud/internal/generator/gomod"
	"gitlab.com/mnm/bud/internal/generator/maingo"
	"gitlab.com/mnm/bud/internal/generator/program"
	"gitlab.com/mnm/bud/internal/generator/public"
	"gitlab.com/mnm/bud/internal/generator/transform"
	"gitlab.com/mnm/bud/internal/generator/view"
	"gitlab.com/mnm/bud/internal/generator/web"
	"gitlab.com/mnm/bud/plugin"

	"gitlab.com/mnm/bud/fsync"
	"gitlab.com/mnm/bud/vfs"

	"gitlab.com/mnm/bud/gen"
	"gitlab.com/mnm/bud/internal/generator/generate"

	"gitlab.com/mnm/bud/go/mod"

	"gitlab.com/mnm/bud/commander"

	"gitlab.com/mnm/bud/log/console"
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
		cli := cli.Command("run", "run the development server")
		cli.Flag("embed", "embed the assets").Bool(&bud.Embed).Default(false)
		cli.Flag("hot", "hot reload the frontend").Bool(&bud.Hot).Default(true)
		cli.Flag("minify", "minify the assets").Bool(&bud.Minify).Default(false)
		cli.Flag("port", "port").Int(&cmd.Port).Default(3000)
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
			cli := cli.Command("v8", "Execute Javascript with V8")
			cli.Arg("eval", "evaluate a script").Strings(&cmd.Eval).Optional()
			cli.Run(cmd.Run)
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

func (c *bud) Generate(dir string) error {
	dirfs := vfs.OS(dir)
	genfs := gen.New(dirfs)
	modFinder := mod.New(mod.WithFS(genfs))
	module, err := modFinder.Find(".")
	if err != nil {
		return err
	}
	parser := parser.New(module)
	injector := di.New(module, parser, di.Map{
		toType("gitlab.com/mnm/bud/gen", "FS"): toType("gitlab.com/mnm/bud/gen", "*FileSystem"),
		toType("gitlab.com/mnm/bud/js", "VM"):  toType("gitlab.com/mnm/bud/js/v8", "*Pool"),
	})
	genfs.Add(map[string]gen.Generator{
		"go.mod": gen.FileGenerator(&gomod.Generator{
			Dir: dir,
			Go: &gomod.Go{
				Version: "1.17",
			},
			Requires: []*gomod.Require{
				{
					Path:    `gitlab.com/mnm/bud`,
					Version: `v0.0.0-20211017185247-da18ff96a31f`,
				},
			},
			// TODO: remove
			Replaces: []*gomod.Replace{
				{
					Old: "gitlab.com/mnm/bud",
					New: "../bud",
				},
				{
					Old: "gitlab.com/mnm/bud-tailwind",
					New: "../bud-tailwind",
				},
			},
		}),
		"bud/plugin": gen.DirGenerator(&plugin.Generator{
			Module: module,
		}),
		"bud/generate/main.go": gen.FileGenerator(&generate.Generator{
			Module: module,
			Embed:  c.Embed,
			Hot:    c.Hot,
			Minify: c.Minify,
		}),
		"bud/generator/generator.go": gen.FileGenerator(&generator.Generator{
			// fill in
		}),
		// TODO: separate the following from the generators to give the generators
		// a chance to add files that are picked up by these compiler plugins.
		"bud/command/command.go": gen.FileGenerator(&command.Generator{
			Module: module,
			Parser: parser,
		}),
		"bud/controller/controller.go": gen.FileGenerator(&controller.Generator{
			Module: module,
		}),
		"bud/transform/transform.go": gen.FileGenerator(&transform.Generator{
			Module: module,
		}),
		"bud/view/view.go": gen.FileGenerator(&view.Generator{
			Module: module,
		}),
		"bud/public/public.go": gen.FileGenerator(&public.Generator{
			Module: module,
			Embed:  c.Embed,
			Minify: c.Minify,
		}),
		"bud/web/web.go": gen.FileGenerator(&web.Generator{
			Module: module,
		}),
		"bud/program/program.go": gen.FileGenerator(&program.Generator{
			Module:   module,
			Injector: injector,
		}),
		"bud/main.go": gen.FileGenerator(&maingo.Generator{
			Module: module,
		}),
	})
	// Sync with the project
	if err := fsync.Dir(module, ".", dirfs, "."); err != nil {
		return err
	}
	// Run generate (if it exists) to support user-defined generators
	generatePath := filepath.Join(module.Directory(), "bud", "generate", "main.go")
	if _, err := os.Stat(generatePath); nil == err {
		if err := gobin.Run(context.Background(), module.Directory(), generatePath); err != nil {
			return err
		}
	}
	return nil
}

func (c *bud) Run(ctx context.Context) error {
	// Find the project directory
	dir, err := mod.FindDirectory(c.Chdir)
	if err != nil {
		return err
	}
	// Generate the code
	if err := c.Generate(dir); err != nil {
		return err
	}
	// Ensure that main.go exists
	mainPath := filepath.Join(dir, "bud", "main.go")
	if _, err := os.Stat(mainPath); err != nil {
		if !errors.Is(err, fs.ErrNotExist) {
			return err
		}
		return fmt.Errorf("unknown command %q", c.Args)
	}
	// Run the command, passing all arguments through
	if err := gobin.Run(ctx, dir, mainPath, c.Args...); err != nil {
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
	dir, err := mod.FindDirectory(c.bud.Chdir)
	if err != nil {
		return err
	}
	// Generate the code
	if err := c.bud.Generate(dir); err != nil {
		return err
	}
	// If bud/main.go doesn't exist, run the welcome server
	mainPath := filepath.Join(dir, "bud", "main.go")
	if _, err := os.Stat(mainPath); err != nil {
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
	// Run the main server. Intentionally use a new background context for running
	// subprocesses because the subprocess should be the one handling the
	// interrupt, not the parent process.
	if err := gobin.Run(context.Background(), dir, mainPath); err != nil {
		return err
	}
	return nil
}

type buildCommand struct {
	bud *bud
}

func (c *buildCommand) Run(ctx context.Context) error {
	// Find the project directory
	dir, err := mod.FindDirectory(c.bud.Chdir)
	if err != nil {
		return err
	}
	// Generate the code
	if err := c.bud.Generate(dir); err != nil {
		return err
	}
	// Verify that bud/main.go exists
	mainPath := filepath.Join(dir, "bud", "main.go")
	if _, err := os.Stat(mainPath); err != nil {
		return err
	}
	// Build the main server
	outPath := filepath.Join(dir, "bud", "main")
	if err := gobin.Build(ctx, dir, mainPath, outPath); err != nil {
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
	modFinder := mod.New()
	module, err := modFinder.Find(c.bud.Chdir)
	if err != nil {
		return err
	}
	parser := parser.New(module)
	fn := &di.Function{
		Hoist: c.Hoist,
	}
	fn.Target, err = c.toImportPath(module, c.Target)
	if err != nil {
		return err
	}
	typeMap := di.Map{}
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
		typeMap[fromDep] = toDep
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
	injector := di.New(module, parser, typeMap)
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
func (c *diCommand) toImportPath(module *mod.Module, importPath string) (string, error) {
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

func (c *diCommand) toDependency(module *mod.Module, dependency string) (di.Dependency, error) {
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
	bud  *bud
	Eval []string
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
	if len(c.Eval) > 0 {
		script := strings.Join(c.Eval, " ")
		return script, nil
	}
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
