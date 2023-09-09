package main

import (
	"os"

	"blog.com/internal/injector"
	"github.com/livebud/bud/pkg/program"
)

func main() {
	os.Exit(program.Run(injector.New()))
}
