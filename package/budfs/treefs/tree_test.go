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
	n.Insert("a", 0, ag)
	bn := n.Insert("b", fs.ModeDir, bg)
	cn := bn.Insert("c", fs.ModeDir, cg)
	cn.Insert("e", 0, eg)
	cn.Insert("f", 0, fg)
	expect := `. generator=<nil> mode=d---------
├── a generator=a mode=----------
└── b generator=b mode=d---------
    └── c generator=c mode=d---------
        ├── e generator=e mode=----------
        └── f generator=f mode=----------
`
	is.Equal(n.Print(), expect)
}

func TestFind(t *testing.T) {
	is := is.New(t)
	n := treefs.New(".")
	n.Insert("a", 0, ag)
	bn := n.Insert("b", fs.ModeDir, bg)
	cn := bn.Insert("c", fs.ModeDir, cg)
	cn.Insert("e", 0, eg)
	cn.Insert("f", 0, fg)
	f, ok := n.Find("a")
	is.True(ok)
	is.Equal(f.Name, "a")
	is.Equal(f.Generator, ag)
	f, ok = n.Find("a/b")
	is.Equal(ok, false)
	is.Equal(f, nil)
	f, ok = n.Find("b/c")
	is.True(ok)
	is.Equal(f.Name, "c")
	is.Equal(f.Generator, cg)
	f, ok = n.Find("b/c/e")
	is.True(ok)
	is.Equal(f.Name, "e")
	is.Equal(f.Generator, eg)
	f, ok = n.Find("b/c/f")
	is.True(ok)
	is.Equal(f.Name, "f")
	is.Equal(f.Generator, fg)
	f, ok = n.Find("b/c/e/f")
	is.Equal(ok, false)
	is.Equal(f, nil)
	// Special case
	f, ok = n.Find(".")
	is.True(ok)
	is.Equal(f.Name, ".")
	is.Equal(f.Generator, nil)
	is.Equal(f.Path(), ".")
}

func TestFindPrefix(t *testing.T) {
	s := is.New(t)
	n := treefs.New(".")
	n.Insert("a", 0, ag)
	bn := n.Insert("b", fs.ModeDir, bg)
	cn := bn.Insert("c", fs.ModeDir, cg)
	cn.Insert("e", 0, eg)
	cn.Insert("f", 0, fg)
	f, prefix, ok := n.FindByPrefix("a")
	s.True(ok)
	s.Equal(f.Name, "a")
	s.Equal(f.Generator, ag)
	s.Equal(prefix, "a")
	s.Equal(f.Path(), "a")
	f, prefix, ok = n.FindByPrefix("a/d")
	s.True(ok)
	s.Equal(f.Name, "a")
	s.Equal(f.Generator, ag)
	s.Equal(prefix, "a")
	s.Equal(f.Path(), "a")
	f, prefix, ok = n.FindByPrefix("b/c/h")
	s.True(ok)
	s.Equal(f.Name, "c")
	s.Equal(f.Generator, cg)
	s.Equal(prefix, "b/c")
	s.Equal(f.Path(), "b/c")
	f, prefix, ok = n.FindByPrefix("c")
	s.Equal(ok, false)
	s.Equal(prefix, "")
	s.Equal(f, nil)
	s.Equal(f.Path(), "")
	// Special case
	f, prefix, ok = n.FindByPrefix(".")
	s.True(ok)
	s.Equal(f.Name, ".")
	s.Equal(f.Generator, nil)
	s.Equal(f.Path(), ".")
}

func TestDelete(t *testing.T) {
	is := is.New(t)
	n := treefs.New(".")
	n.Insert("a", 0, ag)
	bn := n.Insert("b", fs.ModeDir, bg)
	cn := bn.Insert("c", fs.ModeDir, cg)
	cn.Insert("e", 0, eg)
	cn.Insert("f", 0, fg)
	cn, ok := n.Delete("b", "c")
	is.True(ok)
	is.Equal(cn.Name, "c")
	expect := `. generator=<nil> mode=d---------
├── a generator=a mode=----------
└── b generator=b mode=d---------
    ├── e generator=e mode=----------
    └── f generator=f mode=----------
`
	is.Equal(n.Print(), expect)
}
