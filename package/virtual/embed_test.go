package virtual_test

import (
	"testing"

	"github.com/livebud/bud/internal/is"
	"github.com/livebud/bud/package/virtual"
)

func TestEmbed(t *testing.T) {
	is := is.New(t)
	embed := &virtual.Embed{
		Path: "a.txt",
		Data: []byte("a"),
	}
	is.Equal(embed.Data.String(), "\\x61")
}
