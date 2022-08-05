package main

import (
	"context"
	"fmt"
	"io/fs"
	"net/rpc"
	"os"
	"strings"

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
	module, err := gomod.Find(".")
	if err != nil {
		return err
	}
	genfs, err := overlay.Load(log.Discard, module)
	if err != nil {
		return err
	}

	// CustomGenerators()

	conn, err := goplugin.Start("go", "run", "buddy/generator/main.go")
	if err != nil {
		return err
	}
	client := rpc.NewClient(conn)
	defer client.Close()

	// FileSystem()
	// TODO: switch to GenerateDir with a remote filesystem
	// genfs.Mount
	genfs.ServeFile("bud/internal/generator", func(ctx context.Context, fsys overlay.F, file *overlay.File) error {
		rel := strings.TrimPrefix(file.Path(), "bud/internal/generator/")
		var data []byte
		if err := client.Call("generator.Open", rel, &data); err != nil {
			return err
		}
		file.Data = data
		return nil
	})
	code, err := fs.ReadFile(genfs, "bud/internal/generator/tailwind/tailwind.css")
	if err != nil {
		return err
	}
	fmt.Println("got result", string(code))
	return nil
}
