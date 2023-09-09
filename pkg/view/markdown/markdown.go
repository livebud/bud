package markdown

import (
	"context"
	"io"

	"github.com/livebud/bud/pkg/view"
	"github.com/yuin/goldmark"
)

func New() view.Renderer {
	return &renderer{}
}

type renderer struct {
}

func (r *renderer) Render(ctx context.Context, s view.Slot, file view.File, props any) error {
	code, err := io.ReadAll(file)
	if err != nil {
		return err
	}
	if err := goldmark.Convert(code, s); err != nil {
		return err
	}
	return nil
}
