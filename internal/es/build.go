package es

import (
	"fmt"
	"io/fs"
	"strings"

	esbuild "github.com/evanw/esbuild/pkg/api"
)

// TODO: replace with *gomod.Module
func New(absDir string) *Builder {
	return &Builder{
		dir: absDir,
		base: esbuild.BuildOptions{
			AbsWorkingDir: absDir,
			Outdir:        "./",
			Format:        esbuild.FormatIIFE,
			Platform:      esbuild.PlatformBrowser,
			GlobalName:    "bud",
			Bundle:        true,
			Plugins: []esbuild.Plugin{
				httpPlugin(absDir),
				esmPlugin(absDir),
			},
		},
	}
}

type Builder struct {
	dir  string
	base esbuild.BuildOptions
}

func (b *Builder) Directory() string {
	return b.dir
}

type Entrypoint = esbuild.EntryPoint
type Plugin = esbuild.Plugin

type Build struct {
	Entrypoint string
	Plugins    []Plugin
	Minify     bool
}

func (b *Builder) Build(build *Build) ([]byte, error) {
	input := esbuild.BuildOptions{
		EntryPointsAdvanced: []esbuild.EntryPoint{
			{
				InputPath:  build.Entrypoint,
				OutputPath: build.Entrypoint,
			},
		},
		AbsWorkingDir: b.base.AbsWorkingDir,
		Outdir:        b.base.Outdir,
		Format:        b.base.Format,
		Platform:      b.base.Platform,
		GlobalName:    b.base.GlobalName,
		Bundle:        b.base.Bundle,
		Metafile:      b.base.Metafile,
		Plugins:       append(build.Plugins, b.base.Plugins...),
	}
	if build.Minify {
		input.MinifyWhitespace = true
		input.MinifyIdentifiers = true
		input.MinifySyntax = true
	}
	result := esbuild.Build(input)
	if len(result.Errors) > 0 {
		msgs := esbuild.FormatMessages(result.Errors, esbuild.FormatMessagesOptions{
			Color: true,
			Kind:  esbuild.ErrorMessage,
		})
		return nil, fmt.Errorf(strings.Join(msgs, "\n"))
	}
	// Expect exactly 1 output file
	if len(result.OutputFiles) != 1 {
		return nil, fmt.Errorf("expected exactly 1 output file but got %d", len(result.OutputFiles))
	}
	ssrCode := result.OutputFiles[0].Contents
	return ssrCode, nil
}

type Bundle struct {
	Entrypoints []string
	Plugins     []Plugin
	Minify      bool
}

func (b *Builder) Bundle(out fs.FS, bundle *Bundle) error {
	return nil
}
