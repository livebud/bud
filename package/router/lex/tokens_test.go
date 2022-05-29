package lex_test

import (
	"testing"

	"github.com/livebud/bud/internal/is"
	"github.com/livebud/bud/package/router/lex"
)

func tokens(t testing.TB, input string) lex.Tokens {
	l := lex.New(input)
	var tokens lex.Tokens
	for {
		token := l.Next()
		switch token.Type {
		case lex.ErrorToken:
			t.Fatalf("tokens: unexpected error token %+s", token.Value)
		case lex.EndToken:
			return tokens
		default:
			tokens = append(tokens, token)
		}
	}
}

func TestAt(t *testing.T) {
	is := is.New(t)
	toks := tokens(t, "/")
	is.Equal(toks.At(0), "/")
	toks = tokens(t, "/hi")
	is.Equal(toks.At(0), "/")
	is.Equal(toks.At(1), "h")
	is.Equal(toks.At(2), "i")
	is.Equal(toks.At(3), "")
	toks = tokens(t, "/:a")
	is.Equal(toks.At(0), "/")
	is.Equal(toks.At(1), ":a")
	is.Equal(toks.At(2), "")
	toks = tokens(t, "/users/:id")
	is.Equal(toks.At(0), "/")
	is.Equal(toks.At(1), "u")
	is.Equal(toks.At(2), "s")
	is.Equal(toks.At(3), "e")
	is.Equal(toks.At(4), "r")
	is.Equal(toks.At(5), "s")
	is.Equal(toks.At(6), "/")
	is.Equal(toks.At(7), ":id")
	is.Equal(toks.At(8), "")
	toks = tokens(t, "/:noun/:id?")
	is.Equal(toks.At(0), "/")
	is.Equal(toks.At(1), ":noun")
	is.Equal(toks.At(2), "/")
	is.Equal(toks.At(3), ":id?")
	toks = tokens(t, "/:noun/:id*")
	is.Equal(toks.At(0), "/")
	is.Equal(toks.At(1), ":noun")
	is.Equal(toks.At(2), "/")
	is.Equal(toks.At(3), ":id*")
}

func TestSize(t *testing.T) {
	is := is.New(t)
	toks := tokens(t, "/")
	is.Equal(toks.Size(), 1)
	toks = tokens(t, "/hi")
	is.Equal(toks.Size(), 3)
	toks = tokens(t, "/:a")
	is.Equal(toks.Size(), 2)
	toks = tokens(t, "/users/:id")
	is.Equal(toks.Size(), 8)
	toks = tokens(t, "/:noun/:id?")
	is.Equal(toks.Size(), 4)
	toks = tokens(t, "/:noun/:id*")
	is.Equal(toks.Size(), 4)
}
func TestSplit(t *testing.T) {
	is := is.New(t)

	toks := tokens(t, "/")
	is.Equal(len(toks), 1)
	is.Equal(toks.Size(), 1)
	parts := toks.Split(0)
	is.Equal(len(parts), 1)
	is.Equal(parts[0].Size(), 1)
	is.Equal(parts[0].At(0), "/")
	is.Equal(parts[0].At(1), "")

	toks = tokens(t, "/a")
	is.Equal(len(toks), 2)
	is.Equal(toks.Size(), 2)
	parts = toks.Split(1)
	is.Equal(len(parts), 2)
	is.Equal(parts[0].Size(), 1)
	is.Equal(parts[0].At(0), "/")
	is.Equal(parts[0].At(1), "")
	is.Equal(parts[1].Size(), 1)
	is.Equal(parts[1].At(0), "a")
	is.Equal(parts[1].At(1), "")

	toks = tokens(t, "/hi")
	is.Equal(len(toks), 2)
	is.Equal(toks.Size(), 3)
	parts = toks.Split(2)
	is.Equal(len(parts), 2)
	is.Equal(parts[0].Size(), 2)
	is.Equal(parts[0].At(0), "/")
	is.Equal(parts[0].At(1), "h")
	is.Equal(parts[1].Size(), 1)
	is.Equal(parts[1].At(0), "i")

	toks = tokens(t, "/:a")
	is.Equal(len(toks), 2)
	is.Equal(toks.Size(), 2)
	parts = toks.Split(1)
	is.Equal(len(parts), 2)
	is.Equal(parts[0].Size(), 1)
	is.Equal(parts[0].At(0), "/")
	is.Equal(parts[1].Size(), 1)
	is.Equal(parts[1].At(0), ":a")

	toks = tokens(t, "/:id")
	is.Equal(len(toks), 2)
	is.Equal(toks.Size(), 2)
	parts = toks.Split(1)
	is.Equal(len(parts), 2)
	is.Equal(parts[0].Size(), 1)
	is.Equal(parts[0].At(0), "/")
	is.Equal(parts[1].Size(), 1)
	is.Equal(parts[1].At(0), ":id")

	toks = tokens(t, "/users/:id")
	is.Equal(len(toks), 4)
	is.Equal(toks.Size(), 8)
	parts = toks.Split(3)
	is.Equal(len(parts), 2)
	is.Equal(parts[0].Size(), 3)
	is.Equal(parts[0].At(0), "/")
	is.Equal(parts[0].At(1), "u")
	is.Equal(parts[0].At(2), "s")
	is.Equal(parts[0].At(3), "")
	is.Equal(parts[1].Size(), 5)
	is.Equal(parts[1].At(0), "e")
	is.Equal(parts[1].At(1), "r")
	is.Equal(parts[1].At(2), "s")
	is.Equal(parts[1].At(3), "/")
	is.Equal(parts[1].At(4), ":id")
	is.Equal(parts[1].At(5), "")
	// Split again
	parts = parts[1].Split(3)
	is.Equal(len(parts), 2)
	is.Equal(parts[0].Size(), 3)
	is.Equal(parts[0].At(0), "e")
	is.Equal(parts[0].At(1), "r")
	is.Equal(parts[0].At(2), "s")
	is.Equal(parts[0].At(3), "")
	is.Equal(parts[1].Size(), 2)
	is.Equal(parts[1].At(0), "/")
	is.Equal(parts[1].At(1), ":id")
	is.Equal(parts[1].At(2), "")
}
