package overlay_test

import (
	"context"
	"strings"
	"testing"

	"github.com/livebud/bud/internal/testdir"

	"io/fs"

	"github.com/livebud/bud/package/overlay"

	"github.com/livebud/bud/internal/is"
	"github.com/livebud/bud/package/gomod"
)

func TestPlugins(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()
	dir := t.TempDir()
	td := testdir.New(dir)
	td.Modules["github.com/livebud/bud-test-plugin"] = "v0.0.8"
	td.Modules["github.com/livebud/bud-test-nested-plugin"] = "v0.0.5"
	err := td.Write(ctx)
	is.NoErr(err)
	module, err := gomod.Find(dir)
	is.NoErr(err)
	// Load the overlay
	ofs, err := overlay.Load(module)
	is.NoErr(err)
	// Test that we can read files from the overlay
	code, err := fs.ReadFile(ofs, "view/index.svelte")
	is.NoErr(err)
	is.Equal(strings.TrimSpace(string(code)), `<h2>Welcome</h2>`)
	code, err = fs.ReadFile(ofs, "public/admin.css")
	is.NoErr(err)
	is.True(strings.Contains(string(code), `/* admin.css */`))
}
