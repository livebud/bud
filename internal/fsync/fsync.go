package fsync

import (
	"errors"
	"io/fs"
	"path/filepath"
	"strconv"

	"github.com/go-duo/bud/internal/fsync/set"
	"github.com/go-duo/bud/internal/vfs"
)

// TODO: read files during diff so we ensure we're successful before writing.
// TODO: update should compare stamps at the time of writing, not before.

// Dir syncs the source directory from the source filesystem to the target directory
// in the target filesystem
func Dir(sfs fs.FS, sdir string, tfs vfs.ReadWritable, tdir string) error {
	ops, err := diff(sfs, sdir, tfs, tdir)
	if err != nil {
		return err
	}
	return apply(sfs, tfs, ops)
}

type OpType uint8

func (ot OpType) String() string {
	switch ot {
	case CreateType:
		return "create"
	case UpdateType:
		return "update"
	case DeleteType:
		return "delete"
	default:
		return ""
	}
}

const (
	CreateType OpType = iota + 1
	UpdateType
	DeleteType
)

type Op struct {
	Type OpType
	Path string
}

func (o Op) String() string {
	return o.Type.String() + ":" + o.Path
}

func diff(sfs fs.FS, sdir string, tfs vfs.ReadWritable, tdir string) (ops []Op, err error) {
	sourceEntries, err := fs.ReadDir(sfs, sdir)
	if err != nil {
		return nil, err
	}
	targetEntries, err := fs.ReadDir(tfs, tdir)
	if err != nil && !errors.Is(err, fs.ErrNotExist) {
		return nil, err
	}
	sourceSet := set.New(sourceEntries...)
	targetSet := set.New(targetEntries...)
	creates := set.Difference(sourceSet, targetSet)
	deletes := set.Difference(targetSet, sourceSet)
	updates := set.Intersection(sourceSet, targetSet)
	createOps, err := createOps(sfs, sdir, creates.List())
	if err != nil {
		return nil, err
	}
	deleteOps := deleteOps(sdir, deletes.List())
	childOps, err := updateOps(sfs, sdir, tfs, tdir, updates.List())
	if err != nil {
		return nil, err
	}
	ops = append(ops, createOps...)
	ops = append(ops, deleteOps...)
	ops = append(ops, childOps...)
	return ops, nil
}

func createOps(sfs fs.FS, dir string, des []fs.DirEntry) (ops []Op, err error) {
	for _, de := range des {
		path := filepath.Join(dir, de.Name())
		if !de.IsDir() {
			ops = append(ops, Op{CreateType, path})
			continue
		}
		des, err := fs.ReadDir(sfs, path)
		if err != nil {
			return nil, err
		}
		createOps, err := createOps(sfs, path, des)
		if err != nil {
			return nil, err
		}
		ops = append(ops, createOps...)
	}
	return ops, nil
}

func deleteOps(dir string, des []fs.DirEntry) (ops []Op) {
	for _, de := range des {
		path := filepath.Join(dir, de.Name())
		ops = append(ops, Op{DeleteType, path})
		continue
	}
	return ops
}

func updateOps(sfs fs.FS, sdir string, tfs vfs.ReadWritable, tdir string, des []fs.DirEntry) (ops []Op, err error) {
	for _, de := range des {
		sourcePath := filepath.Join(sdir, de.Name())
		targetPath := filepath.Join(tdir, de.Name())
		// Recurse directories
		if de.IsDir() {
			childOps, err := diff(sfs, sourcePath, tfs, targetPath)
			if err != nil {
				return nil, err
			}
			ops = append(ops, childOps...)
			continue
		}
		// Otherwise, check if the file has changed
		sourceStamp, err := stamp(sfs, sourcePath)
		if err != nil {
			return nil, err
		}
		targetStamp, err := stamp(tfs, targetPath)
		if err != nil {
			return nil, err
		}
		if sourceStamp != targetStamp {
			ops = append(ops, Op{UpdateType, sourcePath})
		}
	}
	return ops, nil
}

func apply(sfs fs.FS, tfs vfs.ReadWritable, ops []Op) error {
	for _, op := range ops {
		switch op.Type {
		case CreateType:
			dir := filepath.Dir(op.Path)
			if err := tfs.MkdirAll(dir, 0755); err != nil {
				return err
			}
			data, err := fs.ReadFile(sfs, op.Path)
			if err != nil {
				return err
			}
			if err := tfs.WriteFile(op.Path, data, 0644); err != nil {
				return err
			}
		case UpdateType:
			data, err := fs.ReadFile(sfs, op.Path)
			if err != nil {
				return err
			}
			if err := tfs.WriteFile(op.Path, data, 0644); err != nil {
				return err
			}
		case DeleteType:
			if err := tfs.RemoveAll(op.Path); err != nil {
				return err
			}
		}
	}
	return nil
}

// Stamp the path, returning "" if the file doesn't exist.
// Uses the modtime and size to determine if a file has changed.
func stamp(fsys fs.FS, path string) (stamp string, err error) {
	stat, err := fs.Stat(fsys, path)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return "-1:-1", nil
		}
		return "", err
	}
	mtime := stat.ModTime().UnixNano()
	size := stat.Size()
	stamp = strconv.Itoa(int(size)) + ":" + strconv.Itoa(int(mtime))
	return stamp, nil
}
