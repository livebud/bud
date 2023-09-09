package transpiler

import (
	errors "errors"
	goldmark "github.com/livebud/bud/example/zero/transpiler/goldmark"
	tailwind "github.com/livebud/bud/example/zero/transpiler/tailwind"
	transpiler "github.com/livebud/bud/runtime/transpiler2"
)

// Load the transpiler
func Load(
	tailwind *tailwind.Transpiler,
	goldmark *goldmark.Transpiler,
) Transpiler {
	tr := transpiler.New()
	tr.Add(".gohtml", ".gohtml", tailwind.GohtmlToGohtml)
	tr.Add(".md", ".gohtml", goldmark.MdToGohtml)
	return &proxy{tr}
}

type Transpiler = transpiler.Interface

type proxy struct {
	Transpiler
}

func (p *proxy) Transpile(fromExt, toExt string, code []byte) ([]byte, error) {
	transpiled, err := p.Transpiler.Transpile(fromExt, toExt, code)
	if err != nil {
		if !errors.Is(err, transpiler.ErrNoPath) {
			return nil, err
		}
		return code, nil
	}
	return transpiled, nil
}
