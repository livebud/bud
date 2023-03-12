package cli

import (
	"context"
	"fmt"
	"io/fs"
	"path"
	"path/filepath"

	"github.com/livebud/bud/framework"
	"github.com/livebud/bud/package/virtual"
	"golang.org/x/tools/txtar"
)

type ToolFsTxtar struct {
	Flag *framework.Flag
	Path string
}

func (c *CLI) ToolFsTxtar(ctx context.Context, in *ToolFsTxtar) error {
	// Generate bud files
	generate := &Generate{Flag: in.Flag}
	if err := c.Generate(ctx, generate); err != nil {
		return err
	}

	// Get the current working directory
	abs, err := filepath.Abs(c.Dir)
	if err != nil {
		return err
	}
	fsys := virtual.OS(abs)

	// Walk the directory, adding all the files to the archive
	ar := new(txtar.Archive)
	dir := path.Clean(in.Path)
	err = fs.WalkDir(fsys, dir, func(path string, de fs.DirEntry, err error) error {
		if err != nil {
			return err
		} else if de.IsDir() {
			return nil
		}
		code, err := fs.ReadFile(fsys, path)
		if err != nil {
			return err
		}
		if isBinary(code) {
			return nil
		}
		ar.Files = append(ar.Files, txtar.File{
			Name: path,
			Data: code,
		})
		return nil
	})
	if err != nil {
		return err
	}

	// Print the archive to stdout
	fmt.Fprintln(c.Stdout, string(txtar.Format(ar)))
	return nil
}

// Check if the given byte slice contains any null bytes. Seems to be a good
// enough heuristic for detecting binary files.
func isBinary(code []byte) bool {
	for _, b := range code {
		if b == 0 {
			return true
		}
	}
	return false
}
