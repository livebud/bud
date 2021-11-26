package modcache_test

import (
	"bytes"
	"context"
	"go/build"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"gitlab.com/mnm/bud/vfs"

	"github.com/matryer/is"
	"gitlab.com/mnm/bud/internal/modcache"
)

// Run calls `go run -mod=mod main.go ...`
func goRun(cacheDir, appDir string) (string, string, error) {
	ctx := context.Background()
	cmd := exec.CommandContext(ctx, "go", "run", "-mod=mod", "main.go")
	cmd.Env = append(os.Environ(), "GOMODCACHE="+cacheDir, "GOPRIVATE=*")
	stdout := new(bytes.Buffer)
	cmd.Stdout = stdout
	stderr := new(bytes.Buffer)
	cmd.Stderr = stderr
	cmd.Stdin = os.Stdin
	cmd.Dir = appDir
	err := cmd.Run()
	if stderr.Len() > 0 {
		return "", stderr.String(), nil
	}
	if err != nil {
		return "", "", err
	}
	return stdout.String(), "", nil
}

func TestDirectory(t *testing.T) {
	is := is.New(t)
	dir := modcache.Directory()
	if env := os.Getenv("GOMODCACHE"); env != "" {
		is.Equal(dir, env)
	} else {
		is.Equal(dir, filepath.Join(build.Default.GOPATH, "pkg", "mod"))
	}
}

func TestWriteModule(t *testing.T) {
	is := is.New(t)
	cacheDir := t.TempDir()
	modCache := modcache.New(cacheDir)
	err := modCache.Write(modcache.Modules{
		"mod.test/one@v0.0.1": modcache.Files{
			"const.go": "package one\n\nconst Answer = 42",
		},
		"mod.test/one@v0.0.2": modcache.Files{
			"const.go": "package one\n\nconst Answer = 43",
		},
	})
	is.NoErr(err)
	dir, err := modCache.ResolveDirectory("mod.test/one", "v0.0.2")
	is.NoErr(err)
	is.Equal(dir, modCache.Directory("mod.test", "one@v0.0.2"))
	// Now verify that the go commands don't try downloading mod.test/one
	appDir := t.TempDir()
	err = vfs.Write(appDir, vfs.Map{
		"go.mod": `
			module app.com

			require (
				mod.test/one v0.0.2
			)
		`,
		"main.go": `
			package main

			import (
				"fmt"
				"mod.test/one"
			)

			func main() {
				fmt.Print(one.Answer)
			}
		`,
	})
	is.NoErr(err)
	stdout, stderr, err := goRun(cacheDir, appDir)
	is.NoErr(err)
	is.Equal(stderr, "")
	is.Equal(stdout, "43")
}

func TestResolveDirectoryFromCache(t *testing.T) {
	is := is.New(t)
	modCache := modcache.Default()
	is.True(modCache.Directory() != "")
	dir, err := modCache.ResolveDirectory("github.com/matryer/is", "v1.4.0")
	is.NoErr(err)
	is.Equal(dir, modCache.Directory("github.com", "matryer", "is@v1.4.0"))
}
