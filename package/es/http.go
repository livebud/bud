package es

import (
	"fmt"
	"io"
	"net/http"
	"net/url"

	esbuild "github.com/evanw/esbuild/pkg/api"
)

const httpNamespace = "http-url"

func ExternalHTTP() esbuild.Plugin {
	return esbuild.Plugin{
		Name: "http",
		Setup: func(build esbuild.PluginBuild) {
			// Externalize any HTTP URLs
			build.OnResolve(esbuild.OnResolveOptions{Filter: `^https?://`}, func(args esbuild.OnResolveArgs) (result esbuild.OnResolveResult, err error) {
				result.Path = args.Path
				result.External = true
				return result, nil
			})
		},
	}
}

func HTTP(client *http.Client) esbuild.Plugin {
	return esbuild.Plugin{
		Name: "http",
		Setup: func(build esbuild.PluginBuild) {
			// Intercept import paths starting with "http:" and "https:" so esbuild
			// doesn't attempt to map them to a file system location. Tag them with
			// the "http-url" namespace to associate them with this plugin.
			build.OnResolve(esbuild.OnResolveOptions{Filter: `^https?://`}, func(args esbuild.OnResolveArgs) (result esbuild.OnResolveResult, err error) {
				result.Path = args.Path
				result.Namespace = httpNamespace
				return result, nil
			})

			// We also want to intercept all import paths inside downloaded files and
			// resolve them against the original URL. All of these files will be in
			// the "http-url" namespace. Make sure to keep the newly resolved URL in
			// the "http-url" namespace so imports inside it will also be resolved as
			// URLs recursively.
			build.OnResolve(esbuild.OnResolveOptions{Filter: ".*", Namespace: httpNamespace}, func(args esbuild.OnResolveArgs) (result esbuild.OnResolveResult, err error) {
				importer, ok := args.PluginData.(string)
				if !ok {
					return result, fmt.Errorf("expected plugin data for %q to be a string but got %v", args.Path, args.PluginData)
				}
				base, err := url.Parse(importer)
				if err != nil {
					return result, err
				}
				relative, err := url.Parse(args.Path)
				if err != nil {
					return result, err
				}
				result.Path = base.ResolveReference(relative).String()
				result.Namespace = httpNamespace
				return result, nil
			})

			// When a URL is loaded, we want to actually download the content from the
			// internet. We use plugin data to pass the final URL (after redirects)
			// back to the future resolvers.
			build.OnLoad(esbuild.OnLoadOptions{Filter: ".*", Namespace: httpNamespace}, func(args esbuild.OnLoadArgs) (result esbuild.OnLoadResult, err error) {
				res, err := client.Get(args.Path)
				if err != nil {
					return result, err
				}
				defer res.Body.Close()
				body, err := io.ReadAll(res.Body)
				if err != nil {
					return result, err
				}
				if res.StatusCode != 200 {
					return result, fmt.Errorf("unexpected response for GET %q: %s. %s", args.Path, res.Status, string(body))
				}
				contents := string(body)
				result.Contents = &contents
				result.ResolveDir = "/" + res.Request.URL.String()
				result.PluginData = res.Request.URL.String()
				return result, nil
			})
		},
	}
}
