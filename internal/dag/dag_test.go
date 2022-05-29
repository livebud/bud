package dag_test

import (
	"fmt"
	"testing"

	"github.com/livebud/bud/internal/dag"
	"github.com/livebud/bud/internal/is"
)

func TestRemove(t *testing.T) {
	graph := dag.New()
	graph.Link("go.mod", "duo/main.go")
	graph.Link("duo/view/_ssr.js", "node_modules/react-dom/server.browser.js")
	graph.Link("duo/view/_ssr.js", "node_modules/object-assign/index.js")
	graph.Link("duo/view/_ssr.js", "modules/youtube/index.ts")
	graph.Link("duo/view/_ssr.js", "view/index.jsx")
	graph.Link("duo/view/_ssr.js", "modules/uid/index.ts")
	graph.Link("duo/view/_ssr.js", "node_modules/react/cjs/react.development.js")
	graph.Link("duo/view/_ssr.js", "node_modules/react/index.js")
	graph.Link("duo/view/_ssr.js", "node_modules/react-dom/cjs/react-dom-server.browser.development.js")
	graph.Link("duo/main.go", "duo/program/program.go")
	graph.Link("duo/program/program.go", "duo/web/web.go")
	fmt.Println(graph.String())
	graph.Remove(graph.Parents("modules/uid/index.ts")...)
	fmt.Println(graph.String())
}

func TestShortestPath(t *testing.T) {
	is := is.New(t)
	graph := dag.New()
	graph.Link(".svelte", ".svelte")
	graph.Link(".svelte", ".js")
	graph.Link(".md", ".svelte")
	graph.Link(".md", ".mdx")
	graph.Link(".mdx", ".jsx")
	graph.Link(".jsx", ".js")
	graph.Link(".jsx", ".jsx")
	// digraph g {
	// 	".md" -> ".svelte"
	// 	".md" -> ".mdx"
	// 	".mdx" -> ".jsx"
	// 	".jsx" -> ".js"
	// 	".svelte" -> ".js"
	// }
	nodes, err := graph.ShortestPath(".md", ".js")
	is.NoErr(err)
	is.Equal(len(nodes), 3)
	is.Equal(nodes[0], ".md")
	is.Equal(nodes[1], ".svelte")
	is.Equal(nodes[2], ".js")
}
func TestShortestPathSingle(t *testing.T) {
	is := is.New(t)
	graph := dag.New()
	graph.Link(".svelte", ".svelte")
	graph.Link(".svelte", ".js")
	nodes, err := graph.ShortestPath(".svelte", ".js")
	is.NoErr(err)
	is.Equal(len(nodes), 2)
	is.Equal(nodes[0], ".svelte")
	is.Equal(nodes[1], ".js")
}
func TestShortestPathNone(t *testing.T) {
	is := is.New(t)
	graph := dag.New()
	graph.Link(".md", ".svelte")
	graph.Link(".mdx", ".jsx")
	nodes, err := graph.ShortestPath(".svelte", ".jsx")
	is.Equal(err.Error(), `dag: no path between ".svelte" and ".jsx"`)
	is.Equal(nodes, nil)
}
func TestShortestPathOf(t *testing.T) {
	is := is.New(t)
	graph := dag.New()
	graph.Link(".svelte", ".svelte")
	graph.Link(".svelte", ".js")
	graph.Link(".md", ".svelte")
	graph.Link(".mdx", ".jsx")
	graph.Link(".jsx", ".jsx")
	nodes, err := graph.ShortestPathOf(".md", []string{".jsx", ".js"})
	is.NoErr(err)
	is.Equal(len(nodes), 3)
	is.Equal(nodes[0], ".md")
	is.Equal(nodes[1], ".svelte")
	is.Equal(nodes[2], ".js")
}

func TestShortestPathOfNone(t *testing.T) {
	is := is.New(t)
	graph := dag.New()
	graph.Link(".svelte", ".svelte")
	graph.Link(".svelte", ".js")
	graph.Link(".md", ".svelte")
	graph.Link(".mdx", ".jsx")
	graph.Link(".jsx", ".jsx")
	nodes, err := graph.ShortestPathOf(".md", []string{".jsx", ".mdx"})
	is.Equal(err.Error(), `dag: no path between ".md" and [.jsx .mdx]`)
	is.Equal(nodes, nil)
}
