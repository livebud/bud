package public

import (
	"io/fs"
	"path"
	"strings"

	"gitlab.com/mnm/bud/internal/bail"
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
	state.Files = l.loadFiles()
	if len(state.Files) == 0 {
		return nil, fs.ErrNotExist
	}
	l.imports.AddStd("errors", "io", "io/fs", "net/http", "path", "time")
	// l.imports.AddStd("fmt")
	l.imports.AddNamed("middleware", "gitlab.com/mnm/bud/package/middleware")
	l.imports.AddNamed("overlay", "gitlab.com/mnm/bud/package/overlay")
	state.Imports = l.imports.List()
	return state, nil
}

func (l *loader) loadFiles() (files []*File) {
	files = l.loadFilesFrom("public", ".")
	return files
}

func (l *loader) loadFilesFrom(root, dir string) (files []*File) {
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
			files = append(files, l.loadFilesFrom(root, filePath)...)
			continue
		}
		fullPath := path.Join(root, filePath)
		file := &File{
			Path: fullPath,
			Root: root,
		}
		if l.flag.Embed {
			file.Data = l.loadData(fullPath)
			file.Mode = "0644" // TODO: configurable
		}
		files = append(files, file)
	}
	return files
}

const lowerHex = "0123456789abcdef"

// Based on:
// https://github.com/go-bindata/go-bindata/blob/26949cc13d95310ffcc491c325da869a5aafce8f/stringwriter.go#L18-L36
func (l *loader) loadData(filePath string) string {
	data, err := fs.ReadFile(l.fsys, filePath)
	if err != nil {
		l.Bail(err)
	}
	if len(data) == 0 {
		return ""
	}
	s := new(strings.Builder)
	buf := []byte(`\x00`)
	for _, b := range data {
		buf[2] = lowerHex[b/16]
		buf[3] = lowerHex[b%16]
		s.Write(buf)
	}
	return s.String()
}
