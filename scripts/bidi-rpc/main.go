package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/livebud/bud/internal/extrafile"

	"github.com/livebud/bud/package/socket"

	"github.com/livebud/bud/internal/current"

	"github.com/livebud/bud/package/log/console"
)

func main() {
	if err := run(); err != nil {
		console.Error(err.Error())
	}
}

func run() error {
	if value := os.Getenv("CHILD"); value != "" {
		return child()
	}
	return parent()
}

func parent() error {
	dir, err := current.Directory()
	if err != nil {
		return err
	}
	ln, err := socket.Listen(":0")
	if err != nil {
		return err
	}
	defer ln.Close()
	file, err := ln.File()
	if err != nil {
		return err
	}
	defer file.Close()
	cmd := exec.Command("go", "run", "main.go")
	cmd.Env = append(os.Environ(), "CHILD=1")
	cmd.Dir = dir
	extrafile.Inject(&cmd.ExtraFiles, &cmd.Env, "SOCKET", file)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	return cmd.Run()
}

func child() error {
	files := extrafile.Load("SOCKET")
	if len(files) == 0 {
		return fmt.Errorf("no files passed through")
	}
	ln, err := socket.From(files[0])
	if err != nil {
		return err
	}
	fmt.Println(ln.Addr().String())
	fmt.Println("called child!")
	return nil
}
