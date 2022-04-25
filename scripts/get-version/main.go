package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/livebud/bud/internal/npm"
	"github.com/livebud/bud/package/commander"
	"github.com/livebud/bud/package/gomod"
	"github.com/livebud/bud/package/log/console"
)

func main() {
	if err := run(); err != nil {
		console.Error(err.Error())
		os.Exit(1)
	}
	os.Exit(0)
}

func run() error {
	cmd := new(Command)
	cli := commander.New("get-version")
	cli.Arg("key").String(&cmd.Key)
	cli.Run(cmd.Run)
	return cli.Parse(context.Background(), os.Args[1:])
}

type Command struct {
	Key string
}

func (c *Command) Run(ctx context.Context) error {
	dir, err := gomod.Absolute(".")
	if err != nil {
		return err
	}
	data, err := os.ReadFile(filepath.Join(dir, "package.json"))
	if err != nil {
		return err
	}
	pkg := new(npm.Package)
	if err := json.Unmarshal(data, pkg); err != nil {
		return err
	}
	switch {
	case c.Key == "version":
		fmt.Fprint(os.Stdout, pkg.Version)
		return nil
	case strings.HasPrefix(c.Key, "dependencies."):
		dependency := strings.TrimPrefix(c.Key, "dependencies.")
		version := pkg.Dependencies[dependency]
		if version == "" {
			return fmt.Errorf("version not found for %q", c.Key)
		}
		fmt.Fprint(os.Stdout, version)
		return nil
	case strings.HasPrefix(c.Key, "devDependencies."):
		dependency := strings.TrimPrefix(c.Key, "devDependencies.")
		version := pkg.DevDependencies[dependency]
		if version == "" {
			return fmt.Errorf("version not found for %q", c.Key)
		}
		fmt.Fprint(os.Stdout, version)
		return nil
	default:
		return fmt.Errorf("key not support %q yet", c.Key)
	}
}
