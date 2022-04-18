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
	"runtime"
	"strings"
	"testing/fstest"
	"time"

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
		Modules:     modcache.Modules{},
		NodeModules: map[string]string{},
		BFiles:      map[string][]byte{},
		Files:       map[string]string{},
	}
}

type Dir struct {
	Files       map[string]string // String files (convenient)
	BFiles      map[string][]byte // Byte files (for images and binaries)
	Modules     modcache.Modules  // name@version[path[data]]
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
	currentDir, err := dirname()
	if err != nil {
		return nil, err
	}
	// Replace bud
	budDir, err := gomod.Absolute(currentDir)
	if err != nil {
		return nil, err
	}
	modFile.AddReplace("github.com/livebud/bud", "", budDir, "")
	// Merge the go modules in
	if len(d.Modules) > 0 {
		fsys, err := modcache.WriteFS(d.Modules)
		if err != nil {
			return nil, err
		}
		// Add requires to go.mod
		for pv := range d.Modules {
			path, version, err := modcache.SplitPathVersion(pv)
			if err != nil {
				return nil, err
			}
			// Add require to the go.mod
			if err := modFile.AddRequire(path, version); err != nil {
				return nil, err
			}
		}
		if err := merge(mapfs, fsys, ".mod"); err != nil {
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
				tarPath, err := packLiveBud(budDir)
				if err != nil {
					return nil, err
				}
				stat, err := os.Stat(tarPath)
				if err != nil {
					return nil, err
				}
				data, err := os.ReadFile(tarPath)
				if err != nil {
					return nil, err
				}
				mapfs[".npm/livebud.tgz"] = &fstest.MapFile{
					Data:    data,
					ModTime: stat.ModTime(),
					Mode:    stat.Mode(),
					Sys:     stat.Sys,
				}
				continue
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
	mapfs["go.mod"] = &fstest.MapFile{
		Data:    formatted,
		ModTime: time.Now(),
		Mode:    0644,
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
	if _, ok := fsys[".mod"]; ok {
		modCache = modcache.New(filepath.Join(dir, ".mod"))
	}
	modCacheDir, err := filepath.Abs(modCache.Directory())
	if err != nil {
		return err
	}

	eg, ctx := errgroup.WithContext(context.Background())
	// Download modules that aren't in the module cache
	if _, ok := fsys["go.mod"]; ok {
		// Use "all" to extract cached directories into GOMODCACHE, so there's not
		// a "go: downloading ..." step during go build.
		// TODO: speed this call up, it takes around 60s right now
		cmd := exec.CommandContext(ctx, "go", "mod", "download", "-modcacherw", "all")
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

	// Copy livebud.tgz into node_modules and install any dependencies
	if _, ok := fsys[".npm/livebud.tgz"]; ok {
		cmd := exec.CommandContext(ctx, "npm", "install", "--loglevel=error", ".npm/livebud.tgz")
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
	// Delete .mod/cache/vcs, it's slow to unzip and unnecessary
	if os.RemoveAll(filepath.Join(dir, ".mod", "cache", "vcs")); err != nil {
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

func ModCache(dir string) *modcache.Cache {
	modDir := filepath.Join(dir, ".mod")
	if _, err := os.Stat(modDir); err != nil {
		return modcache.Default()
	}
	return modcache.New(modDir)
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

// dirname gets the directory of this file
func dirname() (string, error) {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return "", errors.New("unable to get the current filename")
	}
	return filepath.Dir(filename), nil
}
