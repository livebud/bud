package create

import (
	"context"
	_ "embed"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
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
		"livebud": "0.0.0",
		"svelte":  "3.44.1",
	}
	code, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		return err
	}
	if err := os.WriteFile(filepath.Join(dir, "package.json"), code, 0644); err != nil {
		return err
	}
	npm, err := exec.LookPath("npm")
	if err != nil {
		return err
	}
	cmd := exec.CommandContext(ctx, npm, "install", "--loglevel=error", "--no-progress")
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
	// if c.Link {
	// 	cmd := exec.CommandContext(ctx, npm, "link", "--loglevel=error", "livebud")
	// 	cmd.Dir = dir
	// 	cmd.Stderr = os.Stderr
	// 	cmd.Env = []string{
	// 		"HOME=" + os.Getenv("HOME"),
	// 		"PATH=" + os.Getenv("PATH"),
	// 		"NO_COLOR=1",
	// 	}
	// 	if err := cmd.Run(); err != nil {
	// 		return err
	// 	}

	// }
	return nil
}
