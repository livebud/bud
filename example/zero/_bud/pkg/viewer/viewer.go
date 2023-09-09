package viewer

import (
	context "context"
	fmt "fmt"
	transpiler "github.com/livebud/bud/example/zero/bud/pkg/transpiler"
	gohtml "github.com/livebud/bud/example/zero/viewer/gohtml"
	router "github.com/livebud/bud/package/router"
	view "github.com/livebud/bud/runtime/view"
	fs "io/fs"
)

// Load the viewer
// TODO: turn most of this code into a runtime library
func New(
	transpiler transpiler.Transpiler,
	gohtml *gohtml.Viewer,
) Viewer {
	return &viewer{
		transpiler: transpiler,
		viewers: map[string]view.Viewer{
			".gohtml": gohtml,
		},
		accepts: []string{
			".gohtml",
		},
	}
}

type Viewer = view.Viewer

type viewer struct {
	transpiler transpiler.Transpiler
	viewers    map[string]view.Viewer
	accepts    []string
}

var _ Viewer = (*viewer)(nil)

func (v *viewer) Register(r *router.Router, pages []*view.Page) {
	for _, viewer := range v.viewers {
		viewer.Register(r, pages)
	}
}

func (v *viewer) Render(ctx context.Context, fsys fs.FS, page *view.Page, propMap view.PropMap) ([]byte, error) {
	viewer, ok := v.viewers[page.Ext]
	if ok {
		return viewer.Render(ctx, fsys, page, propMap)
	}
	// TODO: don't choose best when embedded
	ext, err := v.transpiler.Best(page.Ext, v.accepts)
	if err != nil {
		return nil, fmt.Errorf("viewer: unable to render %q. %w", page.Path, err)
	}
	viewer, ok = v.viewers[ext]
	if ok {
		return viewer.Render(ctx, fsys, page, propMap)
	}
	return nil, fmt.Errorf("viewer: unable to find acceptable viewer to render %q", page.Path)
}

func (v *viewer) RenderError(ctx context.Context, fsys fs.FS, page *view.Page, propMap view.PropMap, err error) []byte {
	viewer, ok := v.viewers[page.Error.Ext]
	if ok {
		return viewer.RenderError(ctx, fsys, page, propMap, err)
	}
	// TODO: don't choose best when embedded
	ext, err := v.transpiler.Best(page.Error.Ext, v.accepts)
	if err != nil {
		msg := fmt.Sprintf("viewer: unable to find extension to render error page %q for error %s", page.Error.Path, err)
		return []byte(msg)
	}
	viewer, ok = v.viewers[ext]
	if ok {
		return viewer.RenderError(ctx, fsys, page, propMap, err)
	}
	msg := fmt.Sprintf("viewer: unable to find acceptable viewer to render error page %q for error %s", page.Error.Path, err)
	return []byte(msg)
}

func (v *viewer) Bundle(ctx context.Context, fsys fs.FS, pages view.Pages, embeds view.Embeds) error {
	exts := map[string]map[string]*view.Page{}
	for _, page := range pages {
		if _, ok := v.viewers[page.Ext]; ok {
			if exts[page.Ext] == nil {
				exts[page.Ext] = map[string]*view.Page{}
			}
			exts[page.Ext][page.Path] = page
			continue
		}
		ext, err := v.transpiler.Best(page.Ext, v.accepts)
		if err != nil {
			return fmt.Errorf("viewer: unable find viewer to bundle %q. %w", page.Path, err)
		}
		if exts[ext] == nil {
			exts[ext] = map[string]*view.Page{}
		}
		exts[ext][page.Path] = page
	}
	// TODO: consider parallelizing this
	for ext, pages := range exts {
		viewer, ok := v.viewers[ext]
		if !ok {
			return fmt.Errorf("viewer: unable find viewer for %q", ext)
		}
		if err := viewer.Bundle(ctx, fsys, pages, embeds); err != nil {
			return fmt.Errorf("viewer: unable to bundle %q. %w", ext, err)
		}
	}
	return nil
}
