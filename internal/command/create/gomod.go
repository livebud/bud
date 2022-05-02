package create

import (
	"context"
	_ "embed"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/livebud/bud/internal/version"

	"github.com/Bowery/prompt"
	"github.com/livebud/bud/internal/gotemplate"
	"github.com/livebud/bud/package/gomod"
)

//go:embed gomod.gotext
var goMod string

func (c *Command) generateGoMod(ctx context.Context, dir string) error {
	generator, err := gotemplate.Parse("go.mod", goMod)
	if err != nil {
		return err
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
	var state struct {
		Name     string
		Version  string
		Requires []*Require
		Replaces []*Replace
	}
	wd, err := os.Getwd()
	if err != nil {
		return err
	}
	absPath := filepath.Join(wd, c.Dir)
	// Get the module name
	state.Name = gomod.Infer(absPath)
	if state.Name == "" {
		state.Name, err = prompt.Basic("Module name? (e.g. github.com/me/app)", true)
		if err != nil {
			return err
		}
	}
	// Get the Go version
	state.Version = strings.TrimPrefix(goVersion(runtime.Version()), "go")
	// Add the required dependencies
	if version.Bud != "latest" {
		state.Requires = []*Require{
			{
				Import:  "github.com/livebud/bud",
				Version: "v" + version.Bud,
			},
		}
	} else {
		// Link to local copy
		state.Requires = []*Require{
			{
				Import:  "github.com/livebud/bud",
				Version: "v0.0.0",
			},
		}
		budModule, err := findBudModule()
		if err != nil {
			return err
		}
		state.Replaces = []*Replace{
			{
				From: "github.com/livebud/bud",
				To:   budModule.Directory(),
			},
		}
	}
	code, err := generator.Generate(state)
	if err != nil {
		return err
	}
	if err := os.WriteFile(filepath.Join(dir, "go.mod"), code, 0644); err != nil {
		return err
	}
	// Download the dependencies in go modules to GOMODCACHE
	cmd := exec.Command("go", "mod", "download", "all")
	cmd.Env = os.Environ()
	cmd.Dir = dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}

// Version can be
func goVersion(version string) string {
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
