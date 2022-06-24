package svelte

import (
	"github.com/livebud/bud/framework/transform/transformrt"
)

func NewTransformable(compiler *Compiler) *Transformable {
	return &Transformable{
		From: ".svelte",
		To:   ".js",
		For: transformrt.Platforms{
			// DOM transform (browser)
			transformrt.PlatformDOM: func(file *transformrt.File) error {
				dom, err := compiler.DOM(file.Path(), file.Code)
				if err != nil {
					return err
				}
				file.Code = []byte(dom.JS)
				return nil
			},

			// SSR transform (server)
			transformrt.PlatformSSR: func(file *transformrt.File) error {
				ssr, err := compiler.SSR(file.Path(), file.Code)
				if err != nil {
					return err
				}
				file.Code = []byte(ssr.JS)
				return nil
			},
		},
	}
}

type Transformable = transformrt.Transformable
