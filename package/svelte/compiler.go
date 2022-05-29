package svelte

//go:generate go run github.com/evanw/esbuild/cmd/esbuild compiler.ts --format=iife --global-name=__svelte__ --bundle --platform=node --inject:shimssr.ts --external:url --outfile=compiler.js --log-level=warning

import (
	"encoding/json"
	"fmt"

	_ "embed"

	"github.com/livebud/bud/package/js"
)

// compiler.js is used to compile .svelte files into JS & CSS
//go:embed compiler.js
var compiler string

func Load(vm js.VM) (*Compiler, error) {
	if err := vm.Script("svelte/compiler.js", compiler); err != nil {
		return nil, err
	}
	// TODO make dev configurable
	return &Compiler{vm, true}, nil
}

type Compiler struct {
	VM  js.VM
	Dev bool
}

type SSR struct {
	JS  string
	CSS string
}

// Compile server-rendered code
func (c *Compiler) SSR(path string, code []byte) (*SSR, error) {
	expr := fmt.Sprintf(`;__svelte__.compile({ "path": %q, "code": %q, "target": "ssr", "dev": %t, "css": false })`, path, code, c.Dev)
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
	expr := fmt.Sprintf(`;__svelte__.compile({ "path": %q, "code": %q, "target": "dom", "dev": %t, "css": true })`, path, code, c.Dev)
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
