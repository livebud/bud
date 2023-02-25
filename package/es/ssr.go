package es

import (
	"fmt"
	"path/filepath"
	"strings"

	esbuild "github.com/evanw/esbuild/pkg/api"
	"github.com/livebud/bud/package/gomod"
	"github.com/livebud/bud/package/virtual"
)

func New(module *gomod.Module) *Builder {
	absDir := module.Directory()
	return &Builder{
		module: module,
		base: esbuild.BuildOptions{
			AbsWorkingDir: absDir,
			Outdir:        "./",
			Plugins: []esbuild.Plugin{
				httpPlugin(),
				// npmPlugin(),
				// esmPlugin(absDir),
			},
		},
	}
}

type Builder struct {
	module *gomod.Module
	base   esbuild.BuildOptions
}

type Entrypoint = esbuild.EntryPoint
type Plugin = esbuild.Plugin

type Mode uint8

const (
	ModeDOM Mode = iota
	ModeSSR
)

type Build struct {
	Entrypoint string
	Mode       Mode
	Plugins    []Plugin
	Minify     bool
}

func (b *Builder) Build(build *Build) ([]byte, error) {
	input := esbuild.BuildOptions{
		EntryPointsAdvanced: []esbuild.EntryPoint{
			{
				InputPath:  build.Entrypoint,
				OutputPath: strings.TrimSuffix(build.Entrypoint, ".js"),
			},
		},
		AbsWorkingDir: b.base.AbsWorkingDir,
		Outdir:        b.base.Outdir,
		Metafile:      b.base.Metafile,
		Plugins:       append(build.Plugins, b.base.Plugins...),

		// Always bundle, use plugins to granularly mark files as external
		Bundle: true,
	}

	// Switch configuration based on build mode
	switch build.Mode {
	case ModeDOM:
		input.Format = esbuild.FormatESModule
		input.Platform = esbuild.PlatformBrowser
		// Add "import" condition to support svelte/interna
		// https://esbuild.github.io/api/#how-conditions-wor
		input.Conditions = []string{"browser", "default", "import"}
	case ModeSSR:
		input.Format = esbuild.FormatIIFE
		input.Platform = esbuild.PlatformNode
		input.GlobalName = "bud"
	}

	// Handle minification
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

func (b *Builder) Bundle(fsys virtual.FS, bundle *Bundle) error {
	entries := make([]esbuild.EntryPoint, len(bundle.Entrypoints))
	for i, entry := range bundle.Entrypoints {
		entries[i] = esbuild.EntryPoint{
			InputPath:  entry,
			OutputPath: strings.TrimSuffix(entry, ".js"),
		}
	}
	input := esbuild.BuildOptions{
		EntryPointsAdvanced: entries,
		AbsWorkingDir:       b.base.AbsWorkingDir,
		Outdir:              b.base.Outdir,
		Format:              b.base.Format,
		Platform:            b.base.Platform,
		GlobalName:          b.base.GlobalName,
		Bundle:              b.base.Bundle,
		Metafile:            b.base.Metafile,
		Plugins:             append(bundle.Plugins, b.base.Plugins...),
	}
	if bundle.Minify {
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
		return fmt.Errorf(strings.Join(msgs, "\n"))
	}
	for _, file := range result.OutputFiles {
		dir, err := filepath.EvalSymlinks(b.module.Directory())
		if err != nil {
			return err
		}
		relPath, err := filepath.Rel(dir, file.Path)
		if err != nil {
			return err
		}
		if err := fsys.MkdirAll(filepath.Dir(relPath), 0755); err != nil {
			return err
		}
		if err := fsys.WriteFile(relPath, file.Contents, 0644); err != nil {
			return err
		}
	}
	return nil
}
