package gohtml_test

import (
	"bytes"
	"io/fs"
	"testing"
	"testing/fstest"

	"github.com/livebud/bud/pkg/slots"
	"github.com/livebud/bud/pkg/view"
	"github.com/livebud/bud/pkg/view/gohtml"
	"github.com/matryer/is"
)

func memFS(files map[string]string) fs.FS {
	fsys := fstest.MapFS{}
	for path, data := range files {
		fsys[path] = &fstest.MapFile{Data: []byte(data)}
	}
	return fsys
}

func TestRender(t *testing.T) {
	is := is.New(t)
	fsys := memFS(map[string]string{
		"index.gohtml": `Hello {{ .Planet }}!`,
	})
	gohtml := gohtml.New(fsys)
	buf := new(bytes.Buffer)
	err := gohtml.Render(buf, "index.gohtml", &view.Data{
		Props: map[string]interface{}{
			"Planet": "Earth",
		},
	})
	is.NoErr(err)
	is.Equal(buf.String(), "Hello Earth!")
}

func TestSlot(t *testing.T) {
	is := is.New(t)
	fsys := memFS(map[string]string{
		"layout.gohtml": `<html><head><title>{{.Title}}</title></head><body>{{slot}}</body></html>`,
	})
	views := gohtml.New(fsys)
	buf := new(bytes.Buffer)
	slots := slots.New()
	slots.Write([]byte("<h2>Hello World!</h2>"))
	slots.Close()
	err := views.Render(buf, "layout.gohtml", &view.Data{
		Props: map[string]interface{}{
			"Title": "Hello",
		},
		Slots: slots.Next(),
	})
	is.NoErr(err)
	is.Equal(buf.String(), "<html><head><title>Hello</title></head><body><h2>Hello World!</h2></body></html>")
}

func TestAttrs(t *testing.T) {
	is := is.New(t)
	fsys := memFS(map[string]string{
		"index.gohtml": `Hello {{ with attr "settings" }}{{ .color }}{{ end }} {{ attr "Planet" }}!`,
	})
	gohtml := gohtml.New(fsys)
	buf := new(bytes.Buffer)
	err := gohtml.Render(buf, "index.gohtml", &view.Data{
		Attrs: map[string]any{
			"Planet": "Earth",
			"settings": map[string]any{
				"color": "Blue",
			},
		},
	})
	is.NoErr(err)
	is.Equal(buf.String(), "Hello Blue Earth!")
}

// func TestHeadSlot(t *testing.T) {
// 	is := is.New(t)
// 	fsys := memFS(map[string]string{
// 		"layout.gohtml": `
// 			<html lang="en">
// 				<head>
// 					{{ head }}
// 				</head>
// 				<body>
// 					{{ slot }}
// 				</body>
// 			</html>
// 		`,
// 		"index.gohtml": `
// 			{{ head "<title>Hello {{.Planet}}!</title>" }}
// 			<h1>Hello {{.Planet}}!</h1>
// 		`,
// 	})
// 	gohtml := gohtml.New(fsys)
// 	buf := new(bytes.Buffer)
// 	// layoutSlots := slot.Mock
// 	props := map[string]interface{}{
// 		"Planet": "Earth",
// 	}
// 	err := gohtml.Render(buf, "layout.gohtml", &view.Data{
// 		Props: props,
// 		// Slots: layoutSlots,
// 	})

// 	// err := gohtml.Render(buf, "index.gohtml", &view.Data{
// 	// 	Attrs: map[string]any{
// 	// 		"Planet": "Earth",
// 	// 		"settings": map[string]any{
// 	// 			"color": "Blue",
// 	// 		},
// 	// 	},
// 	// })
// 	is.NoErr(err)
// 	is.Equal(buf.String(), "Hello Blue Earth!")
// }
