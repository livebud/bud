package transformrt

import esbuild "github.com/evanw/esbuild/pkg/api"

type DOM struct {
	Map *Map
}

var _ Transformer = (*DOM)(nil)

func (d *DOM) Transform(fromPath, toPath string, code []byte) ([]byte, error) {
	return d.Map.DOM.Transform(fromPath, toPath, code)
}

func (d *DOM) Plugins() []esbuild.Plugin {
	return d.Map.DOM.Plugins()
}
