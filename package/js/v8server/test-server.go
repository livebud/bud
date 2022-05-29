//go:build ignore
// +build ignore

package main

import (
	"fmt"
	"os"

	"github.com/livebud/bud/package/js/v8server"
)

func main() {
	if err := v8server.Serve(); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
	os.Exit(0)
}
