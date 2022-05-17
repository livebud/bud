package testplugin

import (
	"fmt"

	"github.com/livebud/bud/internal/current"
	"github.com/livebud/bud/package/gomod"
	"golang.org/x/mod/module"
)

var testPlugin struct {
	path, version string
}

func Plugin() (*module.Version, error) {
	dir, err := current.Directory()
	if err != nil {
		return nil, err
	}
	module, err := gomod.Find(dir)
	if err != nil {
		return nil, err
	}
	tp := module.File().Require("github.com/livebud/bud-test-plugin")
	if tp == nil {
		return nil, fmt.Errorf("testplugin: unable to find required test plugin in Bud's go.mod")
	}
	return tp, nil
}

func NestedPlugin() (*module.Version, error) {
	dir, err := current.Directory()
	if err != nil {
		return nil, err
	}
	module, err := gomod.Find(dir)
	if err != nil {
		return nil, err
	}
	tp := module.File().Require("github.com/livebud/bud-test-nested-plugin")
	if tp == nil {
		return nil, fmt.Errorf("testplugin: unable to find required nested plugin in Bud's go.mod")
	}
	return tp, nil
}
