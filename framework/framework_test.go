package framework_test

import (
	"testing"

	"github.com/livebud/bud/framework"
	"github.com/livebud/bud/internal/is"
)

func TestString(t *testing.T) {
	is := is.New(t)
	f := framework.Flag{
		Embed:  true,
		Minify: true,
		Hot:    false,
	}
	flags := f.Flags()
	is.Equal(flags[0], "--embed=true")
	is.Equal(flags[1], "--minify=true")
	is.Equal(flags[2], "--hot=false")
}
