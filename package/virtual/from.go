package virtual

import (
	"io"
	"io/fs"
)

// From a file to a virtual file
func From(file fs.File) (entry *File, err error) {
	// Get the stats
	stat, err := file.Stat()
	if err != nil {
		return nil, err
	}
	// Copy the directory data over
	if stat.IsDir() {
		return fromDir(file, stat)
	}
	return fromFile(file, stat)
}

func fromDir(file fs.File, stat fs.FileInfo) (entry *File, err error) {
	vdir := &File{
		Path:    stat.Name(),
		ModTime: stat.ModTime(),
		Mode:    stat.Mode(),
	}
	if dir, ok := file.(fs.ReadDirFile); ok {
		des, err := dir.ReadDir(-1)
		if err != nil {
			return nil, err
		}
		vdir.Entries = append(vdir.Entries, des...)
	}
	return vdir, nil
}

func fromFile(file fs.File, stat fs.FileInfo) (entry *File, err error) {
	// Read the data fully
	data, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}
	return &File{
		Path:    stat.Name(),
		Data:    data,
		ModTime: stat.ModTime(),
		Mode:    stat.Mode(),
	}, nil
}
