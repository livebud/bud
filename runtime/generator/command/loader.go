package command

import (
	"errors"
	"fmt"
	"io/fs"
	"path/filepath"

	"github.com/matthewmueller/text"

	"github.com/livebud/bud/internal/bail"
	"github.com/livebud/bud/internal/imports"
	"github.com/livebud/bud/internal/valid"
	"github.com/livebud/bud/package/gomod"
	"github.com/livebud/bud/package/parser"
)

// Load state
func Load(fsys fs.FS, module *gomod.Module, parser *parser.Parser) (*State, error) {
	loader := &loader{
		fsys:    fsys,
		imports: imports.New(),
		parser:  parser,
		module:  module,
	}
	return loader.Load()
}

type loader struct {
	bail.Struct
	fsys    fs.FS
	imports *imports.Set
	module  *gomod.Module
	parser  *parser.Parser
}

// Load the command state
func (l *loader) Load() (state *State, err error) {
	defer l.Recover(&err)
	state = new(State)
	// Add initial imports
	l.imports.AddStd("context")
	l.imports.AddNamed("commander", "github.com/livebud/bud/package/commander")
	// Load the commands
	state.Command = l.loadRoot("command")
	if !state.Command.Runnable && len(state.Command.Subs) == 0 {
		return nil, fs.ErrNotExist
	}
	// Load the imports
	state.Imports = l.imports.List()
	return state, nil
}

// Load the root command, which is unfortunately a special case
func (l *loader) loadRoot(base string) *Command {
	command := new(Command)
	command.Slug = imports.AssumedName(l.module.Import())
	// If a generated web server is present, then the root command is runnable.
	if _, err := fs.Stat(l.fsys, "bud/.app/web/web.go"); nil == err {
		l.imports.AddStd("os")
		l.imports.AddNamed("web", l.module.Import("bud", ".app", "web"))
		command.Runnable = true
	}
	des, err := fs.ReadDir(l.fsys, base)
	if err != nil {
		// Return the build/run command without any subcommands
		if errors.Is(err, fs.ErrNotExist) {
			return command
		}
		l.Bail(err)
	}
	for _, de := range des {
		if !de.IsDir() || !valid.Dir(de.Name()) {
			continue
		}
		sub := l.loadSub(base, de.Name())
		if sub == nil {
			continue
		}
		command.Subs = append(command.Subs, sub)
	}
	return command
}

// Load the subcommand
func (l *loader) loadSub(base, dir string) *Command {
	commandDir := filepath.Join(base, dir)
	command := new(Command)
	// Get the command name
	command.Name = filepath.Base(dir)
	// Load the import
	importPath := l.module.Import(filepath.SplitList(commandDir)...)
	importName := l.imports.Add(importPath)
	command.Import = &imports.Import{
		Name: importName,
		Path: importPath,
	}
	// Parse the command directory
	pkg, err := l.parser.Parse(commandDir)
	if err != nil {
		l.Bail(err)
	}
	// Find the Command struct, ignore if it doesn't exist
	stct := pkg.Struct("Command")
	if stct == nil {
		return nil
	}
	// Check that the command is runnable
	command.Runnable = isRunnable(stct)
	// Gather the fields
	for _, field := range stct.PublicFields() {
		tags, err := field.Tags()
		if err != nil {
			l.Bail(err)
		}
		// Is a dependency
		if len(tags) == 0 && !parser.IsBuiltin(field.Type()) {
			dep := l.loadCommandDep(field)
			command.Deps = append(command.Deps, dep)
			continue
		}
		// Is a flag
		if tags.Has("flag") || tags.Has("short") {
			flag := l.loadCommandFlag(field, tags)
			command.Flags = append(command.Flags, flag)
			continue
		}
		// Is an arg
		if tags.Has("arg") {
			arg := l.loadCommandArg(field, tags)
			command.Args = append(command.Args, arg)
			continue
		}
		l.Bail(fmt.Errorf("command: %q has an unacceptable type %q", command.Name, field.Type()))
	}

	// Read the subdirectories
	des, err := fs.ReadDir(l.fsys, commandDir)
	if err != nil {
		l.Bail(err)
	}
	// Load the subcommands
	for _, de := range des {
		if !de.IsDir() || !valid.Dir(de.Name()) {
			continue
		}
		sub := l.loadSub(base, filepath.Join(dir, de.Name()))
		if sub == nil {
			continue
		}
		sub.Parents = append(sub.Parents, command.Name)
		command.Subs = append(command.Subs, sub)
	}
	return command
}

func (l *loader) loadCommandDep(field *parser.Field) *Dep {
	def, err := parser.Definition(field.Type())
	if err != nil {
		l.Bail(fmt.Errorf("command: %w", err))
	}
	importPath, err := def.Package().Import()
	if err != nil {
		l.Bail(fmt.Errorf("command: %w", err))
	}
	importName := l.imports.Add(importPath)
	// Change original import name to new import name if needed
	fieldType := parser.Requalify(field.Type(), importName)
	return &Dep{
		Import: &imports.Import{
			Path: importPath,
			Name: importName,
		},
		Type: fieldType.String(),
		Name: field.Name(),
	}
}

// Load the command flag
func (l *loader) loadCommandFlag(field *parser.Field, tags parser.Tags) *Flag {
	flag := new(Flag)
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
			l.Bail(fmt.Errorf("command: short flag must be exactly 1 byte long, not %q", short))
		}
		flag.Short = short[0]
	}
	// Set the flag type
	// TODO: add more types
	switch t := field.Type().String(); t {
	case "string", "bool", "int":
		flag.Type = t
	default:
		l.Bail(fmt.Errorf("command: flag must be a string, bool or int, not %q", t))
	}
	// Set the usage
	flag.Help = tags.Get("help")
	// Set the default
	flag.Default = tags.Get("default")
	return flag
}

// Load the command arg
func (l *loader) loadCommandArg(field *parser.Field, tags parser.Tags) *Arg {
	arg := new(Arg)
	// Set the arg name
	arg.Name = tags.Get("arg")
	// Set the arg type
	switch t := field.Type().String(); t {
	case "string", "bool", "int":
		arg.Type = t
	default:
		err := fmt.Errorf("command: arg must be a string, bool or int not %q", t)
		l.Bail(err)
	}
	// Set the usage
	arg.Help = tags.Get("help")
	return arg
}

// Returns true if the command is runnable
func isRunnable(stct *parser.Struct) bool {
	run := stct.Method("Run")
	if run == nil {
		return false
	}
	params := run.Params()
	if len(params) != 1 {
		return false
	}
	// Ensure the first param is context.Context
	typeName := parser.TypeName(params[0].Type())
	return typeName == "Context"
}
