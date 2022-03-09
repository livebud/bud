package mainfile

import "gitlab.com/mnm/bud/pkg/gomod"

// ForApp creates a main file for the app
func ForApp(module *gomod.Module) *Main {
	return &Main{module.Import("bud/.app/program")}
}
