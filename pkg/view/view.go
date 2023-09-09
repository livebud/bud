package view

import (
	"context"
	"fmt"
	"io"
	"io/fs"
	"path"
	"path/filepath"
	"strings"
)

var ErrNotFound = fmt.Errorf("not found")

type View interface {
	Key() string
	Path() string
	Render(ctx context.Context, s Slot, props any) error
}

type Page struct {
	View   View // Entry
	Frames []View
	Layout View
	Error  View
}

type Finder interface {
	FindPage(key string) (*Page, error)
	FindView(key string) (View, error)
}

type File interface {
	io.Reader
	Path() string
	Stat() (fs.FileInfo, error)
}

type Renderer interface {
	Render(ctx context.Context, s Slot, file File, props any) error
}

type Slot interface {
	io.Reader
	io.Writer
}

func New(fsys fs.FS, renderers map[string]Renderer) *Viewer {
	return &Viewer{fsys, renderers}
}

type vfile struct {
	path string
	file fs.File
}

var _ File = (*vfile)(nil)

func (f *vfile) Path() string {
	return f.path
}

func (f *vfile) Read(p []byte) (n int, err error) {
	return f.file.Read(p)
}

func (f *vfile) Stat() (fs.FileInfo, error) {
	return f.file.Stat()
}

type view struct {
	key      string
	path     string
	fsys     fs.FS
	renderer Renderer
}

var _ View = (*view)(nil)

func (v *view) Key() string {
	return v.key
}

func (v *view) Path() string {
	return v.path
}

func (v *view) Render(ctx context.Context, s Slot, props any) error {
	file, err := v.fsys.Open(v.path)
	if err != nil {
		return err
	}
	defer file.Close()
	return v.renderer.Render(ctx, s, &vfile{v.path, file}, props)
}

type inherited struct {
	Layout map[string]View
	Frames map[string][]View
	Error  map[string]View
}

type Viewer struct {
	fsys      fs.FS
	renderers map[string]Renderer
}

var _ Finder = (*Viewer)(nil)

func (v *Viewer) FindPage(key string) (*Page, error) {
	pages, err := v.findPages()
	if err != nil {
		return nil, err
	}
	page, ok := pages[key]
	if !ok {
		// Try without /index
		page, ok = pages[strings.TrimSuffix(key, "/index")]
		if !ok {
			return nil, fmt.Errorf("%q page %w", key, ErrNotFound)
		}
	}
	return page, nil
}

func (v *Viewer) FindView(key string) (View, error) {
	paths, err := fs.Glob(v.fsys, key+".*")
	if err != nil {
		return nil, err
	}
	if len(paths) == 0 {
		return nil, fmt.Errorf("%q view %w", key, ErrNotFound)
	} else if len(paths) > 1 {
		return nil, fmt.Errorf("%q matches too many paths %v", key, paths)
	}
	renderer, err := v.findRenderer(filepath.Ext(paths[0]))
	if err != nil {
		return nil, err
	}
	return &view{
		key:      key,
		path:     paths[0],
		fsys:     v.fsys,
		renderer: renderer,
	}, nil
}

func (v *Viewer) findPages() (map[string]*Page, error) {
	pages := make(map[string]*Page)
	inherited := &inherited{
		Layout: make(map[string]View),
		Frames: make(map[string][]View),
		Error:  make(map[string]View),
	}
	if err := v.findPagesInDir(pages, inherited, "."); err != nil {
		return nil, err
	}
	return pages, nil
}

func (v *Viewer) findRenderer(ext string) (Renderer, error) {
	renderer, ok := v.renderers[ext]
	if !ok {
		return nil, fmt.Errorf("view: %q renderer %w", ext, ErrNotFound)
	}
	return renderer, nil
}

func (v *Viewer) findPagesInDir(pages map[string]*Page, inherited *inherited, dir string) error {
	des, err := fs.ReadDir(v.fsys, dir)
	if err != nil {
		return err
	}

	// First pass: look for layouts, frames and errors
	for _, de := range des {
		if de.IsDir() {
			continue
		}
		ext := filepath.Ext(de.Name())
		extless := de.Name()[:len(de.Name())-len(ext)]
		switch extless {
		case "layout":
			renderer, err := v.findRenderer(ext)
			if err != nil {
				return err
			}
			inherited.Layout[ext] = &view{
				key:      path.Join(dir, extless),
				path:     path.Join(dir, de.Name()),
				fsys:     v.fsys,
				renderer: renderer,
			}
		case "frame":
			renderer, err := v.findRenderer(ext)
			if err != nil {
				return err
			}
			inherited.Frames[ext] = append(inherited.Frames[ext], &view{
				key:      path.Join(dir, extless),
				path:     path.Join(dir, de.Name()),
				fsys:     v.fsys,
				renderer: renderer,
			})
		case "error":
			renderer, err := v.findRenderer(ext)
			if err != nil {
				return err
			}
			inherited.Error[ext] = &view{
				key:      path.Join(dir, extless),
				path:     path.Join(dir, de.Name()),
				fsys:     v.fsys,
				renderer: renderer,
			}
		}
	}

	// Second pass: go through pages
	for _, de := range des {
		if de.IsDir() {
			continue
		}
		ext := filepath.Ext(de.Name())
		extless := de.Name()[:len(de.Name())-len(ext)]
		switch extless {
		case "layout", "frame", "error":
			continue
		default:
			key := path.Join(dir, extless)
			renderer, err := v.findRenderer(ext)
			if err != nil {
				return err
			}
			pages[key] = &Page{
				View: &view{
					key:      key,
					path:     path.Join(dir, de.Name()),
					fsys:     v.fsys,
					renderer: renderer,
				},
				Layout: inherited.Layout[ext],
				Frames: inherited.Frames[ext],
				Error:  inherited.Error[ext],
			}
		}
	}

	// Third pass: go through directories
	for _, de := range des {
		if !de.IsDir() {
			continue
		}
		if err := v.findPagesInDir(pages, inherited, path.Join(dir, de.Name())); err != nil {
			return err
		}
	}

	return nil
}
