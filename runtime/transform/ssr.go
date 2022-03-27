package transform

import esbuild "github.com/evanw/esbuild/pkg/api"

type SSR struct {
	Map *Map
}

var _ Transformer = (*SSR)(nil)

func (d *SSR) Transform(fromPath, toPath string, code []byte) ([]byte, error) {
	return d.Map.SSR.Transform(fromPath, toPath, code)
}

func (d *SSR) Plugins() []esbuild.Plugin {
	return d.Map.SSR.Plugins()
}
