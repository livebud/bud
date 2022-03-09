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
		tw.Write([]byte("\t\t" + cmd.c.name))
		if cmd.c.usage != "" {
			tw.Write([]byte("\t" + dim() + cmd.c.usage + reset()))
		}
		tw.Write([]byte("\n"))
	}
	if err := tw.Flush(); err != nil {
		return "", err
	}
	return strings.TrimSpace(buf.String()), nil
}

func (g *generateCommand) Args() (args []string) {
	for i, arg := range g.c.args {
		// TODO: differentiate between required and optional args
		if i == 0 && len(g.c.commands) > 0 {
			args = append(args, "<command|"+arg.Name+">")
			continue
		}
		args = append(args, "<"+arg.Name+">")
	}
	if len(args) == 0 && len(g.c.commands) > 0 {
		args = append(args, "[command]")
	}
	return args
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

func (flags generateFlags) Usage() (string, error) {
	buf := new(bytes.Buffer)
	tw := tabwriter.NewWriter(buf, 0, 0, 2, ' ', 0)
	for _, flag := range flags {
		tw.Write([]byte("\t\t"))
		if flag.f.short != 0 {
			tw.Write([]byte("-" + string(flag.f.short) + ", "))
		}
		tw.Write([]byte("--" + flag.f.name))
		if flag.f.usage != "" {
			tw.Write([]byte("\t"))
			tw.Write([]byte(dim() + flag.f.usage + reset()))
		}
		tw.Write([]byte("\n"))
	}
	if err := tw.Flush(); err != nil {
		return "", err
	}
	return strings.TrimSpace(buf.String()), nil
}
