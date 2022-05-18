package goimport_test

import (
	"testing"
	"testing/fstest"

	"github.com/livebud/bud/internal/goimport"

	"github.com/matryer/is"
)

func TestDir(t *testing.T) {
	is := is.New(t)
	fsys := fstest.MapFS{
		"a1.go": &fstest.MapFile{Data: []byte(`package a`)},
		"a2.go": &fstest.MapFile{Data: []byte(`package a`)},
		// Ignore tests
		"a_test.go": &fstest.MapFile{Data: []byte(`package a_test`)},
		// Ignore different tags
		"a3.go": &fstest.MapFile{Data: []byte(`
			//go:build ignore
			package main
			func main() {}
		`)},
	}
	importer := goimport.New(fsys)
	pkg, err := importer.Import(".")
	is.NoErr(err)
	is.Equal(len(pkg.GoFiles), 2)
	is.Equal(pkg.GoFiles[0], `a1.go`)
	is.Equal(pkg.GoFiles[1], `a2.go`)
}
