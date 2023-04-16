package es

import (
	"encoding/json"
	"fmt"
	"strings"

	esbuild "github.com/evanw/esbuild/pkg/api"
	"github.com/livebud/bud/framework"
	"github.com/livebud/bud/package/log"
)

type Builder interface {
	Serve(serve *Serve) (*File, error)
	Bundle(bundle *Bundle) ([]File, error)
}

func New(flag *framework.Flag, log log.Log) Builder {
	return &builder{flag, log}
}

type builder struct {
	flag *framework.Flag
	log  log.Log
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
	AbsDir   string
	Entry    string
	Plugins  []esbuild.Plugin
	Platform Platform
}

func (b *builder) serveOptions(serve *Serve) esbuild.BuildOptions {
	switch serve.Platform {
	case DOM:
		return b.dom(serve.AbsDir, []string{serve.Entry}, serve.Plugins)
	default:
		return b.ssr(serve.AbsDir, []string{serve.Entry}, serve.Plugins)
	}
}

var ErrNotRelative = fmt.Errorf("es: entry must be relative")

func (b *builder) Serve(serve *Serve) (*File, error) {
	if !isRelativeEntry(serve.Entry) {
		return nil, fmt.Errorf("%w %q", ErrNotRelative, serve.Entry)
	}
	result := esbuild.Build(b.serveOptions(serve))
	// Check if there were errors
	if result.Errors != nil {
		return nil, &Error{result.Errors}
	} else if len(result.OutputFiles) == 0 {
		return nil, fmt.Errorf("es: no output files")
	}
	// Return the first file
	file := result.OutputFiles[0]
	return &file, nil
}

type Bundle struct {
	AbsDir   string
	Entries  []string
	Plugins  []esbuild.Plugin
	Platform Platform
}

func (b *builder) bundleOptions(bundle *Bundle) esbuild.BuildOptions {
	switch bundle.Platform {
	case DOM:
		return b.dom(bundle.AbsDir, bundle.Entries, bundle.Plugins)
	default:
		return b.ssr(bundle.AbsDir, bundle.Entries, bundle.Plugins)
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
		return nil, &Error{result.Errors}
	} else if len(result.OutputFiles) == 0 {
		return nil, fmt.Errorf("es: no output files")
	}
	// Return the first file
	return result.OutputFiles, nil
}

const outDir = "./"
const globalName = "bud"

// SSR creates a server-rendered preset
func (b *builder) ssr(absDir string, entries []string, plugins []esbuild.Plugin) esbuild.BuildOptions {
	options := esbuild.BuildOptions{
		EntryPoints:   entries,
		Plugins:       plugins,
		AbsWorkingDir: absDir,
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
func (b *builder) dom(absDir string, entries []string, plugins []esbuild.Plugin) esbuild.BuildOptions {
	options := esbuild.BuildOptions{
		EntryPoints:   entries,
		Plugins:       plugins,
		AbsWorkingDir: absDir,
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

func isRelativeEntry(entry string) bool {
	return strings.HasPrefix(entry, "./")
}

type Error struct {
	messages []esbuild.Message
}

func (e *Error) Error() string {
	errors := esbuild.FormatMessages(e.messages, esbuild.FormatMessagesOptions{
		Color: true,
	})
	return strings.Join(errors, "\n\n")
}

func (e *Error) Errors() []error {
	errors := make([]error, len(e.messages))
	for i, message := range e.messages {
		errors[i] = errorMessage(message)
	}
	return errors
}

type errorMessage esbuild.Message

func (e errorMessage) Error() string {
	return e.Text
}

// TODO: wrap esbuild.Message to use lowercase field names
func (e errorMessage) MarshalJSON() ([]byte, error) {
	return json.Marshal((esbuild.Message)(e))
}
