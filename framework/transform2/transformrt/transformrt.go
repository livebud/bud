package transformrt

import (
	"io/fs"
	"strings"

	"github.com/livebud/bud/package/log"
	"github.com/livebud/bud/package/trpipe"
)

type File = trpipe.File

type Transform struct {
	From string
	To   string
	Func func(file *File) error
}

func Load(log log.Log, transforms ...*Transform) *Transformer {
	pipeline := trpipe.New(log)
	for _, t := range transforms {
		pipeline.Add(t.From, t.To, t.Func)
	}
	return &Transformer{pipeline}
}

type Transformer struct {
	pipeline *trpipe.Pipeline
}

func (t *Transformer) Transform(fsys fs.FS, fpath string) ([]byte, error) {
	fromPath, toExt := parsePath(fpath)
	code, err := fs.ReadFile(fsys, fromPath)
	if err != nil {
		return nil, err
	}
	return t.pipeline.Run(fromPath, toExt, code)
}

func parsePath(fpath string) (fromPath, toExt string) {
	root, rest := splitRoot(fpath)
	return rest, "." + root
}

func splitRoot(dir string) (root, rest string) {
	parts := strings.SplitN(dir, "/", 2)
	if len(parts) == 1 {
		return parts[0], ""
	}
	return parts[0], parts[1]
}
