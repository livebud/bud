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
	target := path.Join(subpaths...)
	if target == "" {
		target = "."
	}
	log = log.Field("fn", "sync")
	log.Debug("syncing")
	now := time.Now()

	ops, err := diff(log, from, to, target)
	if err != nil {
		return err
	}
	err = apply(log, to, ops)
	log.Field("duration", time.Since(now)).Debug("synced")
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
}

func (o syncOp) String() string {
	return o.Type.String() + ":" + o.Path
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
	log = log.Field("fn", "diff")
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
		log.Debug("op", op)
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
			ops = append(ops, syncOp{createType, fpath, data})
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
		ops = append(ops, syncOp{deleteType, fpath, nil})
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
				ops = append(ops, syncOp{deleteType, fpath, nil})
				continue
			}
			return nil, err
		}
		ops = append(ops, syncOp{updateType, fpath, data})
	}
	return ops, nil
}

func apply(log log.Log, to FS, ops []syncOp) error {
	log = log.Field("fn", "apply")
	for _, op := range ops {
		switch op.Type {
		case createType:
			log.Debug("creating", op.Path)
			dir := filepath.Dir(op.Path)
			if err := to.MkdirAll(dir, 0755); err != nil {
				return err
			}
			if err := to.WriteFile(op.Path, op.Data, 0644); err != nil {
				return err
			}
		case updateType:
			log.Debug("updating", op.Path)
			if err := to.WriteFile(op.Path, op.Data, 0644); err != nil {
				return err
			}
		case deleteType:
			log.Debug("removing", op.Path)
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
