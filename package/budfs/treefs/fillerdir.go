package treefs

import (
	"errors"
	"fmt"
	"io/fs"

	"github.com/livebud/bud/package/virtual"
)

type fillerDir struct {
	node *Node
}

func (f *fillerDir) Generate(target string) (fs.File, error) {
	path := f.node.Path()
	// Filler directories must be exact matches with the target, otherwise we'll
	// create files that aren't supposed to exist.
	if target != path {
		return nil, fmt.Errorf("treefs: path doesn't match target in filler directory %s != %s", path, target)
	}
	children := f.node.Children()
	var entries []fs.DirEntry
	// TODO: run in parallel
	for _, child := range children {
		de := &dirEntry{child}
		// Stat to ensure the file exists before adding it as a directory entry
		if _, err := de.Info(); err != nil {
			if errors.Is(err, fs.ErrNotExist) {
				continue
			}
			return nil, err
		}
		entries = append(entries, de)
	}
	return virtual.New(&virtual.Dir{
		Path:    path,
		Mode:    fs.ModeDir,
		Entries: entries,
	}), nil
}

type dirEntry struct {
	node *Node
}

var _ fs.DirEntry = (*dirEntry)(nil)

func (e *dirEntry) Name() string {
	return e.node.name
}

func (e *dirEntry) IsDir() bool {
	return e.node.mode.IsDir()
}

func (e *dirEntry) Type() fs.FileMode {
	return e.node.mode
}

func (e *dirEntry) Info() (fs.FileInfo, error) {
	value := e.node.generator
	if value == nil {
		value = &fillerDir{e.node}
	}
	file, err := value.Generate(e.node.Path())
	if err != nil {
		return nil, err
	}
	return file.Stat()
}
