package mainfile

import "gitlab.com/mnm/bud/pkg/gomod"

// ForCLI creates a main file for the project CLI
func ForCLI(module *gomod.Module) *Main {
	return &Main{module.Import("bud/.cli/program")}
}
