package public

import (
	"errors"
	"io/fs"
	"path"

	"github.com/livebud/bud/framework"
	"github.com/livebud/bud/package/vfs"

	"github.com/livebud/bud/internal/bail"
	"github.com/livebud/bud/internal/embed"
	"github.com/livebud/bud/internal/embedded"
	"github.com/livebud/bud/internal/imports"
)

func Load(fsys fs.FS, flag *framework.Flag) (*State, error) {
	loader := &loader{
		fsys:    fsys,
		flag:    flag,
		imports: imports.New(),
	}
	return loader.Load()
}

type loader struct {
	bail.Struct
	flag    *framework.Flag
	fsys    fs.FS
	imports *imports.Set
}

// Load the command state
func (l *loader) Load() (state *State, err error) {
	defer l.Recover(&err)
	state = new(State)
	state.Flag = l.flag
	// Default imports
	l.imports.AddNamed("virtual", "github.com/livebud/bud/package/virtual")
	l.imports.AddNamed("middleware", "github.com/livebud/bud/package/middleware")
	l.imports.AddNamed("publicrt", "github.com/livebud/bud/framework/public/publicrt")
	// Load embeds
	if l.flag.Embed {
		state.Embeds = l.loadEmbedsFrom(".", ".")
	}
	// Load default public files. Out of convenience, these defaults are embedded
	// regardless of flag.Embed
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
		data, err := fs.ReadFile(l.fsys, fullPath)
		if err != nil {
			l.Bail(err)
		}
		files = append(files, &embed.File{
			Path: filePath,
			Data: data,
		})
	}
	return files
}

func (l *loader) loadDefaults() (files []*embed.File) {
	// Add a public favicon if it doesn't exist
	if err := vfs.Exist(l.fsys, "favicon.ico"); err != nil {
		if !errors.Is(err, fs.ErrNotExist) {
			l.Bail(err)
		}
		files = append(files, &embed.File{
			Path: "favicon.ico",
			Data: embedded.Favicon(),
		})
	}
	return files
}
