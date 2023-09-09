package cli

import (
	"fmt"
	"path/filepath"

	"github.com/livebud/bud/pkg/mod"
)

func Default(mod *mod.Module) *CLI {
	name := filepath.Base(mod.Directory())
	help := fmt.Sprintf("%s CLI", name)
	return New(name, help)
}
