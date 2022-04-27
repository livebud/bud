package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/livebud/bud/package/gomod"

	"github.com/livebud/bud/package/commander"
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
	cli := commander.New("test-cache-diff")
	cli.Arg("left").String(&cmd.Left)
	cli.Arg("right").String(&cmd.Right)
	cli.Run(cmd.Run)
	return cli.Parse(context.Background(), os.Args[1:])
}

type Command struct {
	Left  string
	Right string
}

func (c *Command) Run(ctx context.Context) error {
	dir, err := gomod.Absolute(".")
	if err != nil {
		return err
	}
	leftMap, err := fillMap(dir, c.Left)
	if err != nil {
		return err
	}
	rightMap, err := fillMap(dir, c.Right)
	if err != nil {
		return err
	}
	fmt.Println(len(leftMap), len(rightMap))
	for path, lkey := range leftMap {
		rkey, ok := rightMap[path]
		if !ok {
			fmt.Println(c.Right, "missing", path)
			continue
		}
		if lkey != rkey {
			fmt.Println(path, lkey, "!=", rkey)
			continue
		}
	}
	for path, _ := range rightMap {
		_, ok := leftMap[path]
		if !ok {
			fmt.Println(c.Left, "missing", path)
			continue
		}
	}
	return nil
}

func fillMap(dir, path string) (map[string]string, error) {
	m := map[string]string{}
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	s := bufio.NewScanner(f)
	for s.Scan() {
		line := s.Text()
		if !strings.HasPrefix(line, "HASH[testInputs]:") {
			continue
		}
		parts := strings.Split(line, " ")
		for i, part := range parts {
			if strings.HasPrefix(part, dir) {
				m[part] = parts[i+1]
				break
			}
		}
	}
	if err := s.Err(); err != nil {
		return nil, err
	}
	return m, nil
}
