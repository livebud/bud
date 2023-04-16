package es

import (
	"path"

	esbuild "github.com/evanw/esbuild/pkg/api"
)

// ExternalNodeModules externalizes node_modules
func ExternalNodeModules(prefix string) esbuild.Plugin {
	return esbuild.Plugin{
		Name: "external_node_modules",
		Setup: func(epb esbuild.PluginBuild) {
			// Regexp doesn't start with "." or "/" or "\"
			epb.OnResolve(esbuild.OnResolveOptions{Filter: `^[^\.\/\\]`}, func(args esbuild.OnResolveArgs) (result esbuild.OnResolveResult, err error) {
				result.Path = path.Join(prefix, args.Path)
				result.External = true
				return result, nil
			})
		},
	}
}
