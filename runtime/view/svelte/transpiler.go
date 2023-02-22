package svelte

import (
	"context"

	"github.com/livebud/bud/package/log"
	"github.com/livebud/transpiler"
)

func NewTranspiler(compiler *Compiler, log log.Log) *Transpiler {
	return &Transpiler{compiler, log}
}

type Transpiler struct {
	compiler *Compiler
	log      log.Log
}

func (t *Transpiler) SvelteToSsrJs(file *transpiler.File) error {
	t.log.Debug("svelte: transpiling %s to .ssr.js", file.Path())
	ssr, err := t.compiler.SSR(context.TODO(), file.Path(), file.Data)
	if err != nil {
		return err
	}
	file.Data = []byte(ssr.JS)
	return nil
}

func (t *Transpiler) SvelteToDomJs(file *transpiler.File) error {
	t.log.Debug("svelte: transpiling %s to .dom.js", file.Path())
	dom, err := t.compiler.DOM(context.TODO(), file.Path(), file.Data)
	if err != nil {
		return err
	}
	file.Data = []byte(dom.JS)
	return nil
}
