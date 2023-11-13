package ssr

import (
	"encoding/json"
	"fmt"
	"io"
	"io/fs"

	"github.com/dop251/goja"
	"github.com/livebud/bud/internal/js"
	"github.com/livebud/bud/pkg/view"
)

func New(fsys fs.FS, liveUrl string) *Viewer {
	return &Viewer{
		fsys:    fsys,
		liveUrl: liveUrl,
	}
}

type Viewer struct {
	fsys    fs.FS
	liveUrl string
}

var _ view.Viewer = (*Viewer)(nil)

func (v *Viewer) Render(w io.Writer, path string, data *view.Data) error {
	code, err := fs.ReadFile(v.fsys, path)
	if err != nil {
		return err
	}
	// TODO: optimize
	program, err := goja.Compile(path, string(code), false)
	if err != nil {
		return err
	}
	vm := js.New()
	_, err = vm.RunProgram(program)
	if err != nil {
		return err
	}
	props, err := json.Marshal(data.Props)
	if err != nil {
		return err
	}
	result, err := vm.RunString(fmt.Sprintf("bud.render(%s, { liveUrl: %q })", props, v.liveUrl))
	if err != nil {
		return err
	}
	var ssr struct {
		HTML  string          `json:"html"`
		Heads json.RawMessage `json:"heads"`
	}
	err = json.Unmarshal([]byte(result.String()), &ssr)
	if err != nil {
		return err
	}
	if data.Slots != nil {
		if ssr.Heads != nil {
			data.Slots.Slot("heads").Write(ssr.Heads)
		}
	}
	_, err = w.Write([]byte(ssr.HTML))
	return err
}
