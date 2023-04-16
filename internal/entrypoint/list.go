package entrypoint

import (
	"fmt"
	"io/fs"
	"path"
	"path/filepath"
	"sort"
	"strings"

	"github.com/livebud/bud/package/valid"

	"github.com/matthewmueller/text"
)

// List the views
func List(fsys fs.FS, paths ...string) ([]*View, error) {
	// Build a tree of reserved views (layout, frames, error)
	tree, err := buildTree(fsys, path.Clean(path.Join(paths...)))
	if err != nil {
		return nil, err
	}
	// Turn the tree of views into a list of views
	views, err := listViews(fsys, tree, path.Clean(path.Join(paths...)))
	if err != nil {
		return nil, err
	}
	// Sort the views
	sort.Slice(views, func(i, j int) bool {
		return views[i].Page < views[j].Page
	})
	return views, nil
}

func buildTree(fsys fs.FS, dir string) (*tree, error) {
	fis, err := fs.ReadDir(fsys, dir)
	if err != nil {
		return nil, err
	}
	tree := &tree{
		error:   map[string]Path{},
		layout:  map[string]Path{},
		frame:   map[string]Path{},
		subtree: map[string]*tree{},
	}
	for _, fi := range fis {
		name := fi.Name()
		fullpath := path.Join(dir, name)
		if fi.IsDir() {
			if !valid.Dir(name) {
				continue
			}
			subtree, err := buildTree(fsys, fullpath)
			if err != nil {
				return nil, err
			}
			tree.subtree[name] = subtree
			continue
		}
		ext := path.Ext(name)
		switch extless(name) {
		case "Error":
			tree.error[ext] = Path(fullpath)
		case "Layout":
			tree.layout[ext] = Path(fullpath)
		case "Frame":
			tree.frame[ext] = Path(fullpath)
		}
	}
	return tree, nil
}

// reservedTree of reserved views
type tree struct {
	error   map[string]Path
	layout  map[string]Path
	frame   map[string]Path
	subtree map[string]*tree
}

func (t *tree) Error(dir, ext string) Path {
	root, rest := splitRoot(dir)
	if subtree, ok := t.subtree[root]; ok {
		if error := subtree.Error(rest, ext); error != "" {
			return error
		}
	}
	if error, ok := t.error[ext]; ok {
		return error
	}
	return ""
}

func (t *tree) Layout(dir, ext string) Path {
	root, rest := splitRoot(dir)
	if subtree, ok := t.subtree[root]; ok {
		if layout := subtree.Layout(rest, ext); layout != "" {
			return layout
		}
	}
	if layout, ok := t.layout[ext]; ok {
		return layout
	}
	return ""
}

func (t *tree) Frames(dir, ext string) (frames []Path) {
	if frame, ok := t.frame[ext]; ok {
		frames = append(frames, frame)
	}
	root, rest := splitRoot(dir)
	if subtree, ok := t.subtree[root]; ok {
		frames = append(frames, subtree.Frames(rest, ext)...)
	}
	return frames
}

func splitRoot(dir string) (root, rest string) {
	parts := strings.SplitN(dir, "/", 2)
	if len(parts) == 1 {
		return parts[0], ""
	}
	return parts[0], parts[1]
}

func listViews(fsys fs.FS, tree *tree, dir string) (views []*View, err error) {
	fis, err := fs.ReadDir(fsys, dir)
	if err != nil {
		return nil, err
	}
	for _, fi := range fis {
		name := fi.Name()
		fullpath := path.Join(dir, name)
		if fi.IsDir() {
			if !valid.Dir(name) {
				continue
			}
			subviews, err := listViews(fsys, tree, fullpath)
			if err != nil {
				return nil, err
			}
			views = append(views, subviews...)
			continue
		}
		if !valid.View(name) {
			continue
		}
		ext := path.Ext(name)
		// TODO: remove this constraint after we have sufficient testing
		if ext != ".svelte" {
			continue
		}
		views = append(views, &View{
			Page:   Path(fullpath),
			Client: client(fullpath),
			Route:  route(dir, name),
			Frames: tree.Frames(dir, ext),
			Layout: tree.Layout(dir, ext),
			Error:  tree.Error(dir, ext),
			Type:   strings.TrimPrefix(ext, "."),
			Hot:    ":35729", // TODO: configurable
		})
	}
	return views, nil
}

// Generate the IDs for a nested route
// TODO: consolidate with the function in internal/generator/action/loader.go.
func routeDir(dir string) string {
	segments := strings.Split(dir, "/")
	path := new(strings.Builder)
	for i := 0; i < len(segments); i++ {
		if i%2 != 0 {
			path.WriteString("/")
			path.WriteString(":" + text.Snake(text.Singular(segments[i-1])) + "_id")
			path.WriteString("/")
		}
		path.WriteString(text.Snake(segments[i]))
	}
	if path.Len() == 0 {
		return ""
	}
	return path.String()
}

// Path is the route to the action
func route(dir, name string) string {
	dir = strings.TrimPrefix(strings.TrimPrefix(dir, "view"), "/")
	if dir == "." {
		dir = ""
	}
	dir = routeDir(dir)
	action := extless(name)
	switch action {
	case "show":
		return "/" + path.Join(dir, ":id")
	case "new":
		return "/" + path.Join(dir, "new")
	case "edit":
		return "/" + path.Join(dir, ":id", "edit")
	case "index":
		return "/" + dir
	default:
		return "/" + path.Join(dir, text.Lower(text.Snake(action)))
	}
}

// The client path's entrypoint
//
//	e.g. view/index.jsx => bud/view/_index.jsx
//
// We use the _ because this is a generated entrypoint
// that
func client(name string) string {
	dir, path := filepath.Split(name)
	return fmt.Sprintf("bud/%s_%s.js", dir, path)
}
