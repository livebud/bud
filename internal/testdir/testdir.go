package testdir

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"testing/fstest"
	"time"

	"github.com/livebud/bud/internal/dirhash"
	"github.com/livebud/bud/internal/gitignore"

	"github.com/livebud/bud/internal/current"

	"github.com/livebud/bud/internal/fstree"

	"golang.org/x/sync/errgroup"

	"github.com/livebud/bud/internal/snapshot"

	"github.com/livebud/bud/internal/npm"
	"github.com/livebud/bud/package/gomod"
	"github.com/livebud/bud/package/modcache"
	"golang.org/x/mod/modfile"
)

const goMod = `
	module app.com

	require (
		github.com/livebud/bud v0.0.0
	)
`

func New(dir string) *Dir {
	return &Dir{
		dir:         dir,
		Backup:      false, // TODO: re-enable and store snapshots in ./bud/tmp
		Skip:        func(string, bool) bool { return false },
		Modules:     map[string]string{},
		NodeModules: map[string]string{},
		BFiles:      map[string][]byte{},
		Files:       map[string]string{},
	}
}

type Dir struct {
	dir         string
	Backup      bool
	Skip        func(name string, isDir bool) (skip bool)
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

// Dir returns the directory
func (d *Dir) Path(subpaths ...string) string {
	return filepath.Join(append([]string{d.dir}, subpaths...)...)
}

// Hash returns a file hash of our mapped file system
func (d *Dir) Hash() (string, error) {
	// Map out the filesystem
	fsys, err := d.mapfs()
	if err != nil {
		return "", err
	}
	// Compute a hash of the original filesystem
	return snapshot.Hash(fsys)
}

func (d *Dir) writeAll(fsys fs.FS, to string) error {
	return fs.WalkDir(fsys, ".", func(path string, de fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		toPath := filepath.Join(to, path)
		if de.IsDir() {
			mode := de.Type()
			if mode == fs.ModeDir {
				mode = fs.FileMode(0755)
			}
			if err := os.MkdirAll(toPath, mode); err != nil {
				return err
			}
			return nil
		}
		data, err := fs.ReadFile(fsys, path)
		if err != nil {
			return err
		}
		mode := de.Type()
		if mode == 0 {
			mode = fs.FileMode(0644)
		}
		if err := os.WriteFile(toPath, data, mode); err != nil {
			return err
		}
		return nil
	})
}

// Write testdir into dir
func (d *Dir) Write(ctx context.Context) error {
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
	if d.Backup {
		cachedFS, err := snapshot.Restore(hash)
		if nil == err {
			// Write the cache to dir
			if err := d.writeAll(cachedFS, d.dir); err != nil {
				return err
			}
		} else if !errors.Is(err, fs.ErrNotExist) {
			return err
		}
	}
	// Write all the files in the filesystem out
	if err := d.writeAll(fsys, d.dir); err != nil {
		return err
	}
	// Load the module cache
	modCache := modcache.Default()
	modCacheDir, err := filepath.Abs(modCache.Directory())
	if err != nil {
		return err
	}
	eg, ctx := errgroup.WithContext(ctx)
	// Download modules that aren't in the module cache
	if _, ok := fsys["go.mod"]; ok {
		// Use "all" to extract cached directories into GOMODCACHE, so there's not
		// a "go: downloading ..." step during go build.
		cmd := exec.CommandContext(ctx, "go", "mod", "download", "all")
		cmd.Dir = d.dir
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
		// Avoid symlinking in node_modules/.bin to avoid symlink issues downstream
		cmd := exec.CommandContext(ctx, "npm", "install", "--loglevel=error", "--no-bin-links")
		cmd.Dir = d.dir
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
	if d.Backup {
		if err := snapshot.Backup(hash, os.DirFS(d.dir)); err != nil {
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
	liveBudModules := path.Join("node_modules", "livebud")
	return filepath.WalkDir(liveBudDir, func(fpath string, de fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		stat, err := de.Info()
		if err != nil {
			return err
		}
		// Ignore symlinks since they cause write issues later and so far aren't
		// relevant.
		if stat.Mode()&fs.ModeSymlink != 0 {
			return nil
		}
		relPath, err := filepath.Rel(liveBudDir, fpath)
		if err != nil {
			return err
		}
		nodePath := path.Join(liveBudModules, relPath)
		mapfs[nodePath] = &fstest.MapFile{
			ModTime: stat.ModTime(),
			Mode:    stat.Mode(),
			Sys:     stat.Sys(),
		}
		if stat.IsDir() {
			return nil
		}
		data, err := os.ReadFile(fpath)
		if err != nil {
			return err
		}
		mapfs[nodePath].Data = data
		return nil
	})
}

// Exists returns nil if path exists. Should be called after Write.
func (d *Dir) Exists(paths ...string) error {
	eg := new(errgroup.Group)
	for _, path := range paths {
		path := path
		eg.Go(func() error { return d.exist(path) })
	}
	return eg.Wait()
}

func (d *Dir) exist(path string) error {
	if _, err := d.Stat(path); err != nil {
		return err
	}
	return nil
}

// NotExists returns nil if path doesn't exist. Should be called after Write.
func (d *Dir) NotExists(paths ...string) error {
	eg := new(errgroup.Group)
	for _, path := range paths {
		path := path
		eg.Go(func() error { return d.notExist(path) })
	}
	return eg.Wait()
}

func (d *Dir) notExist(path string) error {
	if _, err := d.Stat(path); err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return nil
		}
		return err
	}
	return fmt.Errorf("%s exists: %w", path, fs.ErrExist)
}

// Stat return the file info of path. Should be called after Write.
func (d *Dir) Stat(path string) (fs.FileInfo, error) {
	return os.Stat(filepath.Join(d.dir, path))
}
