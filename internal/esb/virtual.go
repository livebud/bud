package esb

import (
	"regexp"

	esbuild "github.com/evanw/esbuild/pkg/api"
)

type OnLoadArgs = esbuild.OnLoadArgs
type OnLoadResult = esbuild.OnLoadResult

// Virtual creates a virtual file that can be imported as an entry. The
// loader function is called when the file is imported and is like any other
// ESBuild loader.
//
// Note: This plugin doesn't handle resolving relative paths to virtual files.
// See the test for an example of this limitation.
func Virtual(path string, loader func(OnLoadArgs) (OnLoadResult, error)) esbuild.Plugin {
	escapedPath := regexp.QuoteMeta(path)
	return esbuild.Plugin{
		Name: "virtual",
		Setup: func(epb esbuild.PluginBuild) {
			epb.OnResolve(esbuild.OnResolveOptions{Filter: "^" + escapedPath + "$"}, func(args esbuild.OnResolveArgs) (result esbuild.OnResolveResult, err error) {
				result.Path = args.Path
				result.Namespace = escapedPath
				return result, nil
			})
			epb.OnLoad(esbuild.OnLoadOptions{Filter: escapedPath, Namespace: escapedPath}, loader)
		},
	}
}
