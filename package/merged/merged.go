package merged

import (
	"errors"
	"fmt"
	"io"
	"sort"
	"time"

	"io/fs"

	"github.com/livebud/bud/internal/errs"
)

// Merge the filesystems together
// TODO: give the filesystems names
func Merge(fileSystems ...fs.FS) *FS {
	return &FS{fileSystems}
}

type FS struct {
	fileSystems []fs.FS
}

// Open finds the first path in fileSystems
func (f *FS) Open(path string) (fs.File, error) {
	if !fs.ValidPath(path) {
		return nil, &fs.PathError{
			Op:   "open",
			Path: path,
			Err:  fs.ErrInvalid,
		}
	}
	var dirs []dir
	notExists := &notExists{path: path}
	for _, fsys := range f.fileSystems {
		file, err := fsys.Open(path)
		if err != nil {
			if errors.Is(err, fs.ErrNotExist) {
				notExists.errors = append(notExists.errors, err)
				continue
			}
			// Fail fast if it's anything other than a not found.
			return nil, err
		}
		stat, err := file.Stat()
		if err != nil {
			return nil, err
		}
		// Files always have priority. If it's a file, return right away.
		if !stat.IsDir() {
			return file, nil
		}
		dirs = append(dirs, dir{file, stat})
	}
	// If we didn't find any files within the filesystem,
	// return the missing file errors.
	if len(dirs) == 0 {
		return nil, notExists
	}
	return f.mergeDir(path, dirs)
}

// notExistsError is a collection of all errors while attempting to open a file
// in one of the filesystems
type notExists struct {
	path   string
	errors []error
}

// Error implements error and joins all the errors together
func (n *notExists) Error() string {
	return fmt.Errorf("merged: open %q. %w", n.path, errs.Join(n.errors...)).Error()
}

// This type of error should be an fs.ErrNotExist because all underlying errors
// are also fs.ErrNotExists
func (*notExists) Unwrap() error {
	return fs.ErrNotExist
}

type dir struct {
	fs.File
	fs.FileInfo
}

// Creates and returns a new pseudo-directory "File" that contains the contents
// of both files a and b. Both a and b must be directories at the same
// specified path in m.A and m.B, respectively. Closes files a and b before
// returning, since they aren't needed by the MergedDirectory pseudo-file.
func (f *FS) mergeDir(path string, dirs []dir) (fs.File, error) {
	// Initialize merged
	merged := &mergeDir{
		name:    baseName(path),
		mode:    dirs[0].Mode(),    // use the first directory's mode
		modTime: dirs[0].ModTime(), // use the first directory's modtime
		size:    dirs[0].Size(),    // use the first directory's size
		sys:     dirs[0].Sys(),     // use the first directory's sys
	}
	// Ignore conflicting names
	conflicts := map[string]bool{}
	// Loop over the directory
	for _, dir := range dirs {
		defer dir.Close()
		d, ok := dir.File.(fs.ReadDirFile)
		if !ok {
			return nil, fmt.Errorf("merged: directories doesn't implement ReadDirFile")
		}
		// Read all the dir entries
		des, err := d.ReadDir(-1)
		if err != nil {
			return nil, err
		}
		// Add the entries from directory A first.
		for _, de := range des {
			name := de.Name()
			// Conflicts are ignored. File systems are in prioritized order
			if conflicts[name] {
				continue
			}
			merged.entries = append(merged.entries, de)
			conflicts[name] = true
		}
	}
	// Sort all the entries in alphabetical order
	sort.Slice(merged.entries, func(i, j int) bool {
		return merged.entries[i].Name() < merged.entries[j].Name()
	})
	return merged, nil
}

// This is the key component of this library. It represents a directory that is
// present in both filesystems. Implements the fs.File, fs.DirEntry, and
// fs.FileInfo interfaces.
type mergeDir struct {
	// The path to this directory in both FSs
	name    string
	mode    fs.FileMode
	modTime time.Time
	size    int64
	sys     interface{}
	entries []fs.DirEntry
	// The next entry to return with ReadDir.
	readOffset int
}

func (d *mergeDir) Name() string {
	return d.name
}

func (d *mergeDir) Mode() fs.FileMode {
	return d.mode
}

func (d *mergeDir) ModTime() time.Time {
	return d.modTime
}

func (d *mergeDir) IsDir() bool {
	return true
}

func (d *mergeDir) Sys() interface{} {
	return d.sys
}

func (d *mergeDir) Stat() (fs.FileInfo, error) {
	return d, nil
}

func (d *mergeDir) Info() (fs.FileInfo, error) {
	return d, nil
}

func (d *mergeDir) Type() fs.FileMode {
	return d.mode.Type()
}

func (d *mergeDir) Size() int64 {
	return d.size
}

func (d *mergeDir) Read(data []byte) (int, error) {
	return 0, fmt.Errorf("%s is a directory", d.name)
}

func (d *mergeDir) Close() error {
	// Note: Do *not* clear the rest of the fields here, since the
	// mergeDir also serves as a DirEntry or FileInfo, which must be
	// able to outlive the File itself being closed.
	d.entries = nil
	d.readOffset = 0
	return nil
}

func (d *mergeDir) ReadDir(n int) ([]fs.DirEntry, error) {
	if d.readOffset >= len(d.entries) {
		if n <= 0 {
			// A special case required by the FS interface.
			return nil, nil
		}
		return nil, io.EOF
	}
	startEntry := d.readOffset
	var endEntry int
	if n <= 0 {
		endEntry = len(d.entries)
	} else {
		endEntry = startEntry + n
	}
	if endEntry > len(d.entries) {
		endEntry = len(d.entries)
	}
	entries := d.entries[startEntry:endEntry]
	d.readOffset = endEntry
	return entries, nil
}

// Returns the final element of the path. The path must be valid according to
// the rules of fs.ValidPath.
func baseName(path string) string {
	d := []byte(path)
	i := len(d)
	if i <= 1 {
		return path
	}
	i--
	for i >= 0 {
		if d[i] == '/' {
			break
		}
		i--
	}
	return string(d[i+1:])
}
