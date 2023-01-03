package nodemodules

import (
	"bytes"
	"fmt"
	"path/filepath"
	"regexp"
	"strings"

	esbuild "github.com/evanw/esbuild/pkg/api"
	"github.com/livebud/bud/internal/esmeta"
	"github.com/livebud/bud/package/genfs"
	"github.com/livebud/bud/package/gomod"
)

func New(module *gomod.Module) *Generator {
	return &Generator{module}
}

type Generator struct {
	module *gomod.Module
}

// ServeFile serves node modules on demand
func (g *Generator) ServeFile(fsys genfs.FS, file *genfs.File) error {
	// If the name starts with node_modules, trim it to allow esbuild to do
	// the resolving. e.g. node_modules/timeago.js => timeago.js
	entryPoint := trimEntrypoint(file.Target())
	result := esbuild.Build(esbuild.BuildOptions{
		EntryPoints:   []string{entryPoint},
		AbsWorkingDir: g.module.Directory(),
		Format:        esbuild.FormatESModule,
		Platform:      esbuild.PlatformBrowser,
		// Add "import" condition to support svelte/internal
		// https://esbuild.github.io/api/#how-conditions-work
		Conditions: []string{"browser", "default", "import"},
		Metafile:   true,
		Bundle:     true,
		Plugins: []esbuild.Plugin{
			domExternalizePlugin(),
		},
	})
	if len(result.Errors) > 0 {
		msgs := esbuild.FormatMessages(result.Errors, esbuild.FormatMessagesOptions{
			Color:         true,
			Kind:          esbuild.ErrorMessage,
			TerminalWidth: 80,
		})
		return fmt.Errorf(strings.Join(msgs, "\n"))
	}
	content := result.OutputFiles[0].Contents
	// Replace require statements and updates the path on imports
	code := replaceDependencyPaths(content)
	file.Data = code
	// Link the dependencies
	metafile, err := esmeta.Parse(result.Metafile)
	if err != nil {
		return err
	}
	// Watch the dependencies for changes
	if err := fsys.Watch(metafile.Dependencies()...); err != nil {
		return err
	}
	return nil
}

// Transforms the dom file imports into including the "__LIVEBUD_EXTERNAL__:" prefix
// TODO: dedupe with dom
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

// TODO: dedupe with dom
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

// TODO: dedupe with dom
func isNodeModule(path string) bool {
	switch path[0] {
	case '.', '/', '\\':
		return false
	default:
		return true
	}
}

// TODO: dedupe with dom
var reImport = regexp.MustCompile(`([A-Z_a-z$][A-Z_a-z0-9]*)?\(?"(__LIVEBUD_EXTERNAL__:([^"]+))"\)?`)
var importBytes = []byte(`import`)

// This function rewrites require statements and updates the path on imports
// TODO: dedupe with dom
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

// TODO: dedupe with dom
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

// TODO: dedupe with dom
func importStatement(identifier, name string) string {
	return fmt.Sprintf(`import %s from "/bud/node_modules/%s"`+"\n", identifier, name)
}
