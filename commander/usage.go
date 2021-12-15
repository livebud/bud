package commander

import (
	"bytes"
	"sort"
	"strings"
	"text/tabwriter"
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

type generateCommands []*generateCommand

func (cmds generateCommands) Usage() (string, error) {
	buf := new(bytes.Buffer)
	tw := tabwriter.NewWriter(buf, 0, 0, 2, ' ', 0)
	for _, cmd := range cmds {
		tw.Write([]byte("\t\t" + cmd.c.name + "\t" + dim() + cmd.c.usage + reset() + "\n"))
	}
	if err := tw.Flush(); err != nil {
		return "", err
	}
	return strings.TrimSpace(buf.String()), nil
}

func (g *generateCommand) Commands() (commands generateCommands) {
	commands = make(generateCommands, len(g.c.commands))
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

func (g *generateCommand) Flags() (flags generateFlags) {
	flags = make(generateFlags, len(g.c.flags))
	for i, flag := range g.c.flags {
		flags[i] = &generateFlag{flag}
	}
	// Sort by name
	sort.Slice(flags, func(i, j int) bool {
		if hasShort(flags[i]) == hasShort(flags[j]) {
			// Both have shorts or don't have shorts, so sort by name
			return flags[i].f.name < flags[j].f.name
		}
		// Shorts above non-shorts
		return flags[i].f.short > flags[j].f.short
	})
	return flags
}

func hasShort(flag *generateFlag) bool {
	return flag.f.short != 0
}

type generateFlag struct {
	f *Flag
}

func (g *generateFlag) Name() string {
	return g.f.name
}

type generateFlags []*generateFlag

func (flags generateFlags) hasShortFlags() bool {
	for _, flag := range flags {
		if flag.f.short != 0 {
			return true
		}
	}
	return false
}

func (flags generateFlags) Usage() (string, error) {
	buf := new(bytes.Buffer)
	tw := tabwriter.NewWriter(buf, 0, 0, 2, ' ', 0)
	for _, flag := range flags {
		tw.Write([]byte("\t\t"))
		if flag.f.short != 0 {
			tw.Write([]byte("-" + string(flag.f.short) + ", "))
		}
		tw.Write([]byte("--" + flag.f.name))
		tw.Write([]byte("\t"))
		tw.Write([]byte(dim() + flag.f.usage + reset()))
		tw.Write([]byte("\n"))
	}
	if err := tw.Flush(); err != nil {
		return "", err
	}
	return strings.TrimSpace(buf.String()), nil
}
