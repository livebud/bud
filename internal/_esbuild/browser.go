package esbuild

import (
	"io/fs"
	"os"
	"path/filepath"

	esbuild "github.com/evanw/esbuild/pkg/api"
	"github.com/go-duo/bud/internal/entrypoint"
	"github.com/go-duo/bud/internal/gotemplate"
	"github.com/go-duo/bud/svelte"
)

// EntrypointDOM is a plugin for loading client-side entrypoints
// e.g. bud/view/$page.{jsx,svelte}
func BrowserEntryPlugin(fsys fs.FS, dir string, template gotemplate.Template) esbuild.Plugin {
	return esbuild.Plugin{
		Name: "dom",
		Setup: func(epb esbuild.PluginBuild) {
			epb.OnResolve(esbuild.OnResolveOptions{Filter: `^bud\/view\/(?:[A-Za-z\-0-9]+\/)*_[A-Za-z\-0-9]+\.(svelte|jsx)$`}, func(args esbuild.OnResolveArgs) (result esbuild.OnResolveResult, err error) {
				result.Namespace = "dom"
				result.Path = args.Path
				return result, nil
			})
			epb.OnLoad(esbuild.OnLoadOptions{Filter: `.*`, Namespace: "dom"}, func(args esbuild.OnLoadArgs) (result esbuild.OnLoadResult, err error) {
				view, err := entrypoint.FindByClient(fsys, filepath.Clean(args.Path))
				if err != nil {
					return result, err
				}
				data, err := template.Generate(view)
				if err != nil {
					return result, err
				}
				contents := string(data)
				result.ResolveDir = dir
				result.Contents = &contents
				result.Loader = esbuild.LoaderJS
				return result, nil
			})
		},
	}
}

// Transforms the dom file imports into including the "__LIVEBUD_EXTERNAL__:" prefix
func BrowserExternalPlugin() esbuild.Plugin {
	return esbuild.Plugin{
		Name: "externalize-dom",
		Setup: func(epb esbuild.PluginBuild) {
			epb.OnResolve(esbuild.OnResolveOptions{Filter: ".*"}, func(args esbuild.OnResolveArgs) (result esbuild.OnResolveResult, err error) {
				// Externalize node modules
				if args.Importer != "" && isNodeModule(args.Path) {
					result.Path = "__LIVEBUD_EXTERNAL__:" + args.Path
					result.External = true
					return result, nil
				}
				// Don't externalize the entry file or any local files
				return result, nil
			})
		},
	}
}

func isNodeModule(path string) bool {
	switch path[0] {
	case '.', '/', '\\':
		return false
	default:
		return true
	}
}

// Transform svelte files
func BrowserSvelteToJSPlugin(svelte *svelte.Compiler) esbuild.Plugin {
	return esbuild.Plugin{
		Name: "svelte_transform",
		Setup: func(epb esbuild.PluginBuild) {
			// Load svelte files. Add import if not present
			epb.OnLoad(esbuild.OnLoadOptions{Filter: `\.svelte$`}, func(args esbuild.OnLoadArgs) (result esbuild.OnLoadResult, err error) {
				code, err := os.ReadFile(args.Path)
				if err != nil {
					return result, err
				}
				dom, err := svelte.DOM(args.Path, code)
				if err != nil {
					return result, err
				}
				result.ResolveDir = filepath.Dir(args.Path)
				result.Contents = &dom.JS
				result.Loader = esbuild.LoaderJS
				return result, nil
			})
		},
	}
}
