package commander

import (
	"bytes"
	"sort"
	"text/template"
)

func generateUsage(template *template.Template, c *Command) (string, error) {
	buf := new(bytes.Buffer)
	if err := template.Execute(buf, &generateCommand{c}); err != nil {
		return "", err
	}
	return buf.String(), nil
}

type generateCommand struct {
	c *Command
}

func (g *generateCommand) Name() string {
	return g.c.name
}

func (g *generateCommand) Usage() string {
	return g.c.usage
}

func (g *generateCommand) Commands() (commands []*generateCommand) {
	commands = make([]*generateCommand, len(g.c.commands))
	i := 0
	for _, cmd := range g.c.commands {
		commands[i] = &generateCommand{cmd}
		i++
	}
	// Sort by name
	sort.Slice(commands, func(i, j int) bool {
		return commands[i].c.name < commands[j].c.name
	})
	return commands
}

func (g *generateCommand) Flags() (flags []*generateFlag) {
	flags = make([]*generateFlag, len(g.c.flags))
	for i, flag := range g.c.flags {
		flags[i] = &generateFlag{flag}
	}
	return flags
}

type generateFlag struct {
	f *Flag
}

func (g *generateFlag) Name() string {
	return g.f.name
}

func (g *generateFlag) Usage() string {
	return g.f.usage
}

func (g *generateFlag) Short() string {
	if g.f.short == 0 {
		return ""
	}
	return string(g.f.short)
}
