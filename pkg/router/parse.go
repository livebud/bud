package router

import "gitlab.com/mnm/bud/pkg/router/lex"

func Parse(route string) (tokens []lex.Token) {
	lexer := lex.New(route)
	for token := lexer.Next(); token.Type != lex.EndToken; {
		tokens = append(tokens, token)
	}
	return tokens
}
