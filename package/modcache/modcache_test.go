package modcache_test

import (
	"bytes"
	"context"
	"go/build"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/livebud/bud/internal/testdir"

	"github.com/livebud/bud/internal/is"
	"github.com/livebud/bud/package/modcache"
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
	dir := modcache.Default().Directory()
	if env := os.Getenv("GOMODCACHE"); env != "" {
		is.Equal(dir, env)
	} else {
		is.Equal(dir, filepath.Join(build.Default.GOPATH, "pkg", "mod"))
	}
}

func TestResolveDirectoryFromCache(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()
	dir := t.TempDir()
	td := testdir.New(dir)
	td.Modules["github.com/livebud/bud-test-plugin"] = "v0.0.8"
	is.NoErr(td.Write(ctx))
	modCache := modcache.Default()
	dir, err := modCache.ResolveDirectory("github.com/livebud/bud-test-plugin", "v0.0.8")
	is.NoErr(err)
	is.Equal(dir, modCache.Directory(`github.com/livebud`, `bud-test-plugin@v0.0.8`))
}

// func TestWriteModule(t *testing.T) {
// 	is := is.New(t)
// 	cacheDir := t.TempDir()
// 	modCache := modcache.New(cacheDir)
// 	err := modCache.Write(modcache.Modules{
// 		"mod.test/one@v0.0.1": modcache.Files{
// 			"const.go": "package one\n\nconst Answer = 42",
// 		},
// 		"mod.test/one@v0.0.2": modcache.Files{
// 			"const.go": "package one\n\nconst Answer = 43",
// 		},
// 	})
// 	is.NoErr(err)
// 	dir, err := modCache.ResolveDirectory("mod.test/one", "v0.0.2")
// 	is.NoErr(err)
// 	is.Equal(dir, modCache.Directory("mod.test", "one@v0.0.2"))
// 	// Now verify that the go commands don't try downloading mod.test/one
// 	appDir := t.TempDir()
// 	err = vfs.Write(appDir, vfs.Map{
// 		"go.mod": []byte(`
// 			module app.com

// 			require (
// 				mod.test/one v0.0.2
// 			)
// 		`),
// 		"main.go": []byte(`
// 			package main

// 			import (
// 				"fmt"
// 				"mod.test/one"
// 			)

// 			func main() {
// 				fmt.Print(one.Answer)
// 			}
// 		`),
// 	})
// 	is.NoErr(err)
// 	stdout, stderr, err := goRun(cacheDir, appDir)
// 	is.NoErr(err)
// 	is.Equal(stderr, "")
// 	is.Equal(stdout, "43")
// }

// func exists(fsys fs.FS, paths ...string) error {
// 	for _, path := range paths {
// 		if _, err := fs.Stat(fsys, path); err != nil {
// 			return err
// 		}
// 	}
// 	return nil
// }

// func TestWriteModuleFS(t *testing.T) {
// 	is := is.New(t)
// 	cacheDir := t.TempDir()
// 	fsys, err := modcache.WriteFS(modcache.Modules{
// 		"mod.test/one@v0.0.1": modcache.Files{
// 			"const.go": "package one\n\nconst Answer = 42",
// 		},
// 		"mod.test/one@v0.0.2": modcache.Files{
// 			"const.go": "package one\n\nconst Answer = 43",
// 		},
// 	})
// 	is.NoErr(err)
// 	err = exists(fsys,
// 		"cache/download/mod.test/one/@v/v0.0.1.info",
// 		"cache/download/mod.test/one/@v/v0.0.1.zip",
// 		"cache/download/mod.test/one/@v/v0.0.1.mod",
// 		"cache/download/mod.test/one/@v/v0.0.1.ziphash",
// 		"cache/download/mod.test/one/@v/v0.0.2.info",
// 		"cache/download/mod.test/one/@v/v0.0.2.zip",
// 		"cache/download/mod.test/one/@v/v0.0.2.mod",
// 		"cache/download/mod.test/one/@v/v0.0.2.ziphash",
// 		"mod.test/one@v0.0.1/const.go",
// 		"mod.test/one@v0.0.1/go.mod",
// 		"mod.test/one@v0.0.2/const.go",
// 		"mod.test/one@v0.0.2/go.mod",
// 	)
// 	is.NoErr(err)
// 	err = dsync.Dir(fsys, ".", vfs.OS(cacheDir), ".")
// 	is.NoErr(err)
// 	modCache := modcache.New(cacheDir)
// 	dir, err := modCache.ResolveDirectory("mod.test/one", "v0.0.2")
// 	is.NoErr(err)
// 	is.Equal(dir, modCache.Directory("mod.test", "one@v0.0.2"))
// 	// Now verify that the go commands don't try downloading mod.test/one
// 	appDir := t.TempDir()
// 	err = vfs.Write(appDir, vfs.Map{
// 		"go.mod": []byte(`
// 			module app.com

// 			require (
// 				mod.test/one v0.0.2
// 			)
// 		`),
// 		"main.go": []byte(`
// 			package main

// 			import (
// 				"fmt"
// 				"mod.test/one"
// 			)

// 			func main() {
// 				fmt.Print(one.Answer)
// 			}
// 		`),
// 	})
// 	is.NoErr(err)
// 	stdout, stderr, err := goRun(cacheDir, appDir)
// 	is.NoErr(err)
// 	is.Equal(stderr, "")
// 	is.Equal(stdout, "43")
// }

// func TestExportImport(t *testing.T) {
// 	is := is.New(t)
// 	cacheDir := t.TempDir()
// 	modCache := modcache.New(cacheDir)
// 	err := modCache.Write(map[string]modcache.Files{
// 		"github.com/livebud/bud-tailwind@v0.0.1": modcache.Files{
// 			"public/tailwind/preflight.css": `/* tailwind */`,
// 		},
// 	})
// 	is.NoErr(err)
// 	dir, err := modCache.ResolveDirectory("github.com/livebud/bud-tailwind", "v0.0.1")
// 	is.NoErr(err)
// 	is.Equal(dir, modCache.Directory("github.com/livebud/bud-tailwind@v0.0.1"))
// 	cacheDir2 := t.TempDir()
// 	modCache2 := modcache.New(cacheDir2)
// 	// Verify modcache2 doesn't have the module
// 	dir, err = modCache2.ResolveDirectory("github.com/livebud/bud-tailwind", "v0.0.1")
// 	is.Equal(dir, "")
// 	is.True(errors.Is(err, fs.ErrNotExist))
// 	// Export to a new location
// 	tmpDir := t.TempDir()
// 	err = modCache.Export(tmpDir)
// 	is.NoErr(err)
// 	// Import from new location
// 	err = modCache2.Import(tmpDir)
// 	is.NoErr(err)
// 	// Try again
// 	dir, err = modCache2.ResolveDirectory("github.com/livebud/bud-tailwind", "v0.0.1")
// 	is.NoErr(err)
// 	is.Equal(dir, modCache2.Directory("github.com/livebud/bud-tailwind@v0.0.1"))
// }
