package esb

import (
	"encoding/json"
	"regexp"

	esbuild "github.com/evanw/esbuild/pkg/api"
)

func Env(path string, env any) esbuild.Plugin {
	escapedPath := regexp.QuoteMeta(path)
	return esbuild.Plugin{
		Name: "env",
		Setup: func(epb esbuild.PluginBuild) {
			epb.OnResolve(esbuild.OnResolveOptions{Filter: "^" + escapedPath + "$"}, func(args esbuild.OnResolveArgs) (result esbuild.OnResolveResult, err error) {
				result.Path = args.Path
				result.Namespace = escapedPath
				return result, nil
			})
			epb.OnLoad(esbuild.OnLoadOptions{Filter: escapedPath, Namespace: escapedPath}, func(args esbuild.OnLoadArgs) (result esbuild.OnLoadResult, err error) {
				env, err := json.Marshal(env)
				if err != nil {
					return result, err
				}
				result.Loader = esbuild.LoaderJSON
				contents := string(env)
				result.Contents = &contents
				return result, nil
			})
		},
	}
}
