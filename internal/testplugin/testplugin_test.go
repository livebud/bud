package testplugin_test

import (
	"testing"

	"github.com/livebud/bud/internal/current"
	"github.com/livebud/bud/package/gomod"
	"github.com/matryer/is"

	"github.com/livebud/bud/internal/testplugin"
)

func currentModule() (*gomod.Module, error) {
	dir, err := current.Directory()
	if err != nil {
		return nil, err
	}
	return gomod.Find(dir)
}

func TestPlugin(t *testing.T) {
	is := is.New(t)
	module, err := testplugin.Plugin()
	is.NoErr(err)
	is.Equal(module.Path, "github.com/livebud/bud-test-plugin")
	is.True(module.Version != "")
}

func TestNestedPlugin(t *testing.T) {
	is := is.New(t)
	module, err := testplugin.NestedPlugin()
	is.NoErr(err)
	is.Equal(module.Path, "github.com/livebud/bud-test-nested-plugin")
	is.True(module.Version != "")
}
