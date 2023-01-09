package versions

import (
	"context"
	"os"
	"os/exec"

	"github.com/livebud/bud/package/gomod"
)

// AlignRuntime ensures that the CLI and runtime versions are aligned.
// If they're not aligned, the CLI will correct the go.mod file to align them.
func AlignRuntime(ctx context.Context, module *gomod.Module, budVersion string) error {
	modfile := module.File()
	// Do nothing for the latest version
	if budVersion == "latest" {
		// If the module file already replaces bud, don't do anything.
		if modfile.Replace(`github.com/livebud/bud`) != nil {
			return nil
		}
		// Best effort attempt to replace bud with the latest version.
		budModule, err := gomod.FindBudModule()
		if err != nil {
			return nil
		}
		// Replace bud with the local version if we found it.
		if err := modfile.AddReplace("github.com/livebud/bud", "", budModule.Directory(), ""); err != nil {
			return err
		}
		// Write the go.mod file back to disk.
		if err := os.WriteFile(module.Directory("go.mod"), modfile.Format(), 0644); err != nil {
			return err
		}
		return nil
	}
	target := "v" + budVersion
	require := modfile.Require("github.com/livebud/bud")
	// We're good, the CLI matches the runtime version
	if require != nil && require.Version == target {
		return nil
	}
	// Otherwise, update go.mod to match the CLI's version
	if err := modfile.AddRequire("github.com/livebud/bud", target); err != nil {
		return err
	}
	if err := os.WriteFile(module.Directory("go.mod"), modfile.Format(), 0644); err != nil {
		return err
	}
	// Run `go mod download`
	cmd := exec.CommandContext(ctx, "go", "mod", "download")
	cmd.Dir = module.Directory()
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = os.Environ()
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}
