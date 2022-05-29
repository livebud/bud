package scan_test

import (
	"errors"
	"io/fs"
	"testing"

	"github.com/livebud/bud/internal/is"
	"github.com/livebud/bud/internal/scan"

	"github.com/livebud/bud/package/vfs"
)

func TestControllerScan(t *testing.T) {
	is := is.New(t)
	fsys := vfs.Memory{
		"controller/controller.go":         &vfs.File{Data: []byte(``)},
		"controller/users/users.go":        &vfs.File{Data: []byte(``)},
		"controller/users/admin/admin.go":  &vfs.File{Data: []byte(``)},
		"controller/posts/posts.go":        &vfs.File{Data: []byte(``)},
		"controller/messages/_messages.go": &vfs.File{Data: []byte(``)},
		"controller/about":                 &vfs.File{Mode: fs.ModeDir},
	}
	subfs, err := fs.Sub(fsys, "controller")
	is.NoErr(err)
	scanner := scan.Controllers(subfs)
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

func TestControllerScanEmpty(t *testing.T) {
	is := is.New(t)
	fsys := vfs.Memory{}
	subfs, err := fs.Sub(fsys, "controller")
	is.NoErr(err)
	scanner := scan.Controllers(subfs)
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

func TestControllerScanError(t *testing.T) {
	is := is.New(t)
	fsys := &errFS{}
	scanner := scan.Controllers(fsys)
	n := 0
	for scanner.Scan() {
		n++
	}
	is.Equal(n, 0)
	is.True(scanner.Err() != nil)
	is.True(errors.Is(scanner.Err(), fs.ErrInvalid))
}
