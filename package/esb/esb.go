package esb

import (
	"errors"
	"fmt"
	"io/fs"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"

	esbuild "github.com/evanw/esbuild/pkg/api"
	"github.com/livebud/bud/framework"
	"github.com/livebud/bud/internal/once"
)

// SSR creates a server-rendered preset
func SSR(flag *framework.Flag, entries ...string) esbuild.BuildOptions {
	options := esbuild.BuildOptions{
		EntryPoints: entries,
		Outdir:      "./",
		Format:      esbuild.FormatIIFE,
		Platform:    esbuild.PlatformNeutral,
		GlobalName:  "bud",
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
func DOM(flag *framework.Flag, entries ...string) esbuild.BuildOptions {
	options := esbuild.BuildOptions{
		EntryPoints: entries,
		Outdir:      "./",
		Format:      esbuild.FormatESModule,
		Platform:    esbuild.PlatformBrowser,
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

// Only create the temp dir once
var tempDir = func() func(prefix string) (string, error) {
	var once once.String
	return func(prefix string) (string, error) {
		return once.Do(func() (string, error) {
			return os.MkdirTemp("", prefix)
		})
	}
}()

// Serve a single file
func Serve(transport http.RoundTripper, fsys fs.FS, options esbuild.BuildOptions) (*esbuild.OutputFile, error) {
	// Build from a scratch directory to reduce file-system influence
	dir, err := tempDir("bud_esb_*")
	if err != nil {
		return nil, fmt.Errorf("es: unable to create scratch dir. %w", err)
	}
	options.AbsWorkingDir = dir
	options.Plugins = append(options.Plugins, virtualPlugin(fsys))
	// Run esbuild
	result := esbuild.Build(options)
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

// Bundle a group of files
func Bundle(options esbuild.BuildOptions) ([]esbuild.OutputFile, error) {
	// Run esbuild
	result := esbuild.Build(options)
	// Check if there were errors
	if result.Errors != nil {
		errors := esbuild.FormatMessages(result.Errors, esbuild.FormatMessagesOptions{
			Kind: esbuild.ErrorMessage,
		})
		return nil, fmt.Errorf("es: %s", strings.Join(errors, "\n"))
	}
	return result.OutputFiles, nil
}

func resolvePath(fsys fs.FS, fpath string) (string, error) {
	stat, err := fs.Stat(fsys, fpath)
	if err != nil {
		if !errors.Is(err, fs.ErrNotExist) {
			return "", fmt.Errorf("es: unable to stat %q. %w", fpath, err)
		}
		// Read the directory, it might be an extension-less import.
		// e.g. "react-dom/server" where "react-dom/server.js" exists
		dir := path.Dir(fpath)
		des, err := fs.ReadDir(fsys, dir)
		if err != nil {
			return "", fmt.Errorf("es: unable to read dir %q. %w", dir, err)
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
		return "", fmt.Errorf("es: unable to resolve %q. %w", fpath, err)
	}
	// Handle reading the index file from a main directory
	// e.g. "react" resolving to "react/index.js"
	if stat.IsDir() {
		des, err := fs.ReadDir(fsys, fpath)
		if err != nil {
			return "", fmt.Errorf("es: unable to read dir %q. %w", fpath, err)
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
		return "", fmt.Errorf("es: unable to resolve %q. %w", fpath, err)
	}
	return fpath, nil
}

func virtualPlugin(fsys fs.FS) esbuild.Plugin {
	return esbuild.Plugin{
		Name: "virtual",
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
				result.Namespace = "virtual"
				result.Path = resolved
				return result, nil
			})
			// Resolve node_modules
			build.OnResolve(esbuild.OnResolveOptions{Filter: `^[@a-z0-9]`}, func(args esbuild.OnResolveArgs) (result esbuild.OnResolveResult, err error) {
				fpath := path.Join("node_modules", args.Path)
				resolved, err := resolvePath(fsys, fpath)
				if err != nil {
					if errors.Is(err, fs.ErrNotExist) {
						// Continue to the next resolver
						return result, nil
					}
					return result, err
				}
				result.Namespace = "virtual"
				result.Path = resolved
				return result, nil
			})
			// Load virtual files from the virtual filesystem
			build.OnLoad(esbuild.OnLoadOptions{Filter: ".*", Namespace: "virtual"}, func(args esbuild.OnLoadArgs) (result esbuild.OnLoadResult, err error) {
				code, err := fs.ReadFile(fsys, args.Path)
				if err != nil {
					return result, fmt.Errorf("es: unable to read file %q. %w", args.Path, err)
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
					return result, fmt.Errorf("es: unknown loader for %q", args.Path)
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
