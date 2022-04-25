package main

import (
	"context"
	_ "embed"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/livebud/bud/internal/gotemplate"
	"github.com/livebud/bud/package/commander"
	"github.com/livebud/bud/package/gomod"
	"github.com/livebud/bud/package/log/console"
)

//go:embed changelog.gotext
var changelog string

var generator = gotemplate.MustParse("changelog.gotext", changelog)

func main() {
	if err := run(); err != nil {
		console.Error(err.Error())
		os.Exit(1)
	}
	os.Exit(0)
}

func run() error {
	cmd := new(Command)
	cli := commander.New("generate-changelog")
	cli.Arg("version").String(&cmd.Version)
	cli.Run(cmd.Run)
	return cli.Parse(context.Background(), os.Args[1:])
}

type Command struct {
	Version string
}

type State struct {
	Notes string
}

func (c *Command) Run(ctx context.Context) error {
	dir, err := gomod.Absolute(".")
	if err != nil {
		return err
	}
	data, err := os.ReadFile(filepath.Join(dir, "Changelog.md"))
	if err != nil {
		return err
	}
	contents := string(data)
	marker := "\n## " + c.Version
	startIndex := strings.Index(contents, marker)
	if startIndex < 0 {
		return fmt.Errorf("generate-changelog: Unable to find version %s", c.Version)
	}
	notes := contents[startIndex+len(marker):]
	endIndex := strings.Index(notes, "\n## ")
	if endIndex > 0 {
		notes = notes[0:endIndex]
	}
	changelog, err := generator.Generate(&State{
		Notes: notes,
	})
	if err != nil {
		return err
	}
	fmt.Fprintln(os.Stdout, strings.TrimSpace(string(changelog)))
	return nil
}
