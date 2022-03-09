package commandg_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"gitlab.com/mnm/bud/framework/commandg"
	"gitlab.com/mnm/bud/framework/commandp"
	"gitlab.com/mnm/bud/package/overlay"

	"github.com/matryer/is"
	"gitlab.com/mnm/bud/internal/testdir"
	"gitlab.com/mnm/bud/pkg/gomod"
	goparse "gitlab.com/mnm/bud/pkg/parser"
)

func parse(dir string, code []byte) (*goparse.Package, error) {
	now := time.Now()
	td := testdir.New()
	if err := td.Write(dir); err != nil {
		return nil, err
	}
	module, err := gomod.Find(dir)
	if err != nil {
		return nil, err
	}
	overlay, err := overlay.Load(module)
	if err != nil {
		return nil, err
	}
	fmt.Println(time.Since(now))
	goparser := goparse.New(overlay, module)
	return goparser.Parse()
}

func TestEmpty(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()
	dir := t.TempDir()
	code, err := commandg.Generate(ctx, &commandp.State{})
	is.NoErr(err)
	pkg, err := parse(dir, code)
	is.NoErr(err)
	is.Equal(pkg.Name(), "command")
}
