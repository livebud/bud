package run

import (
	"context"
	"os/exec"

	"github.com/livebud/bud/framework"
	"github.com/livebud/bud/internal/command"
	"github.com/livebud/bud/package/socket"
)

func New(bud *command.Bud, web, hot socket.Listener) *Command {
	return &Command{
		bud:  bud,
		web:  web,
		hot:  hot,
		Flag: new(framework.Flag),
	}
}

type Command struct {
	bud *command.Bud

	// Passed in for testing
	web socket.Listener // Can be nil
	hot socket.Listener // Can be nil

	// Flags
	Flag   *framework.Flag
	Listen string
}

func (c *Command) Run(ctx context.Context) (err error) {
	// Setup dependencies
	module, err := c.bud.Module()
	if err != nil {
		return err
	}
	log, err := c.bud.Logger()
	if err != nil {
		return err
	}

	log.Info("running run...")

	{ // 1. Listen
		// Bind the web listener to the listen address
		if c.web == nil && !disabled(c.Listen) {
			c.web, err = socket.Listen(c.Listen)
			if err != nil {
				return err
			}
		}
		// Bind the hot listener to the hot address
		if c.hot == nil && !disabled(c.Flag.Hot) {
			c.hot, err = socket.Listen(c.Flag.Hot)
			if err != nil {
				return err
			}
		}
		// TODO: Setup the Bud pipe
	}

	// { // 2. Compile
	// 	budCompiler := c.bud.Compiler(log, module)
	// 	err := budCompiler.Compile(ctx, &compiler.Flag{
	// 		Minify: c.Minify,
	// 		Embed:  c.Embed,
	// 		Hot:    c.Hot,
	// 	})
	// 	if err != nil {
	// 		return err
	// 	}
	// }

	{ // 3. Start the application
		cmd := exec.Command("bud/app")
		cmd.Dir = module.Directory()
		cmd.Env = c.bud.Env
		cmd.Stdin = c.bud.Stdin
		cmd.Stdout = c.bud.Stdout
		cmd.Stderr = c.bud.Stderr
		if err := cmd.Start(); err != nil {
			// TODO: watch for changes and retry
			return err
		}
	}

	{ // 4. Watch for changes, recompile and start.

	}

	// 1. Listen
	// 2. Compile
	//   a. Generate cli
	//   	 i. Generate bud/internal/cli
	//     ii. Build bud/cli
	//     iii. Run bud/cli
	//   b. Generate app
	//     i. Generate bud/internal/app
	//     ii. Build bud/app
	// 3. Start
	// 4. Watch
	// (...Compile, Start)
	return nil
}

// func (c *Command) listen() (web, hot socket.Listener, err error) {
// 	return c.WebListener, c.HotListener
// }

// // Bind the web listener to the listen address
// func (c *Command) listenWeb() (socket.Listener, error) {
// 	if c.WebListener != nil {
// 		return c.WebListener, nil
// 	} else if disabled(c.Listen) {
// 		return nil, nil
// 	}
// 	listener, err := socket.Listen(c.Listen)
// 	if err != nil {
// 		return nil, err
// 	}

// }

// func (c *Command) listenHot() (socket.Listener, error) {

// }

func disabled(flag string) bool {
	switch flag {
	case "false", "0":
		return true
	default:
		return false
	}
}
