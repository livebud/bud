package testdir

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io/fs"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
	"testing/fstest"
	"time"

	"github.com/livebud/bud/internal/dirhash"
	"github.com/livebud/bud/internal/gitignore"

	"github.com/livebud/bud/internal/current"

	"github.com/livebud/bud/internal/fstree"

	"golang.org/x/sync/errgroup"

	"github.com/livebud/bud/internal/dsync"
	"github.com/livebud/bud/internal/snapshot"

	"github.com/livebud/bud/internal/npm"
	"github.com/livebud/bud/package/gomod"
	"github.com/livebud/bud/package/modcache"
	"github.com/livebud/bud/package/vfs"
	"golang.org/x/mod/modfile"
)

const goMod = `
	module app.com

	require (
		github.com/livebud/bud v0.0.0
	)
`

func New() *Dir {
	return &Dir{
		Modules:     map[string]string{},
		NodeModules: map[string]string{},
		BFiles:      map[string][]byte{},
		Files:       map[string]string{},
	}
}

type Dir struct {
	Files       map[string]string // String files (convenient)
	BFiles      map[string][]byte // Byte files (for images and binaries)
	Modules     map[string]string // name[version]
	NodeModules map[string]string // name[version]
}

func merge(mapfs fstest.MapFS, fsys fs.FS, base ...string) error {
	basePath := path.Join(base...)
	return fs.WalkDir(fsys, ".", func(filePath string, de fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		fi, err := de.Info()
		if err != nil {
			return err
		}
		fullPath := path.Join(basePath, filePath)
		mapfs[fullPath] = &fstest.MapFile{
			ModTime: fi.ModTime(),
			Mode:    fi.Mode(),
		}
		if de.IsDir() {
			return nil
		}
		data, err := fs.ReadFile(fsys, filePath)
		if err != nil {
			return err
		}
		mapfs[fullPath].Data = data
		return nil
	})
}

func (d *Dir) mapfs() (fstest.MapFS, error) {
	mapfs := fstest.MapFS{}
	// Loop over files
	for path, data := range d.Files {
		mapfs[path] = &fstest.MapFile{
			Data:    []byte(data),
			ModTime: time.Now(),
			Mode:    0644,
		}
	}
	// Loop over byte files
	for path, data := range d.BFiles {
		mapfs[path] = &fstest.MapFile{
			Data:    data,
			ModTime: time.Now(),
			Mode:    0644,
		}
	}
	// Build up go.mod automatically
	modFile, err := modfile.Parse("go.mod", []byte(goMod), nil)
	if err != nil {
		return nil, err
	}
	currentDir, err := current.Directory()
	if err != nil {
		return nil, err
	}
	// Replace bud
	budDir, err := gomod.Absolute(currentDir)
	if err != nil {
		return nil, err
	}
	modFile.AddReplace("github.com/livebud/bud", "", budDir, "")
	// Add requires to go.mod
	for path, version := range d.Modules {
		if err := modFile.AddRequire(path, version); err != nil {
			return nil, err
		}
	}
	// Add node_modules
	if len(d.NodeModules) > 0 {
		nodePackage := &npm.Package{
			Name:         "testdir",
			Version:      "0.0.0",
			Dependencies: map[string]string{},
		}
		for name, version := range d.NodeModules {
			if name == "livebud" && version == "*" {
				if err := copyLiveBud(mapfs, budDir); err != nil {
					return nil, err
				}
			}
			nodePackage.Dependencies[name] = version
		}
		// Marshal into a package.json file
		pkg, err := json.MarshalIndent(nodePackage, "", "  ")
		if err != nil {
			return nil, err
		}
		// Add the package.json
		mapfs["package.json"] = &fstest.MapFile{
			Data:    append(pkg, '\n'),
			ModTime: time.Now(),
			Mode:    0644,
		}
	}
	// Add a formatted go.mod
	formatted, err := modFile.Format()
	if err != nil {
		return nil, err
	}
	// Hack to ensure changes to replaced modules trigger snapshot changes
	// TODO: cleanup this mess
	var hashes []byte
	for _, rep := range modFile.Replace {
		var skips []func(name string, isDir bool) bool
		// Handle
		if rep.New.Path == budDir {
			skips = append(skips, gitignore.From(budDir))
		}
		hash, err := dirhash.Hash(os.DirFS(rep.New.Path), dirhash.WithSkip(skips...))
		if err != nil {
			return nil, err
		}
		hashes = append(hashes, []byte(hash)...)
	}
	mapfs["go.mod"] = &fstest.MapFile{
		Data:    formatted,
		ModTime: time.Now(),
		Mode:    0644,
	}
	if len(hashes) > 0 {
		mapfs["go.mod"].Data = append(mapfs["go.mod"].Data, append([]byte("// "), hashes...)...)
	}
	return mapfs, nil
}

type Option func(o *option)

type option struct {
	backup bool
	skips  []func(name string, isDir bool) (skip bool)
}

func WithBackup(backup bool) Option {
	return func(o *option) {
		o.backup = backup
	}
}

func WithSkip(skips ...func(name string, isDir bool) (skip bool)) Option {
	return func(o *option) {
		o.skips = skips
	}
}

func (d *Dir) Hash() (string, error) {
	// Map out the filesystem
	fsys, err := d.mapfs()
	if err != nil {
		return "", err
	}
	// Compute a hash of the original filesystem
	return snapshot.Hash(fsys)
}

// Write testdir into dir
func (d *Dir) Write(dir string, options ...Option) error {
	// Load the options
	opt := &option{
		backup: true,
	}
	for _, option := range options {
		option(opt)
	}
	// Map out the filesystem
	fsys, err := d.mapfs()
	if err != nil {
		return err
	}
	// Compute a hash of the original filesystem
	hash, err := snapshot.Hash(fsys)
	if err != nil {
		return err
	}
	// Try restoring from cache
	if opt.backup {
		cachedFS, err := snapshot.Restore(hash)
		if nil == err {
			return dsync.Dir(cachedFS, ".", vfs.OS(dir), ".", dsync.WithSkip(opt.skips...))
		} else if !errors.Is(err, fs.ErrNotExist) {
			return err
		}
	}
	if err := dsync.Dir(fsys, ".", vfs.OS(dir), ".", dsync.WithSkip(opt.skips...)); err != nil {
		return err
	}
	// Load the module cache
	modCache := modcache.Default()
	modCacheDir, err := filepath.Abs(modCache.Directory())
	if err != nil {
		return err
	}
	eg, ctx := errgroup.WithContext(context.Background())
	// Download modules that aren't in the module cache
	if _, ok := fsys["go.mod"]; ok {
		// Use "all" to extract cached directories into GOMODCACHE, so there's not
		// a "go: downloading ..." step during go build.
		cmd := exec.CommandContext(ctx, "go", "mod", "download", "all")
		cmd.Dir = dir
		cmd.Stderr = os.Stderr
		cmd.Stdout = os.Stdout
		cmd.Env = []string{
			"HOME=" + os.Getenv("HOME"),
			"PATH=" + os.Getenv("PATH"),
			"GOPATH=" + os.Getenv("GOPATH"),
			"GOMODCACHE=" + modCacheDir,
			"NO_COLOR=1",
			// TODO: remove once we can write a sum file to the modcache
			"GOPRIVATE=*",
		}
		eg.Go(func() error {
			err := cmd.Run()
			return err
		})
	}

	if _, ok := fsys["package.json"]; ok {
		cmd := exec.CommandContext(ctx, "npm", "install", "--loglevel=error")
		cmd.Dir = dir
		cmd.Stderr = os.Stderr
		cmd.Env = []string{
			"HOME=" + os.Getenv("HOME"),
			"PATH=" + os.Getenv("PATH"),
			"NO_COLOR=1",
		}
		eg.Go(func() error {
			err := cmd.Run()
			return err
		})
	}

	// Wait for both commands to finish
	if err := eg.Wait(); err != nil {
		return err
	}
	// Backing the snapshot up
	if opt.backup {
		if err := snapshot.Backup(hash, os.DirFS(dir)); err != nil {
			return err
		}
	}
	return nil
}

func Tree(dir string) (string, error) {
	tree, err := fstree.Walk(os.DirFS(dir))
	if err != nil {
		return "", err
	}
	return tree.String(), nil
}

func copyLiveBud(mapfs fstest.MapFS, budDir string) error {
	liveBudDir := filepath.Join(budDir, "livebud")
	return filepath.WalkDir(liveBudDir, func(path string, de fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		stat, err := de.Info()
		if err != nil {
			return err
		}
		relPath, err := filepath.Rel(liveBudDir, path)
		if err != nil {
			return err
		}
		mapfs[relPath] = &fstest.MapFile{
			ModTime: stat.ModTime(),
			Mode:    stat.Mode(),
			Sys:     stat.Sys(),
		}
		if stat.IsDir() {
			return nil
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		mapfs[relPath].Data = data
		return nil
	})
}

func packLiveBud(budDir string) (string, error) {
	liveBudDir := filepath.Join(budDir, "livebud")
	tmpDir, err := ioutil.TempDir("", "testdir-livebud-*")
	if err != nil {
		return "", err
	}
	cmd := exec.Command("npm", "pack", liveBudDir, "--loglevel=error")
	cmd.Dir = tmpDir
	cmd.Stderr = os.Stderr
	tarName := new(bytes.Buffer)
	cmd.Stdout = tarName
	cmd.Env = []string{
		"HOME=" + os.Getenv("HOME"),
		"PATH=" + os.Getenv("PATH"),
		"NO_COLOR=1",
	}
	if err := cmd.Run(); err != nil {
		return "", err
	}
	tarPath := filepath.Join(tmpDir, strings.TrimSpace(tarName.String()))
	return tarPath, nil
}
