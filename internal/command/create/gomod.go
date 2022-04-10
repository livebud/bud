package create

import (
	"context"
	_ "embed"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/Bowery/prompt"
	"gitlab.com/mnm/bud/internal/gotemplate"
	"gitlab.com/mnm/bud/package/gomod"
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
	state.Version = strings.TrimPrefix(runtime.Version(), "go")
	// Add the required dependencies
	state.Requires = []*Require{
		{
			Import:  "gitlab.com/mnm/bud",
			Version: "v0.0.0",
		},
	}
	if c.Link {
		budModule, err := gomod.Find(wd)
		if err != nil {
			return err
		}
		state.Replaces = []*Replace{
			{
				From: "gitlab.com/mnm/bud",
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
	return nil
}
