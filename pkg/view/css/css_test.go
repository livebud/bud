package css_test

import (
	"testing"

	"github.com/livebud/bud/internal/testdir"
	"github.com/livebud/bud/pkg/mod"
	"github.com/livebud/bud/pkg/view/css"
	"github.com/matryer/is"
	"github.com/matthewmueller/diff"
)

func TestCompileCSS(t *testing.T) {
	is := is.New(t)
	dir := t.TempDir()
	is.NoErr(testdir.WriteFiles(dir, map[string]string{
		"index.css": `
			body { background-color: red; }
		`,
	}))
	module, err := mod.Find(dir)
	is.NoErr(err)
	css := css.New(module)
	sheet, err := css.Compile("./index.css")
	is.NoErr(err)
	diff.TestContent(t, string(sheet.Contents), `
		/* index.css */
		body {
		  background-color: red;
		}
	`)
}

func TestCompileImportedCSS(t *testing.T) {
	is := is.New(t)
	dir := t.TempDir()
	is.NoErr(testdir.WriteFiles(dir, map[string]string{
		"Header.css": `
			.Header { background-color: green; }
		`,
		"index.css": `
			@import "./Header.css";
			body { background-color: red; }
		`,
	}))
	module, err := mod.Find(dir)
	is.NoErr(err)
	css := css.New(module)
	sheet, err := css.Compile("./index.css")
	is.NoErr(err)
	diff.TestContent(t, string(sheet.Contents), `
		/* Header.css */
		.Header {
		  background-color: green;
		}

		/* index.css */
		body {
		  background-color: red;
		}
	`)
}
