package ssr

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	_ "embed"

	esbuild "github.com/evanw/esbuild/pkg/api"
	"gitlab.com/mnm/bud/bfs"
	"gitlab.com/mnm/bud/internal/entrypoint"
	"gitlab.com/mnm/bud/internal/gotemplate"
	"gitlab.com/mnm/bud/transform"
)

func Generator(osfs fs.FS, dir string, transformer *transform.Transformer) bfs.Generator {
	plugins := append([]esbuild.Plugin{
		ssrPlugin(osfs, dir),
		ssrRuntimePlugin(osfs, dir),
		jsxPlugin(osfs, dir),
		jsxRuntimePlugin(osfs, dir),
		jsxTransformPlugin(osfs, dir),
		sveltePlugin(osfs, dir),
		svelteRuntimePlugin(osfs, dir),
	}, transformer.Node.Plugins()...)
	return bfs.GenerateFile(func(f bfs.FS, file *bfs.File) error {
		result := esbuild.Build(esbuild.BuildOptions{
			EntryPointsAdvanced: []esbuild.EntryPoint{
				{
					InputPath:  "./bud/view/_ssr.js",
					OutputPath: "./bud/view/_ssr",
				},
			},
			AbsWorkingDir: dir,
			Outdir:        "./",
			Format:        esbuild.FormatIIFE,
			Platform:      esbuild.PlatformBrowser,
			GlobalName:    "bud",
			JSXFactory:    "__budReact__.createElement",
			JSXFragment:   "__budReact__.Fragment",
			Bundle:        true,
			Metafile:      true,
			Plugins:       plugins,
		})
		if len(result.Errors) > 0 {
			msgs := esbuild.FormatMessages(result.Errors, esbuild.FormatMessagesOptions{
				Color:         true,
				Kind:          esbuild.ErrorMessage,
				TerminalWidth: 80,
			})
			return fmt.Errorf(strings.Join(msgs, "\n"))
		}
		// Expect exactly 1 output file
		if len(result.OutputFiles) != 1 {
			return fmt.Errorf("expected exactly 1 output file but got %d", len(result.OutputFiles))
		}
		// if err := esmeta.Link2(dfs, result.Metafile); err != nil {
		// 	return nil, err
		// }
		// TODO: remove WriteEvent and externalize actual file contents so we only
		// need to watch directory changes.
		file.Watch("bud/view/**/*.{svelte,jsx}", bfs.CreateEvent|bfs.RemoveEvent|bfs.WriteEvent)
		file.Write(result.OutputFiles[0].Contents)
		return nil
	})
}

//go:embed ssr.gotext
var ssrTemplate string

// ssrGenerator
var ssrGenerator = gotemplate.MustParse("ssr.gotext", ssrTemplate)

// Generate the bud/view/_ssr.js file
func ssrPlugin(osfs fs.FS, dir string) esbuild.Plugin {
	return esbuild.Plugin{
		Name: "ssr",
		Setup: func(epb esbuild.PluginBuild) {
			epb.OnResolve(esbuild.OnResolveOptions{Filter: `^\.\/bud\/view\/_ssr.js$`}, func(args esbuild.OnResolveArgs) (result esbuild.OnResolveResult, err error) {
				result.Namespace = "ssr"
				result.Path = args.Path
				return result, nil
			})
			epb.OnLoad(esbuild.OnLoadOptions{Filter: `.*`, Namespace: "ssr"}, func(args esbuild.OnLoadArgs) (result esbuild.OnLoadResult, err error) {
				views, err := entrypoint.List(osfs)
				if err != nil {
					return result, err
				}
				code, err := ssrGenerator.Generate(map[string]interface{}{
					"Views": views,
				})
				if err != nil {
					return result, err
				}
				contents := string(code)
				result.ResolveDir = dir
				result.Contents = &contents
				result.Loader = esbuild.LoaderJS
				return result, nil
			})
		},
	}
}

//go:embed ssr.ts
var ssrRuntime string

// Generate the bud/view/_ssr_runtime.ts file imported in bud/view/_ssr.js
func ssrRuntimePlugin(osfs fs.FS, dir string) esbuild.Plugin {
	return esbuild.Plugin{
		Name: "ssr_runtime",
		Setup: func(epb esbuild.PluginBuild) {
			epb.OnResolve(esbuild.OnResolveOptions{Filter: `^./bud/view/_ssr_runtime.ts$`}, func(args esbuild.OnResolveArgs) (result esbuild.OnResolveResult, err error) {
				result.Namespace = "ssr_runtime"
				result.Path = args.Path
				return result, nil
			})
			epb.OnLoad(esbuild.OnLoadOptions{Filter: `.*`, Namespace: "ssr_runtime"}, func(args esbuild.OnLoadArgs) (result esbuild.OnLoadResult, err error) {
				result.ResolveDir = dir
				result.Contents = &ssrRuntime
				result.Loader = esbuild.LoaderTS
				return result, nil
			})
		},
	}
}

//go:embed jsx.gotext
var jsxTemplate string

var jsxGenerator = gotemplate.MustParse("jsx.gotext", jsxTemplate)

// Generate the jsx entry file: bud/view/$page.jsx
func jsxPlugin(osfs fs.FS, dir string) esbuild.Plugin {
	return esbuild.Plugin{
		Name: "jsx",
		Setup: func(epb esbuild.PluginBuild) {
			epb.OnResolve(esbuild.OnResolveOptions{Filter: `^\./bud/view/.*\.jsx$`}, func(args esbuild.OnResolveArgs) (result esbuild.OnResolveResult, err error) {
				result.Path = args.Path
				result.Namespace = "jsx"
				return result, nil
			})
			epb.OnLoad(esbuild.OnLoadOptions{Filter: `.*`, Namespace: "jsx"}, func(args esbuild.OnLoadArgs) (result esbuild.OnLoadResult, err error) {
				view, err := entrypoint.FindByPage(osfs, strings.Trim(filepath.Clean(args.Path), "bud/"))
				if err != nil {
					return result, err
				}
				code, err := jsxGenerator.Generate(view)
				if err != nil {
					return result, err
				}
				contents := string(code)
				result.ResolveDir = dir
				result.Contents = &contents
				result.Loader = esbuild.LoaderJSX
				return result, nil
			})
		},
	}
}

//go:embed jsx.ts
var jsxRuntime string

// Generate the jsx runtime for the entry files
func jsxRuntimePlugin(osfs fs.FS, dir string) esbuild.Plugin {
	return esbuild.Plugin{
		Name: "jsx_runtime",
		Setup: func(epb esbuild.PluginBuild) {
			epb.OnResolve(esbuild.OnResolveOptions{Filter: `^\./bud/view/_jsx\.ts$`}, func(args esbuild.OnResolveArgs) (result esbuild.OnResolveResult, err error) {
				result.Path = args.Path
				result.Namespace = "jsx_runtime"
				return result, nil
			})
			epb.OnLoad(esbuild.OnLoadOptions{Filter: `.*`, Namespace: "jsx_runtime"}, func(args esbuild.OnLoadArgs) (result esbuild.OnLoadResult, err error) {
				result.ResolveDir = dir
				result.Contents = &jsxRuntime
				result.Loader = esbuild.LoaderTS
				return result, nil
			})
		},
	}
}

func jsxTransformPlugin(osfs fs.FS, dir string) esbuild.Plugin {
	return esbuild.Plugin{
		Name: "jsx_transform",
		Setup: func(epb esbuild.PluginBuild) {
			// Load jsx files. Add import if not present
			epb.OnLoad(esbuild.OnLoadOptions{Filter: `\.jsx$`}, func(args esbuild.OnLoadArgs) (result esbuild.OnLoadResult, err error) {
				code, err := os.ReadFile(args.Path)
				if err != nil {
					return result, err
				}
				contents := string(code)
				contents = `import * as __budReact__ from "react"` + "\n\n" + contents
				result.ResolveDir = filepath.Dir(args.Path)
				result.Contents = &contents
				result.Loader = esbuild.LoaderJSX
				return result, nil
			})
		},
	}
}

//go:embed svelte.gotext
var svelteTemplate string

var svelteGenerator = gotemplate.MustParse("svelte.gotext", svelteTemplate)

// Generate the svelte entry file: bud/view/$page.svelte
func sveltePlugin(osfs fs.FS, dir string) esbuild.Plugin {
	return esbuild.Plugin{
		Name: "svelte",
		Setup: func(epb esbuild.PluginBuild) {
			epb.OnResolve(esbuild.OnResolveOptions{Filter: `^\./bud/view/.*\.svelte$`}, func(args esbuild.OnResolveArgs) (result esbuild.OnResolveResult, err error) {
				result.Path = args.Path
				result.Namespace = "svelte"
				return result, nil
			})
			epb.OnLoad(esbuild.OnLoadOptions{Filter: `.*`, Namespace: "svelte"}, func(args esbuild.OnLoadArgs) (result esbuild.OnLoadResult, err error) {
				view, err := entrypoint.FindByPage(osfs, strings.Trim(filepath.Clean(args.Path), "bud/"))
				if err != nil {
					return result, err
				}
				code, err := svelteGenerator.Generate(view)
				if err != nil {
					return result, err
				}
				contents := string(code)
				result.ResolveDir = dir
				result.Contents = &contents
				result.Loader = esbuild.LoaderJSX
				return result, nil
			})
		},
	}
}

//go:embed svelte.ts
var svelteRuntime string

// Generate the svelte runtime for the entry files
func svelteRuntimePlugin(osfs fs.FS, dir string) esbuild.Plugin {
	return esbuild.Plugin{
		Name: "svelte_runtime",
		Setup: func(epb esbuild.PluginBuild) {
			epb.OnResolve(esbuild.OnResolveOptions{Filter: `^\./bud/view/_svelte\.ts$`}, func(args esbuild.OnResolveArgs) (result esbuild.OnResolveResult, err error) {
				result.Path = args.Path
				result.Namespace = "svelte_runtime"
				return result, nil
			})
			epb.OnLoad(esbuild.OnLoadOptions{Filter: `.*`, Namespace: "svelte_runtime"}, func(args esbuild.OnLoadArgs) (result esbuild.OnLoadResult, err error) {
				result.ResolveDir = dir
				result.Contents = &svelteRuntime
				result.Loader = esbuild.LoaderTS
				return result, nil
			})
		},
	}
}
