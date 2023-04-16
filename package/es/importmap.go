package es

import (
	"regexp"

	esbuild "github.com/evanw/esbuild/pkg/api"
	"github.com/livebud/bud/package/log"
)

// ImportMap rewrites imports to a different path. If the import path is a URL,
// you'll need to also add the HTTP plugin.
// TODO: turn into a single `OnResolve` function
func ImportMap(log log.Log, imports map[string]string) esbuild.Plugin {
	return importMap(log, imports, false)
}

// ExternalImportMap rewrites imports to a different path and marks them as external
func ExternalImportMap(log log.Log, imports map[string]string) esbuild.Plugin {
	return importMap(log, imports, true)
}

func importMap(log log.Log, imports map[string]string, external bool) esbuild.Plugin {
	return esbuild.Plugin{
		Name: "import_map",
		Setup: func(epb esbuild.PluginBuild) {
			for from, to := range imports {
				from, to := from, to
				reFrom := regexp.QuoteMeta(from)
				if from[len(from)-1] == '/' {
					epb.OnResolve(esbuild.OnResolveOptions{Filter: `^` + reFrom}, func(args esbuild.OnResolveArgs) (result esbuild.OnResolveResult, err error) {
						log.Debugf("rewrote import %q to %q", args.Path, imports[from]+args.Path[len(from):])
						result.Path = imports[from] + args.Path[len(from):]
						result.Namespace = httpNamespace
						if external {
							result.External = true
						}
						return result, nil
					})
					continue
				}
				epb.OnResolve(esbuild.OnResolveOptions{Filter: `^` + reFrom + `$`}, func(args esbuild.OnResolveArgs) (result esbuild.OnResolveResult, err error) {
					log.Debugf("rewrote import %q to %q", args.Path, to)
					result.Path = to
					result.Namespace = httpNamespace
					if external {
						result.External = true
					}
					return result, nil
				})
			}
		},
	}
}
