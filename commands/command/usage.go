package command

import (
	"bytes"
	_ "embed"
	"strings"
	"text/tabwriter"
)

// var usageGenerator = template.Must(template.New("usage").Funcs(colors).Parse(usage))

// type usageState struct {
// 	Name     string
// 	Flags    usageFlags
// 	Commands usageCommands
// 	Args     usageArgs
// }

// func generateUsage(template *template.Template, c *Command) (string, error) {
// 	buf := new(bytes.Buffer)
// 	if err := template.Execute(buf, &generateCommand{c}); err != nil {
// 		return "", err
// 	}
// 	return buf.String(), nil
// }

// type generateCommand struct {
// 	c *Command
// }

// func (g *generateCommand) Name() string {
// 	return g.c.Name
// }

// type generateCommands []*generateCommand

// func (cmds generateCommands) Usage() (string, error) {
// 	buf := new(bytes.Buffer)
// 	tw := tabwriter.NewWriter(buf, 0, 0, 2, ' ', 0)
// 	for _, cmd := range cmds {
// 		tw.Write([]byte("\t\t" + cmd.c.Name))
// 		if cmd.c.Help != "" {
// 			tw.Write([]byte("\t" + dim() + cmd.c.Help + reset()))
// 		}
// 		tw.Write([]byte("\n"))
// 	}
// 	if err := tw.Flush(); err != nil {
// 		return "", err
// 	}
// 	return strings.TrimSpace(buf.String()), nil
// }

// func (g *generateCommand) Args() (args []string) {
// 	for i, arg := range g.c.Args {
// 		// TODO: differentiate between required and optional args
// 		if i == 0 && len(g.c.Commands) > 0 {
// 			args = append(args, "<command|"+arg.Name+">")
// 			continue
// 		}
// 		args = append(args, "<"+arg.Name+">")
// 	}
// 	if len(args) == 0 && len(g.c.Commands) > 0 {
// 		args = append(args, "[command]")
// 	}
// 	return args
// }

// func (g *generateCommand) Commands() (commands generateCommands) {
// 	commands = make(generateCommands, len(g.c.Commands))
// 	// i := 0
// 	// for _, cmd := range g.c.Commands {
// 	// 	commands[i] = &generateCommand{cmd}
// 	// 	i++
// 	// }
// 	// // Sort by name
// 	// sort.Slice(commands, func(i, j int) bool {
// 	// 	return commands[i].c.Name < commands[j].c.Name
// 	// })
// 	return commands
// }

// func (g *generateCommand) Flags() (flags generateFlags) {
// 	flags = make(generateFlags, len(g.c.Flags))
// 	for i, flag := range g.c.Flags {
// 		flags[i] = &generateFlag{flag}
// 	}
// 	// Sort by name
// 	sort.Slice(flags, func(i, j int) bool {
// 		if hasShort(flags[i]) == hasShort(flags[j]) {
// 			// Both have shorts or don't have shorts, so sort by name
// 			return flags[i].f.Name < flags[j].f.Name
// 		}
// 		// Shorts above non-shorts
// 		return flags[i].f.Short > flags[j].f.Short
// 	})
// 	return flags
// }

// func hasShort(flag *generateFlag) bool {
// 	return flag.f.Short != 0
// }

type usageFlag struct {
	f *Flag
}

func (g *usageFlag) Name() string {
	return g.f.Name
}

type usageFlags []*usageFlag

func (flags usageFlags) Usage() (string, error) {
	buf := new(bytes.Buffer)
	tw := tabwriter.NewWriter(buf, 0, 0, 2, ' ', 0)
	for _, flag := range flags {
		tw.Write([]byte("\t\t"))
		if flag.f.Short != 0 {
			tw.Write([]byte("-" + string(flag.f.Short) + ", "))
		}
		tw.Write([]byte("--" + flag.f.Name))
		if flag.f.Help != "" {
			tw.Write([]byte("\t"))
			tw.Write([]byte(dim() + flag.f.Help + reset()))
		}
		tw.Write([]byte("\n"))
	}
	if err := tw.Flush(); err != nil {
		return "", err
	}
	return strings.TrimSpace(buf.String()), nil
}
