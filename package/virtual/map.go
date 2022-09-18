package virtual

import (
	"io/fs"
	"path"
	"time"
)

type Map map[string]*File

var _ FS = (Map)(nil)

func (m Map) Open(path string) (fs.File, error) {
	if !fs.ValidPath(path) {
		return nil, &fs.PathError{Op: "open", Path: path, Err: fs.ErrInvalid}
	}
	file, ok := m[path]
	if !ok {
		return nil, fs.ErrNotExist
	}
	// Found a file or (empty) directory
	file.Path = path
	if file.IsDir() {
		return &entryDir{&Dir{file.Path, file.Mode, file.ModTime, nil}, 0}, nil
	}
	return &entryFile{file, 0}, nil
}

// Mkdir create a directory.
func (m Map) MkdirAll(path string, perm fs.FileMode) error {
	m[path] = &File{path, nil, perm | fs.ModeDir, time.Time{}}
	return nil
}

// WriteFile writes a file
func (m Map) WriteFile(path string, data []byte, perm fs.FileMode) error {
	m[path] = &File{path, data, perm, time.Time{}}
	return nil
}

// Remove removes a path
func (m Map) RemoveAll(path string) error {
	delete(m, path)
	return nil
}

// Sub returns a submap
func (m Map) Sub(dir string) (FS, error) {
	return &subMap{dir, m}, nil
}

type subMap struct {
	dir string
	m   Map
}

func (s *subMap) Open(filepath string) (fs.File, error) {
	return s.m.Open(path.Join(s.dir, filepath))
}

// Mkdir create a directory.
func (s *subMap) MkdirAll(filepath string, perm fs.FileMode) error {
	return s.m.MkdirAll(path.Join(s.dir, filepath), perm)
}

// WriteFile writes a file
func (s *subMap) WriteFile(filepath string, data []byte, perm fs.FileMode) error {
	return s.m.WriteFile(path.Join(s.dir, filepath), data, perm)
}

// Remove removes a path
func (s *subMap) RemoveAll(filepath string) error {
	return s.m.RemoveAll(path.Join(s.dir, filepath))
}

// Sub returns a submap
func (s *subMap) Sub(dir string) (FS, error) {
	return &subMap{path.Join(s.dir, dir), s.m}, nil
}
