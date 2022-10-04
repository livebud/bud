package scaffold

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/livebud/bud/internal/gotemplate"
	"github.com/livebud/bud/package/vfs"
	"github.com/otiai10/copy"
)

type MapFS = vfs.Memory

func OSFS(dir string) vfs.ReadWritable {
	return vfs.OS(dir)
}

type Scaffolding interface {
	Scaffold(fsys vfs.ReadWritable) error
}

func Template(path, template string, state interface{}) Scaffolding {
	return &templateFile{path, template, state}
}

type templateFile struct {
	path     string
	template string
	state    interface{}
}

func (t *templateFile) Scaffold(fsys vfs.ReadWritable) error {
	template, err := gotemplate.Parse(t.path, t.template)
	if err != nil {
		return err
	}
	code, err := template.Generate(t.state)
	if err != nil {
		return err
	}
	if err := fsys.MkdirAll(filepath.Dir(t.path), 0755); err != nil {
		return err
	}
	return fsys.WriteFile(t.path, code, 0644)
}

func JSON(path string, state interface{}) Scaffolding {
	return &jsonFile{path, state}
}

type jsonFile struct {
	path  string
	state interface{}
}

func (j *jsonFile) Scaffold(fsys vfs.ReadWritable) error {
	code, err := json.MarshalIndent(j.state, "", "  ")
	if err != nil {
		return err
	}
	if err := fsys.MkdirAll(filepath.Dir(j.path), 0755); err != nil {
		return err
	}
	return fsys.WriteFile(j.path, code, 0644)
}

func File(path string, data []byte) Scaffolding {
	return &file{path, data}
}

type file struct {
	path string
	data []byte
}

func (f *file) Scaffold(fsys vfs.ReadWritable) error {
	if err := fsys.MkdirAll(filepath.Dir(f.path), 0755); err != nil {
		return err
	}
	return fsys.WriteFile(f.path, f.data, 0644)
}

func Scaffold(fsys vfs.ReadWritable, scaffoldings ...Scaffolding) error {
	// TODO: make concurrency safe then refactor to use errgroup.
	for _, s := range scaffoldings {
		if err := s.Scaffold(fsys); err != nil {
			return err
		}
	}
	return nil
}

func Write(fsys fs.FS, to string) error {
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
			return os.MkdirAll(toPath, mode)
		}
		data, err := fs.ReadFile(fsys, path)
		if err != nil {
			return err
		}
		mode := de.Type()
		if mode == 0 {
			mode = fs.FileMode(0644)
		}
		if err := ensureNotExist(toPath); err != nil {
			return err
		}
		return os.WriteFile(toPath, data, mode)
	})
}

// Move a `from` directory to a `to` directory.
func Move(fromDir, toDir string) error {
	if err := ensureNotExist(toDir); err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(toDir), 0755); err != nil {
		return fmt.Errorf("scaffold: unable to make the directory. %w", err)
	}
	// Try moving from the temporary directory to the destination directory.
	// If `from` is on a different partition than `to`, the underlying os.Rename
	// can fail with an "invalid cross-device link" error. If this occurs we'll
	// fallback to copying the files over recursively.
	if err := os.Rename(fromDir, toDir); err != nil {
		// If it's not an invalid cross-device link error, return the error
		if !isInvalidCrossLink(err) {
			return fmt.Errorf("scaffold: unable to rename. %w", err)
		}
		// Fallback to copying files recursively
		return copy.Copy(fromDir, toDir)
	}
	return nil
}

// Returns an error if the filepath exists
func ensureNotExist(path string) error {
	if _, err := os.Stat(path); err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return nil
		}
		return err
	}
	return fmt.Errorf("scaffold: %q already exists: %w", path, fs.ErrExist)
}

func isInvalidCrossLink(err error) bool {
	return strings.Contains(err.Error(), "invalid cross-device link")
}

// Command creates a runnable command with common defaults.
func Command(dir, name string, args ...string) *exec.Cmd {
	cmd := exec.Command(name, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = os.Environ()
	cmd.Dir = dir
	return cmd
}
