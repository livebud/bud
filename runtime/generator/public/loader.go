package public

import (
	_ "embed"
	"errors"
	"io/fs"
	"path"

	"github.com/livebud/bud/package/vfs"

	"github.com/livebud/bud/internal/bail"
	"github.com/livebud/bud/internal/embed"
	"github.com/livebud/bud/internal/imports"
	"github.com/livebud/bud/package/gomod"
	"github.com/livebud/bud/runtime/command"
)

func Load(flag *command.Flag, fsys fs.FS, module *gomod.Module) (*State, error) {
	loader := &loader{
		fsys:    fsys,
		flag:    flag,
		imports: imports.New(),
		module:  module,
	}
	return loader.Load()
}

type loader struct {
	bail.Struct
	flag    *command.Flag
	fsys    fs.FS
	imports *imports.Set
	module  *gomod.Module
}

// Load the command state
func (l *loader) Load() (state *State, err error) {
	defer l.Recover(&err)
	state = new(State)
	state.Flag = l.flag
	exist, err := vfs.SomeExist(l.fsys, "public", "controller", "view")
	if err != nil {
		return nil, err
	} else if len(exist) == 0 {
		return nil, fs.ErrNotExist
	}
	// Default imports
	l.imports.AddStd("errors", "io", "io/fs", "net/http", "path", "time")
	l.imports.AddNamed("middleware", "github.com/livebud/bud/package/middleware")
	l.imports.AddNamed("overlay", "github.com/livebud/bud/package/overlay")
	// Load embeds
	if exist["public"] && l.flag.Embed {
		state.Embeds = l.loadEmbedsFrom("public", ".")
	}
	// Load default public files
	state.Embeds = append(state.Embeds, l.loadDefaults()...)
	// Add the imports
	state.Imports = l.imports.List()
	return state, nil
}

func (l *loader) loadEmbedsFrom(root, dir string) (files []*embed.File) {
	fullDir := path.Join(root, dir)
	des, err := fs.ReadDir(l.fsys, fullDir)
	if err != nil {
		l.Bail(err)
	}
	for _, de := range des {
		name := de.Name()
		if name[0] == '_' || name[0] == '.' {
			continue
		}
		filePath := path.Join(dir, name)
		if de.IsDir() {
			files = append(files, l.loadEmbedsFrom(root, filePath)...)
			continue
		}
		fullPath := path.Join(root, filePath)
		file := &embed.File{
			Path: fullPath,
		}
		data, err := fs.ReadFile(l.fsys, fullPath)
		if err != nil {
			l.Bail(err)
		}
		file.Data = data
		files = append(files, file)
	}
	return files
}

//go:embed favicon.ico
var favicon []byte

// Default CSS is modern-normalize by Sindre Sorhus
// https://raw.githubusercontent.com/sindresorhus/modern-normalize/v1.1.0/modern-normalize.css
//go:embed default.css
var defaultCSS []byte

func (l *loader) loadDefaults() (files []*embed.File) {
	// Add a public favicon if it doesn't exist
	if err := vfs.Exist(l.fsys, "public/favicon.ico"); err != nil {
		if !errors.Is(err, fs.ErrNotExist) {
			l.Bail(err)
		}
		files = append(files, &embed.File{
			Path: "public/favicon.ico",
			Data: favicon,
		})
	}
	// Add default.css if it doesn't exist
	if err := vfs.Exist(l.fsys, "public/default.css"); err != nil {
		if !errors.Is(err, fs.ErrNotExist) {
			l.Bail(err)
		}
		files = append(files, &embed.File{
			Path: "public/default.css",
			Data: defaultCSS,
		})
	}
	return files
}
