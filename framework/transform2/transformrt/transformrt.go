package transformrt

import (
	"fmt"
	"io/fs"
	"path"
	"strings"

	"github.com/livebud/bud/package/log"
	"github.com/livebud/bud/package/trpipe"
)

type File = trpipe.File

type Transform struct {
	Import string // import path
	From   string
	To     string
	Func   func(file *File) error
}

// func (t *Transform) Key() string {
// 	return t.Import + ":" + t.From + ">" + t.To
// }

func Load(log log.Interface, transforms ...*Transform) *Transformer {
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
	fmt.Println(string(code))
	return t.pipeline.Run(fromPath, toExt, code)
}

func parsePath(fpath string) (fromPath, toExt string) {
	base := path.Base(fpath)
	fileParts := strings.SplitN(base, ".", 2)
	if len(fileParts) < 2 {
		return fpath, ""
	}
	extParts := strings.SplitN(fileParts[1], "..", 2)
	if len(extParts) < 2 {
		return fpath, "." + fileParts[1]
	}
	return path.Join(path.Dir(fpath), fileParts[0]) + "." + extParts[0], "." + extParts[1]
}
