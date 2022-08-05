package main

import (
	"os"

	"github.com/livebud/bud/example/tailwind/generator/tailwind"
	"github.com/livebud/bud/internal/remotefs"
	"github.com/livebud/bud/package/gomod"
	"github.com/livebud/bud/package/goplugin"
	"github.com/livebud/bud/package/log"
	"github.com/livebud/bud/package/log/console"
	"github.com/livebud/bud/package/overlay"
)

func main() {
	if err := run(); err != nil {
		console.Error(err.Error())
		os.Exit(1)
	}
	os.Exit(0)
}

func run() error {
	conn, err := goplugin.Serve("Generate")
	if err != nil {
		return err
	}
	module, err := gomod.Find(".")
	if err != nil {
		return err
	}
	log := log.Discard
	genfs, err := DILoad(log, module)
	if err != nil {
		return err
	}
	remotefs.Serve(genfs, conn)
	return nil
}

// Load comes from DI
func DILoad(log log.Interface, module *gomod.Module) (*FileSystem, error) {
	genfs, err := overlay.Load(log, module)
	if err != nil {
		return nil, err
	}
	generators := &Generators{
		Tailwind: &tailwind.Generator{},
	}
	return Register(genfs, generators), nil
}

type Generators struct {
	Tailwind *tailwind.Generator
}

func Register(genfs *overlay.FileSystem, generators *Generators) *FileSystem {
	genfs.DirGenerator("tailwind", generators.Tailwind)
	return genfs
}

type FileSystem = overlay.FileSystem

// type Generator struct {
// 	fsys fs.FS
// }

// // type File struct {
// // 	Data []byte
// // }

// func (g *Generator) Open(name string, result *[]byte) error {
// 	code, err := fs.ReadFile(g.fsys, name)
// 	if err != nil {
// 		return err
// 	}
// 	*result = code
// 	return nil
// }
