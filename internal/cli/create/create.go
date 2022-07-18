package create

import (
	"context"
	_ "embed"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/livebud/bud/internal/bail"
	"github.com/livebud/bud/internal/cli/bud"
	"github.com/livebud/bud/internal/format"
	"github.com/livebud/bud/internal/scaffold"
	"github.com/livebud/bud/internal/versions"
	mod "github.com/livebud/bud/package/gomod"
	"golang.org/x/mod/modfile"
)

func New(bud *bud.Command, in *bud.Input) *Command {
	return &Command{}
}

type Command struct {
	Log    string
	Dir    string
	Module string
	Dev    bool

	// Private
	bail      bail.Struct
	budModule *mod.Module
	absDir    string
}

type State struct {
	Module  *Module
	Package *Package
}

type Module struct {
	Name     string
	Requires []*Require
	Replaces []*Replace
}

func (m *Module) Version() string {
	version := strings.TrimPrefix(runtime.Version(), "go")
	parts := strings.SplitN(version, ".", 3)
	switch len(parts) {
	case 1:
		return parts[0] + ".0"
	case 2:
		return strings.Join(parts, ".")
	default:
		return strings.Join(parts[0:2], ".")
	}
}

type Require struct {
	Import   string
	Version  string
	Indirect bool
}

type Replace struct {
	From string
	To   string
}

type Package struct {
	Name         string            `json:"name,omitempty"`
	Private      bool              `json:"private,omitempty"`
	Dependencies map[string]string `json:"dependencies,omitempty"`
}

func (c *Command) Run(ctx context.Context) (err error) {
	// Get the absolutely directory
	c.absDir, err = filepath.Abs(c.Dir)
	if err != nil {
		return err
	}
	// If we're linking to the development version of Bud, we need to
	// find Bud's go.mod file.
	if c.Dev {
		c.budModule, err = bud.BudModule()
		if err != nil {
			return err
		}
	}
	// Load the template state
	state, err := c.Load()
	if err != nil {
		return err
	}
	// Scaffold files from the template state
	return c.Scaffold(state)
}

func (c *Command) Load() (state *State, err error) {
	defer c.bail.Recover2(&err, "create")
	state = new(State)
	state.Module = c.loadModule()
	state.Package = c.loadPackage(filepath.Base(c.Dir))
	return state, nil
}

func (c *Command) loadModule() *Module {
	module := new(Module)
	// Get the module path that's passed in as a flag
	module.Name = c.Module
	if module.Name == "" {
		// Try inferring the module name from the directory
		module.Name = mod.Infer(c.absDir)
		if module.Name == "" {
			// Fail that you need to pass in a module path
			c.bail.Bail(format.Errorf(`
			Unable to infer a module name. Try again using the module <path> name.

			For example,
				bud create --module=github.com/my/app %s
		`, c.Dir))
		}
	}
	// Autoquote the module name since
	module.Name = modfile.AutoQuote(module.Name)
	// Add the required runtime
	module.Requires = []*Require{
		{
			Import:  "github.com/livebud/bud",
			Version: c.budVersion(),
		},
	}
	// Link to local copy
	if c.Dev {
		module.Replaces = []*Replace{
			{
				From: "github.com/livebud/bud",
				To:   modfile.AutoQuote(c.budModule.Directory()),
			},
		}
	}
	return module
}

func (c *Command) loadPackage(name string) *Package {
	pkg := new(Package)
	pkg.Name = name
	pkg.Private = true
	pkg.Dependencies = map[string]string{
		"livebud": versions.Bud,
		"svelte":  versions.Svelte,
	}
	return pkg
}

//go:embed gomod.gotext
var gomod string

//go:embed gitignore.gotext
var gitignore string

// Scaffold state into the specified directory
func (c *Command) Scaffold(state *State) error {
	// Scaffold into that directory
	if err := scaffold.Scaffold(scaffold.OSFS(c.absDir),
		scaffold.Template("go.mod", gomod, state.Module),
		scaffold.Template(".gitignore", gitignore, nil),
		scaffold.JSON("package.json", state.Package),
	); err != nil {
		return err
	}
	// Download the dependencies in go.mod to GOMODCACHE
	// Run `go mod download all`
	// TODO: do we need `all`?
	if err := scaffold.Command(c.absDir, "go", "mod", "download", "all").Run(); err != nil {
		return err
	}
	// Install node modules
	if err := scaffold.Command(c.absDir, "npm", "install", "--loglevel=error", "--no-progress", "--save").Run(); err != nil {
		return err
	}
	if c.Dev {
		// Link node modules
		if err := scaffold.Command(c.absDir, "npm", "link", "--loglevel=error", "livebud", c.budModule.Directory("livebud")).Run(); err != nil {
			return err
		}
	}
	return nil
}

func (c *Command) budVersion() string {
	version := versions.Bud
	if c.Dev && version == "latest" {
		return "v0.0.0"
	}
	return "v" + version
}
