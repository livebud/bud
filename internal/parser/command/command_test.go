package command_test

import (
	"fmt"
	"testing"

	"gitlab.com/mnm/bud/package/overlay"
	"gitlab.com/mnm/bud/pkg/parser"

	"gitlab.com/mnm/bud/internal/parser/command"
	"gitlab.com/mnm/bud/internal/testdir"
	"gitlab.com/mnm/bud/pkg/gomod"

	"github.com/matryer/is"
)

func TestEmpty(t *testing.T) {
	is := is.New(t)
	dir := t.TempDir()
	td := testdir.New()
	err := td.Write(dir)
	is.NoErr(err)
	module, err := gomod.Find(dir)
	is.NoErr(err)
	fsys, err := overlay.Load(module)
	is.NoErr(err)
	parser := parser.New(fsys, module)
	commandParser := command.New(module, parser)
	// TODO: fix this
	state, err := commandParser.ParseCLI()
	is.NoErr(err)
	fmt.Println(state)
}
