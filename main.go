package main

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"

	"gitlab.com/mnm/bud/internal/gobin"

	"gitlab.com/mnm/bud/internal/generator/command"
	"gitlab.com/mnm/bud/internal/generator/controller"
	"gitlab.com/mnm/bud/internal/generator/generator"
	"gitlab.com/mnm/bud/internal/generator/gomod"
	"gitlab.com/mnm/bud/internal/generator/maingo"
	"gitlab.com/mnm/bud/internal/generator/web"

	"gitlab.com/mnm/bud/fsync"
	"gitlab.com/mnm/bud/vfs"

	"gitlab.com/mnm/bud/gen"
	"gitlab.com/mnm/bud/internal/generator/generate"

	"gitlab.com/mnm/bud/go/mod"

	"gitlab.com/mnm/bud/commander"

	"gitlab.com/mnm/bud/log/console"
)

func main() {
	if err := do(); err != nil {
		console.Error(err.Error())
		os.Exit(1)
	}
}

func do() error {
	cmd := &bud{}
	cli := commander.New("bud")
	cli.Flag("chdir", "Change the working directory").Short('C').String(&cmd.Chdir).Default(".")

	{
		cmd := &runCommand{bud: cmd}
		cli := cli.Command("run", "run the development server")
		cli.Flag("embed", "embed the assets").Bool(&cmd.Embed).Default(false)
		cli.Flag("hot", "hot reload the frontend").Bool(&cmd.Hot).Default(true)
		cli.Flag("minify", "minify the assets").Bool(&cmd.Minify).Default(false)
		cli.Flag("port", "port").Int(&cmd.Port).Default(3000)
		cli.Run(cmd.Run)
	}

	{
		cmd := &buildCommand{bud: cmd}
		cli := cli.Command("build", "build the production server")
		cli.Flag("embed", "embed the assets").Bool(&cmd.Embed).Default(true)
		cli.Flag("hot", "hot reload the frontend").Bool(&cmd.Hot).Default(false)
		cli.Flag("minify", "minify the assets").Bool(&cmd.Minify).Default(true)
		cli.Run(cmd.Run)
	}

	cli.Arg("command", "custom command").String(&cmd.Custom)
	cli.Run(cmd.Run)

	return cli.Parse(os.Args[1:])
}

type bud struct {
	Chdir  string
	Custom string
}

func (c *bud) Run(ctx context.Context) error {
	modfile, err := mod.Find(c.Chdir)
	if err != nil {
		return err
	}
	fmt.Println(c.Custom+"ing...", modfile.Directory())
	genfs := gen.New(os.DirFS(modfile.Directory()))
	genfs.Add(map[string]gen.Generator{
		"bud/command/main.go": gen.FileGenerator(&command.Generator{
			// fill in
		}),
		"go.mod": gen.FileGenerator(&gomod.Generator{
			Modfile: modfile,
			Go: &gomod.Go{
				Version: "1.17",
			},
			Requires: []*gomod.Require{
				{
					Path:    `gitlab.com/mnm/bud`,
					Version: `v0.0.0-20211017185247-da18ff96a31f`,
				},
			},
			// TODO: remove
			Replaces: []*gomod.Replace{
				{
					Old: "gitlab.com/mnm/bud",
					New: "../../bud",
				},
			},
		}),
	})
	// Sync genfs
	if err := fsync.Dir(genfs, ".", vfs.OS(modfile.Directory()), "."); err != nil {
		return err
	}
	// If bud/command/main.go doesn't exist, run the welcome server
	commandPath := filepath.Join(modfile.Directory(), "bud", "command", "main.go")
	if _, err := os.Stat(commandPath); err != nil {
		if !errors.Is(err, fs.ErrNotExist) {
			return err
		}
		return fmt.Errorf("unknown command %q", c.Custom)
	}
	// Run the command
	// TODO: pass all arguments through
	if err := gobin.Run(ctx, modfile.Directory(), commandPath, c.Custom); err != nil {
		return err
	}
	return nil
}

type runCommand struct {
	bud    *bud
	Embed  bool
	Hot    bool
	Minify bool
	Port   int
	Args   []string
}

func (c *runCommand) Run(ctx context.Context) error {
	modfile, err := mod.Find(c.bud.Chdir)
	if err != nil {
		return err
	}
	fmt.Println("running...", modfile.Directory(), c.Embed, c.Hot, c.Minify)
	genfs := gen.New(os.DirFS(modfile.Directory()))
	genfs.Add(map[string]gen.Generator{
		"bud/generate/main.go": gen.FileGenerator(&generate.Generator{
			Modfile: modfile,
			Embed:   c.Embed,
			Hot:     c.Hot,
			Minify:  c.Minify,
		}),
		"bud/generator/generator.go": gen.FileGenerator(&generator.Generator{
			// fill in
		}),
		"go.mod": gen.FileGenerator(&gomod.Generator{
			Modfile: modfile,
			Go: &gomod.Go{
				Version: "1.17",
			},
			Requires: []*gomod.Require{
				{
					Path:    `gitlab.com/mnm/bud`,
					Version: `v0.0.0-20211017185247-da18ff96a31f`,
				},
			},
			// TODO: remove
			Replaces: []*gomod.Replace{
				{
					Old: "gitlab.com/mnm/bud",
					New: "../../bud",
				},
			},
		}),
		"bud/controller/controller.go": gen.FileGenerator(&controller.Generator{
			// fill in
		}),
		"bud/web/web.go": gen.FileGenerator(&web.Generator{
			// fill in
		}),
		"bud/main.go": gen.FileGenerator(&maingo.Generator{
			// fill in
		}),
	})
	// Sync genfs
	if err := fsync.Dir(genfs, ".", vfs.OS(modfile.Directory()), "."); err != nil {
		return err
	}
	// Run generate (if it exists) to support user-defined generators
	generatePath := filepath.Join(modfile.Directory(), "bud", "generate", "main.go")
	if _, err := os.Stat(generatePath); nil == err {
		if err := gobin.Run(ctx, modfile.Directory(), generatePath); err != nil {
			return err
		}
	}
	// If bud/main.go doesn't exist, run the welcome server
	mainPath := filepath.Join(modfile.Directory(), "bud", "main.go")
	if _, err := os.Stat(mainPath); err != nil {
		if !errors.Is(err, fs.ErrNotExist) {
			return err
		}
		// TODO: improve the welcome server
		address := fmt.Sprintf(":%d", c.Port)
		console.Info("Listening on http://localhost%s", address)
		return http.ListenAndServe(address, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("Welcome!\n"))
		}))
	}
	// Run the main server
	if err := gobin.Run(ctx, modfile.Directory(), mainPath); err != nil {
		return err
	}
	return nil
}

type buildCommand struct {
	bud    *bud
	Embed  bool
	Hot    bool
	Minify bool
}

func (c *buildCommand) Run(ctx context.Context) error {
	modfile, err := mod.Find(c.bud.Chdir)
	if err != nil {
		return err
	}
	fmt.Println("building...", modfile.Directory(), c.Embed, c.Hot, c.Minify)
	genfs := gen.New(os.DirFS(modfile.Directory()))
	genfs.Add(map[string]gen.Generator{
		"bud/generate/main.go": gen.FileGenerator(&generate.Generator{
			Modfile: modfile,
			Embed:   c.Embed,
			Hot:     c.Hot,
			Minify:  c.Minify,
		}),
		"bud/generator/generator.go": gen.FileGenerator(&generator.Generator{
			// fill in
		}),
		"go.mod": gen.FileGenerator(&gomod.Generator{
			Modfile: modfile,
			Go: &gomod.Go{
				Version: "1.17",
			},
			Requires: []*gomod.Require{
				{
					Path:    `gitlab.com/mnm/bud`,
					Version: `v0.0.0-20211017185247-da18ff96a31f`,
				},
			},
			// TODO: remove
			Replaces: []*gomod.Replace{
				{
					Old: "gitlab.com/mnm/bud",
					New: "../../bud",
				},
			},
		}),
		"bud/controller/controller.go": gen.FileGenerator(&controller.Generator{
			// fill in
		}),
		"bud/web/web.go": gen.FileGenerator(&web.Generator{
			// fill in
		}),
		"bud/main.go": gen.FileGenerator(&maingo.Generator{
			// fill in
		}),
	})
	// Sync genfs
	if err := fsync.Dir(genfs, ".", vfs.OS(modfile.Directory()), "."); err != nil {
		return err
	}
	// Run generate (if it exists) to support user-defined generators
	generatePath := filepath.Join(modfile.Directory(), "bud", "generate", "main.go")
	if _, err := os.Stat(generatePath); nil == err {
		if err := gobin.Run(ctx, modfile.Directory(), generatePath); err != nil {
			return err
		}
	}
	// Verify that bud/main.go exists
	mainPath := filepath.Join(modfile.Directory(), "bud", "main.go")
	if _, err := os.Stat(mainPath); err != nil {
		return err
	}
	// Build the main server
	outPath := filepath.Join(modfile.Directory(), "bud", "main")
	if err := gobin.Build(ctx, modfile.Directory(), mainPath, outPath); err != nil {
		return err
	}
	return nil

	return nil
}
