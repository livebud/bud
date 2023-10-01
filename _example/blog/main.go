package main

import (
	"os"

	"github.com/livebud/bud/example/blog/internal/injector"
	"github.com/livebud/bud/program"
)

func main() {
	os.Exit(program.Run(injector.New()))
}
