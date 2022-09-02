package treefs_test

import (
	"io/fs"
	"testing"

	"github.com/livebud/bud/internal/is"
	"github.com/livebud/bud/package/budfs/treefs"
)

type generator struct{ label string }

func (g *generator) Generate(target string) (fs.File, error) {
	return nil, fs.ErrNotExist
}

func (g *generator) String() string {
	return g.label
}

var ag = &generator{"a"}
var bg = &generator{"b"}
var cg = &generator{"c"}
var eg = &generator{"e"}
var fg = &generator{"f"}

func TestInsert(t *testing.T) {
	is := is.New(t)
	n := treefs.New(".")
	n.InsertFile("a", ag)
	bn := n.InsertDir("b", bg)
	cn := bn.InsertDir("c", cg)
	cn.InsertFile("e", eg)
	cn.InsertFile("f", fg)
	expect := `. mode=d---------
├── a generator=a mode=----------
└── b generator=b mode=d---------
    └── c generator=c mode=d---------
        ├── e generator=e mode=----------
        └── f generator=f mode=----------
`
	is.Equal(n.Print(), expect)
}

func TestFiller(t *testing.T) {
	is := is.New(t)
	n := treefs.New(".")
	n.InsertFile("a", ag)
	n.InsertFile("b/c/e", eg)
	n.InsertFile("b/c/f", fg)
	n.InsertDir("b/c", cg)
	expect := `. mode=d---------
├── a generator=a mode=----------
└── b mode=d---------
    └── c generator=c mode=d---------
        ├── e generator=e mode=----------
        └── f generator=f mode=----------
`
	is.Equal(n.Print(), expect)
}

// func TestFind(t *testing.T) {
// 	is := is.New(t)
// 	n := treefs.New(".")
// 	n.InsertFile("a", ag)
// 	bn := n.InsertDir("b", bg)
// 	cn := bn.InsertDir("c", cg)
// 	cn.InsertFile("e", eg)
// 	cn.InsertFile("f", fg)
// 	f, ok := n.Find("a")
// 	is.True(ok)
// 	is.Equal(f.Path(), "a")
// 	f, ok = n.Find("a/b")
// 	is.Equal(ok, false)
// 	is.Equal(f, nil)
// 	f, ok = n.Find("b/c")
// 	is.True(ok)
// 	is.Equal(f.Path(), "b/c")
// 	f, ok = n.Find("b/c/e")
// 	is.True(ok)
// 	is.Equal(f.Path(), "b/c/e")
// 	f, ok = n.Find("b/c/f")
// 	is.True(ok)
// 	is.Equal(f.Path(), "b/c/f")
// 	f, ok = n.Find("b/c/e/f")
// 	is.Equal(ok, false)
// 	is.Equal(f, nil)
// 	// Special case
// 	f, ok = n.Find(".")
// 	is.True(ok)
// 	is.Equal(f.Path(), ".")
// }

func TestFindPrefix(t *testing.T) {
	s := is.New(t)
	n := treefs.New(".")
	n.InsertFile("a", ag)
	bn := n.InsertDir("b", bg)
	cn := bn.InsertDir("c", cg)
	cn.InsertFile("e", eg)
	cn.InsertFile("f", fg)
	f, prefix, ok := n.FindByPrefix("a")
	s.True(ok)
	s.Equal(prefix, "a")
	s.Equal(f.Path(), "a")
	// File generators must be an exact match
	f, prefix, ok = n.FindByPrefix("a/d")
	s.Equal(ok, false)
	s.Equal(f, nil)
	s.Equal(prefix, "")
	f, prefix, ok = n.FindByPrefix("b/c/h")
	s.True(ok)
	s.Equal(prefix, "b/c")
	s.Equal(f.Path(), "b/c")
	f, prefix, ok = n.FindByPrefix("c")
	s.Equal(ok, false)
	s.Equal(f, nil)
	s.Equal(prefix, "")
	// Special case
	f, prefix, ok = n.FindByPrefix(".")
	s.True(ok)
	s.Equal(prefix, ".")
	s.Equal(f.Path(), ".")
}

func TestDelete(t *testing.T) {
	is := is.New(t)
	n := treefs.New(".")
	n.InsertFile("a", ag)
	bn := n.InsertDir("b", bg)
	cn := bn.InsertDir("c", cg)
	cn.InsertFile("e", eg)
	cn.InsertFile("f", fg)
	cn, ok := n.Delete("b", "c")
	is.True(ok)
	is.Equal(cn.Path(), "b/c")
	expect := `. mode=d---------
├── a generator=a mode=----------
└── b generator=b mode=d---------
    ├── e generator=e mode=----------
    └── f generator=f mode=----------
`
	is.Equal(n.Print(), expect)
}

func TestFillerDirNowGeneratorFile(t *testing.T) {
	is := is.New(t)
	n := treefs.New(".")
	n.InsertDir("bud/node_modules", ag)
	n.InsertFile("bud/node_modules/runtime/hot", bg)
	node, prefix, ok := n.FindByPrefix("bud/node_modules/runtime/svelte")
	is.True(ok)
	is.Equal(node.Path(), "bud/node_modules")
	is.Equal(prefix, "bud/node_modules")
	// Check that parent is a directory
	parent, ok := n.Find("bud/node_modules/runtime")
	is.True(ok)
	is.True(parent.Mode().IsDir())
	n.InsertFile("bud/node_modules/runtime", cg)
	// Check that parent is a file
	parent, ok = n.Find("bud/node_modules/runtime")
	is.True(ok)
	is.True(parent.Mode().IsRegular())
}

func TestGeneratorAndDirectory(t *testing.T) {
	is := is.New(t)
	n := treefs.New(".")
	n.InsertDir("bud/node_modules", ag)
	n.InsertFile("bud/node_modules/runtime", bg)
	n.InsertFile("bud/node_modules/runtime/hot", cg)
	node, prefix, ok := n.FindByPrefix("bud/node_modules/runtime")
	is.True(ok)
	is.Equal(node.Path(), "bud/node_modules/runtime")
	is.Equal(prefix, "bud/node_modules/runtime")
	node, prefix, ok = n.FindByPrefix("bud/node_modules/runtime/hot")
	is.True(ok)
	is.Equal(node.Path(), "bud/node_modules/runtime/hot")
	is.Equal(prefix, "bud/node_modules/runtime/hot")
}

func TestPrefixDifferentFromPath(t *testing.T) {
	is := is.New(t)
	n := treefs.New(".")
	n.InsertDir("bud/node_modules", ag)
	n.InsertFile("bud/node_modules/runtime/hot", bg)
	child, ok := n.Find("bud/node_modules")
	is.True(ok)
	node, prefix, ok := child.FindByPrefix("runtime/hot")
	is.True(ok)
	is.Equal(prefix, "runtime/hot")
	is.Equal(node.Path(), "bud/node_modules/runtime/hot")
}
