package genfs

import (
	"io/fs"
	"testing"

	"github.com/livebud/bud/internal/is"
)

type treeGenerator struct{ label string }

func (g *treeGenerator) Generate(target string) (fs.File, error) {
	return nil, fs.ErrNotExist
}

func (g *treeGenerator) String() string {
	return g.label
}

var ag = &treeGenerator{"a"}
var bg = &treeGenerator{"b"}
var cg = &treeGenerator{"c"}
var eg = &treeGenerator{"e"}
var fg = &treeGenerator{"f"}

func TestTreeInsert(t *testing.T) {
	is := is.New(t)
	n := newTree()
	n.Insert("a", modeGen, ag)
	n.Insert("b", modeGen|modeDir, bg)
	n.Insert("b/c", modeGen|modeDir, cg)
	n.Insert("b/c/e", modeGen, eg)
	n.Insert("b/c/f", modeGen, fg)
	expect := `. mode=d-
├── a mode=-g generator=a
└── b mode=dg generator=b
    └── c mode=dg generator=c
        ├── e mode=-g generator=e
        └── f mode=-g generator=f
`
	is.Equal(n.Print(), expect)
}

func TestTreeFiller(t *testing.T) {
	is := is.New(t)
	n := newTree()
	n.Insert("a", modeGen, ag)
	n.Insert("b/c/e", modeGen, eg)
	n.Insert("b/c/f", modeGen, fg)
	n.Insert("b/c", modeGen|modeDir, cg)
	expect := `. mode=d-
├── a mode=-g generator=a
└── b mode=d-
    └── c mode=dg generator=c
        ├── e mode=-g generator=e
        └── f mode=-g generator=f
`
	is.Equal(n.Print(), expect)
}

func TestTreeFindPrefix(t *testing.T) {
	s := is.New(t)
	n := newTree()
	n.Insert("a", modeGen, ag)
	n.Insert("b", modeGen|modeDir, bg)
	n.Insert("b/c", modeGen|modeDir, cg)
	n.Insert("b/c/e", modeGen, eg)
	n.Insert("b/c/f", modeGen, fg)
	f, ok := n.FindPrefix("a")
	s.True(ok)
	s.Equal(f.Path, "a")
	// File generators must be an exact match
	f, ok = n.FindPrefix("a/d")
	s.Equal(ok, false)
	s.Equal(f, nil)
	f, ok = n.FindPrefix("b/c/h")
	s.True(ok)
	s.Equal(f.Path, "b/c")
	f, ok = n.FindPrefix("c")
	s.Equal(ok, true)
	s.Equal(f.Path, ".")
	f, ok = n.FindPrefix("c")
	s.Equal(ok, true)
	s.Equal(f.Path, ".")
	// Special case
	f, ok = n.FindPrefix(".")
	s.True(ok)
	s.Equal(f.Path, ".")
}

func TestTreeDelete(t *testing.T) {
	is := is.New(t)
	n := newTree()
	n.Insert("a", modeGen, ag)
	n.Insert("b", modeGen|modeDir, bg)
	n.Insert("b/c", modeGen|modeDir, cg)
	n.Insert("b/c/e", modeGen, eg)
	n.Insert("b/c/f", modeGen, fg)
	n.Delete("b/c")
	expect := `. mode=d-
├── a mode=-g generator=a
└── b mode=dg generator=b
`
	is.Equal(n.Print(), expect)
}

func TestTreeFillerDirNowGeneratorFile(t *testing.T) {
	is := is.New(t)
	n := newTree()
	n.Insert("bud/node_modules", modeGen|modeDir, ag)
	n.Insert("bud/node_modules/runtime/hot", modeGen, bg)
	node, ok := n.FindPrefix("bud/node_modules/runtime/svelte")
	is.True(ok)
	is.Equal(node.Path, "bud/node_modules/runtime")
	// Check that parent is a directory
	parent, ok := n.Find("bud/node_modules/runtime")
	is.True(ok)
	is.True(parent.Mode.IsDir())
	n.Insert("bud/node_modules/runtime", modeGen, cg)
	// Check that parent is a file
	parent, ok = n.Find("bud/node_modules/runtime")
	is.True(ok)
	is.Equal(parent.Mode, modeGen)
}

func TestTreeGeneratorAndDirectory(t *testing.T) {
	is := is.New(t)
	n := newTree()
	n.Insert("bud/node_modules", modeGen|modeDir, ag)
	n.Insert("bud/node_modules/runtime", modeGen, bg)
	n.Insert("bud/node_modules/runtime/hot", modeGen, cg)
	node, ok := n.FindPrefix("bud/node_modules/runtime")
	is.True(ok)
	is.Equal(node.Path, "bud/node_modules/runtime")
	node, ok = n.FindPrefix("bud/node_modules/runtime/hot")
	is.True(ok)
	is.Equal(node.Path, "bud/node_modules/runtime/hot")
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

// func TestTreeGenerate(t *testing.T) {
// 	is := is.New(t)
// 	n := newTree()
// 	n.Insert("bud/node_modules", modeGen|modeDir, &nodeModule{})
// 	err := fstest.TestFS(n, "bud/node_modules/runtime")
// 	is.NoErr(err)
// }

// func TestTreeRemove(t *testing.T) {
// 	is := is.New(t)
// 	n := newTree()
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

// func TestTreeClear(t *testing.T) {
// 	is := is.New(t)
// 	n := newTree()
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

// func TestTreeReadDir(t *testing.T) {
// 	is := is.New(t)
// 	n := newTree()
// 	count := 0
// 	n.Insert("controller/controller.go", modeGen, generate(func(target string) (fs.File, error) {
// 		count++
// 		return nil, fs.ErrNotExist
// 	}))
// 	node, ok := n.Find("controller")
// 	is.True(ok)
// 	is.True(node.Mode.IsDir())
// 	file, err := node.Generate("controller")
// 	is.NoErr(err)
// 	dir, ok := file.(fs.ReadDirFile)
// 	is.True(ok)
// 	des, err := dir.ReadDir(-1)
// 	is.NoErr(err)
// 	is.Equal(len(des), 1)
// 	is.Equal(des[0].Name(), "controller.go")
// 	is.Equal(count, 0)
// }
