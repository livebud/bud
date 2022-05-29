package transform_test

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"strings"
	"testing"

	esbuild "github.com/evanw/esbuild/pkg/api"
	"github.com/livebud/bud/internal/is"
	"github.com/livebud/bud/internal/testdir"
	"github.com/livebud/bud/runtime/transform"
)

func TestTransform(t *testing.T) {
	is := is.New(t)
	trace := []string{}
	transformer, err := transform.Load([]*transform.Transformable{
		{
			From: ".svelte",
			To:   ".svelte",
			For: transform.Platforms{
				transform.PlatformAll: func(file *transform.File) error {
					trace = append(trace, ".svelte>.svelte")
					is.Equal(file.Path(), "index.svelte")
					file.Code = bytes.ReplaceAll(file.Code, []byte("<h1>"), []byte("<h1 id='link'>"))
					return nil
				},
			},
		},
		{
			From: ".md",
			To:   ".svelte",
			For: transform.Platforms{
				transform.PlatformAll: func(file *transform.File) error {
					trace = append(trace, ".md>.svelte")
					is.Equal(file.Path(), "index.md")
					file.Code = []byte(`<h1>Hi world</h1>`)
					return nil
				},
			},
		},
		{
			From: ".svelte",
			To:   ".js",
			For: transform.Platforms{
				transform.PlatformAll: func(file *transform.File) error {
					trace = append(trace, ".svelte>.js")
					is.Equal(file.Path(), "index.svelte")
					file.Code = []byte(`document.body.innerHTML = "` + string(file.Code) + `"`)
					return nil
				},
			},
		},
	}...)
	is.NoErr(err)
	result, err := transformer.SSR.Transform("index.md", "index.js", []byte(`# Hi world`))
	is.NoErr(err)
	is.Equal(string(result), `document.body.innerHTML = "<h1 id='link'>Hi world</h1>"`)
	is.Equal(len(trace), 3)
	is.Equal(trace[0], ".md>.svelte")
	is.Equal(trace[1], ".svelte>.svelte")
	is.Equal(trace[2], ".svelte>.js")
	trace = []string{}
	result, err = transformer.DOM.Transform("index.md", "index.js", []byte(`# Hi world`))
	is.NoErr(err)
	is.Equal(string(result), `document.body.innerHTML = "<h1 id='link'>Hi world</h1>"`)
	is.Equal(len(trace), 3)
	is.Equal(trace[0], ".md>.svelte")
	is.Equal(trace[1], ".svelte>.svelte")
	is.Equal(trace[2], ".svelte>.js")
}

func TestPlugins(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()
	transformer, err := transform.Load([]*transform.Transformable{
		{
			From: ".svelte",
			To:   ".svelte",
			For: transform.Platforms{
				transform.PlatformAll: func(file *transform.File) error {
					file.Code = bytes.ReplaceAll(file.Code, []byte("<h1>"), []byte("<h1 id='link'>"))
					return nil
				},
			},
		},
		{
			From: ".md",
			To:   ".svelte",
			For: transform.Platforms{
				transform.PlatformAll: func(file *transform.File) error {
					file.Code = []byte(`<h1>Hi world</h1>`)
					return nil
				},
			},
		},
		{
			From: ".svelte",
			To:   ".js",
			For: transform.Platforms{
				transform.PlatformAll: func(file *transform.File) error {
					file.Code = []byte(`export default "` + string(file.Code) + `"`)
					return nil
				},
			},
		},
	}...)
	is.NoErr(err)
	plugins := transformer.SSR.Plugins()
	is.Equal(len(plugins), 2)
	// Create the test dir
	dir := t.TempDir()
	td := testdir.New(dir)
	td.Files["index.js"] = `
		import hello from "./hello.md"
		console.log(hello)
	`
	td.Files["hello.md"] = `
		# Hi world
	`
	is.NoErr(td.Write(ctx))
	result := esbuild.Build(esbuild.BuildOptions{
		EntryPoints:   []string{"index.js"},
		AbsWorkingDir: dir,
		Plugins:       plugins,
		Bundle:        true,
	})
	if len(result.Errors) > 0 {
		msgs := esbuild.FormatMessages(result.Errors, esbuild.FormatMessagesOptions{
			Kind:          esbuild.ErrorMessage,
			TerminalWidth: 80,
		})
		fmt.Fprintln(os.Stderr, strings.TrimSpace(strings.Join(msgs, "\n")))
		is.Fail()
	}
	is.Equal(len(result.OutputFiles), 1)
	output := result.OutputFiles[0]
	contents := string(output.Contents)
	is.True(strings.Contains(contents, `var hello_default = "<h1 id='link'>Hi world</h1>";`))
}

func TestTargets(t *testing.T) {
	is := is.New(t)
	trace := []string{}
	transformer, err := transform.Load([]*transform.Transformable{
		{
			From: ".svelte",
			To:   ".svelte",
			For: transform.Platforms{
				transform.PlatformAll: func(file *transform.File) error {
					trace = append(trace, ".svelte>.svelte")
					is.Equal(file.Path(), "index.svelte")
					file.Code = bytes.ReplaceAll(file.Code, []byte("<h1>"), []byte("<h1 id='link'>"))
					return nil
				},
			},
		},
		{
			From: ".md",
			To:   ".svelte",
			For: transform.Platforms{
				transform.PlatformAll: func(file *transform.File) error {
					trace = append(trace, ".md>.svelte")
					is.Equal(file.Path(), "index.md")
					file.Code = []byte(`<h1>Hi world</h1>`)
					return nil
				},
			},
		},
		{
			From: ".svelte",
			To:   ".js",
			For: transform.Platforms{
				transform.PlatformSSR: func(file *transform.File) error {
					trace = append(trace, ".svelte>.js(ssr)")
					is.Equal(file.Path(), "index.svelte")
					file.Code = []byte(`export default "` + string(file.Code) + `"`)
					return nil
				},
				transform.PlatformDOM: func(file *transform.File) error {
					trace = append(trace, ".svelte>.js(dom)")
					is.Equal(file.Path(), "index.svelte")
					file.Code = []byte(`document.body.innerHTML = "` + string(file.Code) + `"`)
					return nil
				},
			},
		},
	}...)
	is.NoErr(err)
	result, err := transformer.SSR.Transform("index.md", "index.js", []byte(`# Hi world`))
	is.NoErr(err)
	is.Equal(string(result), `export default "<h1 id='link'>Hi world</h1>"`)
	is.Equal(len(trace), 3)
	is.Equal(trace[0], ".md>.svelte")
	is.Equal(trace[1], ".svelte>.svelte")
	is.Equal(trace[2], ".svelte>.js(ssr)")
	trace = []string{}
	result, err = transformer.DOM.Transform("index.md", "index.js", []byte(`# Hi world`))
	is.NoErr(err)
	is.Equal(string(result), `document.body.innerHTML = "<h1 id='link'>Hi world</h1>"`)
	is.Equal(len(trace), 3)
	is.Equal(trace[0], ".md>.svelte")
	is.Equal(trace[1], ".svelte>.svelte")
	is.Equal(trace[2], ".svelte>.js(dom)")
}
