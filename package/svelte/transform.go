package svelte

import (
	"gitlab.com/mnm/bud/runtime/transform"
)

func NewTransformable(compiler *Compiler) *Transformable {
	return &Transformable{
		From: ".svelte",
		To:   ".js",
		For: transform.Platforms{
			// DOM transform (browser)
			transform.PlatformDOM: func(file *transform.File) error {
				dom, err := compiler.DOM(file.Path(), file.Code)
				if err != nil {
					return err
				}
				file.Code = []byte(dom.JS)
				return nil
			},

			// SSR transform (server)
			transform.PlatformSSR: func(file *transform.File) error {
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

type Transformable = transform.Transformable
