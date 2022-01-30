package scan_test

import (
	"errors"
	"io/fs"
	"testing"

	"github.com/matryer/is"
	"gitlab.com/mnm/bud/internal/scan"

	"gitlab.com/mnm/bud/pkg/vfs"
)

func TestActionScan(t *testing.T) {
	is := is.New(t)
	fsys := vfs.Memory{
		"action/action.go":             &vfs.File{Data: []byte(``)},
		"action/users/users.go":        &vfs.File{Data: []byte(``)},
		"action/users/admin/admin.go":  &vfs.File{Data: []byte(``)},
		"action/posts/posts.go":        &vfs.File{Data: []byte(``)},
		"action/messages/_messages.go": &vfs.File{Data: []byte(``)},
		"action/about":                 &vfs.File{Mode: fs.ModeDir},
	}
	subfs, err := fs.Sub(fsys, "action")
	is.NoErr(err)
	scanner := scan.Actions(subfs)
	expect := [...]string{
		".",
		"posts",
		"users",
		"users/admin",
	}
	n := 0
	for scanner.Scan() {
		is.Equal(expect[n], scanner.Text())
		n++
	}
	is.NoErr(scanner.Err())
	is.Equal(n, 4)
}

func TestActionScanEmpty(t *testing.T) {
	is := is.New(t)
	fsys := vfs.Memory{}
	subfs, err := fs.Sub(fsys, "action")
	is.NoErr(err)
	scanner := scan.Actions(subfs)
	n := 0
	for scanner.Scan() {
		n++
	}
	is.Equal(n, 0)
	is.NoErr(scanner.Err())

}

type errFS struct {
}

func (*errFS) Open(name string) (fs.File, error) {
	return nil, fs.ErrInvalid
}

func TestActionScanError(t *testing.T) {
	is := is.New(t)
	fsys := &errFS{}
	scanner := scan.Actions(fsys)
	n := 0
	for scanner.Scan() {
		n++
	}
	is.Equal(n, 0)
	is.True(scanner.Err() != nil)
	is.True(errors.Is(scanner.Err(), fs.ErrInvalid))
}
