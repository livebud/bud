package es

import (
	"fmt"
	"path/filepath"
	"strings"

	esbuild "github.com/evanw/esbuild/pkg/api"
	"github.com/livebud/bud/package/virtual"
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
				OutputPath: strings.TrimSuffix(build.Entrypoint, ".js"),
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
		dir, err := filepath.EvalSymlinks(b.dir)
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
