package public

import (
	"io/fs"
	"path"

	"gitlab.com/mnm/bud/internal/bail"
	"gitlab.com/mnm/bud/internal/embed"
	"gitlab.com/mnm/bud/internal/imports"
	"gitlab.com/mnm/bud/package/gomod"
	"gitlab.com/mnm/bud/runtime/bud"
)

func Load(flag *bud.Flag, fsys fs.FS, module *gomod.Module) (*State, error) {
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
	flag    *bud.Flag
	fsys    fs.FS
	imports *imports.Set
	module  *gomod.Module
}

// Load the command state
func (l *loader) Load() (state *State, err error) {
	defer l.Recover(&err)
	state = new(State)
	state.Flag = l.flag
	if l.flag.Embed {
		state.Embeds = l.loadEmbedsFrom("public", ".")
	}
	l.imports.AddStd("errors", "io", "io/fs", "net/http", "path", "time")
	// l.imports.AddStd("fmt")
	l.imports.AddNamed("middleware", "gitlab.com/mnm/bud/package/middleware")
	l.imports.AddNamed("overlay", "gitlab.com/mnm/bud/package/overlay")
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
