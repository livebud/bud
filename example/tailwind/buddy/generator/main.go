package main

import (
	"io/fs"
	"net/rpc"
	"os"

	"github.com/livebud/bud/example/tailwind/generator/tailwind"
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
	generator, err := Load(log, module)
	server := rpc.NewServer()
	if err := server.RegisterName("generator", generator); err != nil {
		return err
	}
	server.ServeConn(conn)
	return nil
}

// Load comes from DI
func Load(log log.Interface, module *gomod.Module) (*Generator, error) {
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

func Register(genfs *overlay.FileSystem, generators *Generators) *Generator {
	genfs.DirGenerator("tailwind", generators.Tailwind)
	return &Generator{genfs}
}

type Generator struct {
	fsys fs.FS
}

// type File struct {
// 	Data []byte
// }

func (g *Generator) Open(name string, result *[]byte) error {
	code, err := fs.ReadFile(g.fsys, name)
	if err != nil {
		return err
	}
	*result = code
	return nil
}
