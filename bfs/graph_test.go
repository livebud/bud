package bfs

import (
	"testing"

	"github.com/matryer/is"
)

func TestGraphMatch(t *testing.T) {
	is := is.New(t)
	graph := newGraph()
	graph.Link("duo/main.go", "go.mod", WriteEvent|RemoveEvent)
	graph.Link("duo/view", "view/{**,*}.go", CreateEvent|RemoveEvent)
	graph.Link("duo/view", "view/public/public.go", WriteEvent|RemoveEvent)
	nodes := graph.match("view/public/public.go")
	is.Equal(len(nodes), 2)
	is.Equal(nodes[0], "view/public/public.go")
	is.Equal(nodes[1], "view/{**,*}.go")
}

func TestGraphIns(t *testing.T) {
	is := is.New(t)
	graph := newGraph()
	graph.Link("duo/main.go", "go.mod", WriteEvent|RemoveEvent)
	graph.Link("duo/view", "view/{**,*}.{svelte,jsx}", CreateEvent|RemoveEvent)
	graph.Link("duo/view/index.svelte", "view/index.svelte", WriteEvent|RemoveEvent)
	graph.Link("duo/view/index.svelte", "node_modules/svelte/internal/index.js", WriteEvent|RemoveEvent)
	graph.Link("duo/view/about/about.svelte", "view/about/about.svelte", WriteEvent|RemoveEvent)
	graph.Link("duo/view/about/about.svelte", "node_modules/svelte/internal/index.js", WriteEvent|RemoveEvent)
	nodes := graph.Ins("view/index.svelte", 0)
	is.Equal(len(nodes), 0)
	nodes = graph.Ins("view/index.svelte", WriteEvent)
	is.Equal(len(nodes), 1)
	is.Equal(nodes[0], "duo/view/index.svelte")
	nodes = graph.Ins("view/index.svelte", CreateEvent)
	is.Equal(len(nodes), 1)
	is.Equal(nodes[0], "duo/view")
	nodes = graph.Ins("view/index.svelte", WriteEvent|CreateEvent)
	is.Equal(len(nodes), 2)
	is.Equal(nodes[0], "duo/view")
	is.Equal(nodes[1], "duo/view/index.svelte")
	nodes = graph.Ins("view/index.svelte", WriteEvent|CreateEvent|RemoveEvent)
	is.Equal(len(nodes), 2)
	is.Equal(nodes[0], "duo/view")
	is.Equal(nodes[1], "duo/view/index.svelte")
}
func TestGraphDeepIns(t *testing.T) {
	is := is.New(t)
	graph := newGraph()
	graph.Link("duo/main.go", "go.mod", WriteEvent|RemoveEvent)
	graph.Link("duo/main.go", "duo/program/program.go", WriteEvent|RemoveEvent)
	graph.Link("duo/program/program.go", "duo/web/web.go", WriteEvent|RemoveEvent)
	graph.Link("duo/web/web.go", "duo/view/view.go", WriteEvent|RemoveEvent)
	graph.Link("duo/view/view.go", "view/{**,*}.{svelte,jsx}", CreateEvent|RemoveEvent)
	graph.Link("duo/web/web.go", "duo/web/ssr.js", WriteEvent|RemoveEvent)
	graph.Link("duo/web/ssr.js", "view/index.svelte", WriteEvent|RemoveEvent)
	nodes := graph.DeepIns("view/index.svelte", WriteEvent)
	is.Equal(len(nodes), 4)
	is.Equal(nodes[0], "duo/web/ssr.js")
	is.Equal(nodes[1], "duo/web/web.go")
	is.Equal(nodes[2], "duo/program/program.go")
	is.Equal(nodes[3], "duo/main.go")
	nodes = graph.DeepIns("view/edit.svelte", CreateEvent)
	is.Equal(len(nodes), 1)
	is.Equal(nodes[0], "duo/view/view.go")
	nodes = graph.DeepIns("duo/view/view.go", WriteEvent)
	is.Equal(len(nodes), 3)
	is.Equal(nodes[0], "duo/web/web.go")
	is.Equal(nodes[1], "duo/program/program.go")
	is.Equal(nodes[2], "duo/main.go")
}

// func TestGraphRemove(t *testing.T) {
// 	graph := newGraph()
// 	graph.Link("duo/main.go", "go.mod", WriteEvent)
// 	graph.Link("duo/view/_ssr.js", "node_modules/react-dom/server.browser.js", 0)
// 	graph.Link("duo/view/_ssr.js", "node_modules/object-assign/index.js", 0)
// 	graph.Link("duo/view/_ssr.js", "modules/youtube/index.ts", 0)
// 	graph.Link("duo/view/_ssr.js", "view/index.jsx", 0)
// 	graph.Link("duo/view/_ssr.js", "modules/uid/index.ts", 0)
// 	graph.Link("duo/view/_ssr.js", "node_modules/react/cjs/react.development.js", 0)
// 	graph.Link("duo/view/_ssr.js", "node_modules/react/index.js", 0)
// 	graph.Link("duo/view/_ssr.js", "node_modules/react-dom/cjs/react-dom-server.browser.development.js", 0)
// 	graph.Link("duo/main.go", "duo/program/program.go", 0)
// 	graph.Link("duo/program/program.go", "duo/web/web.go", 0)
// 	fmt.Println(graph.String())
// 	// for _, in := range graph.Ins(graph, "modules/uid/index.ts") {
// 	// 	graph.Remove(in)
// 	// }
// 	// fmt.Println(graph.String())
// }
