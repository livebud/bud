package esb

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"net/http"
	"path"
	"path/filepath"
	"strings"

	"github.com/livebud/bud/package/es"
	"github.com/livebud/bud/package/gomod"
	"github.com/livebud/bud/package/log"

	esbuild "github.com/evanw/esbuild/pkg/api"
	"github.com/livebud/bud/framework"
)

func New(log log.Log, module *gomod.Module) *Builder {
	return &Builder{log, module, http.DefaultTransport}
}

type Builder struct {
	log    log.Log
	module *gomod.Module
	rt     http.RoundTripper
}

// SSR creates a server-rendered preset
func SSR(flag *framework.Flag) esbuild.BuildOptions {
	options := esbuild.BuildOptions{
		Outdir:     "./",
		Format:     esbuild.FormatIIFE,
		Platform:   esbuild.PlatformNeutral,
		GlobalName: "bud",
		// Always bundle, use plugins to granularly mark files as external
		Bundle: true,
	}
	// Support minifying
	if flag.Minify {
		options.MinifyWhitespace = true
		options.MinifyIdentifiers = true
		options.MinifySyntax = true
	}
	return options
}

// DOM creates a browser preset
func DOM(flag *framework.Flag) esbuild.BuildOptions {
	options := esbuild.BuildOptions{
		Outdir:   "./",
		Format:   esbuild.FormatESModule,
		Platform: esbuild.PlatformBrowser,
		// Add "import" condition to support svelte/internal
		// https://esbuild.github.io/api/#how-conditions-work
		Conditions: []string{"browser", "default", "import"},
		// Always bundle, use plugins to granularly mark files as external
		Bundle: true,
	}
	// Support minifying
	if flag.Minify {
		options.MinifyWhitespace = true
		options.MinifyIdentifiers = true
		options.MinifySyntax = true
	}
	if flag.Embed {
		options.Splitting = true
	}
	return options
}

// Serve a single file
func (b *Builder) Serve(options esbuild.BuildOptions, entries ...string) (*esbuild.OutputFile, error) {
	options.EntryPoints = entries
	options.Plugins = append(options.Plugins, notFoundPlugin())
	// Run esbuild
	result := esbuild.Build(options)
	// Check if there were errors
	if result.Errors != nil {
		errs := make([]string, len(result.Errors))
		for i, err := range result.Errors {
			errs[i] = `es: ` + err.Text
		}
		return nil, errors.New(strings.Join(errs, "\n"))
	} else if len(result.OutputFiles) == 0 {
		return nil, fmt.Errorf("es: no output files")
	}
	// Return the first file
	file := result.OutputFiles[0]
	return &file, nil
}

// Bundle a group of files
func (b *Builder) Bundle(options esbuild.BuildOptions, entries ...string) ([]esbuild.OutputFile, error) {
	options.EntryPoints = entries
	options.Plugins = append(options.Plugins, notFoundPlugin())
	// Run esbuild
	result := esbuild.Build(options)
	// Check if there were errors
	if result.Errors != nil {
		errs := make([]string, len(result.Errors))
		for i, err := range result.Errors {
			errs[i] = `es: ` + err.Text
		}
		return nil, errors.New(strings.Join(errs, "\n"))
	}
	return result.OutputFiles, nil
}

func resolvePath(fsys fs.FS, fpath string) (string, error) {
	stat, err := fs.Stat(fsys, fpath)
	if err != nil {
		if !errors.Is(err, fs.ErrNotExist) {
			return "", fmt.Errorf("unable to stat %q. %w", fpath, err)
		}
		// Read the directory, it might be an extension-less import.
		// e.g. "react-dom/server" where "react-dom/server.js" exists
		dir := path.Dir(fpath)
		des, err := fs.ReadDir(fsys, dir)
		if err != nil {
			return "", fmt.Errorf("unable to read dir %q. %w", dir, err)
		}
		baseAndDot := path.Base(fpath) + "."
		for _, de := range des {
			if de.IsDir() {
				continue
			}
			name := de.Name()
			if !strings.HasPrefix(name, baseAndDot) {
				continue
			}
			switch path.Ext(name) {
			case ".js", ".mjs", ".cjs", ".jsx", ".ts", ".tsx", ".json":
				fpath = path.Join(dir, name)
				return fpath, nil
			}
		}
		return "", fmt.Errorf("unable to resolve %q. %w", fpath, fs.ErrNotExist)
	}
	// Handle reading the index file from a main directory
	// e.g. "react" resolving to "react/index.js"
	if stat.IsDir() {
		des, err := fs.ReadDir(fsys, fpath)
		if err != nil {
			return "", fmt.Errorf("unable to read dir %q. %w", fpath, err)
		}
		for _, de := range des {
			if de.IsDir() {
				continue
			}
			name := de.Name()
			if strings.HasPrefix(name, "index.") {
				switch path.Ext(name) {
				case ".js", ".mjs", ".cjs", ".jsx", ".ts", ".tsx", ".json":
					fpath = path.Join(fpath, name)
					return fpath, nil
				}
			}
		}
		return "", fmt.Errorf("unable to resolve %q. %w", fpath, fs.ErrNotExist)
	}
	return fpath, nil
}

type packageJSON struct {
	Module  string `json:"module"`
	Main    string `json:"main"`
	Browser string `json:"browser"`
}

func resolveNodeModule(fsys fs.FS, module string) (string, error) {
	fpath := path.Join("node_modules", module)
	sections := strings.SplitN(module, "/", 3)
	hasScope := strings.HasPrefix(sections[0], "@")
	if (hasScope && len(sections) > 2) || (!hasScope && len(sections) > 1) {
		return resolvePath(fsys, fpath)
	}
	pkgCode, err := fs.ReadFile(fsys, path.Join(fpath, "package.json"))
	if err != nil {
		return "", err
	}
	var pkg packageJSON
	if err := json.Unmarshal(pkgCode, &pkg); err != nil {
		return "", err
	}
	switch {
	case pkg.Module != "":
		fpath = path.Join(fpath, pkg.Module)
	case pkg.Main != "":
		fpath = path.Join(fpath, pkg.Main)
	case pkg.Browser != "":
		fpath = path.Join(fpath, pkg.Browser)
	}
	return resolvePath(fsys, fpath)
}

// type FSOption struct {
// 	FS fs.FS

func FS(fsys fs.FS, namespace string) esbuild.Plugin {
	return esbuild.Plugin{
		Name: "scope",
		Setup: func(build esbuild.PluginBuild) {
			// Resolve relative paths that start with a dot or slash
			build.OnResolve(esbuild.OnResolveOptions{Filter: `^\.{0,2}/`}, func(args esbuild.OnResolveArgs) (result esbuild.OnResolveResult, err error) {
				fpath := filepath.Clean(args.Path)
				// Resolve paths relative to the importer
				if args.Importer != "" {
					fpath = path.Join(path.Dir(args.Importer), fpath)
				}
				resolved, err := resolvePath(fsys, fpath)
				if err != nil {
					if errors.Is(err, fs.ErrNotExist) {
						// Continue to the next resolver
						return result, nil
					}
					return result, err
				}
				result.Namespace = namespace
				result.Path = resolved
				return result, nil
			})
			// Resolve node_modules
			build.OnResolve(esbuild.OnResolveOptions{Filter: `^[@a-z0-9]`}, func(args esbuild.OnResolveArgs) (result esbuild.OnResolveResult, err error) {
				resolved, err := resolveNodeModule(fsys, args.Path)
				if err != nil {
					if errors.Is(err, fs.ErrNotExist) {
						// Continue to the next resolver
						return result, nil
					}
					return result, err
				}
				result.Namespace = namespace
				result.Path = resolved
				return result, nil
			})
			// Load virtual files from the virtual filesystem
			build.OnLoad(esbuild.OnLoadOptions{Filter: ".*", Namespace: namespace}, func(args esbuild.OnLoadArgs) (result esbuild.OnLoadResult, err error) {
				code, err := fs.ReadFile(fsys, args.Path)
				if err != nil {
					return result, fmt.Errorf("unable to read file %q. %w", args.Path, err)
				}
				// Set the loader
				switch path.Ext(args.Path) {
				case ".js", ".mjs", ".cjs":
					result.Loader = esbuild.LoaderJS
				case ".jsx":
					result.Loader = esbuild.LoaderJSX
				case ".ts":
					result.Loader = esbuild.LoaderTS
				case ".tsx":
					result.Loader = esbuild.LoaderTSX
				case ".json":
					result.Loader = esbuild.LoaderJSON
				default:
					return result, fmt.Errorf("unknown loader for %q", args.Path)
				}
				// Resolve the directory
				result.ResolveDir = path.Dir(args.Path)
				// Set the contents
				contents := string(code)
				result.Contents = &contents
				// Return the result
				return result, nil
			})
		},
	}
}

func notFoundPlugin() es.Plugin {
	return esbuild.Plugin{
		Name: "not_found",
		Setup: func(build esbuild.PluginBuild) {
			// Resolver that matches anything and just fails to avoid esbuild
			// resolving with AbsWorkingDir.
			build.OnResolve(esbuild.OnResolveOptions{Filter: ".*"}, func(args esbuild.OnResolveArgs) (result esbuild.OnResolveResult, err error) {
				return result, fmt.Errorf("unable to find %q", args.Path)
			})
		},
	}
}
