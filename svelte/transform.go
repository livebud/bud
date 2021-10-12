package svelte

import (
	"github.com/go-duo/bud/transform"
)

func NewTransform(compiler *Compiler) *Transform {
	return &Transform{compiler}
}

type Transform struct {
	compiler *Compiler
}

func (t *Transform) SvelteToJS(file *transform.File) error {
	ssr, err := t.compiler.SSR(file.Path(), file.Code)
	if err != nil {
		return err
	}
	file.Code = []byte(ssr.JS)
	return nil
}

// func svelteTransformPlugin(svelte *svelte.Compiler, transformer transform.Transformer) esbuild.Plugin {
// 	return esbuild.Plugin{
// 		Name: "svelte_transform",
// 		Setup: func(epb esbuild.PluginBuild) {
// 			// Load svelte files. Add import if not present
// 			epb.OnLoad(esbuild.OnLoadOptions{Filter: `\.svelte$`}, func(args esbuild.OnLoadArgs) (result esbuild.OnLoadResult, err error) {
// 				code, err := os.ReadFile(args.Path)
// 				if err != nil {
// 					return result, err
// 				}
// 				ssr, err := svelte.SSR(args.Path, code)
// 				if err != nil {
// 					return result, err
// 				}
// 				result.ResolveDir = filepath.Dir(args.Path)
// 				result.Contents = &ssr.JS
// 				result.Loader = esbuild.LoaderJS
// 				return result, nil
// 			})
// 		},
// 	}
// }

// // Transform svelte files
// func svelteTransformPlugin(svelte *svelte.Compiler) esbuild.Plugin {
// 	return esbuild.Plugin{
// 		Name: "svelte_transform",
// 		Setup: func(epb esbuild.PluginBuild) {
// 			// Load svelte files. Add import if not present
// 			epb.OnLoad(esbuild.OnLoadOptions{Filter: `\.svelte$`}, func(args esbuild.OnLoadArgs) (result esbuild.OnLoadResult, err error) {
// 				code, err := os.ReadFile(args.Path)
// 				if err != nil {
// 					return result, err
// 				}
// 				dom, err := svelte.DOM(args.Path, code)
// 				if err != nil {
// 					return result, err
// 				}
// 				result.ResolveDir = filepath.Dir(args.Path)
// 				result.Contents = &dom.JS
// 				result.Loader = esbuild.LoaderJS
// 				return result, nil
// 			})
// 		},
// 	}
// }
