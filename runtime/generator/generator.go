package generator

import (
	"github.com/livebud/bud/package/genfs"
	"github.com/livebud/bud/package/log"
	"github.com/livebud/bud/package/virtual"
)

type FS = genfs.FS
type File = genfs.File
type Dir = genfs.Dir
type GenerateFile = genfs.GenerateFile
type GenerateDir = genfs.GenerateDir
type FileSystem = genfs.FileSystem

type Extension interface {
	Extend(genfs FileSystem)
}

func NewGenerator(genfs genfs.FileSystem, log log.Log, extensions ...Extension) *Generator {
	for _, extension := range extensions {
		extension.Extend(genfs)
	}
	return &Generator{genfs, log}
}

type Generator struct {
	from genfs.FileSystem
	log  log.Log
}

// Generate files to the bud/ directory
func (g *Generator) Generate(to virtual.FS, subdirs ...string) error {
	return virtual.Sync(g.log, g.from, to, subdirs...)
}

// type Generator struct {
// 	from genfs.FileSystem
// }

// // type Schema struct {
// // 	GenerateFiles map[string]GenerateFile
// // 	GenerateDirs  map[string]GenerateDir
// // }

// // func New(genfs, extensions ...Extension) genfs.FileSystem {
// // 	// Genfs combines both the generators and the real filesystem files into one
// // 	// virtual filesystem. We don't want to the virtual filesystem to know about
// // 	// the bud/ files because otherwise it will never cleanup bud/ files that
// // 	// were once generated but are no longer in the generators.
// // 	//
// // 	// There's one subtly though. We want to keep the gen files, like
// // 	// bud/cmd/gen/main.go, so we don't exclude those.
// // 	dirfs := virtual.Exclude(module, func(path string) bool {
// // 		return exclude(path)
// // 	})
// // 	gen := genfs.New(dag.Discard, dirfs, log)
// // 	gen.Extend(extensions...)
// // 	return gen
// // }

// func exclude(path string) bool {
// 	if isGenPath(path) || isBudChild(path) {
// 		return false
// 	} else if isBudPath(path) {
// 		return true
// 	}
// 	return false
// }

// // Exclude everything in bud/
// func isBudPath(p string) bool {
// 	return strings.HasPrefix(p, "bud/")
// }

// func isBudChild(p string) bool {
// 	return path.Dir(p) == "bud"
// }

// func isGenPath(p string) bool {
// 	return strings.HasPrefix(p, "bud/cmd/gen") ||
// 		strings.HasPrefix(p, "bud/internal/gen") ||
// 		strings.HasPrefix(p, "bud/pkg/gen")
// }

// // type FileSystem struct {
// // 	gen    genfs.FileSystem
// // 	log    log.Log
// // 	module *gomod.Module
// // }

// // func (g *FileSystem) Sync() error {
// // 	return virtual.Sync(g.log, g.gen, g.module, "bud")
// // }
