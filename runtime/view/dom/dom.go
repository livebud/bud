package dom

import (
	"bytes"
	"context"
	_ "embed"
	"fmt"
	"io/fs"
	"path/filepath"
	"regexp"
	"strings"

	"gitlab.com/mnm/bud/package/overlay"

	esbuild "github.com/evanw/esbuild/pkg/api"
	"gitlab.com/mnm/bud/internal/entrypoint"
	"gitlab.com/mnm/bud/internal/gotemplate"
	"gitlab.com/mnm/bud/package/gomod"
	"gitlab.com/mnm/bud/runtime/transform"
)

//go:embed dom.gotext
var template string

// generator
var generator = gotemplate.MustParse("dom.gotext", template)

func Runner(fsys fs.FS, module *gomod.Module, transformer *transform.Transformer) overlay.FileServer {
	plugins := append([]esbuild.Plugin{
		domPlugin(fsys, module),
		domExternalizePlugin(),
	}, transformer.Browser.Plugins()...)
	return overlay.ServeFile(func(ctx context.Context, f overlay.F, file *overlay.File) error {
		// If the name starts with node_modules, trim it to allow esbuild to do
		// the resolving. e.g. node_modules/livebud => livebud
		entryPoint := trimEntrypoint(file.Path())
		result := esbuild.Build(esbuild.BuildOptions{
			EntryPoints:   []string{entryPoint},
			AbsWorkingDir: module.Directory(),
			Format:        esbuild.FormatESModule,
			Platform:      esbuild.PlatformBrowser,
			// Add "import" condition to support svelte/internal
			// https://esbuild.github.io/api/#how-conditions-work
			Conditions: []string{"browser", "default", "import"},
			Metafile:   true,
			Bundle:     true,
			Plugins:    plugins,
		})
		if len(result.Errors) > 0 {
			msgs := esbuild.FormatMessages(result.Errors, esbuild.FormatMessagesOptions{
				Color:         true,
				Kind:          esbuild.ErrorMessage,
				TerminalWidth: 80,
			})
			return fmt.Errorf(strings.Join(msgs, "\n"))
		}
		// if err := esmeta.Link2(dfs, result.Metafile); err != nil {
		// 	return nil, err
		// }
		code := result.OutputFiles[0].Contents
		// Replace require statements and updates the path on imports
		code = replaceDependencyPaths(code)
		file.Data = code
		source := strings.TrimPrefix(file.Path(), "bud/")
		file.Link(source)
		return nil
	})
}

func NodeModules(module *gomod.Module) overlay.FileServer {
	plugins := []esbuild.Plugin{
		domExternalizePlugin(),
	}
	return overlay.ServeFile(func(ctx context.Context, f overlay.F, file *overlay.File) error {
		// If the name starts with node_modules, trim it to allow esbuild to do
		// the resolving. e.g. node_modules/timeago.js => timeago.js
		entryPoint := trimEntrypoint(file.Path())
		result := esbuild.Build(esbuild.BuildOptions{
			EntryPoints:   []string{entryPoint},
			AbsWorkingDir: module.Directory(),
			Format:        esbuild.FormatESModule,
			Platform:      esbuild.PlatformBrowser,
			// Add "import" condition to support svelte/internal
			// https://esbuild.github.io/api/#how-conditions-work
			Conditions: []string{"browser", "default", "import"},
			Metafile:   true,
			Bundle:     true,
			Plugins:    plugins,
		})
		if len(result.Errors) > 0 {
			msgs := esbuild.FormatMessages(result.Errors, esbuild.FormatMessagesOptions{
				Color:         true,
				Kind:          esbuild.ErrorMessage,
				TerminalWidth: 80,
			})
			return fmt.Errorf(strings.Join(msgs, "\n"))
		}
		// if err := esmeta.Link2(dfs, result.Metafile); err != nil {
		// 	return nil, err
		// }
		content := result.OutputFiles[0].Contents
		// Replace require statements and updates the path on imports
		code := replaceDependencyPaths(content)
		file.Data = code
		source := strings.TrimPrefix(file.Path(), "bud/")
		// fmt.Println("linked", file.Path(), "->", source)
		file.Link(source)
		return nil
	})
}

func Builder(fsys fs.FS, module *gomod.Module, transformer *transform.Transformer) overlay.DirGenerator {
	plugins := append([]esbuild.Plugin{
		domPlugin(fsys, module),
	}, transformer.Browser.Plugins()...)
	return overlay.GenerateDir(func(ctx context.Context, f overlay.F, dir *overlay.Dir) error {
		views, err := entrypoint.List(fsys)
		if err != nil {
			return err
		}
		entries := make([]esbuild.EntryPoint, len(views))
		viewDir := filepath.Join("bud", "view") + string(filepath.Separator)
		for i, view := range views {
			entryPath := filepath.Join("bud", toEntry(string(view.Page)))
			outPath := strings.TrimPrefix(entryPath, viewDir)
			entries[i] = esbuild.EntryPoint{
				InputPath:  entryPath,
				OutputPath: outPath,
			}
		}
		// If the name starts with node_modules, trim it to allow esbuild to do
		// the resolving. e.g. node_modules/livebud => livebud
		result := esbuild.Build(esbuild.BuildOptions{
			EntryPointsAdvanced: entries,
			Outdir:              "/",
			AbsWorkingDir:       module.Directory(),
			ChunkNames:          "[name]-[hash]",
			Format:              esbuild.FormatESModule,
			Platform:            esbuild.PlatformBrowser,
			// Add "import" condition to support svelte/internal
			// https://esbuild.github.io/api/#how-conditions-work
			Conditions:        []string{"browser", "default", "import"},
			Metafile:          false,
			Bundle:            true,
			Splitting:         true,
			MinifyIdentifiers: true,
			MinifySyntax:      true,
			MinifyWhitespace:  true,
			Plugins:           plugins,
			Write:             false,
		})
		if len(result.Errors) > 0 {
			msgs := esbuild.FormatMessages(result.Errors, esbuild.FormatMessagesOptions{
				Color:         true,
				Kind:          esbuild.ErrorMessage,
				TerminalWidth: 80,
			})
			return fmt.Errorf(strings.Join(msgs, "\n"))
		}
		for _, outFile := range result.OutputFiles {
			outFile := outFile
			outPath := strings.TrimPrefix(outFile.Path, "/")
			if isEntry(outPath) {
				outPath = strings.TrimSuffix(outPath, ".js")
			}
			dir.GenerateFile(outPath, func(ctx context.Context, f overlay.F, file *overlay.File) error {
				file.Data = outFile.Contents
				return nil
			})
		}
		return nil
	})
}

func toEntry(path string) string {
	dir, base := filepath.Split(path)
	return filepath.Join(dir, "_"+base)
}

func isEntry(path string) bool {
	base := filepath.Base(path)
	return base[0] == '_'
}

func trimEntrypoint(path string) string {
	// Trim up node_modules so esbuild can resolve them, yet they're valid url
	// paths on the frontend.
	// e.g.
	//   /bud/node_modules/livebud/hot => livebud/hot
	//   /bud/node_modules/react => react
	if strings.HasPrefix(path, "bud/node_modules") {
		return strings.TrimPrefix(path, "bud/node_modules/")
	}
	// If the basepath starts with an underscore it could be the entrypoint
	if filepath.Base(path)[0] == '_' {
		return path
	}
	// Trim up /bud from the path so we can map to a valid underlying view file
	// e.g. bud/view/new.js => view/new.js
	if strings.HasPrefix(path, "bud/view") {
		return strings.TrimPrefix(path, "bud/")
	}
	return path
}

// Build the bud/view/$page.{jsx,svelte} client-side entrypoint
func domPlugin(fsys fs.FS, module *gomod.Module) esbuild.Plugin {
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
				code, err := generator.Generate(view)
				if err != nil {
					return result, err
				}
				contents := string(code)
				result.ResolveDir = module.Directory()
				result.Contents = &contents
				result.Loader = esbuild.LoaderJS
				return result, nil
			})
		},
	}
}

// Transforms the dom file imports into including the "__LIVEBUD_EXTERNAL__:" prefix
func domExternalizePlugin() esbuild.Plugin {
	return esbuild.Plugin{
		Name: "dom_resolver",
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

var reImport = regexp.MustCompile(`([A-Z_a-z$][A-Z_a-z0-9]*)?\(?"(__LIVEBUD_EXTERNAL__:([^"]+))"\)?`)
var importBytes = []byte(`import`)

// This function rewrites require statements and updates the path on imports
func replaceDependencyPaths(content []byte) []byte {
	identifiers := map[string]bool{}
	out := new(bytes.Buffer)
	code := new(bytes.Buffer)
	since := 0
	// Submatches: [
	//  (0) matchStart,
	//  (1) matchEnd,
	//  (2) requireOrImportStart,
	//  (3) requireOrImportEnd,
	//  (4) modulePathStart,
	//  (5) modulePathEnd,
	//  (6) moduleNameStart,
	//  (7) moduleNameEnd,
	// ]
	for _, submatches := range reImport.FindAllSubmatchIndex(content, -1) {
		// Write the bytes since the last match
		code.Write(content[since:submatches[0]])
		// Update since with the end of the match
		since = submatches[1]
		// Get the path of the node module
		path := string(content[submatches[6]:submatches[7]])
		// Handle require(...) or import(...)
		var importOrRequire []byte
		if submatches[2] >= 0 && submatches[3] >= 0 {
			importOrRequire = content[submatches[2]:submatches[3]]
		}
		// We have a require(...), replace the whole expression
		if importOrRequire != nil && !bytes.Equal(importOrRequire, importBytes) {
			identifier := "__" + toIdentifier(path) + "$"
			code.WriteString(identifier)
			// Only add this import if we haven't seen this identifier yet
			if !identifiers[identifier] {
				out.WriteString(importStatement(identifier, path))
				identifiers[identifier] = true
			}
			continue
		}
		// Otherwise, we'll just replace the path
		code.Write(content[submatches[0]:submatches[4]])
		code.WriteString("/bud/node_modules/" + path)
		code.Write(content[submatches[5]:submatches[1]])
	}
	// Write the remaining bytes
	code.Write(content[since:])
	// Write code to out
	out.Write(code.Bytes())
	return out.Bytes()
}

func toIdentifier(importPath string) string {
	p := []byte(importPath)
	for i, c := range p {
		switch c {
		case '/', '-', '@', '.':
			p[i] = '_'
		default:
			p[i] = c
		}
	}
	return string(p)
}

func importStatement(identifier, name string) string {
	return fmt.Sprintf(`import %s from "/bud/node_modules/%s"`+"\n", identifier, name)
}
