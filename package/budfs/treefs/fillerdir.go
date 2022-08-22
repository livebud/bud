package treefs

import (
	"fmt"
	"io/fs"

	"github.com/livebud/bud/internal/virtual"
)

type fillerDir struct {
	node *Node
}

func (f *fillerDir) Generate(target string) (fs.File, error) {
	path := f.node.Path()
	// Filler directories must be exact matches with the target, otherwise we'll
	// create files that aren't supposed to exist.
	if target != "." {
		return nil, fmt.Errorf("treefs: path doesn't match target in filler directory %s != %s. %w", path, target, fs.ErrNotExist)
	}
	children := f.node.Children()
	entries := make([]fs.DirEntry, len(children))
	for i, child := range children {
		entries[i] = &dirEntry{child}
	}
	return &virtual.Dir{
		Path:    path,
		Mode:    fs.ModeDir,
		Entries: entries,
	}, nil
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
	file, err := value.Generate(".")
	if err != nil {
		return nil, err
	}
	return file.Stat()
}
