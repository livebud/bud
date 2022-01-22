package mergefs_test

import (
	"fmt"
	"testing"

	"gitlab.com/mnm/bud/2/mod"
	"gitlab.com/mnm/bud/internal/fstree"
	"gitlab.com/mnm/bud/internal/modcache"

	"github.com/matryer/is"
	"github.com/yalue/merged_fs"
	"gitlab.com/mnm/bud/vfs"
)

func TestPlugins(t *testing.T) {
	is := is.New(t)
	appDir := t.TempDir()
	err := vfs.Write(appDir, vfs.Map{
		"go.mod":                 []byte("module app.com\nrequire tailwind.com v1.0.0\nrequire markdown.com v1.0.0"),
		"transform/jade/jade.go": []byte(`package jade`),
		"package.json":           []byte(`{ "name": "app.com" }`),
	})
	is.NoErr(err)
	cacheDir := t.TempDir()
	modCache := modcache.New(cacheDir)
	err = modCache.Write(map[string]map[string]string{
		"tailwind.com@v1.0.0": map[string]string{
			"package.json":                   `{ "name": "tailwind.com" }`,
			"transform/tailwind/tailwind.go": `package tailwind`,
		},
		"markdown.com@v1.0.0": map[string]string{
			"package.json":                   `{ "name": "markdown.com" }`,
			"transform/markdown/markdown.go": `package markdown`,
		},
	})
	is.NoErr(err)
	module, err := mod.Find(appDir, mod.WithModCache(modCache))
	is.NoErr(err)
	tailwind, err := module.Find("tailwind.com")
	is.NoErr(err)
	markdown, err := module.Find("markdown.com")
	is.NoErr(err)
	fsys := merged_fs.NewMergedFS(merged_fs.NewMergedFS(module, tailwind), markdown)
	tree, err := fstree.Walk(fsys)
	is.NoErr(err)
	fmt.Println(tree.String())
	// pkg, err := fs.ReadFile(fsys, "package.json")
	// is.NoErr(err)
	// fmt.Println(string(pkg))
	// is.Equal(tree.String(), dedent.Dedent(`
	// 	.
	// 	├── go.mod
	// 	├── package.json
	// 	└── transform
	// 	    ├── jade
	// 	    │   └── jade.go
	// 	    ├── markdown
	// 	    │   └── markdown.go
	// 	    └── tailwind
	// 	        └── tailwind.go
	// `))
}
