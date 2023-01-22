package transpiler

import (
	"io/fs"
	"path"
	"strings"

	"github.com/livebud/transpiler"
)

type File = transpiler.File
type Interface = transpiler.Interface

var New = transpiler.New

const transpilerDir = `bud/internal/transpiler`

// SplitRoot splits the root directory off a file path.
func SplitRoot(fpath string) (rootDir, remainingPath string) {
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
