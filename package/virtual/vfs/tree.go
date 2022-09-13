package vfs

import (
	"io/fs"
	"path"
	"sort"
	"strings"
	"time"
)

type Tree map[string]*Entry

var _ FS = (Tree)(nil)

func (fsys Tree) Open(path string) (fs.File, error) {
	if !fs.ValidPath(path) {
		return nil, &fs.PathError{Op: "open", Path: path, Err: fs.ErrInvalid}
	}

	// The following logic is based on "testing/fstest".MapFS.Open
	entry, ok := fsys[path]
	if ok {
		if entry.Mode.IsDir() {
			return &entryDir{entry, nil, path, 0}, nil
		}
		return &entryFile{entry, path, 0}, nil
	}

	// Directory, possibly synthesized.
	// Note that file can be nil here: the map need not contain explicit parent directories for all its files.
	// But file can also be non-nil, in case the user wants to set metadata for the directory explicitly.
	// Either way, we need to construct the list of children of this directory.
	var list []fs.DirEntry
	var need = make(map[string]bool)
	if path == "." {
		for fname, entry := range fsys {
			i := strings.Index(fname, "/")
			if i < 0 {
				if fname != "." {
					list = append(list, &dirEntry{entry, fname})
				}
			} else {
				need[fname[:i]] = true
			}
		}
	} else {
		prefix := path + "/"
		for fname, entry := range fsys {
			if strings.HasPrefix(fname, prefix) {
				felem := fname[len(prefix):]
				i := strings.Index(felem, "/")
				if i < 0 {
					list = append(list, &dirEntry{entry, felem})
				} else {
					need[fname[len(prefix):len(prefix)+i]] = true
				}
			}
		}
		// If the directory name is not in the map,
		// and there are no children of the name in the map,
		// then the directory is treated as not existing.
		if entry == nil && list == nil && len(need) == 0 {
			return nil, &fs.PathError{Op: "open", Path: path, Err: fs.ErrNotExist}
		}
	}
	for _, fi := range list {
		delete(need, fi.Name())
	}
	for path := range need {
		dir := &Entry{nil, fs.ModeDir, time.Time{}}
		list = append(list, &dirEntry{dir, path})
	}
	sort.Slice(list, func(i, j int) bool {
		return list[i].Name() < list[j].Name()
	})

	if entry == nil {
		entry = &Entry{nil, fs.ModeDir, time.Time{}}
	}
	return &entryDir{entry, list, path, 0}, nil
}

// Mkdir create a directory.
func (t Tree) MkdirAll(path string, perm fs.FileMode) error {
	t[path] = &Entry{nil, perm | fs.ModeDir, time.Time{}}
	return nil
}

// WriteFile writes a file
// TODO: WriteFile should fail if path.Dir(name) doesn't exist
func (t Tree) WriteFile(path string, data []byte, perm fs.FileMode) error {
	t[path] = &Entry{data, perm, time.Time{}}
	return nil
}

// Remove removes a path
func (t Tree) RemoveAll(path string) error {
	delete(t, path)
	return nil
}

// Sub returns a subtree
func (t Tree) Sub(dir string) (FS, error) {
	return &subTree{dir, t}, nil
}

type subTree struct {
	dir  string
	tree Tree
}

func (s *subTree) Open(filepath string) (fs.File, error) {
	return s.tree.Open(path.Join(s.dir, filepath))
}

// Mkdir create a directory.
func (s *subTree) MkdirAll(filepath string, perm fs.FileMode) error {
	return s.tree.MkdirAll(path.Join(s.dir, filepath), perm)
}

// WriteFile writes a file
func (s *subTree) WriteFile(filepath string, data []byte, perm fs.FileMode) error {
	return s.tree.WriteFile(path.Join(s.dir, filepath), data, perm)
}

// Remove removes a path
func (s *subTree) RemoveAll(filepath string) error {
	return s.tree.RemoveAll(path.Join(s.dir, filepath))
}

// Sub returns a subtree
func (s *subTree) Sub(dir string) (FS, error) {
	return &subTree{path.Join(s.dir, dir), s.tree}, nil
}
