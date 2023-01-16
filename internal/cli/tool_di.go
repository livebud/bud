package cli

import (
	"context"
	"fmt"
	"io/fs"
	"os"
	"strings"

	"github.com/livebud/bud/internal/imports"
	"github.com/livebud/bud/package/di"
	"github.com/livebud/bud/package/gomod"
	"github.com/livebud/bud/package/parser"
)

type ToolDi struct {
	Name         string
	Target       string
	Map          map[string]string
	Dependencies []string
	Externals    []string
	Hoist        bool
	Verbose      bool
}

func (c *CLI) ToolDi(ctx context.Context, in *ToolDi) error {
	log, err := c.loadLog()
	if err != nil {
		return err
	}
	module, err := c.findModule()
	if err != nil {
		return err
	}
	// For the dependency injection CLI, use written files insted of generated
	// files. Note that this was changed due to budfs ignoring the bud/* dir.
	var fsys fs.FS = module
	parser := parser.New(fsys, module)
	fn := &di.Function{
		Hoist: in.Hoist,
	}
	target, err := diToImportPath(module, in.Target)
	if err != nil {
		return err
	}
	fn.Target = target
	fn.Name = in.Name
	fn.Aliases = di.Aliases{}
	// Add the type mapping
	for from, to := range in.Map {
		fromDep, err := diToDependency(module, from)
		if err != nil {
			return err
		}
		toDep, err := diToDependency(module, to)
		if err != nil {
			return err
		}
		fn.Aliases[fromDep] = toDep
	}
	// Add the dependencies
	for _, dependency := range in.Dependencies {
		dep, err := diToDependency(module, dependency)
		if err != nil {
			return err
		}
		fn.Results = append(fn.Results, dep)
	}
	if len(in.Dependencies) > 0 {
		fn.Results = append(fn.Results, &di.Error{})
	}
	// Add the externals
	for _, external := range in.Externals {
		ext, err := diToDependency(module, external)
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
	if in.Verbose {
		fmt.Println(node.Print())
	}
	provider := node.Generate(imports.New(), in.Name, fn.Target)
	fmt.Fprintln(os.Stdout, provider.File())
	return nil
}

// This should handle both stdlib (e.g. "net/http"), directories (e.g. "web"),
// and dependencies
func diToImportPath(module *gomod.Module, importPath string) (string, error) {
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

func diToDependency(module *gomod.Module, dependency string) (di.Dependency, error) {
	i := strings.LastIndex(dependency, ".")
	if i < 0 {
		return nil, fmt.Errorf("di: external must have form '<import>.<type>'. got %q ", dependency)
	}
	importPath, err := diToImportPath(module, dependency[0:i])
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
