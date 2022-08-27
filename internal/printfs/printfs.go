package printfs

import (
	"io/fs"
	"path/filepath"
	"strings"

	"github.com/xlab/treeprint"
)

func New() *Tree {
	return &Tree{
		tree: treeprint.New(),
	}
}

type Tree struct {
	tree treeprint.Tree
}

func (t *Tree) Add(path string) {
	parent := t.tree
	for _, element := range strings.Split(filepath.ToSlash(path), "/") {
		existing := parent.FindByValue(element)
		if existing != nil {
			parent = existing
		} else {
			parent = parent.AddBranch(element)
		}
	}
}

func (t *Tree) String() string {
	return t.tree.String()
}

func Walk(fsys fs.FS) (*Tree, error) {
	tree := New()
	err := fs.WalkDir(fsys, ".", func(path string, de fs.DirEntry, err error) error {
		if err != nil {
			return err
		} else if path == "." {
			return nil
		}
		tree.Add(path)
		return nil
	})
	if err != nil {
		return nil, err
	}
	return tree, nil
}

// func (t *Tree) String2() string {
// 	// Handle special cases
// 	switch len(t.entries) {
// 	case 0:
// 		return ""
// 	case 1:
// 		return strings.Join(t.entries[0].parts, sep)
// 	}
// 	// Sort the entries
// 	sort.Slice(t.entries, func(i, j int) bool {
// 		il := len(t.entries[i].parts)
// 		ij := len(t.entries[j].parts)
// 		lesser := il
// 		if ij < il {
// 			lesser = ij
// 		}
// 		for k := 0; k < lesser; k++ {
// 			if t.entries[i].parts[k] < t.entries[j].parts[k] {
// 				return true
// 			}
// 			if t.entries[i].parts[k] > t.entries[j].parts[k] {
// 				return false
// 			}
// 		}
// 		return il < ij
// 	})

// 	fmt.Println(0, t.entries[0].path)
// 	prevPath := t.entries[0].path
// 	depth := 0
// 	for i := 1; i < len(t.entries); i++ {

// 		currPath := t.entries[1].path
// 		if strings.HasPrefix(currPath, prevPath) {
// 			depth++
// 		} else if !strings.HasPrefix(currPath, filepath.Dir(prevPath)) {
// 			depth--
// 		}
// 		// fmt.Println(prev, curr, depth, t.entries[i].parts)
// 		fmt.Println(depth, t.entries[i].path)
// 		prevPath = currPath
// 		// prevLen = currLen
// 	}

// 	// out := new(strings.Builder)
// 	// prevDepth := len(t.entries[0].parts)
// 	// depth := 0
// 	// for i := 0; i < len(t.entries); i++ {
// 	// 	if i > 0 {
// 	// 		out.WriteByte('\n')
// 	// 	}
// 	// 	entry := t.entries[i]
// 	// 	currDepth := len(entry.parts)
// 	// 	if prevDepth < currDepth {
// 	// 		depth++

// 	// 	} else if prevDepth > currDepth {
// 	// 		depth--
// 	// 	}
// 	// 	// Write the indent
// 	// 	out.WriteString(indent(depth))
// 	// 	if depth >= currDepth {
// 	// 		out.WriteString(entry.parts[currDepth-1])
// 	// 	} else {
// 	// 		out.WriteString(strings.Join(entry.parts[depth:], sep))
// 	// 	}
// 	// 	// Write the /
// 	// 	if entry.isDir {
// 	// 		out.WriteString(sep)
// 	// 	}
// 	// 	prevDepth = currDepth
// 	// }
// 	return ""
// 	// return out.String()
// }

// type entry struct {
// 	path  string
// 	parts []string
// 	isDir bool
// }

// func indent(n int) string {
// 	out := ""
// 	for i := 0; i < n; i++ {
// 		out += "  "
// 	}
// 	return out
// }
