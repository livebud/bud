package expand_test

import (
	"context"
	"os"
	"testing"

	"gitlab.com/mnm/bud/pkg/buddy"

	"github.com/matryer/is"
	"gitlab.com/mnm/bud/internal/buddy/expand"
	"gitlab.com/mnm/bud/internal/testdir"
)

func TestEmpty(t *testing.T) {
	is := is.New(t)
	td := testdir.New()
	dir := "_tmp"
	is.NoErr(os.RemoveAll(dir))
	err := td.Write(dir)
	is.NoErr(err)
	bud, err := buddy.Load(dir)
	is.NoErr(err)
	ctx := context.Background()
	err = bud.Expand(ctx, &expand.Input{})
	is.NoErr(err)
}
