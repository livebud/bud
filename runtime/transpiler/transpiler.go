package transpiler

import (
	"errors"
	"fmt"
	"io/fs"
	"path"
	"strings"

	"github.com/livebud/bud/package/genfs"
	"github.com/livebud/transpiler"
)

type File = transpiler.File

type Transpiler interface {
	Transpile(fromPath, toExt string) ([]byte, error)
}

func NewServer() *Server {
	return &Server{transpiler.New()}
}

type Server struct {
	tr *transpiler.Transpiler
}

func (s *Server) Add(fromExt, toExt string, transpile func(file *File) error) {
	s.tr.Add(fromExt, toExt, transpile)
}

func (s *Server) Serve(fsys genfs.FS, file *genfs.File) error {
	toExt, inputPath := splitRoot(file.Relative())
	data, err := fs.ReadFile(fsys, inputPath)
	if err != nil {
		return err
	}
	output, err := s.tr.Transpile(inputPath, toExt, data)
	if err != nil {
		return err
	}
	file.Data = output
	return nil
}

// Aliasing allows us to target the transpiler filesystem directly
type FS = fs.FS

func NewTester(fsys FS) *Tester {
	return &Tester{fsys, transpiler.New()}
}

type Tester struct {
	fsys fs.FS
	tr   *transpiler.Transpiler
}

var _ Transpiler = (*Tester)(nil)

func (t *Tester) Add(fromExt, toExt string, transpile func(file *File) error) {
	t.tr.Add(fromExt, toExt, transpile)
}

func (t *Tester) Transpile(fromPath, toExt string) ([]byte, error) {
	data, err := fs.ReadFile(t.fsys, fromPath)
	if err != nil {
		return nil, err
	}
	return t.tr.Transpile(fromPath, toExt, data)
}

func NewClient(fsys FS) *Client {
	return &Client{fsys}
}

// Client transpiles files within the filesystem
type Client struct {
	fsys FS
}

var _ Transpiler = (*Client)(nil)

const transpilerDir = `bud/internal/transpiler`

// Transpile a file from one extension to another. This assumes the filesystem
// is serving from the transpiler directory.
func (c *Client) Transpile(fromPath, toExt string) ([]byte, error) {
	return fs.ReadFile(c.fsys, path.Join(transpilerDir, toExt, fromPath))
}

// Proxy tries to transpile a file, but if it doesn't exist, it will just return
// the contents of the original file.
func (c *Client) Proxy(fromPath, toExt string) ([]byte, error) {
	data, err := c.Transpile(fromPath, toExt)
	if err != nil {
		if !errors.Is(err, fs.ErrNotExist) {
			return nil, err
		}
		data, err = fs.ReadFile(c.fsys, fromPath)
		if err != nil {
			return nil, err
		}
	}
	return data, nil
}

func (c *Client) Path(fromExt, toExt string) (hops []string, err error) {
	return nil, fmt.Errorf("transpiler: path not implemented yet")
}

// splitRoot splits the root directory off a file path.
func splitRoot(fpath string) (rootDir, remainingPath string) {
	parts := strings.SplitN(fpath, "/", 2)
	if len(parts) == 1 {
		return "", parts[0]
	}
	return parts[0], parts[1]
}
