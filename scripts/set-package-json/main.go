package main

import (
	"os"
	"path/filepath"

	"github.com/livebud/bud/internal/npm"
	"github.com/livebud/bud/internal/version"
	"github.com/livebud/bud/package/gomod"
	"github.com/livebud/bud/package/log/console"
)

func main() {
	if err := run(); err != nil {
		console.Error(err.Error())
		os.Exit(1)
	}
	os.Exit(0)
}

func run() error {
	dir, err := gomod.Absolute(".")
	if err != nil {
		return err
	}
	// Update the dependencies in ./livebud/package.json
	if err := npm.Set(filepath.Join(dir, "livebud"), map[string]string{
		"dependencies.svelte":              version.Svelte,
		"dependencies.react":               version.React,
		"dependencies.react-dom":           version.React,
		"devDependencies.@types/react":     version.React,
		"devDependencies.@types/react-dom": version.React,
	}); err != nil {
		return err
	}
	// Update the dependencies in .
	if err := npm.Set(dir, map[string]string{
		"devDependencies.svelte":           version.Svelte,
		"devDependencies.react":            version.React,
		"devDependencies.react-dom":        version.React,
		"devDependencies.@types/react":     version.React,
		"devDependencies.@types/react-dom": version.React,
	}); err != nil {
		return err
	}
	return nil
}
