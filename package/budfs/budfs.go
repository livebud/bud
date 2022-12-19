package budfs

import (
	"io/fs"

	"github.com/livebud/bud/package/virtual"

	"github.com/livebud/bud/internal/genfs"
	"github.com/livebud/bud/internal/pathlink"
	"github.com/livebud/bud/package/gomod"
	"github.com/livebud/bud/package/log"
	"github.com/livebud/bud/package/virtual/vcache"
)

type FileGenerator = genfs.FileGenerator
type DirGenerator = genfs.DirGenerator
type FileServer = genfs.FileServer
type File = genfs.File
type Dir = genfs.Dir
type FS = genfs.FS

type FileSystem interface {
	fs.FS
	GenerateFile(path string, fn func(fsys FS, file *File) error)
	FileGenerator(path string, generator FileGenerator)
	GenerateDir(path string, fn func(fsys FS, dir *Dir) error)
	DirGenerator(path string, generator DirGenerator)
	ServeFile(path string, fn func(fsys FS, file *File) error)
	FileServer(path string, server FileServer)
	Sync(paths ...string) error
	Change(paths ...string) error
	Close() error
}

func Load(module *gomod.Module, log log.Log) (*fileSystem, error) {
	cache := vcache.New()
	linker := pathlink.New(log)
	genfs := genfs.New(cache, module, linker, log)
	// parser := parser.New(genfs, module)
	// injector := di.New(genfs, log, module, parser)
	// genfs.FileGenerator("bud/internal/generator/generator.go", generator)
	// genfs.FileGenerator("bud/cmd/afs/main.go", afs.New(exec, injector, log, module))
	return &fileSystem{cache, genfs, linker, log, module}, nil
}

var _ FileSystem = (*fileSystem)(nil)

type fileSystem struct {
	cache vcache.Cache
	genfs.FileSystem
	linker pathlink.Linker
	log    log.Log
	module *gomod.Module
}

func (f *fileSystem) Sync(writable virtual.FS, paths ...string) error {
	return nil
}

func (f *fileSystem) Change(paths ...string) error {
	return nil
}

func (f *fileSystem) Close() error {
	return nil
}
