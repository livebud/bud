package virtual

import (
	"io/fs"
	"path"
	"sort"
	"strings"
	"time"
)

type List []*File

var _ FS = (*List)(nil)

func (fsys List) Open(path string) (fs.File, error) {
	if !fs.ValidPath(path) {
		return nil, &fs.PathError{Op: "open", Path: path, Err: fs.ErrInvalid}
	}
	file, ok := fsys.find(path)
	if !ok {
		return nil, fs.ErrNotExist
	}
	// Found a file or directory
	file.Path = path
	if !file.IsDir() {
		return &openFile{file, 0}, nil
	}
	// The following logic is based on "testing/fstest".MapFS.Open
	// Directory, possibly synthesized.
	// Note that file can be nil here: the map need not contain explicit parent directories for all its files.
	// But file can also be non-nil, in case the user wants to set metadata for the directory explicitly.
	// Either way, we need to construct the list of children of this directory.
	var list []fs.DirEntry
	var need = make(map[string]bool)
	if path == "." {
		for _, file := range fsys {
			fname := file.Path
			i := strings.Index(fname, "/")
			if i < 0 {
				if fname != "." {
					file.Path = fname
					list = append(list, file)
				}
			} else {
				need[fname[:i]] = true
			}
		}
	} else {
		prefix := path + "/"
		for _, file := range fsys {
			fname := file.Path
			if strings.HasPrefix(fname, prefix) {
				felem := fname[len(prefix):]
				i := strings.Index(felem, "/")
				if i < 0 {
					file.Path = felem
					list = append(list, file)
				} else {
					need[fname[len(prefix):len(prefix)+i]] = true
				}
			}
		}
		// If the directory name is not in the map,
		// and there are no children of the name in the map,
		// then the directory is treated as not existing.
		if file == nil && list == nil && len(need) == 0 {
			return nil, &fs.PathError{Op: "open", Path: path, Err: fs.ErrNotExist}
		}
	}
	for _, fi := range list {
		delete(need, fi.Name())
	}
	for path := range need {
		dir := &File{path, nil, fs.ModeDir, time.Time{}, nil}
		list = append(list, dir)
	}
	sort.Slice(list, func(i, j int) bool {
		return list[i].Name() < list[j].Name()
	})
	// Return the synthesized entries as a directory.
	return &openDir{&File{path, nil, fs.ModeDir, time.Time{}, list}, 0}, nil
}

// Mkdir create a directory.
func (fsys *List) MkdirAll(path string, perm fs.FileMode) error {
	file, ok := fsys.find(path)
	if ok {
		if file.IsDir() {
			return nil
		}
		return &fs.PathError{
			Op:   "MkdirAll",
			Path: path,
			Err:  fs.ErrExist,
		}
	}
	*fsys = append(*fsys, &File{path, nil, perm | fs.ModeDir, time.Time{}, nil})
	return nil
}

// WriteFile writes a file
func (fsys *List) WriteFile(path string, data []byte, perm fs.FileMode) error {
	file, ok := fsys.find(path)
	if ok {
		if file.IsDir() {
			return &fs.PathError{
				Op:   "WriteFile",
				Path: path,
				Err:  fs.ErrExist,
			}
		}
		file.Data = data
		file.Mode = perm
		return nil
	}
	*fsys = append(*fsys, &File{path, data, perm, time.Time{}, nil})
	return nil
}

// Remove removes a path
func (fsys *List) RemoveAll(path string) error {
	idx := fsys.indexOf(path)
	if idx < 0 {
		return nil
	}
	*fsys = append((*fsys)[:idx], (*fsys)[idx+1:]...)
	return nil
}

// Sub returns a submap
func (fsys List) Sub(dir string) (FS, error) {
	return &subList{dir, fsys}, nil
}

type subList struct {
	dir string
	m   List
}

func (s *subList) Open(filepath string) (fs.File, error) {
	return s.m.Open(path.Join(s.dir, filepath))
}

// Mkdir create a directory.
func (s *subList) MkdirAll(filepath string, perm fs.FileMode) error {
	return s.m.MkdirAll(path.Join(s.dir, filepath), perm)
}

// WriteFile writes a file
func (s *subList) WriteFile(filepath string, data []byte, perm fs.FileMode) error {
	return s.m.WriteFile(path.Join(s.dir, filepath), data, perm)
}

// Remove removes a path
func (s *subList) RemoveAll(filepath string) error {
	return s.m.RemoveAll(path.Join(s.dir, filepath))
}

// Sub returns a subList
func (s *subList) Sub(dir string) (FS, error) {
	return &subList{path.Join(s.dir, dir), s.m}, nil
}

func (l List) find(path string) (f *File, ok bool) {
	for _, file := range l {
		if file.Path == path {
			return file, true
		}
	}
	return nil, false
}

func (l List) indexOf(path string) (i int) {
	for i, file := range l {
		if file.Path == path {
			return i
		}
	}
	return -1
}
