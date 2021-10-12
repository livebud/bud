package transform_test

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"testing/fstest"

	esbuild "github.com/evanw/esbuild/pkg/api"
	"github.com/go-duo/bud/internal/vfs"
	"github.com/go-duo/bud/transform"
	"github.com/matryer/is"
)

func TestTransform(t *testing.T) {
	is := is.New(t)
	trace := []string{}
	transformer, err := transform.Load([]*transform.Transform{
		{
			From: ".svelte",
			To:   ".svelte",
			Func: func(file *transform.File) error {
				trace = append(trace, ".svelte>.svelte")
				is.Equal(file.Path(), "index.svelte")
				file.Code = bytes.ReplaceAll(file.Code, []byte("<h1>"), []byte("<h1 id='link'>"))
				return nil
			},
		},
		{
			From: ".md",
			To:   ".svelte",
			Func: func(file *transform.File) error {
				trace = append(trace, ".md>.svelte")
				is.Equal(file.Path(), "index.md")
				file.Code = []byte(`<h1>Hi world</h1>`)
				return nil
			},
		},
		{
			From: ".svelte",
			To:   ".js",
			Func: func(file *transform.File) error {
				trace = append(trace, ".svelte>.js")
				is.Equal(file.Path(), "index.svelte")
				file.Code = []byte(`document.body.innerHTML = "` + string(file.Code) + `"`)
				return nil
			},
		},
	})
	is.NoErr(err)
	result, err := transformer.Transform("index.md", "index.js", []byte(`# Hi world`))
	is.NoErr(err)
	is.Equal(string(result), `document.body.innerHTML = "<h1 id='link'>Hi world</h1>"`)
	is.Equal(len(trace), 3)
	is.Equal(trace[0], ".md>.svelte")
	is.Equal(trace[1], ".svelte>.svelte")
	is.Equal(trace[2], ".svelte>.js")
}

func TestPlugins(t *testing.T) {
	is := is.New(t)
	transformer, err := transform.Load([]*transform.Transform{
		{
			From: ".svelte",
			To:   ".svelte",
			Func: func(file *transform.File) error {
				file.Code = bytes.ReplaceAll(file.Code, []byte("<h1>"), []byte("<h1 id='link'>"))
				return nil
			},
		},
		{
			From: ".md",
			To:   ".svelte",
			Func: func(file *transform.File) error {
				file.Code = []byte(`<h1>Hi world</h1>`)
				return nil
			},
		},
		{
			From: ".svelte",
			To:   ".js",
			Func: func(file *transform.File) error {
				file.Code = []byte(`export default "` + string(file.Code) + `"`)
				return nil
			},
		},
	})
	is.NoErr(err)
	plugins := transformer.Plugins()
	is.Equal(len(plugins), 2)
	// Create files in _tmp
	cwd, err := os.Getwd()
	is.NoErr(err)
	absdir := filepath.Join(cwd, "_tmp")
	is.NoErr(os.RemoveAll(absdir))
	defer func() {
		if !t.Failed() {
			is.NoErr(os.RemoveAll(absdir))
		}
	}()
	is.NoErr(vfs.WriteAll(".", absdir, vfs.Memory{
		"index.js": &fstest.MapFile{
			Data: []byte(`
				import hello from "./hello.md"
				console.log(hello)
			`),
		},
		"hello.md": &fstest.MapFile{
			Data: []byte(`# Hi world"`),
		},
	}))
	result := esbuild.Build(esbuild.BuildOptions{
		EntryPoints:   []string{"index.js"},
		AbsWorkingDir: absdir,
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
