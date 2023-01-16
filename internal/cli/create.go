package cli

import (
	"context"
	_ "embed"
	"encoding/json"
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/livebud/bud/internal/gotemplate"

	"github.com/livebud/bud/internal/embedded"
	"github.com/livebud/bud/internal/versions"
	"github.com/livebud/bud/package/gomod"
	"github.com/livebud/bud/package/virtual"
	"golang.org/x/mod/modfile"
)

type Create struct {
	Dir    string
	Module string
	Dev    bool
}

type createModule struct {
	Name      string
	GoVersion string
	Requires  []*createRequire
	Replaces  []*createReplace
}

type createRequire struct {
	Import   string
	Version  string
	Indirect bool
}

type createReplace struct {
	From string
	To   string
}

type createPackage struct {
	Name         string            `json:"name,omitempty"`
	Private      bool              `json:"private,omitempty"`
	Dependencies map[string]string `json:"dependencies,omitempty"`
}

func (c *CLI) Create(ctx context.Context, in *Create) error {
	var files []*virtual.File

	// Get the absolute directory
	absDir, err := filepath.Abs(in.Dir)
	if err != nil {
		return err
	}

	// create go.mod
	gomodFile, err := c.createGoMod(in, absDir)
	if err != nil {
		return err
	}
	files = append(files, gomodFile)

	pkgFile, err := c.createPackageJson(in, absDir)
	if err != nil {
		return err
	}
	files = append(files, pkgFile)

	// Add the .gitignore
	files = append(files, &virtual.File{
		Path: ".gitignore",
		Data: embedded.Gitignore(),
	})

	// Add favicon.ico
	files = append(files, &virtual.File{
		Path: "public/favicon.ico",
		Data: embedded.Favicon(),
	})

	// Sync the files
	for _, file := range files {
		abspath := filepath.Join(absDir, file.Path)
		if err := os.MkdirAll(filepath.Dir(abspath), 0755); err != nil {
			return err
		}
		if err := os.WriteFile(abspath, file.Data, 0644); err != nil {
			return err
		}
	}

	// Download the dependencies in go.mod to GOMODCACHE
	// Run `go mod download all`
	// TODO: do we need `all`?
	if err := c.command(absDir, "go", "mod", "download", "all").Run(); err != nil {
		return err
	}

	// Install node_modules
	npmInstall := c.command(absDir, "npm", "install", "--no-audit", "--loglevel=error", "--no-progress", "--save")
	if err := npmInstall.Run(); err != nil {
		return err
	}

	if in.Dev {
		budModule, err := c.findBudModule()
		if err != nil {
			return err
		}
		// Link node_modules
		npmLink := c.command(absDir, "npm", "link", "--no-audit", "--loglevel=error", "livebud", budModule.Directory("livebud"))
		if err := npmLink.Run(); err != nil {
			return err
		}
	}

	return nil
}

func (c *CLI) createGoMod(in *Create, absDir string) (*virtual.File, error) {
	// Create the module
	module := &createModule{
		Name:      in.Module,
		GoVersion: goVersion(),
	}

	// Update the module name
	if module.Name == "" {
		// Try inferring the module name from the directory
		module.Name = gomod.Infer(absDir)
		// if can't infer then default it to `change.me`
		if module.Name == "" {
			module.Name = "change.me"
		}
	}
	// Autoquote the module name
	module.Name = modfile.AutoQuote(module.Name)

	// Add the required runtime
	runtime := &createRequire{
		Import:  "github.com/livebud/bud",
		Version: versions.Bud,
	}
	if in.Dev && versions.Bud == "latest" {
		runtime.Version = "v0.0.0"
	}
	module.Requires = append(module.Requires, runtime)

	// Link to local copy
	if in.Dev {
		budModule, err := c.findBudModule()
		if err != nil {
			return nil, err
		}
		module.Replaces = append(module.Replaces, &createReplace{
			From: "github.com/livebud/bud",
			To:   modfile.AutoQuote(budModule.Directory()),
		})
	}

	code, err := fs.ReadFile(embedFS, "create/gomod.gotext")
	if err != nil {
		return nil, err
	}
	template, err := gotemplate.Parse("create/gomod.gotext", string(code))
	if err != nil {
		return nil, err
	}
	// Generate go.mod
	data, err := template.Generate(module)
	if err != nil {
		return nil, err
	}
	return &virtual.File{
		Path: "go.mod",
		Data: data,
	}, nil
}

func (c *CLI) createPackageJson(in *Create, absDir string) (*virtual.File, error) {
	// Create the package
	pkg := &createPackage{
		Name:    filepath.Base(absDir),
		Private: true,
		Dependencies: map[string]string{
			"livebud": versions.Bud,
			"svelte":  versions.Svelte,
		},
	}
	// Generate package.json
	data, err := json.MarshalIndent(pkg, "", "  ")
	if err != nil {
		return nil, err
	}
	return &virtual.File{
		Path: "package.json",
		Data: data,
	}, nil
}

func goVersion() string {
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

// func (c *CLI) Run(ctx context.Context) (err error) {
// 	// Get the absolutely directory
// 	c.absDir, err = filepath.Abs(c.Dir)
// 	if err != nil {
// 		return err
// 	}
// 	// If we're linking to the development version of Bud, we need to
// 	// find Bud's go.mod file.
// 	if c.Dev {
// 		c.budModule, err = mod.FindBudModule()
// 		if err != nil {
// 			return err
// 		}
// 	}
// 	// Load the template state
// 	state, err := c.Load()
// 	if err != nil {
// 		return err
// 	}
// 	// Scaffold files from the template state
// 	return c.Scaffold(state)
// }

// // func (c *CLI) load() (state *State, err error) {
// // 	defer c.bail.Recover2(&err, "create")
// // 	state = new(State)
// // 	state.Module = c.loadModule()
// // 	state.Package = c.loadPackage(filepath.Base(c.Dir))
// // 	return state, nil
// // }

// // func (c *CLI) loadModule() *Module {
// // 	module := new(Module)
// // 	// Get the module path that's passed in as a flag
// // 	module.Name = c.Module
// // 	if module.Name == "" {
// // 		// Try inferring the module name from the directory
// // 		module.Name = mod.Infer(c.absDir)
// // 		// if can't infer then default it to `change.me`
// // 		if module.Name == "" {
// // 			module.Name = "change.me"
// // 		}
// // 	}
// // 	// Autoquote the module name since
// // 	module.Name = modfile.AutoQuote(module.Name)
// // 	// Add the required runtime
// // 	module.Requires = []*Require{
// // 		{
// // 			Import:  "github.com/livebud/bud",
// // 			Version: c.budVersion(),
// // 		},
// // 	}
// // 	// Link to local copy
// // 	if c.Dev {
// // 		module.Replaces = []*Replace{
// // 			{
// // 				From: "github.com/livebud/bud",
// // 				To:   modfile.AutoQuote(c.budModule.Directory()),
// // 			},
// // 		}
// // 	}
// // 	return module
// // }

// // func (c *CLI) loadPackage(name string) *Package {
// // 	pkg := new(Package)
// // 	pkg.Name = name
// // 	pkg.Private = true
// // 	pkg.Dependencies = map[string]string{
// // 		"livebud": versions.Bud,
// // 		"svelte":  versions.Svelte,
// // 	}
// // 	return pkg
// // }

// // //go:embed gomod.gotext
// // var gomod string

// // //go:embed gitignore.gotext
// // var gitignore string

// // // Scaffold state into the specified directory
// // func (c *CLI) scaffold(state *State) error {
// // 	// Scaffold into that directory
// // 	if err := scaffold.Scaffold(scaffold.OSFS(c.absDir),
// // 		scaffold.Template("go.mod", gomod, state.Module),
// // 		scaffold.Template(".gitignore", gitignore, nil),
// // 		scaffold.JSON("package.json", state.Package),
// // 		scaffold.File("public/favicon.ico", embedded.Favicon()),
// // 	); err != nil {
// // 		return err
// // 	}
// // 	// Download the dependencies in go.mod to GOMODCACHE
// // 	// Run `go mod download all`
// // 	// TODO: do we need `all`?
// // 	if err := scaffold.Command(c.absDir, "go", "mod", "download", "all").Run(); err != nil {
// // 		return err
// // 	}
// // 	// Install node modules
// // 	if err := scaffold.Command(c.absDir, "npm", "install", "--loglevel=error", "--no-progress", "--save").Run(); err != nil {
// // 		return err
// // 	}
// // 	if c.Dev {
// // 		// Link node modules
// // 		if err := scaffold.Command(c.absDir, "npm", "link", "--loglevel=error", "livebud", c.budModule.Directory("livebud")).Run(); err != nil {
// // 			return err
// // 		}
// // 	}
// // 	return nil
// // }

// // func (c *Command) budVersion() string {
// // 	version := versions.Bud
// // 	if c.Dev && version == "latest" {
// // 		return "v0.0.0"
// // 	}
// // 	return "v" + version
// // }
