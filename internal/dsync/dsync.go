package dsync

import (
	"errors"
	"io/fs"
	"path/filepath"
	"strconv"

	"github.com/livebud/bud/internal/dsync/set"
	"github.com/livebud/bud/package/vfs"
)

type skipFunc = func(name string, isDir bool) bool

type option struct {
	Skip skipFunc
	rel  func(spath string) (string, error)
}

type Option func(o *option)

// Provide a skip function
//
// Note: try to skip as high up in the tree as possible.
// E.g. if the source filesystem doesn't have bud, it will
// delete bud, even if you're skipping bud/generate.
func WithSkip(skips ...skipFunc) Option {
	return func(o *option) {
		o.Skip = composeSkips(skips)
	}
}

func composeSkips(skips []skipFunc) skipFunc {
	return func(name string, isDir bool) bool {
		for _, skip := range skips {
			if skip(name, isDir) {
				return true
			}
		}
		return false
	}
}

func Rel(sdir, tdir string) func(path string) (string, error) {
	return func(spath string) (string, error) {
		rel, err := filepath.Rel(sdir, spath)
		if err != nil {
			return "", err
		}
		return filepath.Join(tdir, rel), nil
	}
}

// Dir syncs the source directory from the source filesystem to the target directory
// in the target filesystem
func Dir(sfs fs.FS, sdir string, tfs vfs.ReadWritable, tdir string, options ...Option) error {
	opt := &option{
		Skip: func(name string, isDir bool) bool { return false },
		rel:  Rel(sdir, tdir),
	}
	for _, option := range options {
		option(opt)
	}
	ops, err := diff(opt, sfs, sdir, tfs, tdir)
	if err != nil {
		return err
	}
	err = apply(sfs, tfs, ops)
	return err
}

// To syncs the "to" directory from the source to target filesystem
func To(sfs fs.FS, tfs vfs.ReadWritable, to string, options ...Option) error {
	return Dir(sfs, to, tfs, to, options...)
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
	Data []byte
}

func (o Op) String() string {
	return o.Type.String() + ":" + o.Path
}

func diff(opt *option, sfs fs.FS, sdir string, tfs vfs.ReadWritable, tdir string) (ops []Op, err error) {
	sourceEntries, err := fs.ReadDir(sfs, sdir)
	if err != nil {
		return nil, err
	}
	targetEntries, err := fs.ReadDir(tfs, tdir)
	if err != nil && !errors.Is(err, fs.ErrNotExist) {
		return nil, err
	}
	// Create the source set from the source entries
	sourceSet := set.NewWithSize(len(sourceEntries))
	for _, de := range sourceEntries {
		// Ensure all sources actually exist
		if _, err := de.Info(); err != nil {
			if errors.Is(err, fs.ErrNotExist) {
				continue
			}
			return nil, err
		}
		sourceSet.Add(de)
	}
	// Create a target set from the target entries
	targetSet := set.NewWithSize(len(targetEntries))
	for _, de := range targetEntries {
		// Ensure all sources actually exist
		if _, err := de.Info(); err != nil {
			if errors.Is(err, fs.ErrNotExist) {
				continue
			}
			return nil, err
		}
		targetSet.Add(de)
	}
	creates := set.Difference(sourceSet, targetSet)
	deletes := set.Difference(targetSet, sourceSet)
	updates := set.Intersection(sourceSet, targetSet)
	createOps, err := createOps(opt, sfs, sdir, creates.List())
	if err != nil {
		return nil, err
	}
	deleteOps, err := deleteOps(opt, sdir, deletes.List())
	if err != nil {
		return nil, err
	}
	childOps, err := updateOps(opt, sfs, sdir, tfs, tdir, updates.List())
	if err != nil {
		return nil, err
	}
	ops = append(ops, createOps...)
	ops = append(ops, deleteOps...)
	ops = append(ops, childOps...)
	return ops, nil
}

func createOps(opt *option, sfs fs.FS, dir string, des []fs.DirEntry) (ops []Op, err error) {
	for _, de := range des {
		if de.Name() == "." {
			continue
		}
		path := filepath.Join(dir, de.Name())
		if opt.Skip(path, de.IsDir()) {
			continue
		}
		if !de.IsDir() {
			data, err := fs.ReadFile(sfs, path)
			if err != nil {
				// Don't error out on files that don't exist
				if errors.Is(err, fs.ErrNotExist) {
					continue
				}
				return nil, err
			}
			rel, err := opt.rel(path)
			if err != nil {
				return nil, err
			}
			ops = append(ops, Op{CreateType, rel, data})
			continue
		}
		des, err := fs.ReadDir(sfs, path)
		if err != nil {
			// Ignore ReadDir that fail when the path doesn't exist
			if errors.Is(err, fs.ErrNotExist) {
				continue
			}
			return nil, err
		}
		createOps, err := createOps(opt, sfs, path, des)
		if err != nil {
			return nil, err
		}
		ops = append(ops, createOps...)
	}
	return ops, nil
}

func deleteOps(opt *option, dir string, des []fs.DirEntry) (ops []Op, err error) {
	for _, de := range des {
		// Don't allow the directory itself to be deleted
		if de.Name() == "." {
			continue
		}
		path := filepath.Join(dir, de.Name())
		if opt.Skip(path, de.IsDir()) {
			continue
		}
		rel, err := opt.rel(path)
		if err != nil {
			return nil, err
		}
		ops = append(ops, Op{DeleteType, rel, nil})
		continue
	}
	return ops, nil
}

func updateOps(opt *option, sfs fs.FS, sdir string, tfs vfs.ReadWritable, tdir string, des []fs.DirEntry) (ops []Op, err error) {
	for _, de := range des {
		if de.Name() == "." {
			continue
		}
		spath := filepath.Join(sdir, de.Name())
		if opt.Skip(spath, de.IsDir()) {
			continue
		}
		tpath := filepath.Join(tdir, de.Name())
		// Recurse directories
		if de.IsDir() {
			childOps, err := diff(opt, sfs, spath, tfs, tpath)
			if err != nil {
				return nil, err
			}
			ops = append(ops, childOps...)
			continue
		}
		// Otherwise, check if the file has changed
		sourceStamp, err := stamp(sfs, spath)
		if err != nil {
			return nil, err
		}
		targetStamp, err := stamp(tfs, tpath)
		if err != nil {
			return nil, err
		}
		// Skip if the source and target are the same
		if sourceStamp == targetStamp {
			continue
		}
		data, err := fs.ReadFile(sfs, spath)
		if err != nil {
			// Don't error out on files that don't exist
			if errors.Is(err, fs.ErrNotExist) {
				continue
			}
			return nil, err
		}
		rel, err := opt.rel(spath)
		if err != nil {
			return nil, err
		}
		ops = append(ops, Op{UpdateType, rel, data})
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
			if err := tfs.WriteFile(op.Path, op.Data, 0644); err != nil {
				return err
			}
		case UpdateType:
			if err := tfs.WriteFile(op.Path, op.Data, 0644); err != nil {
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
	mode := stat.Mode()
	size := stat.Size()
	stamp = strconv.Itoa(int(size)) + ":" + mode.String() + ":" + strconv.Itoa(int(mtime))
	return stamp, nil
}
