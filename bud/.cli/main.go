package main

import (
	context "context"
	program "gitlab.com/mnm/bud/bud/.cli/program"
	os "os"
)

func main() {
	ctx := context.Background()
	exitCode := program.Run(ctx, os.Args[1:]...)
	os.Exit(exitCode)
}
