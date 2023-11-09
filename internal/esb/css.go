package esb

import (
	"fmt"
	"path/filepath"

	esbuild "github.com/evanw/esbuild/pkg/api"
)

var ErrCSSImportInJS = fmt.Errorf("esb: CSS imports are not allowed in non-CSS files")

func DisableCSSImportsInJS() esbuild.Plugin {
	return esbuild.Plugin{
		Name: "disable-css-imports-in-js",
		Setup: func(build esbuild.PluginBuild) {
			build.OnResolve(esbuild.OnResolveOptions{Filter: `.css$`}, func(args esbuild.OnResolveArgs) (result esbuild.OnResolveResult, err error) {
				if filepath.Ext(args.Importer) != ".css" {
					return result, fmt.Errorf("esb: non-CSS file %q is not allowed to import CSS file %q", args.Importer, args.Path)
				}
				return result, nil
			})
		},
	}
}
