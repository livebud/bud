package esb

import (
	"strings"

	esbuild "github.com/evanw/esbuild/pkg/api"
)

type File = esbuild.OutputFile

type Error struct {
	Messages []esbuild.Message
}

func (e *Error) Error() string {
	errors := esbuild.FormatMessages(e.Messages, esbuild.FormatMessagesOptions{
		Color: true,
	})
	return strings.Join(errors, "\n\n")
}
