package modcache_test

import (
	"context"
	"go/build"
	"os"
	"path/filepath"
	"testing"

	"github.com/livebud/bud/internal/testdir"

	"github.com/livebud/bud/internal/is"
	"github.com/livebud/bud/package/modcache"
)

// Run calls `go run -mod=mod main.go ...`
// func goRun(cacheDir, appDir string) (string, string, error) {
// 	ctx := context.Background()
// 	cmd := exec.CommandContext(ctx, "go", "run", "-mod=mod", "main.go")
// 	cmd.Env = append(os.Environ(), "GOMODCACHE="+cacheDir, "GOPRIVATE=*")
// 	stdout := new(bytes.Buffer)
// 	cmd.Stdout = stdout
// 	stderr := new(bytes.Buffer)
// 	cmd.Stderr = stderr
// 	cmd.Stdin = os.Stdin
// 	cmd.Dir = appDir
// 	err := cmd.Run()
// 	if stderr.Len() > 0 {
// 		return "", stderr.String(), nil
// 	}
// 	if err != nil {
// 		return "", "", err
// 	}
// 	return stdout.String(), "", nil
// }

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
	td.Modules["github.com/livebud/bud-test-plugin"] = "v0.0.9"
	is.NoErr(td.Write(ctx))
	modCache := modcache.Default()
	dir, err := modCache.ResolveDirectory("github.com/livebud/bud-test-plugin", "v0.0.9")
	is.NoErr(err)
	is.Equal(dir, modCache.Directory(`github.com/livebud`, `bud-test-plugin@v0.0.9`))
}
