package transpiler

import (
	"errors"
	"io/fs"
	"path"
	"strings"

	"github.com/livebud/bud/package/genfs"
	"github.com/livebud/transpiler"
)

type File = transpiler.File
type Interface = transpiler.Interface

func New() *Transpiler {
	return transpiler.New()
}

type Transpiler = transpiler.Transpiler

const transpilerDir = `bud/internal/transpiler`

// splitRoot splits the root directory off a file path.
func splitRoot(fpath string) (rootDir, remainingPath string) {
	parts := strings.SplitN(fpath, "/", 2)
	if len(parts) == 1 {
		return "", parts[0]
	}
	return parts[0], parts[1]
}

// Aliasing allows us to target the transpiler filesystem directly
type FS = fs.FS

// TranspileFile transpiles a file from one extension to another. It assumes
// the transpiler generator is hooked up and serving from the transpiler
// directory.
func TranspileFile(fsys FS, inputPath, toExt string) ([]byte, error) {
	return fs.ReadFile(fsys, path.Join(transpilerDir, toExt, inputPath))
}

func Serve(tr transpiler.Interface, fsys FS, file *genfs.File) error {
	toExt, inputPath := splitRoot(file.Relative())
	input, err := fs.ReadFile(fsys, inputPath)
	if err != nil {
		return err
	}
	output, err := tr.Transpile(file.Ext(), toExt, input)
	if err != nil {
		return err
	}
	file.Data = output
	return nil
}

// Proxy a filesystem through the transpiler
type Proxy struct {
	FS FS
}

func (p *Proxy) Open(name string) (fs.File, error) {
	file, err := p.FS.Open(path.Join(transpilerDir, path.Ext(name), name))
	if err != nil {
		if !errors.Is(err, fs.ErrNotExist) {
			return nil, err
		}
		return p.FS.Open(name)
	}
	return file, nil
}
