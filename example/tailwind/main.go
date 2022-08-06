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

	// // FileSystem()
	// // TODO: switch to GenerateDir with a remote filesystem
	// // genfs.Mount
	// genfs.ServeFile("bud/internal/generator", func(ctx context.Context, fsys overlay.F, file *overlay.File) error {
	// 	rel := strings.TrimPrefix(file.Path(), "bud/internal/generator/")
	// 	var data []byte
	// 	if err := client.Call("generator.Open", rel, &data); err != nil {
	// 		return err
	// 	}
	// 	file.Data = data
	// 	return nil
	// })
	// err = fs.WalkDir(client, ".", func(path string, d fs.DirEntry, err error) error {
	// 	if err != nil {
	// 		return err
	// 	}
	// 	fmt.Println(path)
	// 	return nil
	// })
	// if err != nil {
	// 	return err
	// }
	// now := time.Now()
	// overlay.Load(log.Discard, client)
	// // code, err := fs.ReadFile(client, "tailwind/tailwind.css")
	// // if err != nil {
	// // 	return err
	// // }
	// // fmt.Println(string(code), time.Since(now))
	// // now = time.Now()
	// des, err := fs.ReadDir(client, "tailwind")
	// if err != nil {
	// 	return err
	// }
	// for _, de := range des {
	// 	fmt.Println(de.Name(), time.Since(now))
	// }
	return nil
}
