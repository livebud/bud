package svelte

import (
	"gitlab.com/mnm/bud/transform"
)

func NewTransformable(compiler *Compiler) *Transformable {
	return &transform.Transformable{
		From: ".svelte",
		To:   ".js",
		Map: transform.Map{
			// Browser transform
			transform.PlatformBrowser: func(file *transform.File) error {
				dom, err := compiler.DOM(file.Path(), file.Code)
				if err != nil {
					return err
				}
				file.Code = []byte(dom.JS)
				return nil
			},

			// Node transform
			transform.PlatformNode: func(file *transform.File) error {
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
