package main

import (
	"fmt"
	"io/fs"
	"os"
	"time"

	"github.com/livebud/bud/internal/remotefs"
	"github.com/livebud/bud/package/goplugin"
	"github.com/livebud/bud/package/log/console"
)

func main() {
	if err := run(); err != nil {
		console.Error(err.Error())
		os.Exit(1)
	}
	os.Exit(0)
}

func run() error {
	// module, err := gomod.Find(".")
	// if err != nil {
	// 	return err
	// }
	// genfs, err := overlay.Load(log.Discard, module)
	// if err != nil {
	// 	return err
	// }

	// CustomGenerators()

	conn, err := goplugin.Start("go", "run", "buddy/generator/main.go")
	if err != nil {
		return err
	}
	client := remotefs.NewClient(conn)
	defer client.Close()

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
	now := time.Now()
	// code, err := fs.ReadFile(client, "tailwind/tailwind.css")
	// if err != nil {
	// 	return err
	// }
	// fmt.Println(string(code), time.Since(now))
	// now = time.Now()
	des, err := fs.ReadDir(client, "tailwind")
	if err != nil {
		return err
	}
	for _, de := range des {
		fmt.Println(de.Name(), time.Since(now))
	}
	return nil
}
