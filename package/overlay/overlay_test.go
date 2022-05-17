package overlay_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/livebud/bud/internal/testplugin"

	"io/fs"

	"github.com/livebud/bud/package/overlay"

	"github.com/livebud/bud/package/gomod"
	"github.com/matryer/is"
)

func TestPlugins(t *testing.T) {
	is := is.New(t)
	// Create a gomod with some dependencies
	dep1, err := testplugin.Plugin()
	is.NoErr(err)
	dep2, err := testplugin.NestedPlugin()
	is.NoErr(err)
	modFile := `
		module app.com
		require ` + dep1.Path + ` ` + dep1.Version + `
		require ` + dep2.Path + ` ` + dep2.Version + `
	`
	appDir := t.TempDir()
	err = os.WriteFile(filepath.Join(appDir, "go.mod"), []byte(modFile), 0644)
	is.NoErr(err)
	module, err := gomod.Find(appDir)
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
