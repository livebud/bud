package svelte

//go:generate esbuild compiler.ts --format=iife --global-name=__svelte__ --bundle --platform=node --inject:shimssr.ts --external:url --outfile=compiler.js --log-level=warning

import (
	"encoding/json"
	"fmt"

	"gitlab.com/mnm/bud/internal/ldflag"

	_ "embed"

	"gitlab.com/mnm/bud/pkg/js"
)

type Input struct {
	VM  js.VM
	Dev bool
}

func New(vm js.VM) *Compiler {
	return &Compiler{vm, ldflag.IsDevelopment()}
}

type Compiler struct {
	VM  js.VM
	Dev bool
}

type SSR struct {
	JS  string
	CSS string
}

// compiler.js is used to compile .svelte files into JS & CSS
//go:embed compiler.js
var compiler string

// Compile server-rendered code
func (c *Compiler) SSR(path string, code []byte) (*SSR, error) {
	expr := fmt.Sprintf(`%s; __svelte__.compile({ "path": %q, "code": %q, "target": "ssr", "dev": %t, "css": false })`, compiler, path, code, c.Dev)
	result, err := c.VM.Eval(path, expr)
	if err != nil {
		return nil, err
	}
	out := new(SSR)
	if err := json.Unmarshal([]byte(result), out); err != nil {
		return nil, err
	}
	return out, nil
}

type DOM struct {
	JS  string
	CSS string
}

// Compile DOM code
func (c *Compiler) DOM(path string, code []byte) (*DOM, error) {
	expr := fmt.Sprintf(`%s; __svelte__.compile({ "path": %q, "code": %q, "target": "dom", "dev": %t, "css": true })`, compiler, path, code, c.Dev)
	result, err := c.VM.Eval(path, expr)
	if err != nil {
		return nil, err
	}
	out := new(DOM)
	if err := json.Unmarshal([]byte(result), out); err != nil {
		return nil, err
	}
	return out, nil
}
