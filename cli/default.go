package cli

import (
	"fmt"
	"path/filepath"

	"github.com/livebud/bud/internal/mod"
)

func Default(mod *mod.Module) *CLI {
	name := filepath.Base(mod.Directory())
	desc := fmt.Sprintf("%s CLI", name)
	return New(name, desc)
}
