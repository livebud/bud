package create

import (
	"context"
	_ "embed"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/livebud/bud/internal/version"
)

func (c *Command) generatePackageJSON(ctx context.Context, dir, name string) error {
	type Dependency struct {
		Name, Version string
	}
	var state struct {
		Name         string            `json:"name,omitempty"`
		Private      bool              `json:"private"`
		Dependencies map[string]string `json:"dependencies,omitempty"`
	}
	state.Name = name
	state.Private = true
	state.Dependencies = map[string]string{
		"livebud": version.Bud,
		"svelte":  version.Svelte,
	}
	code, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		return err
	}
	if err := os.WriteFile(filepath.Join(dir, "package.json"), code, 0644); err != nil {
		return err
	}
	npmPath, err := exec.LookPath("npm")
	if err != nil {
		return err
	}
	cmd := exec.CommandContext(ctx, npmPath, "install", "--loglevel=error", "--no-progress", "--save")
	cmd.Dir = dir
	cmd.Stderr = os.Stderr
	cmd.Env = []string{
		"HOME=" + os.Getenv("HOME"),
		"PATH=" + os.Getenv("PATH"),
		"NO_COLOR=1",
	}
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}
