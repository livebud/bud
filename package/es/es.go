package es

import (
	"fmt"
	"path"
	"strings"

	esbuild "github.com/evanw/esbuild/pkg/api"
	"github.com/livebud/bud/framework"
	"github.com/livebud/bud/package/gomod"
	"github.com/livebud/bud/package/log"
)

type Builder interface {
	Serve(serve *Serve) (*File, error)
	Bundle(bundle *Bundle) ([]File, error)
}

func New(flag *framework.Flag, log log.Log, module *gomod.Module) Builder {
	return &builder{flag, module}
}

type builder struct {
	flag   *framework.Flag
	module *gomod.Module
}

var _ Builder = (*builder)(nil)

type File = esbuild.OutputFile
type Plugin = esbuild.Plugin

type Platform uint8

const (
	DOM Platform = iota
	SSR
)

type Serve struct {
	Entry    string
	Plugins  []esbuild.Plugin
	Platform Platform
}

func (b *builder) serveOptions(serve *Serve) esbuild.BuildOptions {
	switch serve.Platform {
	case DOM:
		return b.dom([]string{serve.Entry}, serve.Plugins)
	default:
		return b.ssr([]string{serve.Entry}, serve.Plugins)
	}
}

var ErrNotRelative = fmt.Errorf("es: entry must be relative")

func (b *builder) Serve(serve *Serve) (*File, error) {
	if !isRelativeEntry(serve.Entry) {
		return nil, fmt.Errorf("%w %q", ErrNotRelative, serve.Entry)
	}
	// Externalize dependencies in development on non-dependency entries
	if serve.Platform == DOM && !b.flag.Embed && !isNodeModuleEntry(serve.Entry) {
		serve.Plugins = append(serve.Plugins, domExternalize())
	}
	result := esbuild.Build(b.serveOptions(serve))
	// Check if there were errors
	if result.Errors != nil {
		errors := esbuild.FormatMessages(result.Errors, esbuild.FormatMessagesOptions{
			Kind: esbuild.ErrorMessage,
		})
		return nil, fmt.Errorf("es: %s", strings.Join(errors, "\n"))
	} else if len(result.OutputFiles) == 0 {
		return nil, fmt.Errorf("es: no output files")
	}
	// Return the first file
	file := result.OutputFiles[0]
	return &file, nil
}

type Bundle struct {
	Entries  []string
	Plugins  []esbuild.Plugin
	Platform Platform
}

func (b *builder) bundleOptions(bundle *Bundle) esbuild.BuildOptions {
	switch bundle.Platform {
	case DOM:
		return b.dom(bundle.Entries, bundle.Plugins)
	default:
		return b.ssr(bundle.Entries, bundle.Plugins)
	}
}

func (b *builder) Bundle(bundle *Bundle) ([]File, error) {
	for _, entry := range bundle.Entries {
		if !isRelativeEntry(entry) {
			return nil, fmt.Errorf("%w %q", ErrNotRelative, entry)
		}
	}
	result := esbuild.Build(b.bundleOptions(bundle))
	// Check if there were errors
	if result.Errors != nil {
		errors := esbuild.FormatMessages(result.Errors, esbuild.FormatMessagesOptions{
			Kind: esbuild.ErrorMessage,
		})
		return nil, fmt.Errorf("es: %s", strings.Join(errors, "\n"))
	} else if len(result.OutputFiles) == 0 {
		return nil, fmt.Errorf("es: no output files")
	}
	// Return the first file
	return result.OutputFiles, nil
}

const outDir = "./"
const globalName = "bud"

// SSR creates a server-rendered preset
func (b *builder) ssr(entries []string, plugins []esbuild.Plugin) esbuild.BuildOptions {
	options := esbuild.BuildOptions{
		EntryPoints:   entries,
		Plugins:       plugins,
		AbsWorkingDir: b.module.Directory(),
		Outdir:        outDir,
		Format:        esbuild.FormatIIFE,
		Platform:      esbuild.PlatformNeutral,
		GlobalName:    globalName,
		// Always bundle, use plugins to granularly mark files as external
		Bundle: true,
	}
	if b.flag.Minify {
		options.MinifyWhitespace = true
		options.MinifyIdentifiers = true
		options.MinifySyntax = true
	}
	return options
}

// DOM creates a dom-rendered preset
func (b *builder) dom(entries []string, plugins []esbuild.Plugin) esbuild.BuildOptions {
	options := esbuild.BuildOptions{
		EntryPoints:   entries,
		Plugins:       plugins,
		AbsWorkingDir: b.module.Directory(),
		Outdir:        outDir,
		Format:        esbuild.FormatESModule,
		Platform:      esbuild.PlatformBrowser,
		// Add "import" condition to support svelte/internal
		// https://esbuild.github.io/api/#how-conditions-work
		Conditions: []string{"browser", "default", "import"},
		// Always bundle, use plugins to granularly mark files as external
		Bundle: true,
	}
	// Support minifying
	if b.flag.Minify {
		options.MinifyWhitespace = true
		options.MinifyIdentifiers = true
		options.MinifySyntax = true
	}
	if b.flag.Embed {
		options.Splitting = true
	}
	return options
}

func domExternalize() esbuild.Plugin {
	return esbuild.Plugin{
		Name: "dom_external_modules",
		Setup: func(epb esbuild.PluginBuild) {
			epb.OnResolve(esbuild.OnResolveOptions{Filter: ".*"}, func(args esbuild.OnResolveArgs) (result esbuild.OnResolveResult, err error) {
				// Externalize node modules
				if args.Importer != "" && isNodeModule(args.Path) {
					result.Path = "/" + path.Join("node_modules", args.Path)
					result.External = true
					return result, nil
				}
				// Don't externalize the entry file or any local files
				return result, nil
			})
		},
	}
}

func isNodeModuleEntry(importPath string) bool {
	importPath = path.Clean(importPath)
	if len(importPath) == 0 {
		return false
	}
	if importPath[0] == '/' {
		return strings.HasPrefix(importPath, "/node_modules/")
	}
	return strings.HasPrefix(importPath, "node_modules/")
}

func isRelativeEntry(entry string) bool {
	return strings.HasPrefix(entry, "./")
}

func isNodeModule(path string) bool {
	switch path[0] {
	case '.', '/', '\\':
		return false
	default:
		return true
	}
}
