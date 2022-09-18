package gobuild

import (
	"context"
	"errors"
	"io"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/livebud/bud/internal/imhash"
	"github.com/livebud/bud/internal/symlink"
	"github.com/livebud/bud/package/gomod"
)

func New(module *gomod.Module) *Builder {
	return &Builder{
		os.Environ(),
		os.Stderr,
		os.Stdin,
		os.Stdout,
		module,
		module.Directory("bud", "cache"),
	}
}

type Builder struct {
	Env      []string
	Stderr   io.Writer
	Stdin    io.Reader
	Stdout   io.Writer
	module   *gomod.Module
	cacheDir string
}

// Build a Go binary and cache it for later use
func (b *Builder) Build(ctx context.Context, mainPath string, outPath string, flags ...string) error {
	hash, err := imhash.Hash(b.module, filepath.Dir(mainPath))
	if err != nil {
		return err
	}
	cachePath := filepath.Join(b.cacheDir, hash)
	exists, err := b.exists(cachePath)
	if err != nil {
		return err
	} else if exists {
		return symlink.Link(cachePath, b.module.Directory(outPath))
	}
	if err := b.build(ctx, mainPath, cachePath, flags...); err != nil {
		return err
	}
	return symlink.Link(cachePath, b.module.Directory(outPath))
}

// Build calls `go build -mod=mod -o main [flags...] main.go`
func (b *Builder) build(ctx context.Context, mainPath string, outPath string, flags ...string) error {
	// Compile the args
	args := append([]string{
		"build",
		"-mod=mod",
		"-o=" + outPath,
	}, flags...)
	args = append(args, mainPath)
	cmd := exec.CommandContext(ctx, "go", args...)
	cmd.Env = append(b.Env,
		"GOMODCACHE="+b.module.ModCache(),
	)
	cmd.Stdout = b.Stdout
	cmd.Stderr = b.Stderr
	cmd.Stdin = b.Stdin
	cmd.Dir = b.module.Directory()
	err := cmd.Run()
	if err != nil {
		return err
	}
	return nil
}

// Check if the path exists
func (b *Builder) exists(path string) (bool, error) {
	if _, err := os.Stat(path); err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
