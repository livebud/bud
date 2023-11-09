package esb_test

import (
	"path/filepath"
	"strings"
	"testing"

	esbuild "github.com/evanw/esbuild/pkg/api"
	"github.com/livebud/bud/internal/esb"
	"github.com/matryer/is"
)

func TestVirtualTopLevel(t *testing.T) {
	is := is.New(t)
	dir := t.TempDir()
	writeFiles(t, dir, map[string]string{
		"view/index.js": `
			export default "view/index.js"
		`,
	})
	loader := func(args esbuild.OnLoadArgs) (result esbuild.OnLoadResult, err error) {
		content := `
			import path from './index.js'
			export default path
		`
		result.Contents = &content
		result.ResolveDir = filepath.Dir(args.Path)
		return result, nil
	}
	file, err := serve(esbuild.BuildOptions{
		AbsWorkingDir: dir,
		EntryPoints:   []string{"./view/index.js"},
		Plugins: []esbuild.Plugin{
			esb.Virtual("./view/index.dom.js", loader),
		},
	})
	is.NoErr(err)
	code := string(file.Contents)
	is.True(strings.Contains(code, `var view_default = "view/index.js";`))
}

func TestVirtualInner(t *testing.T) {
	is := is.New(t)
	dir := t.TempDir()
	writeFiles(t, dir, map[string]string{
		"view/index.js": `
			import path from './index.dom.js'
			export default path
		`,
	})
	loader := func(args esbuild.OnLoadArgs) (result esbuild.OnLoadResult, err error) {
		content := `
			export default "view/index.dom.js"
		`
		result.Contents = &content
		result.ResolveDir = filepath.Dir(args.Path)
		return result, nil
	}
	file, err := serve(esbuild.BuildOptions{
		AbsWorkingDir: dir,
		EntryPoints:   []string{"./view/index.js"},
		Plugins: []esbuild.Plugin{
			esb.Virtual("./index.dom.js", loader),
		},
	})
	is.NoErr(err)
	code := string(file.Contents)
	is.True(strings.Contains(code, `var index_dom_default = "view/index.dom.js";`))
}

func TestVirtualRelativeLimitation(t *testing.T) {
	is := is.New(t)
	dir := t.TempDir()
	writeFiles(t, dir, map[string]string{
		"view/index.js": `
			import path from './index.dom.js'
			export default path
		`,
	})
	loader := func(args esbuild.OnLoadArgs) (result esbuild.OnLoadResult, err error) {
		content := `
			export default "view/index.dom.js"
		`
		result.Contents = &content
		result.ResolveDir = filepath.Dir(args.Path)
		return result, nil
	}
	file, err := serve(esbuild.BuildOptions{
		AbsWorkingDir: dir,
		EntryPoints:   []string{"./view/index.js"},
		Plugins: []esbuild.Plugin{
			esb.Virtual("./view/index.dom.js", loader),
		},
	})
	is.True(err != nil)
	is.True(strings.Contains(err.Error(), `Could not resolve "./index.dom.js"`))
	is.Equal(file, nil)
}
