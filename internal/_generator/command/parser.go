package command

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/matthewmueller/gotext"

	"github.com/livebud/bud/internal/bail"
	"github.com/livebud/bud/internal/imports"
	"github.com/livebud/bud/internal/valid"
	"github.com/livebud/bud/package/di"
	"github.com/livebud/bud/package/gomod"
	goparse "github.com/livebud/bud/package/parser"
	"github.com/matthewmueller/text"
)

type parser struct {
	bail.Struct
	fs       fs.FS
	imports  *imports.Set
	injector *di.Injector
	module   *gomod.Module
	parser   *goparse.Parser
}

func (p *parser) Parse(ctx context.Context) (state *State, err error) {
	defer p.Recover2(&err, "command: unable to parse")
	// Default imports
	p.imports.AddStd("context", "os")
	p.imports.AddNamed("commander", "github.com/livebud/bud/package/commander")
	p.imports.AddNamed("gomod", "github.com/livebud/bud/package/gomod")
	p.imports.AddNamed("run", "github.com/livebud/bud/runtime/command/run")
	p.imports.AddNamed("new_controller", "github.com/livebud/bud/runtime/command/new/controller")
	p.imports.AddNamed("build", "github.com/livebud/bud/runtime/command/build")
	state = new(State)
	state.Provider = p.loadProvider()
	state.Imports = p.imports.List()
	return state, nil
}

func (p *parser) Parse2(ctx context.Context) (state *State, err error) {
	defer p.Recover2(&err, "command: unable to parse")
	// Default imports
	p.imports.AddStd("context")
	state = new(State)
	state.Command = p.loadCommand2(nil, "command", "")
	state.Imports = p.imports.List()
	return state, nil
}

func (p *parser) loadCommand2(parent *Cmd, base, dir string) *Cmd {
	cmd := new(Cmd)
	commandDir := filepath.Join(base, dir)
	// Traverse the subdirectories
	des, err := fs.ReadDir(p.fs, commandDir)
	if err != nil {
		if !errors.Is(err, fs.ErrNotExist) {
			p.Bail(err)
		}
		return cmd
	}
	shouldParse := false
	for _, de := range des {
		// Valid command file
		if !de.IsDir() {
			shouldParse = shouldParse || valid.CommandFile(de.Name())
			continue
		}
		// 	// Valid sub-directory
		// 	if !valid.Dir(de.Name()) {
		// 		continue
		// 	}
		// 	// Valid sub-command dir

		// 	if de.IsDir() && valid.Dir(de.Name()) {
		// 		sub := p.loadCommand(command, base, filepath.Join(dir, de.Name()))
		// 		if sub == nil {
		// 			continue
		// 		}
		// 		if hasNameConflict(command, sub) {
		// 			p.Bail(fmt.Errorf("command and subcommand cannot have conflicting names %s", command.Name))
		// 		}
		// 		command.Subs = append(command.Subs, sub)
		// 	}
	}
	return cmd
}

func (p *parser) loadProvider() *di.Provider {
	provider, err := p.injector.Wire(&di.Function{
		Name:   "loadGenerator",
		Target: p.module.Import("bud/.cli/command"),
		Params: []di.Dependency{
			di.ToType("github.com/livebud/bud/package/gomod", "*Module"),
			di.ToType("context", "Context"),
			di.ToType("github.com/livebud/bud/runtime/command", "*Flag"),
		},
		Results: []di.Dependency{
			di.ToType(p.module.Import("bud/.cli/generator"), "*FileSystem"),
			&di.Error{},
		},
		Aliases: di.Aliases{
			di.ToType("github.com/livebud/bud/package/js", "VM"):          di.ToType("github.com/livebud/bud/package/js/v8client", "*Client"),
			di.ToType("io/fs", "FS"):                                      di.ToType("github.com/livebud/bud/package/overlay", "*FileSystem"),
			di.ToType("github.com/livebud/bud/runtime/transform", "*Map"): di.ToType(p.module.Import("bud/.cli/transform"), "*Map"),
		},
	})
	if err != nil {
		p.Bail(err)
	}
	// Add the imports
	for _, im := range provider.Imports {
		p.imports.AddNamed(im.Name, im.Path)
	}
	return provider
}

// func isValidSubDir(de fs.DirEntry) bool {
// 	return de.IsDir() && valid.Dir(de.Name())
// }

func (p *parser) loadCommand(parent *Cmd, base, dir string) *Cmd {
	commandDir := filepath.Join(base, dir)
	command := new(Cmd)
	// Get the command name
	name := filepath.Base(dir)
	if name == "." {
		name = ""
	}
	command.Name = name
	// Load the import
	importPath := p.module.Import(filepath.SplitList(commandDir)...)
	importName := p.imports.Add(importPath)
	command.Import = &imports.Import{
		Name: importName,
		Path: importPath,
	}
	// Parse the command directory
	pkg, err := p.parser.Parse(commandDir)
	if err != nil {
		p.Bail(err)
	}
	// Find the Command struct, ignore if it doesn't exist
	stct := pkg.Struct("Command")
	if stct == nil {
		return nil
	}
	// Check that the command is runnable
	command.Runnable = stct.Method(gotext.Pascal(name)) != nil
	// Create subcommands for each of the public methods
	for _, method := range stct.PublicMethods() {
		sub := p.loadSub(command, method)
		if sub == nil {
			continue
		}
		if hasNameConflict(command, sub) {
			p.Bail(fmt.Errorf("command and subcommand cannot have conflicting names %s", command.Name))
		}
		// Same import because they're in the same file
		sub.Import = command.Import
		command.Subs = append(command.Subs, sub)
	}
	// Traverse the subdirectories
	des, err := fs.ReadDir(p.fs, commandDir)
	if err != nil {
		p.Bail(err)
	}
	// Load the subcommands
	for _, de := range des {
		if !de.IsDir() || !valid.Dir(de.Name()) {
			continue
		}
		sub := p.loadCommand(command, base, filepath.Join(dir, de.Name()))
		if sub == nil {
			continue
		}
		if hasNameConflict(command, sub) {
			p.Bail(fmt.Errorf("command and subcommand cannot have conflicting names %s", command.Name))
		}
		command.Subs = append(command.Subs, sub)
	}
	return command
}

func hasNameConflict(parent *Cmd, child *Cmd) bool {
	return parent.Runnable && child.Runnable && parent.Name == child.Name
}

// Load subcommand
func (p *parser) loadSub(parent *Cmd, method *goparse.Function) *Cmd {
	sub := &Cmd{
		Name:     method.Name(),
		Parent:   parent,
		Runnable: true,
	}
	for i, param := range method.Params() {
		paramType := param.Type()
		// Load the context
		if p.isContext(i, paramType) {
			sub.Context = true
			continue
		}
		// Builtins get mapped to args
		if builtin := goparse.IsBuiltin(paramType); builtin {
			sub.Args = append(sub.Args, p.loadArg(param))
			continue
		}
		// Structs gets mapped to flags
		def, err := goparse.Definition(paramType)
		if err != nil {
			p.Bail(err)
		}
		stct := def.Package().Struct(def.Name())
		if stct != nil {
			sub.Flags = append(sub.Flags, p.loadFlags(stct)...)
			continue
		}
		p.Bail(fmt.Errorf("unable to handle param %s", param.Type().String()))
	}
	return sub
}

func (p *parser) isContext(nth int, paramType goparse.Type) bool {
	isContext, err := goparse.IsImportType(paramType, "context", "Context")
	if err != nil {
		p.Bail(err)
	}
	if !isContext {
		return false
	}
	// Context is only supported on the first
	if nth > 0 {
		p.Bail(fmt.Errorf("context must be the first argument"))
	}
	return isContext
}

func (p *parser) loadArg(param *goparse.Param) *Arg {
	arg := &Arg{
		Name: param.Name(),
	}
	ts := param.Type().String()
	// Check if optional
	if strings.HasPrefix(ts, "*") {
		arg.Optional = true
		ts = strings.TrimPrefix(ts, "*")
	}
	switch ts {
	case "string", "...string", "bool", "int":
		// use original type string to support optionals
		arg.Type = param.Type().String()
	default:
		p.Bail(fmt.Errorf("command: arg must be a string, ...string, bool or int, not %q", ts))
	}
	return arg
}

func (p *parser) loadFlags(stct *goparse.Struct) (flags []*Flag) {
	for _, field := range stct.PublicFields() {
		flags = append(flags, p.loadFlag(field))
	}
	return flags
}

// Load the command flag
func (p *parser) loadFlag(field *goparse.Field) *Flag {
	flag := new(Flag)
	tags, err := field.Tags()
	if err != nil {
		p.Bail(err)
	}
	// Set the flag name
	flag.Name = tags.Get("flag")
	// Use the field name (flags are the default)
	if flag.Name == "" {
		flag.Name = text.Slug(field.Name())
	}
	// Set the short flag
	short := tags.Get("short")
	if short != "" {
		if len(short) != 1 {
			p.Bail(fmt.Errorf("command: short flag must be exactly 1 byte long, not %q", short))
		}
		flag.Short = short[0]
	}
	ts := field.Type().String()
	// Check if optional
	if strings.HasPrefix(ts, "*") {
		flag.Optional = true
		ts = strings.TrimPrefix(ts, "*")
	}
	// Set the flag type
	// TODO: add more types
	switch ts {
	case "string", "bool", "int":
		// use original type string to support optionals
		flag.Type = field.Type().String()
	default:
		p.Bail(fmt.Errorf("command: flag must be a string, bool or int, not %q", ts))
	}
	// Set the usage
	flag.Help = tags.Get("help")
	// Set the default
	if def := tags.Get("default"); def != "" {
		if flag.Type == "string" {
			def = strconv.Quote(def)
		}
		flag.Default = &def
	}
	return flag
}

// Load the command arg
// func (p *parser) loadCommandArg(field *goparse.Field, tags goparse.Tags) *Arg {
// 	arg := new(Arg)
// 	// Set the arg name
// 	arg.Name = tags.Get("arg")
// 	// Set the arg type
// 	switch t := field.Type().String(); t {
// 	case "string", "bool", "int":
// 		arg.Type = t
// 	default:
// 		err := fmt.Errorf("command: arg must be a string, bool or int not %q", t)
// 		p.Bail(err)
// 	}
// 	return arg
// }

// // Returns true if the command has a run command
// func isRunnable(stct *goparse.Struct) bool {
// 	run := stct.Method("Run")
// 	if run == nil {
// 		return false
// 	}
// 	return true
// 	// params := run.Params()
// 	// if len(params) != 1 {
// 	// 	return false
// 	// }
// 	// // Ensure the first param is context.Context
// 	// typeName := goparse.TypeName(params[0].Type())
// 	// return typeName == "Context"
// }

// // Check if the type is a context
// func isContext(t goparse.Type) (bool, error) {

// }

// func findStruct(t goparse.Type) (bool, error) {
// }
