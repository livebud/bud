package public

import (
	"errors"
	"io/fs"
	"path"

	"gitlab.com/mnm/bud/internal/valid"

	"gitlab.com/mnm/bud/budfs"
	"gitlab.com/mnm/bud/internal/bail"
	"gitlab.com/mnm/bud/internal/imports"
	"gitlab.com/mnm/bud/mod"
)

func Load(bfs budfs.FS, module *mod.Module, embed, minify bool) (*State, error) {
	loader := &loader{
		bfs:     bfs,
		imports: imports.New(),
		module:  module,
		embed:   embed,
		minify:  minify,
	}
	return loader.Load()
}

type loader struct {
	bail.Struct
	bfs     budfs.FS
	imports *imports.Set
	module  *mod.Module
	embed   bool
	minify  bool
}

// Load the command state
func (l *loader) Load() (state *State, err error) {
	defer l.Recover(&err)
	state = new(State)
	state.Embed = l.embed
	state.Files = l.loadFiles()
	if len(state.Files) == 0 {
		return nil, fs.ErrNotExist
	}
	return state, nil
}

func (l *loader) loadFiles() (files []*File) {
	files = l.loadFilesFrom("public", ".")
	files = append(files, l.loadFilesFromPlugins("bud/plugin")...)
	return files
}

func (l *loader) loadFilesFrom(root, dir string) (files []*File) {
	fullPath := path.Join(root, dir)
	des, err := fs.ReadDir(l.bfs, fullPath)
	if err != nil {
		if !errors.Is(err, fs.ErrNotExist) {
			l.Bail(err)
		}
		return files
	}
	for _, de := range des {
		name := de.Name()
		if name[0] == '_' || name[0] == '.' {
			continue
		}
		if de.IsDir() {
			files = append(files, l.loadFilesFrom(root, path.Join(dir, name))...)
			continue
		}
		files = append(files, &File{
			Path: path.Join("bud", fullPath, name),
			Root: root,
		})
	}
	return files
}

func (l *loader) loadFilesFromPlugins(pluginBaseDir string) (files []*File) {
	des, err := fs.ReadDir(l.bfs, pluginBaseDir)
	if err != nil {
		if !errors.Is(err, fs.ErrNotExist) {
			l.Bail(err)
		}
		return files
	}
	for _, de := range des {
		name := de.Name()
		if valid.PluginDir(name) {
			continue
		}
		fullDir := path.Join(pluginBaseDir, name)
		des, err := fs.ReadDir(l.bfs, fullDir)
		if err != nil {
			if !errors.Is(err, fs.ErrNotExist) {
				l.Bail(err)
			}
			return files
		}
		for _, de := range des {
			if de.Name() != "public" {
				continue
			}
			pluginPublicDir := path.Join(fullDir, "public")
			files = append(files, l.loadFilesFromPlugin(pluginPublicDir, ".")...)
		}
	}
	return files
}

func (l *loader) loadFilesFromPlugin(pluginPublicDir, dir string) (files []*File) {
	fullPath := path.Join(pluginPublicDir, dir)
	des, err := fs.ReadDir(l.bfs, fullPath)
	if err != nil {
		l.Bail(err)
		return files
	}
	for _, de := range des {
		name := de.Name()
		if name[0] == '_' || name[0] == '.' {
			continue
		}
		filePath := path.Join(dir, name)
		if de.IsDir() {
			files = append(files, l.loadFilesFromPlugin(pluginPublicDir, filePath)...)
			continue
		}
		files = append(files, &File{
			Path: path.Join("bud", "public", filePath),
			Root: pluginPublicDir,
		})
	}
	return files
}
