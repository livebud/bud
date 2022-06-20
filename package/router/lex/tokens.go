package lex

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

// token type
type token string

// Types of tokens
const (
	PathToken     token = "path"
	SlashToken    token = "slash"
	SlotToken     token = "slot"
	QuestionToken token = "question"
	StarToken     token = "star"
	ErrorToken    token = "error"
	EndToken      token = "end"
)

// Token produced by the lexer
type Token struct {
	Type  token
	Value string
}

func (t Token) String() string {
	if t.Type == EndToken {
		return ""
	}
	return fmt.Sprintf("%s:%q", t.Type, t.Value)
}

// Tokens is a list of tokens
type Tokens []Token

// At returns the individual value at i. For paths, this will be a single
// character, for slots, this will be the whole slot name.
func (tokens Tokens) At(i int) string {
	for _, token := range tokens {
		switch token.Type {
		case PathToken, SlashToken:
			for _, char := range token.Value {
				if i == 0 {
					return string(char)
				}
				i--
			}
		case SlotToken, QuestionToken, StarToken:
			if i == 0 {
				return token.Value
			}
			i--
		}
	}
	return ""
}

// Size returns the number individual values. For paths, this will count the
// number of characters, whereas slots will count as one.
func (tokens Tokens) Size() (n int) {
	for _, token := range tokens {
		switch token.Type {
		case PathToken, SlashToken:
			n += utf8.RuneCountInString(token.Value)
		case SlotToken, QuestionToken, StarToken:
			n++
		}
	}
	return n
}

// Split the token list into two lists of tokens. If we land in the middle of a
// path token then split that token into two path tokens.
func (tokens Tokens) Split(at int) []Tokens {
	for i, token := range tokens {
		switch token.Type {
		case PathToken, SlashToken:
			for j := range token.Value {
				if at != 0 {
					at--
					continue
				}
				left, right := token.Value[:j], token.Value[j:]
				// At the edge
				if left == "" || right == "" {
					if i > 0 && i < len(tokens) {
						return []Tokens{tokens[:i], tokens[i:]}
					}
					return []Tokens{tokens}
				}
				newToken := Token{token.Type, left}
				leftTokens := append(append(Tokens{}, tokens[:i]...), newToken)
				rightTokens := append(Tokens{}, tokens[i:]...)
				rightTokens[0].Value = right
				return []Tokens{leftTokens, rightTokens}
			}
		case SlotToken, QuestionToken, StarToken:
			if at != 0 {
				at--
				continue
			}
			if i > 0 && i < len(tokens) {
				return []Tokens{tokens[:i], tokens[i:]}
			}
			return []Tokens{tokens}
		}
	}
	return []Tokens{tokens}
}

// String returns a string of tokens for testing
func (tokens Tokens) String() string {
	var ts []string
	for _, t := range tokens {
		s := t.String()
		if s == "" {
			continue
		}
		ts = append(ts, s)
	}
	return strings.Join(ts, " ")
}
