package testdir

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/livebud/bud/internal/current"
	"github.com/livebud/bud/internal/npm"
	"github.com/livebud/bud/package/gomod"
	"github.com/livebud/bud/package/log"
	"github.com/livebud/bud/package/log/testlog"
	"github.com/livebud/bud/package/virtual"
	"golang.org/x/mod/modfile"
	"golang.org/x/sync/errgroup"
)

// Directory to copy into, if empty, we use $TMPDIR/bud_testdir_*.
var dirFlag = flag.String("dir", "", "dir to copy into")

func Load() (*Dir, error) {
	dir, err := makeDir(*dirFlag)
	if err != nil {
		return nil, err
	}
	dir, err = filepath.EvalSymlinks(dir)
	if err != nil {
		return nil, err
	}
	return &Dir{
		Files:       map[string]string{},
		Bytes:       map[string][]byte{},
		Modules:     map[string]string{},
		NodeModules: map[string]string{},

		dir:        dir,
		log:        testlog.New(),
		to:         virtual.OS(dir),
		cache:      map[string][]byte{},
		writeFiles: syncFiles,
	}, nil
}

type Dir struct {
	Files       map[string]string
	Bytes       map[string][]byte
	Modules     map[string]string
	NodeModules map[string]string

	dir        string
	log        log.Log
	to         virtual.FS
	cache      map[string][]byte
	writeFiles func(log log.Log, from fs.FS, to virtual.FS) error
}

var _ fs.FS = (*Dir)(nil)
var _ fs.StatFS = (*Dir)(nil)

const goMod = `
	module app.com

	require (
		github.com/livebud/bud v0.0.0
	)
`

func (d *Dir) loadFS() (virtual.Tree, error) {
	tree := virtual.Tree{}

	// Loop over files
	for path, data := range d.Files {
		tree[path] = &virtual.File{
			Data:    []byte(data),
			ModTime: time.Now(),
			Mode:    0644,
		}
	}
	// Loop over byte files
	for path, data := range d.Bytes {
		tree[path] = &virtual.File{
			Data:    data,
			ModTime: time.Now(),
			Mode:    0644,
		}
	}
	// Generate the go.mod automatically
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
	// Add a formatted go.mod
	formatted, err := modFile.Format()
	if err != nil {
		return nil, err
	}
	tree["go.mod"] = &virtual.File{
		Data:    formatted,
		ModTime: time.Now(),
		Mode:    0644,
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
				if err := copyLiveBud(tree, budDir); err != nil {
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
		tree["package.json"] = &virtual.File{
			Data:    append(pkg, '\n'),
			ModTime: time.Now(),
			Mode:    0644,
		}
	}
	return tree, nil
}

func (d *Dir) Write(ctx context.Context) error {
	from, err := d.loadFS()
	if err != nil {
		return err
	}
	if len(d.cache) > 0 {
		newfs := virtual.Tree{}
		for path, file := range from {
			if cached, ok := d.cache[path]; !ok || !bytes.Equal(cached, file.Data) {
				newfs[path] = file
			}
		}
		from = newfs
	}
	if err := d.writeFiles(d.log, from, d.to); err != nil {
		return err
	}
	eg, ctx := errgroup.WithContext(ctx)
	// Download modules that aren't in the module cache
	if _, ok := from["go.mod"]; ok {
		d.log.Debug("testdir: downloading go modules")
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
			"NO_COLOR=1",
		}
		eg.Go(func() error {
			err := cmd.Run()
			return err
		})
	}
	// Download node modules if there's a package.json
	if _, ok := from["package.json"]; ok {
		d.log.Debug("testdir: downloading node_modules")
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
	// Change to copying instead of syncing
	d.writeFiles = copyFiles
	// Cache the files for future reads
	for path, file := range from {
		d.cache[path] = file.Data
	}
	return nil
}

func (d *Dir) Directory() string {
	return d.dir
}

func (d *Dir) Open(name string) (fs.File, error) {
	return d.to.Open(name)
}

func (d *Dir) Stat(name string) (fs.FileInfo, error) {
	return fs.Stat(d.to, name)
}

func (d *Dir) RemoveAll(path string) error {
	return d.to.RemoveAll(path)
}

// Rename is os-specific, so we don't use virtual.FS
func (d *Dir) Rename(from, to string) error {
	return os.Rename(filepath.Join(d.dir, from), filepath.Join(d.dir, to))
}

// Exists returns nil if path exists. Should be called after Write.
func (d *Dir) Exists(paths ...string) error {
	eg := new(errgroup.Group)
	for _, path := range paths {
		path := path
		eg.Go(func() error {
			return d.exist(path)
		})
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
		eg.Go(func() error {
			return d.notExist(path)
		})
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

func makeDir(dirFlag string) (string, error) {
	if dirFlag != "" {
		if !strings.HasPrefix(dirFlag, "_tmp") {
			return "", fmt.Errorf("testdir: dir path must be relative and start with _tmp")
		}
		if strings.Contains(dirFlag, "*") {
			return os.MkdirTemp(".", dirFlag)
		}
		if err := os.MkdirAll(dirFlag, 0755); err != nil {
			return "", err
		}
		return dirFlag, nil
	}
	return os.MkdirTemp("", "bud_testdir_*")
}

// Copy files into a directory.
func copyFiles(log log.Log, from fs.FS, to virtual.FS) error {
	log.Debug("testdir: copying")
	return virtual.Copy(log, from, to)
}

// Sync files into a directory.
func syncFiles(log log.Log, from fs.FS, to virtual.FS) error {
	log.Debug("testdir: syncing")
	return virtual.Sync(log, from, to)
}

// Copy livebud into the node_modules directory.
func copyLiveBud(tree virtual.Tree, budDir string) error {
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
		tree[nodePath] = &virtual.File{
			Path:    nodePath,
			ModTime: stat.ModTime(),
			Mode:    stat.Mode(),
		}
		if stat.IsDir() {
			return nil
		}
		data, err := os.ReadFile(fpath)
		if err != nil {
			return err
		}
		tree[nodePath].Data = data
		return nil
	})
}
