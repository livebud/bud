package glob

import (
	"path/filepath"
	"strings"

	"github.com/gobwas/glob/syntax/lexer"
)

const sep = string(filepath.Separator)

// Base gets the non-magical part of the glob
func Base(pattern string) string {
	parts := strings.Split(pattern, sep)
	var base []string
outer:
	for _, part := range parts {
		lex := lexer.NewLexer(part)
	inner:
		for {
			token := lex.Next()
			switch token.Type {
			case lexer.Text:
				continue
			case lexer.EOF:
				break inner
			default:
				break outer
			}
		}
		base = append(base, part)
	}
	if len(base) == 0 {
		return "."
	} else if len(base) == 1 && base[0] == "" {
		return sep
	}
	return filepath.Clean(strings.Join(base, sep))
}
