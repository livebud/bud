package virtual

import "io/fs"

type Map map[string]Entry

func (m Map) Open(name string) (fs.File, error) {
	entry, ok := m[name]
	if !ok {
		return nil, fs.ErrNotExist
	}
	return entry.Open(), nil
}

// Mkdir create a directory.
// TODO: Mkdir should fail if path.Dir(dirpath) doesn't exist
func (m Map) Mkdir(dirpath string, perm fs.FileMode) error {
	m[dirpath] = &Dir{
		Path: dirpath,
		Mode: perm | fs.ModeDir,
	}
	return nil
}

// WriteFile writes a file
// TODO: WriteFile should fail if path.Dir(name) doesn't exist
func (m Map) WriteFile(name string, data []byte, perm fs.FileMode) error {
	m[name] = &File{
		Path: name,
		Data: data,
		Mode: perm,
	}
	return nil
}

// Remove removes a path
// TODO: Remove should fail if path doesn't exist
func (m Map) Remove(path string) error {
	delete(m, path)
	return nil
}
