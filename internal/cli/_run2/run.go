package run

import (
	"context"
	"fmt"

	"github.com/livebud/bud/internal/budapi"
	"github.com/livebud/bud/internal/cli/bud"
)

func New(cmd *budapi.Command) *Command {
	return &Command{cmd, new(budapi.RunInput)}
}

type Command struct {
	cmd *budapi.Command
	*budapi.RunInput
}

func (c *Command) Run(ctx context.Context) error {
	bud := budapi.New()
	fmt.Println("running???")
	return bud.Run(ctx, c.RunInput)

	// module, err := gomod.Find(c.cfg.Dir)
	// if err != nil {
	// 	return err
	// }
	// // Ensure we have version alignment between the CLI and the runtime
	// // if err := config.EnsureVersionAlignment(ctx, module, versions.Bud); err != nil {
	// // 	return err
	// // }
	// // Setup the logger
	// logLevel, err := log.ParseLevel(c.cfg.Log)
	// if err != nil {
	// 	return err
	// }
	// log := log.New(levelfilter.New(console.New(os.Stderr), logLevel))
	// db, err := dag.Load(log, module.Directory("bud", "bud.db"))
	// if err != nil {
	// 	return err
	// }
	// defer db.Close()
	// pubSub := pubsub.New()
	// v8, err := v8.Load()
	// if err != nil {
	// 	return err
	// }
	// defer v8.Close()
	// devLn, err := socket.Listen(c.cfg.ListenDev)
	// if err != nil {
	// 	return err
	// }
	// defer devLn.Close()
	// devLauncher := &dev.Launcher{
	// 	Listener: devLn,
	// 	Pubsub:   pubSub,
	// 	VM:       v8,
	// }
	// devProcess, err := devLauncher.Launch(ctx)
	// if err != nil {
	// 	return err
	// }
	// defer devProcess.Close()
	// afsLn, err := socket.Listen(c.cfg.ListenAFS)
	// if err != nil {
	// 	return err
	// }
	// defer afsLn.Close()
	// afs, err := afs.Launch()
	// afsLauncher := &afs.Launcher{
	// 	Cache:    db,
	// 	Config:   c.cfg,
	// 	Listener: afsLn,
	// 	Log:      log,
	// 	Module:   module,
	// 	VM:       devProcess,
	// }
	// afsProcess, err := afsLauncher.Launch(ctx)
	// if err != nil {
	// 	return err
	// }
	// defer afsProcess.Close()

	// // app := app.New(c.cfg, afs, dev, db)

	// // app.Generate()

	// // c.cfg.WebLn
	// log.Info("hi")
	// fmt.Println("hi", module.Directory())
	return nil
}
