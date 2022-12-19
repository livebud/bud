package treefs_test

import (
	"io/fs"
	"testing"

	"github.com/livebud/bud/internal/is"
	"github.com/livebud/bud/internal/treefs"
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

const (
	modeGenerator    = treefs.ModeGenerator
	modeGeneratorDir = treefs.ModeDir | treefs.ModeGenerator
)

func TestInsert(t *testing.T) {
	is := is.New(t)
	n := treefs.New(".")
	n.Insert("a", modeGenerator, ag)
	n.Insert("b", modeGeneratorDir, bg)
	n.Insert("b/c", modeGeneratorDir, cg)
	n.Insert("b/c/e", modeGenerator, eg)
	n.Insert("b/c/f", modeGenerator, fg)
	expect := `. mode=d-
├── a mode=-g generator=a
└── b mode=dg generator=b
    └── c mode=dg generator=c
        ├── e mode=-g generator=e
        └── f mode=-g generator=f
`
	is.Equal(n.Print(), expect)
}

func TestFiller(t *testing.T) {
	is := is.New(t)
	n := treefs.New(".")
	n.Insert("a", modeGenerator, ag)
	n.Insert("b/c/e", modeGenerator, eg)
	n.Insert("b/c/f", modeGenerator, fg)
	n.Insert("b/c", modeGeneratorDir, cg)
	expect := `. mode=d-
├── a mode=-g generator=a
└── b mode=d-
    └── c mode=dg generator=c
        ├── e mode=-g generator=e
        └── f mode=-g generator=f
`
	is.Equal(n.Print(), expect)
}

func TestFindPrefix(t *testing.T) {
	s := is.New(t)
	n := treefs.New(".")
	n.Insert("a", modeGenerator, ag)
	n.Insert("b", modeGeneratorDir, bg)
	n.Insert("b/c", modeGeneratorDir, cg)
	n.Insert("b/c/e", modeGenerator, eg)
	n.Insert("b/c/f", modeGenerator, fg)
	f, ok := n.FindByPrefix("a")
	s.True(ok)
	s.Equal(f.Path(), "a")
	// File generators must be an exact match
	f, ok = n.FindByPrefix("a/d")
	s.Equal(ok, false)
	s.Equal(f, nil)
	f, ok = n.FindByPrefix("b/c/h")
	s.True(ok)
	s.Equal(f.Path(), "b/c")
	f, ok = n.FindByPrefix("c")
	s.Equal(ok, true)
	s.Equal(f.Path(), ".")
	f, ok = n.FindByPrefix("c")
	s.Equal(ok, true)
	s.Equal(f.Path(), ".")
	// Special case
	f, ok = n.FindByPrefix(".")
	s.True(ok)
	s.Equal(f.Path(), ".")
}

func TestDelete(t *testing.T) {
	is := is.New(t)
	n := treefs.New(".")
	n.Insert("a", modeGenerator, ag)
	n.Insert("b", modeGeneratorDir, bg)
	n.Insert("b/c", modeGeneratorDir, cg)
	n.Insert("b/c/e", modeGenerator, eg)
	n.Insert("b/c/f", modeGenerator, fg)
	n.Delete("b/c")
	expect := `. mode=d-
├── a mode=-g generator=a
└── b mode=dg generator=b
`
	is.Equal(n.Print(), expect)
}

func TestFillerDirNowGeneratorFile(t *testing.T) {
	is := is.New(t)
	n := treefs.New(".")
	n.Insert("bud/node_modules", modeGeneratorDir, ag)
	n.Insert("bud/node_modules/runtime/hot", modeGenerator, bg)
	node, ok := n.FindByPrefix("bud/node_modules/runtime/svelte")
	is.True(ok)
	is.Equal(node.Path(), "bud/node_modules/runtime")
	// Check that parent is a directory
	parent, ok := n.Get("bud/node_modules/runtime")
	is.True(ok)
	is.True(parent.Mode().IsDir())
	n.Insert("bud/node_modules/runtime", modeGenerator, cg)
	// Check that parent is a file
	parent, ok = n.Get("bud/node_modules/runtime")
	is.True(ok)
	is.True(parent.Mode().IsFile())
}

func TestGeneratorAndDirectory(t *testing.T) {
	is := is.New(t)
	n := treefs.New(".")
	n.Insert("bud/node_modules", modeGeneratorDir, ag)
	n.Insert("bud/node_modules/runtime", modeGenerator, bg)
	n.Insert("bud/node_modules/runtime/hot", modeGenerator, cg)
	node, ok := n.FindByPrefix("bud/node_modules/runtime")
	is.True(ok)
	is.Equal(node.Path(), "bud/node_modules/runtime")
	node, ok = n.FindByPrefix("bud/node_modules/runtime/hot")
	is.True(ok)
	is.Equal(node.Path(), "bud/node_modules/runtime/hot")
}

// type nodeModule struct {
// }

// func (n *nodeModule) Generate(target string) (fs.File, error) {
// 	switch target {
// 	case "bud/node_modules/runtime":
// 		return virtual.New(&virtual.Dir{
// 			Path: "bud/node_modules/runtime",
// 			Mode: 0755 | treefs.ModeDir,
// 		}), nil
// 	case "bud/node_modules":
// 		return virtual.New(&virtual.Dir{
// 			Path: "bud/node_modules",
// 			Mode: 0755 | treefs.ModeDir,
// 			Entries: []fs.DirEntry{
// 				&virtual.DirEntry{
// 					Path: "bud/node_modules/runtime",
// 					Mode: 0755 | treefs.ModeDir,
// 				},
// 			},
// 		}), nil
// 	default:
// 		return nil, fmt.Errorf("error generating %q. %w", target, fs.ErrNotExist)
// 	}
// }

// func TestGenerate(t *testing.T) {
// 	is := is.New(t)
// 	n := treefs.New(".")
// 	n.Insert("bud/node_modules", modeGeneratorDir, &nodeModule{})
// 	err := fstest.TestFS(n, "bud/node_modules/runtime")
// 	is.NoErr(err)
// }

// func TestRemove(t *testing.T) {
// 	is := is.New(t)
// 	n := treefs.New(".")
// 	n.FileGenerator("a", ag)
// 	bn := n.DirGenerator("b", bg)
// 	cn := bn.DirGenerator("c", cg)
// 	cn.FileGenerator("e", eg)
// 	cn.FileGenerator("f", fg)
// 	expect := `. mode=d-
// ├── a generator=a mode=-g
// └── b generator=b mode=d-
//     └── c generator=c mode=d-
//         ├── e generator=e mode=-g
//         └── f generator=f mode=-g
// `
// 	is.Equal(n.Print(), expect)
// 	n.Remove("b")
// 	expect = `. mode=d-
// └── a generator=a mode=-g
// `
// 	is.Equal(n.Print(), expect)
// }

// func TestClear(t *testing.T) {
// 	is := is.New(t)
// 	n := treefs.New(".")
// 	n.FileGenerator("a", ag)
// 	bn := n.DirGenerator("b", bg)
// 	cn := bn.DirGenerator("c", cg)
// 	cn.FileGenerator("e", eg)
// 	cn.FileGenerator("f", fg)
// 	expect := `. mode=d-
// ├── a generator=a mode=-g
// └── b generator=b mode=d-
//     └── c generator=c mode=d-
//         ├── e generator=e mode=-g
//         └── f generator=f mode=-g
// `
// 	is.Equal(n.Print(), expect)
// 	cn.Clear()
// 	expect = `. mode=d-
// ├── a generator=a mode=-g
// └── b generator=b mode=d-
//     └── c generator=c mode=d-
// `
// 	is.Equal(n.Print(), expect)
// }

func TestReadDir(t *testing.T) {
	is := is.New(t)
	n := treefs.New(".")
	count := 0
	n.Insert("controller/controller.go", modeGenerator, treefs.Generate(func(target string) (fs.File, error) {
		count++
		return nil, fs.ErrNotExist
	}))
	node, ok := n.Get("controller")
	is.True(ok)
	is.True(node.Mode().IsDir())
	file, err := node.Generate("controller")
	is.NoErr(err)
	dir, ok := file.(fs.ReadDirFile)
	is.True(ok)
	des, err := dir.ReadDir(-1)
	is.NoErr(err)
	is.Equal(len(des), 1)
	is.Equal(des[0].Name(), "controller.go")
	is.Equal(count, 0)
}
