package oldcli

import (
	"context"
	"flag"
	"fmt"
	"io"
	"reflect"
	"strconv"
	"strings"

	"golang.org/x/sync/errgroup"
)

type Runner interface {
	Run(ctx context.Context) error
}

type Command interface {
	Command(name, help string, runners ...Runner) Command
	Flag(name, help string) *Flag
	Arg(name, help string) *Arg
	Args(name string) *Args
	Run(fn func(ctx context.Context) error)
	Advanced() Command
	Hidden() Command
}

type subcommand struct {
	parent   *subcommand
	name     string
	full     string
	help     string
	commands map[string]*subcommand
	flags    []*Flag
	args     []*Arg
	run      func(ctx context.Context) error
	restArgs *Args // optional, collects the rest of the args
	advanced bool
	hidden   bool
}

var _ Command = (*subcommand)(nil)

func (c *subcommand) Command(name, help string, runners ...Runner) Command {
	full := name
	if c.full != "" {
		full = c.full + ":" + name
	}
	command := &subcommand{
		c,
		name,
		full,
		help,
		map[string]*subcommand{},
		nil,
		nil,
		nil,
		nil,
		false,
		false,
	}
	c.commands[name] = command
	c.reflectRunners(command, runners...)
	return command
}

func (c *subcommand) Flag(name, help string) *Flag {
	flag := &Flag{name, help, 0, nil}
	c.flags = append(c.flags, flag)
	return flag
}

func (c *subcommand) Arg(name, help string) *Arg {
	arg := &Arg{name, help, nil}
	c.args = append(c.args, arg)
	return arg
}

func (c *subcommand) Args(name string) *Args {
	args := &Args{name, nil}
	c.restArgs = args
	return args
}

func (c *subcommand) Run(fn func(ctx context.Context) error) {
	c.run = fn
}

func (c *subcommand) Advanced() Command {
	c.advanced = true
	return c
}

func (c *subcommand) Hidden() Command {
	c.hidden = true
	return c
}

func (c *subcommand) extract(fset *flag.FlagSet, arguments []string) (args []string, err error) {
	for len(arguments) > 0 {
		if err := fset.Parse(arguments); err != nil {
			return nil, err
		}
		if fset.NArg() == 0 {
			return args, nil
		}
		args = append(args, fset.Arg(0))
		arguments = fset.Args()[1:]
	}
	return args, nil
}

func (c *subcommand) parse(ctx context.Context, arguments []string) error {
	fset := flag.NewFlagSet(c.full, flag.ContinueOnError)
	fset.SetOutput(io.Discard)
	for _, flag := range c.flags {
		fset.Var(flag.value, flag.name, flag.help)
		if flag.short != 0 {
			fset.Var(flag.value, string(flag.short), flag.help)
		}
	}
	args, err := c.extract(fset, arguments)
	if err != nil {
		return err
	}
	numArgs := len(c.args)
	for i, arg := range args {
		// Handle variadic arguments
		if i >= numArgs {
			if c.restArgs == nil {
				return fmt.Errorf("%w %s", ErrCommandNotFound, arg)
			}
			for _, arg := range args[i:] {
				if err := c.restArgs.value.Set(arg); err != nil {
					return err
				}
			}
			break
		}
		// Otherwise, set the argument as normal
		if err := c.args[i].value.Set(arg); err != nil {
			return err
		}
	}
	// Verify that all the args have been set or have default values
	if err := verifyArgs(c.args); err != nil {
		return err
	}
	// Also verify rest args if we have any
	if c.restArgs != nil {
		if err := c.restArgs.verify(c.restArgs.Name); err != nil {
			return err
		}
	}
	// Print usage if there's no run function defined
	if c.run == nil {
		return flag.ErrHelp
	}
	// Verify that all the flags have been set or have default values
	if err := verifyFlags(c.flags); err != nil {
		return err
	}
	// Run the command
	return c.run(ctx)
}

type runGroup []func(ctx context.Context) error

func (fns runGroup) Run(ctx context.Context) error {
	if len(fns) == 0 {
		return nil
	} else if len(fns) == 1 {
		return fns[0](ctx)
	}
	eg, ctx := errgroup.WithContext(ctx)
	for _, fn := range fns {
		fn := fn
		eg.Go(func() error {
			return fn(ctx)
		})
	}
	return eg.Wait()
}

func (s *subcommand) reflectRunners(sub Command, runners ...Runner) {
	if len(runners) == 0 {
		return
	}
	runGroup := make(runGroup, len(runners))
	sub.Run(runGroup.Run)
	for i, runner := range runners {
		if err := s.reflectRunner(sub, runner); err != nil {
			// Since this is a developer error that occurs on boot, it's ok to panic
			panic(err)
		}
		runGroup[i] = runner.Run
	}
}

func (s *subcommand) reflectRunner(cli Command, runner Runner) error {
	// Ensure we're working with a pointer to a struct
	ptrStructValue := reflect.ValueOf(runner)
	if ptrStructValue.Kind() != reflect.Ptr {
		return fmt.Errorf("cli: command must be a pointer to a struct")
	}
	structValue := ptrStructValue.Elem()
	if structValue.Kind() != reflect.Struct {
		return fmt.Errorf("cli: command must be a pointer to a struct")
	}

	// Add flags and arguments
	structType := structValue.Type()
	for i := 0; i < structType.NumField(); i++ {
		field := structType.Field(i)
		if field.Tag.Get("flag") != "" {
			if err := s.reflectFlag(cli, field, structValue.Field(i)); err != nil {
				return err
			}
		} else if field.Tag.Get("arg") != "" {
			if err := s.reflectArg(cli, field, structValue.Field(i)); err != nil {
				return err
			}
		}
	}

	return nil
}

// TODO: this is missing support for optional slices and maps
func (s *subcommand) reflectFlag(cli Command, field reflect.StructField, fieldValue reflect.Value) error {
	flagName := field.Tag.Get("flag")
	flag := cli.Flag(flagName, field.Tag.Get("desc"))
	switch field.Type.Kind() {
	// Handle integers
	case reflect.Int:
		intFlag := flag.Int(fieldValue.Addr().Interface().(*int))
		defaultValue := field.Tag.Get("default")
		if defaultValue != "" {
			n, err := strconv.Atoi(defaultValue)
			if err != nil {
				return fmt.Errorf("cli: invalid default value for %q: %w", flagName, err)
			}
			intFlag.Default(n)
		}
		return nil
	// Handle strings
	case reflect.String:
		stringFlag := flag.String(fieldValue.Addr().Interface().(*string))
		defaultValue := field.Tag.Get("default")
		if defaultValue != "" {
			stringFlag.Default(defaultValue)
		}
		return nil
	// Handle booleans
	case reflect.Bool:
		boolFlag := flag.Bool(fieldValue.Addr().Interface().(*bool))
		defaultValue := field.Tag.Get("default")
		if defaultValue != "" {
			b, err := strconv.ParseBool(defaultValue)
			if err != nil {
				return fmt.Errorf("cli: invalid default value %s for %q: %w", defaultValue, flagName, err)
			}
			boolFlag.Default(b)
		}
		return nil
	case reflect.Slice:
		if field.Type.Elem().Kind() != reflect.String {
			return fmt.Errorf("cli: unsupported flag type %q", field.Type.Elem().Kind())
		}
		stringsFlag := flag.Strings(fieldValue.Addr().Interface().(*[]string))
		defaultValue := field.Tag.Get("default")
		if defaultValue != "" {
			values := strings.Split(defaultValue, ",")
			stringsFlag.Default(values...)
		}
		return nil
	case reflect.Map:
		if field.Type.Elem().Kind() != reflect.String {
			return fmt.Errorf("cli: unsupported flag type %q", field.Type.Elem().Kind())
		}
		stringMapFlag := flag.StringMap(fieldValue.Addr().Interface().(*map[string]string))
		defaultValue := field.Tag.Get("default")
		if defaultValue != "" {
			fields := strings.Split(defaultValue, ",")
			keyValues := make(map[string]string)
			for _, field := range fields {
				keyValue := strings.Split(field, ":")
				if len(keyValue) != 2 {
					return fmt.Errorf("cli: invalid default value %s for %q", defaultValue, flagName)
				}
				keyValues[keyValue[0]] = keyValue[1]
			}
			stringMapFlag.Default(keyValues)
		}
		return nil
	case reflect.Ptr:
		elem := field.Type.Elem()
		oflag := flag.Optional()
		switch elem.Kind() {
		case reflect.Int:
			intFlag := oflag.Int(fieldValue.Addr().Interface().(**int))
			defaultValue := field.Tag.Get("default")
			if defaultValue != "" {
				n, err := strconv.Atoi(defaultValue)
				if err != nil {
					return fmt.Errorf("cli: invalid default value for %q: %w", flagName, err)
				}
				intFlag.Default(n)
			}
			return nil
		case reflect.String:
			stringFlag := oflag.String(fieldValue.Addr().Interface().(**string))
			defaultValue := field.Tag.Get("default")
			if defaultValue != "" {
				stringFlag.Default(defaultValue)
			}
			return nil
		case reflect.Bool:
			boolFlag := oflag.Bool(fieldValue.Addr().Interface().(**bool))
			defaultValue := field.Tag.Get("default")
			if defaultValue != "" {
				b, err := strconv.ParseBool(defaultValue)
				if err != nil {
					return fmt.Errorf("cli: invalid default value %s for %q: %w", defaultValue, flagName, err)
				}
				boolFlag.Default(b)
			}
			return nil
		default:
			return fmt.Errorf("cli: unsupported flag type %q", field.Type)
		}
	default:
		return fmt.Errorf("cli: unsupported flag type %q", field.Type)
	}
}

func (s *subcommand) reflectArg(cli Command, structField reflect.StructField, fieldValue reflect.Value) error {
	argName := structField.Tag.Get("arg")
	arg := cli.Arg(argName, structField.Tag.Get("desc"))
	switch structField.Type.Kind() {
	// Handle integers
	case reflect.Int:
		intFlag := arg.Int(fieldValue.Addr().Interface().(*int))
		defaultValue := structField.Tag.Get("default")
		if defaultValue != "" {
			n, err := strconv.Atoi(defaultValue)
			if err != nil {
				return fmt.Errorf("cli: invalid default value for %q: %w", argName, err)
			}
			intFlag.Default(n)
		}
		return nil
	// Handle strings
	case reflect.String:
		stringArg := arg.String(fieldValue.Addr().Interface().(*string))
		defaultValue := structField.Tag.Get("default")
		if defaultValue != "" {
			stringArg.Default(defaultValue)
		}
		return nil
	case reflect.Map:
		if structField.Type.Elem().Kind() != reflect.String {
			return fmt.Errorf("cli: unsupported arg type %q", structField.Type.Elem().Kind())
		}
		stringMapArg := arg.StringMap(fieldValue.Addr().Interface().(*map[string]string))
		defaultValue := structField.Tag.Get("default")
		if defaultValue != "" {
			fields := strings.Split(defaultValue, ",")
			keyValues := make(map[string]string)
			for _, field := range fields {
				keyValue := strings.Split(field, ":")
				if len(keyValue) != 2 {
					return fmt.Errorf("cli: invalid default value %s for %q", defaultValue, argName)
				}
				keyValues[keyValue[0]] = keyValue[1]
			}
			stringMapArg.Default(keyValues)
		}
		return nil
	default:
		return fmt.Errorf("cli: unsupported arg type %q", structField.Type.Kind())
	}
}
