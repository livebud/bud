package cli

import (
	"context"
	_ "embed"
	"encoding/json"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/livebud/bud/framework"
	"github.com/livebud/bud/internal/gotemplate"

	"github.com/livebud/bud/internal/embedded"
	"github.com/livebud/bud/internal/versions"
	"github.com/livebud/bud/package/gomod"
	"github.com/livebud/bud/package/virtual"
	"golang.org/x/mod/modfile"
)

type Create struct {
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

	log, err := c.loadLog()
	if err != nil {
		return err
	}

	// Get the absolute directory
	absDir, err := filepath.Abs(c.Dir)
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
		log.Info("Created: %s", file.Path)
	}

	// Download the dependencies in go.mod to GOMODCACHE
	// Run `go mod download all`
	// TODO: do we need `all`?
	if err := c.command(absDir, "go", "mod", "download", "all").Run(); err != nil {
		return err
	}
	log.Info("Installed: go modules")

	// Install node_modules
	npmInstall := c.command(absDir, "npm", "install", "--no-audit", "--loglevel=error", "--no-progress", "--save")
	// Stdout is ignored because there's still "added 1 package" output despite setting log levels
	npmInstall.Stdout = io.Discard
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
		// Stdout is ignored because there's still "added 1 package" output despite setting log levels
		npmLink.Stdout = io.Discard
		if err := npmLink.Run(); err != nil {
			return err
		}
	}
	log.Info("Installed: node modules")

	if err := c.Generate(ctx, &Generate{
		Flag: &framework.Flag{},
	}); err != nil {
		return err
	}
	log.Info("Generated: bud")

	log.Info("Ready: %s", c.Dir)
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
		Version: "v" + versions.Bud,
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
