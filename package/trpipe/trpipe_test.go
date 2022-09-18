package trpipe_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/livebud/bud/internal/is"
	"github.com/livebud/bud/package/log/testlog"
	"github.com/livebud/bud/package/trpipe"
)

func TestSvelteJSX(t *testing.T) {
	is := is.New(t)
	log := testlog.New()
	pipe := trpipe.New(log)
	trace := []string{}
	pipe.Add(".svelte", ".jsx", func(file *trpipe.File) error {
		trace = append(trace, "svelte->jsx")
		file.Data = []byte(`export default function() { return ` + string(file.Data) + ` }`)
		return nil
	})
	pipe.Add(".svelte", ".svelte", func(file *trpipe.File) error {
		trace = append(trace, "svelte->svelte")
		file.Data = []byte("<main>" + string(file.Data) + "</main>")
		return nil
	})
	result, err := pipe.Run("hello.svelte", ".jsx", []byte("<h1>hi world</h1>"))
	is.NoErr(err)
	is.Equal(string(result), `export default function() { return <main><h1>hi world</h1></main> }`)
	is.Equal(strings.Join(trace, " "), "svelte->svelte svelte->jsx")
}

func TestSvelteSvelte(t *testing.T) {
	is := is.New(t)
	log := testlog.New()
	pipe := trpipe.New(log)
	trace := []string{}
	pipe.Add(".svelte", ".svelte", func(file *trpipe.File) error {
		trace = append(trace, "svelte->svelte")
		file.Data = []byte("<main>" + string(file.Data) + "</main>")
		return nil
	})
	result, err := pipe.Run("hello.svelte", ".svelte", []byte("<h1>hi world</h1>"))
	is.NoErr(err)
	is.Equal(string(result), `<main><h1>hi world</h1></main>`)
	is.Equal(strings.Join(trace, " "), "svelte->svelte")
}

func TestNoPath(t *testing.T) {
	is := is.New(t)
	log := testlog.New()
	pipe := trpipe.New(log)
	pipe.Add(".svelte", ".svelte", func(file *trpipe.File) error {
		file.Data = []byte("<main>" + string(file.Data) + "</main>")
		return nil
	})
	result, err := pipe.Run("hello.svelte", ".jsx", []byte("<h1>hi world</h1>"))
	is.True(err != nil)
	is.In(err.Error(), `trpipe: no path to transform ".svelte" to ".jsx" for "hello.svelte"`)
	is.Equal(result, nil)
}

func TestMultiStep(t *testing.T) {
	is := is.New(t)
	log := testlog.New()
	pipe := trpipe.New(log)
	trace := []string{}
	pipe.Add(".svelte", ".jsx", func(file *trpipe.File) error {
		trace = append(trace, "svelte->jsx")
		file.Data = []byte(`export default function() { return ` + string(file.Data) + ` }`)
		return nil
	})
	pipe.Add(".svelte", ".svelte", func(file *trpipe.File) error {
		trace = append(trace, "svelte->svelte")
		file.Data = []byte("<main>" + string(file.Data) + "</main>")
		return nil
	})
	pipe.Add(".md", ".svelte", func(file *trpipe.File) error {
		trace = append(trace, "md->svelte")
		file.Data = bytes.TrimPrefix(file.Data, []byte("# "))
		file.Data = []byte("<h1>" + string(file.Data) + "</h1>")
		return nil
	})
	result, err := pipe.Run("hello.md", ".jsx", []byte("# hi world"))
	is.NoErr(err)
	is.Equal(string(result), `export default function() { return <main><h1>hi world</h1></main> }`)
	is.Equal(strings.Join(trace, " "), "md->svelte svelte->svelte svelte->jsx")
}
