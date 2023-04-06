package viewer

import (
	"io/fs"
	"path"
	"path/filepath"
)

func NewFinder(fsys fs.FS) Finder {
	return &finder{fsys}
}

type finder struct {
	fsys fs.FS
}

func (f *finder) Find(root string) (pages Pages, err error) {
	return Find(f.fsys)
}

type Finder interface {
	Find(root string) (Pages, error)
}

// Find pages
func Find(fsys fs.FS) (pages map[Key]*Page, err error) {
	pages = make(map[Key]*Page)
	inherited := &inherited{
		Layout: make(map[ext]*View),
		Frames: make(map[ext][]*View),
		Error:  make(map[ext]*View),
	}
	if err := find(fsys, pages, inherited, "."); err != nil {
		return nil, err
	}
	return pages, nil
}

type ext = string

type inherited struct {
	Layout map[ext]*View
	Frames map[ext][]*View
	Error  map[ext]*View
}

func find(fsys fs.FS, pages map[Key]*Page, inherited *inherited, dir string) error {
	des, err := fs.ReadDir(fsys, dir)
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
			inherited.Layout[ext] = &View{
				Path: path.Join(dir, de.Name()),
				Key:  path.Join(dir, extless),
				Ext:  ext,
			}
		case "frame":
			inherited.Frames[ext] = append(inherited.Frames[ext], &View{
				Path: path.Join(dir, de.Name()),
				Key:  path.Join(dir, extless),
				Ext:  ext,
			})
		case "error":
			inherited.Error[ext] = &View{
				Path: path.Join(dir, de.Name()),
				Key:  path.Join(dir, extless),
				Ext:  ext,
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
			pages[key] = &Page{
				View: &View{
					Path: path.Join(dir, de.Name()),
					Key:  key,
					Ext:  ext,
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
		if err := find(fsys, pages, inherited, de.Name()); err != nil {
			return err
		}
	}

	return nil
}
