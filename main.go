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
	"gitlab.com/mnm/bud/internal/modcache"
	"gitlab.com/mnm/bud/internal/parser"
	v8 "gitlab.com/mnm/bud/js/v8"

	"gitlab.com/mnm/bud/internal/generator/command"
	"gitlab.com/mnm/bud/internal/generator/controller"
	"gitlab.com/mnm/bud/internal/generator/generator"
	"gitlab.com/mnm/bud/internal/generator/gomod"
	"gitlab.com/mnm/bud/internal/generator/maingo"
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
		console.Error(err.Error())
		os.Exit(1)
	}
}

func do() error {
	cmd := new(bud)
	cli := commander.New("bud")
	cli.Flag("chdir", "Change the working directory").Short('C').String(&cmd.Chdir).Default(".")

	{
		cmd := &runCommand{bud: cmd}
		cli := cli.Command("run", "run the development server")
		cli.Flag("embed", "embed the assets").Bool(&cmd.Embed).Default(false)
		cli.Flag("hot", "hot reload the frontend").Bool(&cmd.Hot).Default(true)
		cli.Flag("minify", "minify the assets").Bool(&cmd.Minify).Default(false)
		cli.Flag("port", "port").Int(&cmd.Port).Default(3000)
		cli.Run(cmd.Run)
	}

	{
		cmd := &buildCommand{bud: cmd}
		cli := cli.Command("build", "build the production server")
		cli.Flag("embed", "embed the assets").Bool(&cmd.Embed).Default(true)
		cli.Flag("hot", "hot reload the frontend").Bool(&cmd.Hot).Default(false)
		cli.Flag("minify", "minify the assets").Bool(&cmd.Minify).Default(true)
		cli.Run(cmd.Run)
	}

	{
		cli := cli.Command("tool", "extra tools")

		{ // bud tool di
			cmd := &diCommand{bud: cmd}
			cli := cli.Command("di", "dependency injection generator")
			cli.Flag("dependency", "generate dependency provider").Short('d').Strings(&cmd.Dependencies)
			cli.Flag("external", "mark dependency as external").Short('e').Strings(&cmd.Externals).Optional()
			cli.Flag("map", "map interface types to concrete types").Short('m').StringMap(&cmd.Map).Optional()
			cli.Flag("target", "target import path").Short('t').String(&cmd.Target)
			cli.Flag("hoist", "hoist dependencies that depend on externals").Bool(&cmd.Hoist).Default(false)
			cli.Flag("verbose", "verbose logging").Short('v').Bool(&cmd.Verbose).Default(false)
			cli.Run(cmd.Run)
		}

		{ // bud tool v8
			cmd := &v8Command{bud: cmd}
			cli := cli.Command("v8", "Execute Javascript with V8")
			cli.Arg("eval", "evaluate a script").Strings(&cmd.Eval).Optional()
			cli.Run(cmd.Run)
		}
	}

	cli.Arg("command", "custom command").String(&cmd.Custom)
	cli.Run(cmd.Run)

	return cli.Parse(os.Args[1:])
}

type bud struct {
	Chdir  string
	Custom string
}

func (c *bud) Run(ctx context.Context) error {
	module := mod.New(modcache.Default())
	modfile, err := module.Find(c.Chdir)
	if err != nil {
		return err
	}
	fmt.Println(c.Custom+"ing...", modfile.Directory())
	genfs := gen.New(os.DirFS(modfile.Directory()))
	genfs.Add(map[string]gen.Generator{
		"bud/command/main.go": gen.FileGenerator(&command.Generator{
			// fill in
		}),
		"go.mod": gen.FileGenerator(&gomod.Generator{
			Modfile: modfile,
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
			},
		}),
	})
	// Sync genfs
	if err := fsync.Dir(genfs, ".", vfs.OS(modfile.Directory()), "."); err != nil {
		return err
	}
	// If bud/command/main.go doesn't exist, run the welcome server
	commandPath := filepath.Join(modfile.Directory(), "bud", "command", "main.go")
	if _, err := os.Stat(commandPath); err != nil {
		if !errors.Is(err, fs.ErrNotExist) {
			return err
		}
		return fmt.Errorf("unknown command %q", c.Custom)
	}
	// Run the command
	// TODO: pass all arguments through
	if err := gobin.Run(ctx, modfile.Directory(), commandPath, c.Custom); err != nil {
		return err
	}
	return nil
}

type runCommand struct {
	bud    *bud
	Embed  bool
	Hot    bool
	Minify bool
	Port   int
	Args   []string
}

func (c *runCommand) Run(ctx context.Context) error {
	module := mod.New(modcache.Default())
	modFile, err := module.Find(c.bud.Chdir)
	if err != nil {
		return err
	}
	parser := parser.New(module)
	injector := di.New(modFile, parser, di.Map{})
	genfs := gen.New(os.DirFS(modFile.Directory()))
	genfs.Add(map[string]gen.Generator{
		"go.mod": gen.FileGenerator(&gomod.Generator{
			Modfile: modFile,
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
			Modfile: modFile,
		}),
		"bud/generate/main.go": gen.FileGenerator(&generate.Generator{
			Modfile: modFile,
			Embed:   c.Embed,
			Hot:     c.Hot,
			Minify:  c.Minify,
		}),
		"bud/generator/generator.go": gen.FileGenerator(&generator.Generator{
			// fill in
		}),
		// TODO: separate the following from the generators to give the generators
		// a chance to add files that are picked up by these compiler plugins.
		"bud/command/command.go": gen.FileGenerator(&command.Generator{
			Modfile:  modFile,
			Injector: injector,
		}),
		"bud/controller/controller.go": gen.FileGenerator(&controller.Generator{
			Modfile: modFile,
		}),
		"bud/transform/transform.go": gen.FileGenerator(&transform.Generator{
			Modfile: modFile,
		}),
		"bud/view/view.go": gen.FileGenerator(&view.Generator{
			Modfile: modFile,
		}),
		"bud/public/public.go": gen.FileGenerator(&public.Generator{
			Modfile: modFile,
			Embed:   c.Embed,
			Minify:  c.Minify,
		}),
		"bud/web/web.go": gen.FileGenerator(&web.Generator{
			Modfile: modFile,
		}),
		"bud/main.go": gen.FileGenerator(&maingo.Generator{
			Modfile: modFile,
		}),
	})
	// Sync genfs
	if err := fsync.Dir(genfs, ".", vfs.OS(modFile.Directory()), "."); err != nil {
		return err
	}
	// Intentionally use a different context for running subprocesses because
	// the subprocess should be the one handling the interrupt, not the parent
	// process.
	runCtx, cancel := context.WithCancel(context.Background())
	defer cancel()
	// Run generate (if it exists) to support user-defined generators
	generatePath := filepath.Join(modFile.Directory(), "bud", "generate", "main.go")
	if _, err := os.Stat(generatePath); nil == err {
		if err := gobin.Run(runCtx, modFile.Directory(), generatePath); err != nil {
			return err
		}
	}
	// If bud/main.go doesn't exist, run the welcome server
	mainPath := filepath.Join(modFile.Directory(), "bud", "main.go")
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
	// Run the main server
	if err := gobin.Run(runCtx, modFile.Directory(), mainPath); err != nil {
		return err
	}
	return nil
}

type buildCommand struct {
	bud    *bud
	Embed  bool
	Hot    bool
	Minify bool
}

func (c *buildCommand) Run(ctx context.Context) error {
	module := mod.New(modcache.Default())
	modFile, err := module.Find(c.bud.Chdir)
	if err != nil {
		return err
	}
	parser := parser.New(module)
	injector := di.New(modFile, parser, di.Map{})
	fmt.Println("building...", modFile.Directory(), c.Embed, c.Hot, c.Minify)
	genfs := gen.New(os.DirFS(modFile.Directory()))
	genfs.Add(map[string]gen.Generator{
		"bud/generate/main.go": gen.FileGenerator(&generate.Generator{
			Modfile: modFile,
			Embed:   c.Embed,
			Hot:     c.Hot,
			Minify:  c.Minify,
		}),
		"bud/generator/generator.go": gen.FileGenerator(&generator.Generator{
			// fill in
		}),
		"bud/command/command.go": gen.FileGenerator(&command.Generator{
			Modfile:  modFile,
			Injector: injector,
		}),
		"go.mod": gen.FileGenerator(&gomod.Generator{
			Modfile: modFile,
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
		"bud/controller/controller.go": gen.FileGenerator(&controller.Generator{
			// Fill in
		}),
		"bud/web/web.go": gen.FileGenerator(&web.Generator{
			Modfile: modFile,
		}),
		"bud/main.go": gen.FileGenerator(&maingo.Generator{
			// fill in
		}),
	})
	// Sync genfs
	if err := fsync.Dir(genfs, ".", vfs.OS(modFile.Directory()), "."); err != nil {
		return err
	}
	// Run generate (if it exists) to support user-defined generators
	generatePath := filepath.Join(modFile.Directory(), "bud", "generate", "main.go")
	if _, err := os.Stat(generatePath); nil == err {
		if err := gobin.Run(ctx, modFile.Directory(), generatePath); err != nil {
			return err
		}
	}
	// Verify that bud/main.go exists
	mainPath := filepath.Join(modFile.Directory(), "bud", "main.go")
	if _, err := os.Stat(mainPath); err != nil {
		return err
	}
	// Build the main server
	outPath := filepath.Join(modFile.Directory(), "bud", "main")
	if err := gobin.Build(ctx, modFile.Directory(), mainPath, outPath); err != nil {
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
	module := mod.New(modcache.Default())
	modfile, err := module.Find(c.bud.Chdir)
	if err != nil {
		return err
	}
	parser := parser.New(module)
	fn := &di.Function{
		Hoist: c.Hoist,
	}
	fn.Target, err = c.toImportPath(modfile, c.Target)
	if err != nil {
		return err
	}
	typeMap := di.Map{}
	// Add the type mapping
	for from, to := range c.Map {
		fromDep, err := c.toDependency(modfile, from)
		if err != nil {
			return err
		}
		toDep, err := c.toDependency(modfile, to)
		if err != nil {
			return err
		}
		typeMap[fromDep] = toDep
	}
	// Add the dependencies
	for _, dependency := range c.Dependencies {
		dep, err := c.toDependency(modfile, dependency)
		if err != nil {
			return err
		}
		fn.Results = append(fn.Results, dep)
	}
	// Add the externals
	for _, external := range c.Externals {
		ext, err := c.toDependency(modfile, external)
		if err != nil {
			return err
		}
		fn.Params = append(fn.Params, ext)
	}
	injector := di.New(modfile, parser, typeMap)
	node, err := injector.Load(fn)
	if err != nil {
		return err
	}
	if c.Verbose {
		fmt.Println(node.Print())
	}
	provider := node.Generate(fn.Target)
	fmt.Println(provider.File("Load"))
	return nil
}

// This should handle both stdlib (e.g. "net/http"), directories (e.g. "web"),
// and dependencies
func (c *diCommand) toImportPath(modfile *mod.File, importPath string) (string, error) {
	importPath = strings.Trim(importPath, "\"")
	maybeDir := modfile.Directory(importPath)
	if _, err := os.Stat(maybeDir); err == nil {
		importPath, err = modfile.ResolveImport(maybeDir)
		if err != nil {
			return "", fmt.Errorf("di: unable to resolve import %s because %+s", importPath, err)
		}
	}
	return importPath, nil
}

func (c *diCommand) toDependency(modfile *mod.File, dependency string) (di.Dependency, error) {
	i := strings.LastIndex(dependency, ".")
	if i < 0 {
		return nil, fmt.Errorf("di: external must have form '<import>.<type>'. got %q ", dependency)
	}
	importPath, err := c.toImportPath(modfile, dependency[0:i])
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
	fmt.Println(result)
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
