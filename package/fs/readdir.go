// Copyright 2020 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fs

import (
	"context"
	"errors"
	"io/fs"
	"sort"
)

// ReadDirFS is the interface implemented by a file system
// that provides an optimized implementation of ReadDir.
type ReadDirFS interface {
	FS

	// ReadDir reads the named directory
	// and returns a list of directory entries sorted by filename.
	ReadDirContext(ctx context.Context, name string) ([]fs.DirEntry, error)
}

// ReadDir reads the named directory
// and returns a list of directory entries sorted by filename.
//
// If fs implements ReadDirFS, ReadDir calls fs.ReadDir.
// Otherwise ReadDir calls fs.Open and uses ReadDir and Close
// on the returned file.
func ReadDir(ctx context.Context, fsys FS, name string) ([]fs.DirEntry, error) {
	if fsys, ok := fsys.(ReadDirFS); ok {
		return fsys.ReadDirContext(ctx, name)
	}
	if fsys, ok := fsys.(OpenFS); ok {
		return readDir(ctx, fsys, name)
	}
	return fs.ReadDir(fsys, name)
}

func readDir(ctx context.Context, fsys OpenFS, name string) ([]fs.DirEntry, error) {
	file, err := fsys.OpenContext(ctx, name)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	dir, ok := file.(ReadDirFile)
	if !ok {
		return nil, &fs.PathError{Op: "readdir", Path: name, Err: errors.New("not implemented")}
	}

	list, err := dir.ReadDir(-1)
	sort.Slice(list, func(i, j int) bool { return list[i].Name() < list[j].Name() })
	return list, err
}
