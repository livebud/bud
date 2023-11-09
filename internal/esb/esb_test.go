package esb_test

import (
	"fmt"
	"strings"
	"testing"

	esbuild "github.com/evanw/esbuild/pkg/api"
	"github.com/livebud/bud/internal/esb"
	"github.com/livebud/bud/internal/testdir"
)

func isIn(t testing.TB, str, substr string) {
	t.Helper()
	if !strings.Contains(str, substr) {
		t.Fatalf("%q does not contain %q", str, substr)
	}
}

func writeFiles(t testing.TB, dir string, files map[string]string) {
	err := testdir.WriteFiles(dir, files)
	if err != nil {
		t.Fatal(err)
	}
}

// Example serve function
func serve(opts esbuild.BuildOptions) (*esbuild.OutputFile, error) {
	for _, entry := range opts.EntryPoints {
		if !isRelativeEntry(entry) {
			return nil, fmt.Errorf("entry must be relative %q", entry)
		}
	}
	for _, entry := range opts.EntryPointsAdvanced {
		if !isRelativeEntry(entry.InputPath) {
			return nil, fmt.Errorf("entry must be relative %q", entry.InputPath)
		}
	}
	// Set default options
	if opts.Outdir == "" {
		opts.Outdir = "./"
	}
	if opts.Format == esbuild.FormatDefault {
		opts.Format = esbuild.FormatESModule
	}
	if opts.Platform == esbuild.PlatformDefault {
		opts.Platform = esbuild.PlatformBrowser
	}
	if len(opts.Conditions) == 0 {
		// Add "import" condition to support svelte/internal
		// https://esbuild.github.io/api/#how-conditions-work
		opts.Conditions = []string{"browser", "default", "import"}
	}
	// Always bundle, use plugins to granularly mark files as external
	opts.Bundle = true
	result := esbuild.Build(opts)
	// Check if there were errors
	if result.Errors != nil {
		return nil, &esb.Error{result.Errors}
	} else if len(result.OutputFiles) == 0 {
		return nil, fmt.Errorf("esb: no output files")
	}
	return &result.OutputFiles[0], nil
}

func isRelativeEntry(entry string) bool {
	return strings.HasPrefix(entry, "./")
}
