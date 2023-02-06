package es

import (
	"io"
	"net/http"

	esbuild "github.com/evanw/esbuild/pkg/api"
)

func httpPlugin(absDir string) esbuild.Plugin {
	return esbuild.Plugin{
		Name: "bud-http",
		Setup: func(epb esbuild.PluginBuild) {
			epb.OnResolve(esbuild.OnResolveOptions{Filter: `^http[s]?://`}, func(args esbuild.OnResolveArgs) (result esbuild.OnResolveResult, err error) {
				result.Namespace = "http"
				result.Path = args.Path
				return result, nil
			})
			epb.OnLoad(esbuild.OnLoadOptions{Filter: `.*`, Namespace: `http`}, func(args esbuild.OnLoadArgs) (result esbuild.OnLoadResult, err error) {
				res, err := http.Get(args.Path)
				if err != nil {
					return result, err
				}
				defer res.Body.Close()
				body, err := io.ReadAll(res.Body)
				if err != nil {
					return result, err
				}
				contents := string(body)
				result.ResolveDir = absDir
				result.Contents = &contents
				result.Loader = esbuild.LoaderJS
				return result, nil
			})
		},
	}
}

func esmPlugin(absDir string) esbuild.Plugin {
	return esbuild.Plugin{
		Name: "bud-esm",
		Setup: func(epb esbuild.PluginBuild) {
			epb.OnResolve(esbuild.OnResolveOptions{Filter: `^[a-z0-9@]`}, func(args esbuild.OnResolveArgs) (result esbuild.OnResolveResult, err error) {
				result.Namespace = "esm"
				result.Path = args.Path
				return result, nil
			})
			epb.OnLoad(esbuild.OnLoadOptions{Filter: `.*`, Namespace: `esm`}, func(args esbuild.OnLoadArgs) (result esbuild.OnLoadResult, err error) {
				res, err := http.Get("https://esm.sh/" + args.Path)
				if err != nil {
					return result, err
				}
				defer res.Body.Close()
				body, err := io.ReadAll(res.Body)
				if err != nil {
					return result, err
				}
				contents := string(body)
				result.ResolveDir = absDir
				result.Contents = &contents
				result.Loader = esbuild.LoaderJS
				return result, nil
			})
		},
	}
}
