package transpiler_test

import (
	"bytes"
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/livebud/bud/internal/is"
	"github.com/livebud/bud/package/transpiler"
)

func TestTranspileSvelteToJSX(t *testing.T) {
	is := is.New(t)
	tr := transpiler.New()
	paths := []string{}
	trace := []string{}
	tr.Add(".svelte", ".jsx", func(ctx context.Context, file *transpiler.File) error {
		paths = append(paths, file.Path())
		trace = append(trace, "svelte->jsx")
		file.Data = []byte(`export default function() { return ` + string(file.Data) + ` }`)
		return nil
	})
	tr.Add(".svelte", ".svelte", func(ctx context.Context, file *transpiler.File) error {
		paths = append(paths, file.Path())
		trace = append(trace, "svelte->svelte")
		file.Data = []byte("<main>" + string(file.Data) + "</main>")
		return nil
	})
	ctx := context.Background()
	result, err := tr.Transpile(ctx, "hello.svelte", ".jsx", []byte("<h1>hi world</h1>"))
	is.NoErr(err)
	is.Equal(string(result), `export default function() { return <main><h1>hi world</h1></main> }`)
	is.Equal(strings.Join(trace, " "), "svelte->svelte svelte->jsx")
	is.Equal(strings.Join(paths, " "), "hello.svelte hello.svelte")
	hops, err := tr.Path(".svelte", ".jsx")
	is.NoErr(err)
	is.Equal(len(hops), 2)
	is.Equal(hops[0], ".svelte")
	is.Equal(hops[1], ".jsx")
}

func TestSvelteSvelte(t *testing.T) {
	is := is.New(t)
	tr := transpiler.New()
	trace := []string{}
	paths := []string{}
	tr.Add(".svelte", ".svelte", func(ctx context.Context, file *transpiler.File) error {
		paths = append(paths, file.Path())
		trace = append(trace, "svelte->svelte")
		file.Data = []byte("<main>" + string(file.Data) + "</main>")
		return nil
	})
	ctx := context.Background()
	result, err := tr.Transpile(ctx, "hello.svelte", ".svelte", []byte("<h1>hi world</h1>"))
	is.NoErr(err)
	is.Equal(string(result), `<main><h1>hi world</h1></main>`)
	is.Equal(strings.Join(trace, " "), "svelte->svelte")
	is.Equal(strings.Join(paths, " "), "hello.svelte")
	hops, err := tr.Path(".svelte", ".svelte")
	is.NoErr(err)
	is.Equal(len(hops), 1)
	is.Equal(hops[0], ".svelte")
}

func TestNoExt(t *testing.T) {
	is := is.New(t)
	tr := transpiler.New()
	tr.Add(".svelte", ".svelte", func(ctx context.Context, file *transpiler.File) error {
		file.Data = []byte("<main>" + string(file.Data) + "</main>")
		return nil
	})
	ctx := context.Background()
	result, err := tr.Transpile(ctx, "hello.svelte", ".jsx", []byte("<h1>hi world</h1>"))
	is.True(err != nil)
	is.True(errors.Is(err, transpiler.ErrNoPath))
	is.Equal(result, nil)
	hops, err := tr.Path(".svelte", ".jsx")
	is.True(err != nil)
	is.True(errors.Is(err, transpiler.ErrNoPath))
	is.Equal(hops, nil)
}

func TestNoPath(t *testing.T) {
	is := is.New(t)
	tr := transpiler.New()
	tr.Add(".svelte", ".svelte", func(ctx context.Context, file *transpiler.File) error {
		file.Data = []byte("<main>" + string(file.Data) + "</main>")
		return nil
	})
	tr.Add(".jsx", ".jsx", func(ctx context.Context, file *transpiler.File) error {
		file.Data = []byte("<main>" + string(file.Data) + "</main>")
		return nil
	})
	ctx := context.Background()
	result, err := tr.Transpile(ctx, "hello.svelte", ".jsx", []byte("<h1>hi world</h1>"))
	is.True(err != nil)
	is.True(errors.Is(err, transpiler.ErrNoPath))
	is.Equal(result, nil)
	hops, err := tr.Path(".svelte", ".jsx")
	is.True(err != nil)
	is.True(errors.Is(err, transpiler.ErrNoPath))
	is.Equal(hops, nil)
}

func TestMultiStep(t *testing.T) {
	is := is.New(t)
	tr := transpiler.New()
	trace := []string{}
	path := []string{}
	tr.Add(".jsx", ".jsx", func(ctx context.Context, file *transpiler.File) error {
		path = append(path, file.Path())
		trace = append(trace, "jsx->jsx")
		file.Data = []byte("/* some prelude */ " + string(file.Data))
		return nil
	})
	tr.Add(".svelte", ".jsx", func(ctx context.Context, file *transpiler.File) error {
		path = append(path, file.Path())
		trace = append(trace, "svelte->jsx")
		file.Data = []byte(`export default function() { return ` + string(file.Data) + ` }`)
		return nil
	})
	tr.Add(".svelte", ".svelte", func(ctx context.Context, file *transpiler.File) error {
		path = append(path, file.Path())
		trace = append(trace, "1:svelte->svelte")
		file.Data = []byte("<main>" + string(file.Data) + "</main>")
		return nil
	})
	tr.Add(".svelte", ".svelte", func(ctx context.Context, file *transpiler.File) error {
		path = append(path, file.Path())
		trace = append(trace, "2:svelte->svelte")
		file.Data = []byte("<div>" + string(file.Data) + "</div>")
		return nil
	})
	tr.Add(".md", ".svelte", func(ctx context.Context, file *transpiler.File) error {
		path = append(path, file.Path())
		trace = append(trace, "md->svelte")
		file.Data = bytes.TrimPrefix(file.Data, []byte("# "))
		file.Data = []byte("<h1>" + string(file.Data) + "</h1>")
		return nil
	})
	ctx := context.Background()
	result, err := tr.Transpile(ctx, "hello.md", ".jsx", []byte("# hi world"))
	is.NoErr(err)
	is.Equal(strings.Join(trace, " "), "md->svelte 1:svelte->svelte 2:svelte->svelte svelte->jsx jsx->jsx")
	is.Equal(strings.Join(path, " "), "hello.md hello.svelte hello.svelte hello.svelte hello.jsx")
	is.Equal(string(result), `/* some prelude */ export default function() { return <div><main><h1>hi world</h1></main></div> }`)
	hops, err := tr.Path(".md", ".jsx")
	is.NoErr(err)
	is.Equal(len(hops), 3)
	is.Equal(hops[0], ".md")
	is.Equal(hops[1], ".svelte")
	is.Equal(hops[2], ".jsx")
}

func TestTranspileSSRJS(t *testing.T) {
	is := is.New(t)
	tr := transpiler.New()
	tr.Add(".svelte", ".ssr.js", func(ctx context.Context, file *transpiler.File) error {
		file.Data = []byte(`module.exports = "` + string(file.Data) + `"`)
		return nil
	})
	ctx := context.Background()
	code, err := tr.Transpile(ctx, "hello.svelte", ".ssr.js", []byte("<h1>hello</h1>"))
	is.NoErr(err)
	is.Equal(string(code), `module.exports = "<h1>hello</h1>"`)
	code, err = tr.Transpile(ctx, "hello.svelte", ".ssr.js", []byte("<h1>world</h1>"))
	is.NoErr(err)
	is.Equal(string(code), `module.exports = "<h1>world</h1>"`)
}

func TestTranspileNoPath(t *testing.T) {
	is := is.New(t)
	tr := transpiler.New()
	ctx := context.Background()
	code, err := tr.Transpile(ctx, "hello.jsx", ".ssr.js", []byte("<h1>hello</h1>"))
	is.True(err != nil)
	is.True(errors.Is(err, transpiler.ErrNoPath))
	is.Equal(code, nil)
}

func TestTranspileBest(t *testing.T) {
	is := is.New(t)
	tr := transpiler.New()
	trace := []string{}
	path := []string{}
	tr.Add(".jsx", ".js", func(ctx context.Context, file *transpiler.File) error {
		path = append(path, file.Path())
		trace = append(trace, "jsx->js")
		file.Data = []byte("// " + string(file.Data))
		return nil
	})
	tr.Add(".jsx", ".jsx", func(ctx context.Context, file *transpiler.File) error {
		path = append(path, file.Path())
		trace = append(trace, "jsx->jsx")
		file.Data = []byte("/* some prelude */ " + string(file.Data))
		return nil
	})
	tr.Add(".svelte", ".jsx", func(ctx context.Context, file *transpiler.File) error {
		path = append(path, file.Path())
		trace = append(trace, "svelte->jsx")
		file.Data = []byte(`export default function() { return ` + string(file.Data) + ` }`)
		return nil
	})
	tr.Add(".svelte", ".svelte", func(ctx context.Context, file *transpiler.File) error {
		path = append(path, file.Path())
		trace = append(trace, "1:svelte->svelte")
		file.Data = []byte("<main>" + string(file.Data) + "</main>")
		return nil
	})
	tr.Add(".svelte", ".svelte", func(ctx context.Context, file *transpiler.File) error {
		path = append(path, file.Path())
		trace = append(trace, "2:svelte->svelte")
		file.Data = []byte("<div>" + string(file.Data) + "</div>")
		return nil
	})
	tr.Add(".md", ".svelte", func(ctx context.Context, file *transpiler.File) error {
		path = append(path, file.Path())
		trace = append(trace, "md->svelte")
		file.Data = bytes.TrimPrefix(file.Data, []byte("# "))
		file.Data = []byte("<h1>" + string(file.Data) + "</h1>")
		return nil
	})
	tr.Add(".md", ".jsx", func(ctx context.Context, file *transpiler.File) error {
		path = append(path, file.Path())
		trace = append(trace, "md->jsx")
		file.Data = bytes.TrimPrefix(file.Data, []byte("# "))
		file.Data = []byte("<h1>" + string(file.Data) + "</h1>")
		return nil
	})
	tr.Add(".png", ".jpg", func(ctx context.Context, file *transpiler.File) error {
		path = append(path, file.Path())
		trace = append(trace, "png->jpg")
		file.Data = []byte("jpg")
		return nil
	})
	// Prefer jsx because same hops as svelte
	ext, err := tr.Best(".md", []string{".jsx", ".svelte", ".js"})
	is.NoErr(err)
	is.Equal(ext, ".jsx")
	// Prefer svelte because same hops as jsx
	ext, err = tr.Best(".md", []string{".svelte", ".jsx", ".js"})
	is.NoErr(err)
	is.Equal(ext, ".svelte")
	// Svelte has less hops than JS
	ext, err = tr.Best(".md", []string{".js", ".svelte", ".jsx"})
	is.NoErr(err)
	is.Equal(ext, ".svelte")

	ext, err = tr.Best(".scss", []string{".css"})
	is.True(err != nil)
	is.True(errors.Is(err, transpiler.ErrNoPath))
	is.Equal(ext, "")

	ext, err = tr.Best(".jpg", []string{".min.jpg"})
	is.True(err != nil)
	is.True(errors.Is(err, transpiler.ErrNoPath))
	is.Equal(ext, "")
}
