package virtual

import (
	"errors"
	"io/fs"
	"path"
	"path/filepath"
	"sort"
	"strconv"
	"time"

	"github.com/livebud/bud/package/log"
)

// Sync files from one filesystem to another at subpath
func Sync(log log.Log, from fs.FS, to FS, subpaths ...string) error {
	log.Debug("virtual: syncing")
	target := path.Join(subpaths...)
	if target == "" {
		target = "."
	}
	now := time.Now()

	ops, err := diff(log, from, to, target)
	if err != nil {
		return err
	}
	err = apply(log, to, ops)
	log.Debug("virtual: synced in", time.Since(now))
	return err
}

type syncType uint8

const (
	createType syncType = iota + 1
	updateType
	deleteType
)

func (t syncType) String() string {
	switch t {
	case createType:
		return "create"
	case updateType:
		return "update"
	case deleteType:
		return "delete"
	default:
		return ""
	}
}

type syncOp struct {
	Type syncType
	Path string
	Data []byte
	Mode fs.FileMode
}

func (o syncOp) String() string {
	return o.Type.String() + " " + o.Path + " " + o.Mode.String()
}

func newSet(des []fs.DirEntry) set {
	s := make(set, len(des))
	for _, de := range des {
		s[de.Name()] = de
	}
	return s
}

type set map[string]fs.DirEntry

func (source set) Difference(target set) (des []fs.DirEntry) {
	for name, de := range source {
		if _, ok := target[name]; !ok {
			des = append(des, de)
		}
	}
	sort.Slice(des, func(i, j int) bool {
		return des[i].Name() < des[j].Name()
	})
	return des
}

func (source set) Intersection(target set) (des []fs.DirEntry) {
	for name, de := range source {
		if _, ok := target[name]; ok {
			des = append(des, de)
		}
	}
	sort.Slice(des, func(i, j int) bool {
		return des[i].Name() < des[j].Name()
	})
	return des
}

func diff(log log.Log, from fs.FS, to FS, dir string) (ops []syncOp, err error) {
	sourceEntries, err := fs.ReadDir(from, dir)
	if err != nil {
		return nil, err
	}
	targetEntries, err := fs.ReadDir(to, dir)
	if err != nil && !errors.Is(err, fs.ErrNotExist) {
		return nil, err
	}
	// Create the source set from the source entries
	sourceSet := newSet(sourceEntries)
	// Create a target set from the target entries
	targetSet := newSet(targetEntries)
	// Compute the operations
	creates := sourceSet.Difference(targetSet)
	deletes := targetSet.Difference(sourceSet)
	updates := sourceSet.Intersection(targetSet)
	createOps, err := createOps(log, from, dir, creates)
	if err != nil {
		return nil, err
	}
	deleteOps, err := deleteOps(log, dir, deletes)
	if err != nil {
		return nil, err
	}
	childOps, err := updateOps(log, from, to, dir, updates)
	if err != nil {
		return nil, err
	}
	ops = append(ops, createOps...)
	ops = append(ops, deleteOps...)
	ops = append(ops, childOps...)
	for _, op := range ops {
		log.Debug("virtual: op", op)
	}
	return ops, nil
}

func createOps(log log.Log, from fs.FS, dir string, des []fs.DirEntry) (ops []syncOp, err error) {
	for _, de := range des {
		if de.Name() == "." {
			continue
		}
		fpath := path.Join(dir, de.Name())
		if !de.IsDir() {
			data, err := fs.ReadFile(from, fpath)
			if err != nil {
				// Don't error out on files that don't exist
				if errors.Is(err, fs.ErrNotExist) {
					continue
				}
				return nil, err
			}
			ops = append(ops, syncOp{createType, fpath, data, de.Type()})
			continue
		}
		des, err := fs.ReadDir(from, fpath)
		if err != nil {
			// Ignore ReadDir that fail when the path doesn't exist
			if errors.Is(err, fs.ErrNotExist) {
				continue
			}
			return nil, err
		}
		createOps, err := createOps(log, from, fpath, des)
		if err != nil {
			return nil, err
		}
		ops = append(ops, createOps...)
	}
	return ops, nil
}

func deleteOps(log log.Log, dir string, des []fs.DirEntry) (ops []syncOp, err error) {
	for _, de := range des {
		// Don't allow the directory itself to be deleted
		if de.Name() == "." {
			continue
		}
		fpath := path.Join(dir, de.Name())
		ops = append(ops, syncOp{deleteType, fpath, nil, de.Type()})
		continue
	}
	return ops, nil
}

func updateOps(log log.Log, from fs.FS, to FS, dir string, des []fs.DirEntry) (ops []syncOp, err error) {
	for _, de := range des {
		if de.Name() == "." {
			continue
		}
		fpath := path.Join(dir, de.Name())
		// Recurse directories
		if de.IsDir() {
			childOps, err := diff(log, from, to, fpath)
			if err != nil {
				return nil, err
			}
			ops = append(ops, childOps...)
			continue
		}
		// Otherwise, check if the file has changed
		sourceStamp, err := stamp(from, fpath)
		if err != nil {
			return nil, err
		}
		targetStamp, err := stamp(to, fpath)
		if err != nil {
			return nil, err
		}
		// Skip if the source and target are the same
		if sourceStamp == targetStamp {
			continue
		}
		data, err := fs.ReadFile(from, fpath)
		if err != nil {
			// Don't error out on files that don't exist
			if errors.Is(err, fs.ErrNotExist) {
				// The file no longer exists, delete it
				ops = append(ops, syncOp{deleteType, fpath, nil, de.Type()})
				continue
			}
			return nil, err
		}
		ops = append(ops, syncOp{updateType, fpath, data, de.Type()})
	}
	return ops, nil
}

func apply(log log.Log, to FS, ops []syncOp) error {
	for _, op := range ops {
		switch op.Type {
		case createType:
			dir := filepath.Dir(op.Path)
			// TODO: create ops for new directories too and maintain original
			// permission bits.
			mode := fs.FileMode(0755 | fs.ModeDir)
			log.Debug("virtual: creating", dir, mode)
			if err := to.MkdirAll(dir, mode); err != nil {
				return err
			}
			// Many of the virtual filesystems don't set a mode. Copying these to an
			// actual filesystem will cause permission errors, so we'll use common
			// permissions when not explicitly set.
			if op.Mode == 0 {
				op.Mode = 0644
			}
			log.Debug("virtual: creating", op.Path, op.Mode)
			if err := to.WriteFile(op.Path, op.Data, op.Mode); err != nil {
				return err
			}
		case updateType:
			// Many of the virtual filesystems don't set a mode. Copying these to an
			// actual filesystem will cause permission errors, so we'll use common
			// permissions when not explicitly set.
			if op.Mode == 0 {
				op.Mode = 0644
			}
			log.Debug("virtual: updating", op.Path, op.Mode)
			if err := to.WriteFile(op.Path, op.Data, op.Mode); err != nil {
				return err
			}
		case deleteType:
			log.Debug("virtual: removing", op.Path)
			if err := to.RemoveAll(op.Path); err != nil {
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
