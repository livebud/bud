package scaffolder

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/otiai10/copy"

	"github.com/livebud/bud/internal/gotemplate"
	"golang.org/x/sync/errgroup"
)

// type Task interface {
// 	Run(ctx context.Context) error
// 	String() string
// }

func Load() (*Scaffold, error) {
	tmpDir, err := ioutil.TempDir("", "bud-create-*")
	if err != nil {
		return nil, err
	}
	return &Scaffold{tmpDir}, nil
}

type Scaffold struct {
	dir string
}

type Generator interface {
	Generate() error
}

type Generators = map[Generator]interface{}

func (s *Scaffold) Template(path, template string, state interface{}) Generator {
	return &templateGenerator{filepath.Join(s.dir, path), template, state}
}

type templateGenerator struct {
	path     string
	template string
	state    interface{}
}

func (g *templateGenerator) Generate() error {
	template, err := gotemplate.Parse(g.path, g.template)
	if err != nil {
		return err
	}
	code, err := template.Generate(g.state)
	if err != nil {
		return err
	}
	return os.WriteFile(g.path, code, 0644)
}

func (s *Scaffold) JSON(path string, state interface{}) Generator {
	return &jsonGenerator{filepath.Join(s.dir, path), state}
}

type jsonGenerator struct {
	path  string
	state interface{}
}

func (g *jsonGenerator) Generate() error {
	code, err := json.MarshalIndent(g.state, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(g.path, code, 0644)
}

// Generate templates
func (s *Scaffold) Generate(generators ...Generator) error {
	eg := new(errgroup.Group)
	for _, generator := range generators {
		generator := generator
		eg.Go(func() error { return generator.Generate() })
	}
	return eg.Wait()
}

// Command creates a runnable command with common defaults.
func (s *Scaffold) Command(name string, args ...string) *exec.Cmd {
	cmd := exec.Command(name, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = os.Environ()
	cmd.Dir = s.dir
	return cmd
}

// Move the scaffolding to a destination directory.
func (s *Scaffold) Move(toDir string) error {
	if err := ensureNotExist(toDir); err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(toDir), 0755); err != nil {
		return err
	}
	// Try moving from the temporary directory to the destination directory.
	// If `from` is on a different partition than `to`, the underlying os.Rename
	// can fail with an "invalid cross-device link" error. If this occurs we'll
	// fallback to copying the files over recursively.
	if err := os.Rename(s.dir, toDir); err != nil {
		// If it's not an invalid cross-device link error, return the error
		if !isInvalidCrossLink(err) {
			return err
		}
		// Fallback to copying files recursively
		return copy.Copy(s.dir, toDir)
	}
	return nil
}

func ensureNotExist(dir string) error {
	if _, err := os.Stat(dir); err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return nil
		}
		return err
	}
	return fmt.Errorf("scaffold: refusing to overwrite. %q already exists", dir)
}

func isInvalidCrossLink(err error) bool {
	return strings.Contains(err.Error(), "invalid cross-device link")
}
