package transpiler

import (
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

// TranspileFile transpiles a file from one extension to another. It assumes
// the transpiler generator is hooked up and serving from the transpiler
// directory.
func TranspileFile(fsys fs.FS, inputPath, toExt string) ([]byte, error) {
	return fs.ReadFile(fsys, path.Join(transpilerDir, toExt, inputPath))
}

func Serve(tr transpiler.Interface, fsys fs.FS, file *genfs.File) error {
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
