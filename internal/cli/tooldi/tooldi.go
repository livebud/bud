package tooldi

import (
	"context"
	"fmt"
	"io/fs"
	"os"
	"strings"

	"github.com/livebud/bud/internal/cli/bud"
	"github.com/livebud/bud/internal/imports"
	"github.com/livebud/bud/package/di"
	"github.com/livebud/bud/package/gomod"
	"github.com/livebud/bud/package/parser"
)

func New(bud *bud.Command, in *bud.Input) *Command {
	return &Command{bud: bud, in: in}
}

type Command struct {
	bud          *bud.Command
	in           *bud.Input
	Name         string
	Target       string
	Map          map[string]string
	Dependencies []string
	Externals    []string
	Hoist        bool
	Verbose      bool
}

func (c *Command) Run(ctx context.Context) error {
	log, err := bud.Log(c.in.Stderr, c.bud.Log)
	if err != nil {
		return err
	}
	module, err := bud.Module(c.bud.Dir)
	if err != nil {
		return err
	}
	// For the dependency injection CLI, use written files insted of generated
	// files. Note that this was changed due to budfs ignoring the bud/* dir.
	var fsys fs.FS = module
	parser := parser.New(fsys, module)
	fn := &di.Function{
		Hoist: c.Hoist,
	}
	target, err := c.toDependency(module, c.Target)
	if err != nil {
		return err
	}
	fn.Target = target.ImportPath()
	fn.Name = target.TypeName()
	fn.Aliases = di.Aliases{}
	// Add the type mapping
	for from, to := range c.Map {
		fromDep, err := c.toDependency(module, from)
		if err != nil {
			return err
		}
		toDep, err := c.toDependency(module, to)
		if err != nil {
			return err
		}
		fn.Aliases[fromDep] = toDep
	}
	// Add the dependencies
	for _, dependency := range c.Dependencies {
		dep, err := c.toDependency(module, dependency)
		if err != nil {
			return err
		}
		fn.Results = append(fn.Results, dep)
	}
	if len(c.Dependencies) > 0 {
		fn.Results = append(fn.Results, &di.Error{})
	}
	// Add the externals
	for _, external := range c.Externals {
		ext, err := c.toDependency(module, external)
		if err != nil {
			return err
		}
		fn.Params = append(fn.Params, &di.Param{
			Import: ext.ImportPath(),
			Type:   ext.TypeName(),
		})
	}
	injector := di.New(module, log, module, parser)
	node, err := injector.Load(fn)
	if err != nil {
		return err
	}
	if c.Verbose {
		fmt.Println(node.Print())
	}
	provider := node.Generate(imports.New(), c.Name, fn.Target)
	fmt.Fprintln(os.Stdout, provider.File())
	return nil
}

// This should handle both stdlib (e.g. "net/http"), directories (e.g. "web"),
// and dependencies
func (c *Command) toImportPath(module *gomod.Module, importPath string) (string, error) {
	importPath = strings.Trim(importPath, "\"")
	maybeDir := module.Directory(importPath)
	if _, err := os.Stat(maybeDir); err == nil {
		importPath, err = module.ResolveImport(maybeDir)
		if err != nil {
			return "", fmt.Errorf("di: unable to resolve import %s because %+s", importPath, err)
		}
	}
	return importPath, nil
}

func (c *Command) toDependency(module *gomod.Module, dependency string) (di.Dependency, error) {
	i := strings.LastIndex(dependency, ".")
	if i < 0 {
		return nil, fmt.Errorf("di: external must have form '<import>.<type>'. got %q ", dependency)
	}
	importPath, err := c.toImportPath(module, dependency[0:i])
	if err != nil {
		return nil, err
	}
	dataType := dependency[i+1:]
	// Create the dependency
	return &di.Type{
		Import: importPath,
		Type:   dataType,
	}, nil
}
