package color_test

import (
	"fmt"
	"testing"

	"github.com/livebud/bud/internal/color"
	"github.com/matryer/is"
)

func TestColor(t *testing.T) {
	is := is.New(t)
	color := color.New()
	is.True(color.Enabled())
	is.Equal(color.Blue("hello"), "\x1b[34mhello\x1b[0m")
	is.Equal(color.Red("hello"), "\x1b[31mhello\x1b[0m")
}

func TestIgnore(t *testing.T) {
	is := is.New(t)
	color := color.Ignore()
	is.True(!color.Enabled())
	is.Equal(color.Blue("hello"), "hello")
	is.Equal(color.Red("hello"), "hello")
}

func ExampleIgnore() {
	color := color.Ignore()
	fmt.Println(color.Blue("hello"))
	fmt.Println(color.Red("hello"))
	// Output:
	// hello
	// hello
}
