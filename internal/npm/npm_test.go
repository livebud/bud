package npm_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/livebud/bud/internal/npm"
	"github.com/matryer/is"
)

func exists(t testing.TB, path string) {
	t.Helper()
	if _, err := os.Stat(path); err != nil {
		t.Fatalf("expected %q to exist: %s", path, err)
	}
}

func TestInstallSvelte(t *testing.T) {
	is := is.New(t)
	dir := t.TempDir()
	ctx := context.Background()
	err := npm.Install(ctx, dir, "svelte@3.42.3", "uid@2.0.0")
	is.NoErr(err)
	exists(t, filepath.Join(dir, "node_modules", "svelte", "package.json"))
	exists(t, filepath.Join(dir, "node_modules", "uid", "package.json"))
	exists(t, filepath.Join(dir, "node_modules", "svelte", "internal", "index.js"))
}

func TestInstallReact(t *testing.T) {
	is := is.New(t)
	dir := t.TempDir()
	ctx := context.Background()
	err := npm.Install(ctx, dir, "react@18.2.0", "react-dom@18.2.0")
	is.NoErr(err)
	exists(t, filepath.Join(dir, "node_modules", "react", "package.json"))
	exists(t, filepath.Join(dir, "node_modules", "react-dom", "package.json"))
}

func TestInstallStripe(t *testing.T) {
	is := is.New(t)
	dir := t.TempDir()
	ctx := context.Background()
	err := npm.Install(ctx, dir, "@stripe/stripe-js@2.1.11")
	is.NoErr(err)
	exists(t, filepath.Join(dir, "node_modules", "@stripe", "stripe-js", "package.json"))
}
