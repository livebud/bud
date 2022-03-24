package entrypoint

import (
	"fmt"
	"io/fs"
	"path/filepath"
)

func FindByPage(fsys fs.FS, page string) (view *View, err error) {
	views, err := List(fsys, "view")
	if err != nil {
		return nil, err
	}
	page = filepath.Clean(page)
	for _, view := range views {
		if string(view.Page) == page {
			return view, nil
		}
	}
	return nil, fmt.Errorf("unable to find view by page %q", page)
}

func FindByClient(fsys fs.FS, client string) (view *View, err error) {
	views, err := List(fsys, "view")
	if err != nil {
		return nil, err
	}
	for _, view := range views {
		if view.Client == client {
			return view, nil
		}
	}
	return nil, fmt.Errorf("unable to find view by client path %q", client)
}
