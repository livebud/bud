package terminal

import (
	"bytes"
	_ "embed"

	terminal "github.com/buildkite/terminal-to-html"
)

// Pre tag
func Pre(data []byte) []byte {
	var b bytes.Buffer
	b.WriteString(`<pre class="term-container">`)
	b.Write(terminal.Render(data))
	b.WriteString(`</pre>`)
	return b.Bytes()
}

//go:embed terminal.css
var css []byte

// CSS that needs to be rendered
func CSS() []byte {
	return css
}
