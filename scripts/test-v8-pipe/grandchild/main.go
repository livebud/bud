package main

import (
	"context"
	"fmt"
	"log"

	"github.com/livebud/bud/internal/extrafile"
	"github.com/livebud/bud/package/js/v8client"
	"github.com/livebud/bud/package/socket"
)

func run(ctx context.Context) error {
	fmt.Println("calling grandchild")
	// files := extrafile.Load("V8")
	v8client, err := v8client.Load(ctx)
	if err != nil {
		return err
	}
	// v8client := v8client.New(files[0], files[1])
	result, err := v8client.Eval("script.js", "__svelte__ + 2")
	if err != nil {
		return err
	}
	fmt.Println(result)
	appFile := extrafile.Load("APP")
	listener, err := socket.From(appFile[0])
	if err != nil {
		return err
	}
	fmt.Println("got listener", listener)
	return nil
}

func main() {
	if err := run(context.Background()); err != nil {
		log.Fatal(err)
	}
}
