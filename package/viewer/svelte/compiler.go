package svelte

//go:generate go run github.com/evanw/esbuild/cmd/esbuild compiler.ts --format=iife --global-name=__svelte__ --bundle --platform=node --inject:compiler_shim.ts --external:url --outfile=compiler.js --log-level=warning

import (
	"context"
	"encoding/json"
	"fmt"

	_ "embed"

	"github.com/livebud/bud/framework"
	"github.com/livebud/js"
)

// compiler.js is used to compile .svelte files into JS & CSS
//
//go:embed compiler.js
var compiler string

func Load(flag *framework.Flag, js js.VM) (*Compiler, error) {
	if _, err := js.Evaluate(context.TODO(), "svelte/compiler.js", compiler); err != nil {
		return nil, err
	}
	return &Compiler{flag, js}, nil
}

type Compiler struct {
	flag *framework.Flag
	js   js.VM
}

type SSR struct {
	JS  string
	CSS string
}

// Compile server-rendered code
func (c *Compiler) SSR(ctx context.Context, path string, code []byte) (*SSR, error) {
	isDev := !c.flag.Minify
	expr := fmt.Sprintf(`;__svelte__.compile({ "path": %q, "code": %q, "target": "ssr", "dev": %t, "css": false })`, path, code, isDev)
	result, err := c.js.Evaluate(ctx, path, expr)
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
	isDev := !c.flag.Minify
	expr := fmt.Sprintf(`;__svelte__.compile({ "path": %q, "code": %q, "target": "dom", "dev": %t, "css": false })`, path, code, isDev)
	result, err := c.js.Evaluate(ctx, path, expr)
	if err != nil {
		return nil, err
	}
	out := new(DOM)
	if err := json.Unmarshal([]byte(result), out); err != nil {
		return nil, err
	}
	return out, nil
}
