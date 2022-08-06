package main

import (
	"fmt"
	"io/fs"
	"os"

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
	conn, err := goplugin.Start("go", "run", "buddy/generator/main.go")
	if err != nil {
		return err
	}
	remotefs := remotefs.NewClient(conn)
	defer remotefs.Close()

	// Load the overlay
	module, err := gomod.Find(".")
	if err != nil {
		return err
	}
	genfs, err := overlay.Load(log.Discard, module)
	if err != nil {
		return err
	}
	genfs.Mount("bud/generator", remotefs)

	data, err := fs.ReadFile(genfs, "bud/generator/tailwind/tailwind.css")
	if err != nil {
		return err
	}

	fmt.Println("tailwind/tailwind.css:", string(data))

	return nil
}
