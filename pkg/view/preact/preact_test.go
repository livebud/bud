package preact_test

import (
	"context"
	"strings"
	"testing"

	"github.com/dop251/goja"
	"github.com/livebud/bud/internal/js"
	"github.com/livebud/bud/internal/npm"
	"github.com/livebud/bud/internal/testdir"
	"github.com/livebud/bud/pkg/mod"
	"github.com/livebud/bud/pkg/view/preact"
	"github.com/matryer/is"
)

func TestCompileSSR(t *testing.T) {
	is := is.New(t)
	dir := t.TempDir()
	ctx := context.Background()
	is.NoErr(npm.Install(ctx, dir, "preact@10.18.1", "preact-render-to-string@6.2.2"))
	is.NoErr(testdir.WriteFiles(dir, map[string]string{
		"index.jsx": `
			import { createElement } from "preact";
			export default function() {
				return <h1>Hello World!</h1>
			}
		`,
	}))
	module, err := mod.Find(dir)
	is.NoErr(err)
	preact := preact.New(module)
	entry, err := preact.CompileSSR("./index.jsx")
	is.NoErr(err)
	program, err := goja.Compile("index.jsx", string(entry.Contents), false)
	is.NoErr(err)
	vm := js.New()
	_, err = vm.RunProgram(program)
	is.NoErr(err)
	result, err := vm.RunString("bud.render({})")
	is.NoErr(err)
	is.Equal(result.String(), `{"html":"<h1>Hello World!</h1>","head":""}`)
}

func TestCompileSSRWithHead(t *testing.T) {
	is := is.New(t)
	dir := t.TempDir()
	ctx := context.Background()
	is.NoErr(npm.Install(ctx, dir, "preact@10.18.1", "preact-render-to-string@6.2.2"))
	is.NoErr(testdir.WriteFiles(dir, map[string]string{
		"index.jsx": `
			import { createElement } from "preact";
			export default function(props, context) {
				context.head.push(<title>hello</title>)
				return <h1>Hello World!</h1>
			}
		`,
	}))
	module, err := mod.Find(dir)
	is.NoErr(err)
	preact := preact.New(module)
	entry, err := preact.CompileSSR("./index.jsx")
	is.NoErr(err)
	program, err := goja.Compile("index.jsx", string(entry.Contents), false)
	is.NoErr(err)
	vm := js.New()
	_, err = vm.RunProgram(program)
	is.NoErr(err)
	result, err := vm.RunString("bud.render({})")
	is.NoErr(err)
	is.Equal(result.String(), `{"html":"<h1>Hello World!</h1>","head":"<title>hello</title>"}`)
}

func TestCompileDOM(t *testing.T) {
	is := is.New(t)
	dir := t.TempDir()
	ctx := context.Background()
	is.NoErr(npm.Install(ctx, dir, "preact@10.18.1"))
	is.NoErr(testdir.WriteFiles(dir, map[string]string{
		"index.jsx": `
			import { createElement } from "preact";
			export default function(props, context) {
				context.head.push(<title>hello</title>)
				return <h1>Hello World!</h1>
			}
		`,
	}))
	module, err := mod.Find(dir)
	is.NoErr(err)
	preact := preact.New(module)
	entry, err := preact.CompileDOM("./index.jsx")
	is.NoErr(err)
	is.True(strings.Contains(string(entry.Contents), `var target = document.getElementById(".bud") || document.body;`))
	is.True(strings.Contains(string(entry.Contents), `document.getElementById(".bud_props")?.textContent || "{}"`))
}
