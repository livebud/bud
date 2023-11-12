package virt

import (
	"io/fs"
	"path/filepath"
	"strings"

	"github.com/xlab/treeprint"
)

// Print out a virtual filesystem.
func Print(fsys fs.FS) (string, error) {
	tree, err := printFS(fsys)
	if err != nil {
		return "", err
	}
	return tree.String(), nil
}

func newPrinter() *printer {
	return &printer{
		tree: treeprint.New(),
	}
}

type printer struct {
	tree treeprint.Tree
}

func (t *printer) Add(path string) {
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

func (t *printer) String() string {
	return t.tree.String()
}

func printFS(fsys fs.FS) (*printer, error) {
	printer := newPrinter()
	err := fs.WalkDir(fsys, ".", func(path string, de fs.DirEntry, err error) error {
		if err != nil {
			return err
		} else if path == "." {
			return nil
		}
		printer.Add(path)
		return nil
	})
	if err != nil {
		return nil, err
	}
	return printer, nil
}
