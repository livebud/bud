// Copyright 2020 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fs

import (
	"context"
	"io/fs"
)

// A StatFS is a file system with a Stat method.
type StatFS interface {
	FS

	// Stat returns a FileInfo describing the file.
	// If there is an error, it should be of type *PathError.
	StatContext(ctx context.Context, name string) (fs.FileInfo, error)
}

// Stat returns a FileInfo describing the named file from the file system.
//
// If fs implements StatFS, Stat calls fs.Stat.
// Otherwise, Stat opens the file to stat it.
func Stat(ctx context.Context, fsys FS, name string) (fs.FileInfo, error) {
	if fsys, ok := fsys.(StatFS); ok {
		return fsys.StatContext(ctx, name)
	}
	if fsys, ok := fsys.(OpenFS); ok {
		return stat(ctx, fsys, name)
	}
	return fs.Stat(fsys, name)
}

func stat(ctx context.Context, fsys OpenFS, name string) (fs.FileInfo, error) {
	file, err := fsys.OpenContext(ctx, name)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	return file.Stat()
}
