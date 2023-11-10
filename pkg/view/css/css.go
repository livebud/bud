package css

import (
	"fmt"
	"io"
	"path/filepath"

	esbuild "github.com/evanw/esbuild/pkg/api"
	"github.com/livebud/bud/internal/esb"
	"github.com/livebud/bud/pkg/mod"
	"github.com/livebud/bud/pkg/view"
)

func New(module *mod.Module) *Viewer {
	return &Viewer{module}
}

type Viewer struct {
	module *mod.Module
}

var _ view.Viewer = (*Viewer)(nil)

func (v *Viewer) Render(w io.Writer, path string, data *view.Data) error {
	cssFile, err := v.Compile(path)
	if err != nil {
		return err
	}
	w.Write(cssFile.Contents)
	return nil
}

func (v *Viewer) Compile(entryPath string) (*esbuild.OutputFile, error) {
	options := esbuild.BuildOptions{
		AbsWorkingDir: v.module.Directory(),
		EntryPoints:   []string{entryPath},
		// These are supplanted by the plugin below, in the future we may want to
		// not externalize assets under a certain size.
		Loader: map[string]esbuild.Loader{
			".jpeg": esbuild.LoaderDataURL,
			".jpg":  esbuild.LoaderDataURL,
			".png":  esbuild.LoaderDataURL,
			".svg":  esbuild.LoaderDataURL,
			// TODO: fonts
		},
		// This is needed now to make sure CSS nesting works correctly.
		Engines: []esbuild.Engine{
			{Name: esbuild.EngineChrome, Version: "119"},
			{Name: esbuild.EngineEdge, Version: "118"},
			{Name: esbuild.EngineFirefox, Version: "118"},
			{Name: esbuild.EngineSafari, Version: "17"},
		},
		Plugins: []esbuild.Plugin{
			{
				Name: "externalize-css-assets",
				Setup: func(epb esbuild.PluginBuild) {
					epb.OnResolve(esbuild.OnResolveOptions{Filter: "^."}, func(args esbuild.OnResolveArgs) (result esbuild.OnResolveResult, err error) {
						if filepath.Ext(args.Importer) != ".css" || filepath.Ext(args.Path) == ".css" {
							return result, nil
						}
						result.Path = args.Path
						result.External = true
						return result, nil
					})
				},
			},
		},
		Bundle: true,
	}
	result := esbuild.Build(options)
	if result.Errors != nil {
		return nil, &esb.Error{Messages: result.Errors}
	} else if len(result.OutputFiles) == 0 {
		return nil, fmt.Errorf("esb: no output files")
	}
	return &result.OutputFiles[0], nil
}
