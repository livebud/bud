package printfs

import (
	"io/fs"
	"path/filepath"
	"strings"

	"github.com/xlab/treeprint"
)

func New() *Tree {
	return &Tree{
		tree: treeprint.New(),
	}
}

type Tree struct {
	tree treeprint.Tree
}

func (t *Tree) Add(path string) {
	parent := t.tree
	for _, element := range strings.Split(filepath.ToSlash(path), "/") {
		existing := parent.FindByValue(element)
		if existing != nil {
			parent = existing
		} else {
			parent = parent.AddBranch(element)
		}
	}
}

func (t *Tree) String() string {
	return t.tree.String()
}

func Walk(fsys fs.FS) (*Tree, error) {
	tree := New()
	err := fs.WalkDir(fsys, ".", func(path string, de fs.DirEntry, err error) error {
		if err != nil {
			return err
		} else if path == "." {
			return nil
		}
		tree.Add(path)
		return nil
	})
	if err != nil {
		return nil, err
	}
	return tree, nil
}

func Print(fsys fs.FS, dir string) (string, error) {
	tree := New()
	// Set the top-node
	tree.tree.SetValue(dir)
	subfs, err := fs.Sub(fsys, dir)
	if err != nil {
		return "", err
	}
	// Only walk the sub-tree
	err = fs.WalkDir(subfs, ".", func(path string, de fs.DirEntry, err error) error {
		if err != nil {
			return err
		} else if path == "." {
			return nil
		}
		tree.Add(path)
		return nil
	})
	if err != nil {
		return "", err
	}
	return tree.String(), nil
}
