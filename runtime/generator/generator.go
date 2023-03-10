package generator

import (
	"path"
	"strings"

	"github.com/livebud/bud/internal/dag"

	"github.com/livebud/bud/package/genfs"
	"github.com/livebud/bud/package/gomod"
	"github.com/livebud/bud/package/log"
	"github.com/livebud/bud/package/virtual"
)

type FS = genfs.FS
type File = genfs.File
type GenerateFile = genfs.GenerateFile

type Schema struct {
	GenerateFiles map[string]GenerateFile
}

func New(log log.Log, module *gomod.Module, schema *Schema) *Generator {
	// Genfs combines both the generators and the real filesystem files into one
	// virtual filesystem. We don't want to the virtual filesystem to know about
	// the bud/ files because otherwise it will never cleanup bud/ files that
	// were once generated but are no longer in the generators.
	//
	// There's one subtly though. We want to keep the gen files, like
	// bud/cmd/gen/main.go, so we don't exclude those.
	dirfs := virtual.Exclude(module, func(path string) bool {
		return exclude(path)
	})
	gen := genfs.New(dag.Discard, dirfs, log)
	for path, generator := range schema.GenerateFiles {
		gen.FileGenerator(path, generator)
	}
	return &Generator{gen, log, module}
}

func exclude(path string) bool {
	if isGenPath(path) || isBudChild(path) {
		return false
	} else if isBudPath(path) {
		return true
	}
	return false
}

// Exclude everything in bud/
func isBudPath(p string) bool {
	return strings.HasPrefix(p, "bud/")
}

func isBudChild(p string) bool {
	return path.Dir(p) == "bud"
}

func isGenPath(p string) bool {
	return strings.HasPrefix(p, "bud/cmd/gen") ||
		strings.HasPrefix(p, "bud/internal/gen") ||
		strings.HasPrefix(p, "bud/pkg/gen")
}

type Generator struct {
	gen    genfs.FileSystem
	log    log.Log
	module *gomod.Module
}

func (g *Generator) Generate() error {
	return virtual.Sync(g.log, g.gen, g.module, "bud")
}
