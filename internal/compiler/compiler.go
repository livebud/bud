package compiler

import (
	"context"
	"io"

	"github.com/livebud/bud/package/gomod"
	"github.com/livebud/bud/package/log"
)

// Bud compiler
type Bud struct {
	// Dependencies
	Module *gomod.Module
	Log    log.Interface

	// Passed into bud/cli
	Env    []string
	Stdin  io.Reader
	Stdout io.Writer
	Stderr io.Writer
}

// Flag contains the compilation configuration
type Flag struct {
	Embed  bool
	Minify bool
	Hot    string
}

// 2. Compile
//   a. Generate cli
//   	 i. Generate bud/internal/cli
//     ii. Build bud/cli
//     iii. Run bud/cli
//   b. Generate app
//     i. Generate bud/internal/app
//     ii. Build bud/app
func (b *Bud) Compile(ctx context.Context, flag *Flag) error {
	return nil
}
