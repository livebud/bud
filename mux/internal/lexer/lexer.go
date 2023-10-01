package lexer

import (
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/livebud/bud/mux/internal/token"
)

type state = func(l *Lexer) token.Type

func New(input string) *Lexer {
	l := &Lexer{
		input:  input,
		states: []state{initialState},
	}
	l.step()
	return l
}

func Lex(input string) []token.Token {
	l := New(input)
	var tokens []token.Token
	for l.Next() {
		tokens = append(tokens, l.Token)
	}
	return tokens
}

// Print the input as tokens
func Print(input string) string {
	tokens := Lex(input)
	stoken := make([]string, len(tokens))
	for i, token := range tokens {
		stoken[i] = token.String()
	}
	return strings.Join(stoken, " ")
}

type Lexer struct {
	Token token.Token // Current token
	input string      // Input string
	start int         // Index to the start of the current token
	end   int         // Index to the end of the current token
	cp    rune        // Code point being considered
	next  int         // Index to the next rune to be considered
	err   string      // Error message for an error token

	states []state // Stack of states
	peaked []token.Token
}

func (l *Lexer) nextToken() token.Token {
	l.start = l.end
	tokenType := l.states[len(l.states)-1](l)
	t := token.Token{
		Type:  tokenType,
		Start: l.start,
		Text:  l.text(),
	}
	if tokenType == token.Error {
		t.Text = l.err
		l.err = ""
	}
	return t
}

func (l *Lexer) Next() bool {
	if len(l.peaked) > 0 {
		l.Token = l.peaked[0]
		l.peaked = l.peaked[1:]
	} else {
		l.Token = l.nextToken()
	}
	return l.Token.Type != token.End
}

func (l *Lexer) Peak(nth int) token.Token {
	if len(l.peaked) >= nth {
		return l.peaked[nth-1]
	}
	for i := len(l.peaked); i < nth; i++ {
		l.peaked = append(l.peaked, l.nextToken())
	}
	return l.peaked[nth-1]
}

// Use -1 to indicate the end of the file
const eof = -1

// Step advances the lexer to the next token
func (l *Lexer) step() {
	codePoint, width := utf8.DecodeRuneInString(l.input[l.next:])
	if width == 0 {
		codePoint = eof
	}
	l.cp = codePoint
	l.end = l.next
	l.next += width
}

func (l *Lexer) text() string {
	return l.input[l.start:l.end]
}

func (l *Lexer) stepUntil(rs ...rune) bool {
	for {
		if l.cp == eof {
			return false
		}
		for _, r := range rs {
			if l.cp == r {
				return true
			}
		}
		l.step()
	}
}

func (l *Lexer) pushState(state state) {
	l.states = append(l.states, state)
}

func (l *Lexer) popState() {
	l.states = l.states[:len(l.states)-1]
}

func (l *Lexer) errorf(msg string, args ...interface{}) token.Type {
	l.err = fmt.Sprintf(msg, args...)
	return token.Error
}

func initialState(l *Lexer) token.Type {
	switch l.cp {
	case eof:
		return token.End
	case '/':
		l.step()
		l.pushState(pathState)
		return token.Slash
	default:
		l.stepUntil('/')
		return l.errorf(`path must start with a slash /`)
	}
}

func pathState(l *Lexer) token.Type {
	switch {
	case l.cp == eof:
		l.popState()
		return token.End
	case l.cp == '/':
		l.step()
		return token.Slash
	case l.cp == '{':
		l.step()
		l.pushState(slotState)
		return token.OpenCurly
	case isPathChar(l.cp):
		l.step()
		for isPathChar(l.cp) {
			l.step()
		}
		return token.Path
	}
	// Skip forward for the error
	for {
		l.step()
		if l.cp == eof || l.cp == '/' || l.cp == '{' || isPathChar(l.cp) {
			break
		}
	}
	return l.errorf("unexpected character '%s' in path", l.text())
}

func slotState(l *Lexer) token.Type {
	switch {
	case l.cp == eof:
		l.popState()
		return l.errorf("unclosed slot")
	case isSlotChar(l.cp):
		l.step()
		for isSlotChar(l.cp) {
			l.step()
		}
		l.popState()
		l.pushState(slotModifierState)
		return token.Slot
	default:
		l.step()
		return l.errorf("slot can't start with '%s'", l.text())
	}
}

func slotModifierState(l *Lexer) token.Type {
	switch {
	case l.cp == eof:
		l.popState()
		return l.errorf("unclosed slot")
	case l.cp == '}':
		l.step()
		l.popState()
		return token.CloseCurly
	case l.cp == '?':
		l.step()
		l.popState()
		l.pushState(slotCloseState)
		return token.Question
	case l.cp == '*':
		l.step()
		l.popState()
		l.pushState(slotCloseState)
		return token.Star
	case l.cp == '|':
		l.step()
		l.pushState(slotRegexpState)
		return token.Pipe
	default:
		l.step()
		l.popState()
		l.pushState(slotCloseState)
		return l.errorf("invalid character '%s' in slot", l.text())
	}
}

func slotCloseState(l *Lexer) token.Type {
	switch l.cp {
	case eof:
		l.popState()
		return l.errorf("unclosed slot")
	case '}':
		l.step()
		l.popState()
		return token.CloseCurly
	}
	if l.stepUntil('}') {
		return l.errorf(`expected '}' but got '%s'`, l.text())
	}
	l.popState()
	return l.errorf("unclosed slot")
}

func slotRegexpState(l *Lexer) token.Type {
	depth := 0
loop:
	for l.cp != eof {
		switch l.cp {
		case '{':
			l.step()
			depth++
		case '}':
			if depth > 0 {
				depth--
				l.step()
				continue loop
			}
			l.popState()
			return token.Regexp
		default:
			l.step()
		}
	}
	l.popState()
	return token.End
}

func isLowerLetter(r rune) bool {
	return unicode.IsLetter(r) && unicode.IsLower(r)
}

func isNumber(r rune) bool {
	return unicode.IsNumber(r)
}

func isDash(r rune) bool {
	return r == '-'
}

func isUnderscore(r rune) bool {
	return r == '_'
}

func isPeriod(r rune) bool {
	return r == '.'
}

func isLowerAlpha(r rune) bool {
	return 'a' <= r && r <= 'z'
}

func isDigit(r rune) bool {
	return '0' <= r && r <= '9'
}

func isPathChar(r rune) bool {
	return isLowerLetter(r) || isNumber(r) || isDash(r) || isUnderscore(r) || isPeriod(r)
}

func isSlotChar(r rune) bool {
	return isLowerAlpha(r) || isDigit(r) || isUnderscore(r)
}
