package svelte

//go:generate go run github.com/evanw/esbuild/cmd/esbuild compiler.ts --format=iife --global-name=__svelte__ --bundle --platform=node --inject:compiler_shim.ts --external:url --outfile=compiler.js --log-level=warning

import (
	"context"
	"encoding/json"
	"fmt"

	_ "embed"

	"github.com/livebud/js"
)

// compiler.js is used to compile .svelte files into JS & CSS
//
//go:embed compiler.js
var compiler string

func Load(ctx context.Context, js js.VM) (*Compiler, error) {
	if _, err := js.Evaluate(ctx, "svelte/compiler.js", compiler); err != nil {
		return nil, err
	}
	return &Compiler{js, true}, nil
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
func (c *Compiler) SSR(ctx context.Context, path string, code []byte) (*SSR, error) {
	expr := fmt.Sprintf(`;__svelte__.compile({ "path": %q, "code": %q, "target": "ssr", "dev": %t, "css": false })`, path, code, c.Dev)
	result, err := c.VM.Evaluate(ctx, path, expr)
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
func (c *Compiler) DOM(ctx context.Context, path string, code []byte) (*DOM, error) {
	expr := fmt.Sprintf(`;__svelte__.compile({ "path": %q, "code": %q, "target": "dom", "dev": %t, "css": true })`, path, code, c.Dev)
	result, err := c.VM.Evaluate(ctx, path, expr)
	if err != nil {
		return nil, err
	}
	out := new(DOM)
	if err := json.Unmarshal([]byte(result), out); err != nil {
		return nil, err
	}
	return out, nil
}
