package public

import (
	"io/fs"
	"path"
	"strings"

	"github.com/livebud/bud/internal/valid"

	"github.com/livebud/bud/framework"
	"github.com/livebud/bud/package/finder"

	"github.com/livebud/bud/internal/bail"
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
	paths, err := finder.Find(l.fsys, "public/**", func(fullpath string, isDir bool) (entries []string) {
		if isDir {
			return nil
		}
		if valid.PublicFile(path.Base(fullpath)) {
			entries = append(entries, fullpath)
		}
		return entries
	})
	if err != nil {
		return nil, err
	} else if len(paths) == 0 {
		return nil, fs.ErrNotExist
	}
	state = new(State)
	// Load the files from paths
	state.Files = l.loadFiles(paths)
	// Default imports
	l.imports.AddNamed("virtual", "github.com/livebud/bud/package/virtual")
	l.imports.AddNamed("publicrt", "github.com/livebud/bud/framework/public/publicrt")
	l.imports.AddNamed("router", "github.com/livebud/bud/package/router")
	l.imports.AddNamed("http", "net/http")
	l.imports.AddNamed("fs", "io/fs")
	// Add the imports
	state.Imports = l.imports.List()
	return state, nil
}

func (l *loader) loadFiles(paths []string) (files []*File) {
	for _, path := range paths {
		files = append(files, l.loadFile(path))
	}
	return files
}

func (l *loader) loadFile(path string) *File {
	file := new(File)
	file.Path = path
	file.Route = strings.TrimPrefix(path, "public")
	if l.flag.Embed {
		data, err := fs.ReadFile(l.fsys, path)
		if err != nil {
			l.Bail(err)
		}
		file.Data = data
	}
	return file
}
