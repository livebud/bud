package generator

import (
	"fmt"
	"strings"

	"github.com/livebud/bud/internal/imports"
)

type State struct {
	Imports        []*imports.Import
	FileGenerators []*CodeGenerator
	FileServers    []*CodeGenerator
	GenerateDirs   []*CodeGenerator
	ServeFiles     []*CodeGenerator
}

type Type string

type CodeGenerator struct {
	Import *imports.Import
	Path   string // Path that triggers the generator (e.g. "bud/cmd/app/main.go")
	Camel  string
}

func (c *CodeGenerator) Method() (string, error) {
	switch {
	case strings.HasPrefix(c.Path, "bud/internal"):
		return "Generate", nil
	case strings.HasPrefix(c.Path, "bud/cmd"):
		return "GenerateCmd", nil
	case strings.HasPrefix(c.Path, "bud/pkg"):
		return "GeneratePkg", nil
	default:
		return "", fmt.Errorf("generator: unexpected generator path %q", c.Path)
	}
}
