package react

import (
	"io/fs"

	"github.com/livebud/bud/package/viewer"
	"github.com/livebud/bud/runtime/transpiler"
	"github.com/livebud/js"
)

func New(fsys fs.FS, js js.VM, transpiler transpiler.Interface, pages viewer.Pages) *Viewer {
	return &Viewer{fsys, js, pages, transpiler}
}

type Viewer struct {
	fsys       fs.FS
	js         js.VM
	pages      viewer.Pages
	transpiler transpiler.Interface
}
