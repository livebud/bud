package generate

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"path"
	"path/filepath"
	"strings"

	"github.com/livebud/bud/framework/afs"
	generator "github.com/livebud/bud/framework/generator2"
	"github.com/livebud/bud/package/di"
	"github.com/livebud/bud/package/genfs"
	"github.com/livebud/bud/package/gomod"
	"github.com/livebud/bud/package/log"
	"github.com/livebud/bud/package/parser"
	"github.com/livebud/bud/package/remotefs"

	"github.com/livebud/bud/framework"
	"github.com/livebud/bud/internal/cli/bud"
	"github.com/livebud/bud/internal/dsync"
	"github.com/livebud/bud/internal/exe"
	"github.com/livebud/bud/internal/fscache"
)

// New command for bud generate
func New(bud *bud.Command, in *bud.Input) *Command {
	return &Command{
		bud: bud,
		in:  in,
		Flag: &framework.Flag{
			Env:    in.Env,
			Stderr: in.Stderr,
			Stdin:  in.Stdin,
			Stdout: in.Stdout,
		},
	}
}

// Command for running bud generate
type Command struct {
	bud  *bud.Command
	in   *bud.Input
	Flag *framework.Flag
	Args []string
}

// Run the generate command
func (c *Command) Run(ctx context.Context) error {
	log, err := bud.Log(c.in.Stderr, c.bud.Log)
	if err != nil {
		return err
	}
	module, err := gomod.Find(c.bud.Dir)
	if err != nil {
		return err
	}
	exec := &exe.Template{
		Dir:    module.Directory(),
		Stdin:  c.in.Stdin,
		Stdout: c.in.Stdout,
		Stderr: c.in.Stderr,
		Env:    c.in.Env,
	}
	afs, err := loadAFS(exec, c.Flag, log, module)
	if err != nil {
		return err
	}
	defer afs.Close()
	return afs.Generate(ctx, c.Args...)
}

func loadAFS(exec *exe.Template, flag *framework.Flag, log log.Log, module *gomod.Module) (*AFS, error) {
	// Create the generator filesystem
	cache := fscache.Discard
	genfs := genfs.New(cache, module, log)
	// Load the parser
	parser := parser.New(genfs, module)
	// Load the injector
	injector := di.New(genfs, log, module, parser)
	// Setup the initial file generators
	gen := generator.New(log, module, parser)
	state, err := generator.Load(genfs, log, module, parser)
	if err != nil {
		return nil, err
	}
	// TODO: mount all of the generators that proxy to the remote filesystem
	fmt.Println("got state", state)
	afs := afs.New(exec, flag, injector, log, module)
	genfs.FileGenerator("bud/internal/generator/generator.go", gen)
	genfs.FileGenerator("bud/cmd/afs/main.go", afs)
	// Return the new afs controller
	return &AFS{cache, exec, flag, genfs, log, module, nil}, nil
}

type AFS struct {
	cache   fscache.Cache
	exec    *exe.Template
	flag    *framework.Flag
	genfs   genfs.FileSystem
	log     log.Log
	module  *gomod.Module
	process *remotefs.Process // nil when afs is not running
}

func only(dirs ...string) func(name string, isDir bool) bool {
	m := map[string]bool{}
	for _, dir := range dirs {
		for dir != "." {
			m[dir] = true
			dir = path.Dir(dir)
		}
	}
	hasPrefix := func(name string) bool {
		for _, dir := range dirs {
			if strings.HasPrefix(name, dir) {
				return true
			}
		}
		return false
	}
	return func(name string, isDir bool) bool {
		return !m[name] && !hasPrefix(name)
	}
}

func (a *AFS) Generate(ctx context.Context, dirs ...string) error {
	skips := []func(name string, isDir bool) bool{}
	// Skip hidden files and directories
	skips = append(skips, func(name string, isDir bool) bool {
		base := filepath.Base(name)
		return base[0] == '_' || base[0] == '.'
	})
	// Skip files that aren't within these directories
	if len(dirs) > 0 {
		skips = append(skips, only(dirs...))
	}
	// Sync the files
	if err := dsync.To(a.genfs, a.module, "bud", dsync.WithSkip(skips...)); err != nil {
		if !errors.Is(err, fs.ErrNotExist) {
			return err
		}
	}
	// Check if the main.go file exists, otherwise return early
	if _, err := fs.Stat(a.genfs, "bud/cmd/afs/main.go"); err != nil {
		if !errors.Is(err, fs.ErrNotExist) {
			return err
		}
		return nil
	}
	// Build the main.go file
	cmd := a.exec.Command("go", "build",
		"-o", "bud/afs",
		"-mod=mod",
		"bud/cmd/afs/main.go",
	)
	if err := cmd.Run(); err != nil {
		return err
	}
	// Start the application file server
	if a.process != nil {
		if err := a.process.Restart(ctx); err != nil {
			return err
		}
	} else {
		// TODO: pass flags through
		process, err := remotefs.Start(ctx, a.exec, "bud/afs")
		if err != nil {
			return err
		}
		a.process = process
	}
	// Sync the remote file server
	if err := dsync.To(a.process, a.module, "bud", dsync.WithSkip(skips...)); err != nil {
		return err
	}
	return nil
}

func includeBinary(dir []string) bool {
	if len(dir) == 0 {
		return true
	}
	for _, d := range dir {
		if d == "bud/afs" {
			return true
		}
	}
	return false
}

func needsBinary(dirs []string) bool {
	return false
}

func (a *AFS) Reload(paths ...string) error {
	return nil
}

func (a *AFS) Close() error {
	if a.process != nil {
		if err := a.process.Close(); err != nil {
			return err
		}
		a.process = nil
	}
	return nil
}

func newApp(afs *AFS, exec *exe.Template, log log.Log) *App {
	return &App{afs, exec, log}
}

type App struct {
	afs  *AFS
	exec *exe.Template
	log  log.Log
}

func (a *App) Run() error {
	return nil
}

func (a *App) Watch() error {
	return nil
}

func (a *App) Build() error {
	return nil
}
