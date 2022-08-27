package lex

import (
	"fmt"
	"unicode"
	"unicode/utf8"
)

// New lexer
func New(route string) Lexer {
	l := &lexer{
		tokenCh: make(chan Token),
		input:   route,
	}
	go l.run()
	return l
}

// Lexer interface
type Lexer interface {
	Next() Token
}

// Lexer for turning routes into tokens
//
// Example of lexer state:
// /:firstName
//
//	      ^ pos is at "N", width 1 ("N" is 1 byte wide)
//	^ start is at ":"
type lexer struct {
	tokenCh chan Token
	input   string // buffer containing the full route
	start   int    // start position of the new token
	pos     int    // current position in the token stream
	width   int    // width of the current rune
}

// Token gets the next token
func (l *lexer) Next() Token {
	return <-l.tokenCh
}

// Run the lexer
func (l *lexer) run() {
	for state := lexSlash; state != nil; {
		state = state(l)
	}
	close(l.tokenCh)
}

const end = 0

// peek returns but does not consume the next rune in the input.
func (l *lexer) peek() rune {
	r := l.step()
	l.backup()
	return r
}

// Step to the next rune
func (l *lexer) step() (r rune) {
	if l.pos >= len(l.input) {
		l.width = 0
		return end
	}
	r, w := utf8.DecodeRuneInString(l.input[l.pos:])
	l.width = w
	l.pos += l.width
	return r
}

func (l *lexer) emit(t token) {
	value := l.input[l.start:l.pos]
	l.tokenCh <- Token{Type: t, Value: value}
	l.start = l.pos
}

// backup steps back one rune. Can only be called once per call of next.
func (l *lexer) backup() {
	l.pos -= l.width
}

func (l *lexer) errorf(format string, args ...interface{}) stateFn {
	value := fmt.Sprintf(format, args...)
	l.tokenCh <- Token{
		Type:  ErrorToken,
		Value: value,
	}
	return nil
}

// State function
type stateFn func(l *lexer) stateFn

func lexSlash(l *lexer) stateFn {
	r := l.step()
	if r != '/' {
		return l.errorf(`route %q: must start with a slash "/"`, l.input)
	}
	l.emit(SlashToken)
	return lexText
}

// Lex the route
func lexText(l *lexer) stateFn {
	switch r := l.step(); {
	case r == '/':
		l.emit(SlashToken)
		if r := l.peek(); r == end {
			return l.errorf(`route %q: remove the slash "/" at the end`, l.input)
		}
		return lexText
	case r == ':':
		// l.ignore()
		return lexSlotFirst
	case r == end:
		l.emit(EndToken)
		return nil
	case r == '?' || r == '*':
		return l.errorf("route %q: unexpected modifier %q", l.input, string(r))
	case unicode.IsUpper(r):
		return l.errorf("route %q: uppercase letters are not allowed %q", l.input, string(r))
	case isPath(r):
		return lexPath
	default:
		return l.errorf("route %q: invalid character %q", l.input, string(r))
	}
}

func isPath(r rune) bool {
	switch r {
	case ' ', ':', '/', '*', '?', end:
		return false
	}
	if unicode.IsUpper(r) {
		return false
	}
	return unicode.IsPrint(r)
}

func lexPath(l *lexer) stateFn {
	r := l.step()
	for isPath(r) {
		r = l.step()
	}
	l.backup()
	l.emit(PathToken)
	return lexText
}

// Slot first must be [a-z]
// TODO: consider allowing any unicode letter or digit
func isSlotFirst(r rune) bool {
	switch r {
	case 'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm',
		'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z':
		return true
	}
	return false
}

func lexSlotFirst(l *lexer) stateFn {
	r := l.step()
	// Next character after : can't be end
	if r == end {
		return l.errorf(`route %q: missing slot name after ":"`, l.input)
	}
	if !isSlotFirst(r) {
		return l.errorf(`route %q: first letter after ":" must be a lowercase Latin letter`, l.input)
	}
	return lexSlotRest
}

// Slot rest must be [a-z0-9_]+
// TODO: consider allowing any unicode letter or digit
func isSlotRest(r rune) bool {
	switch r {
	case '_', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9',
		'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm',
		'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z':
		return true
	}
	return false
}

// Loop over remaining slots
func lexSlotRest(l *lexer) stateFn {
	r := l.step()
	for isSlotRest(r) {
		r = l.step()
	}
	// Slot can't have uppercase letters
	if unicode.IsUpper(r) {
		return l.errorf(`route %q: uppercase letters are not allowed %q`, l.input, string(r))
	}
	// After the slot name
	switch r {
	case '?':
		// Support optional modifiers
		l.emit(QuestionToken)
		return lexQuestion
	case '*':
		// Support wildcard modifiers
		l.emit(StarToken)
		return lexStar
	case '.', '/', end:
		// Valid post-slot values
		// TODO: There are probably some other characters that should be allowed.
		l.backup()
		l.emit(SlotToken)
		return lexText
	default:
		// All other slot values should be invalid
		return l.errorf(`route %q: invalid slot character %q`, l.input, string(r))
	}
}

func lexQuestion(l *lexer) stateFn {
	// Expect End after
	switch r := l.step(); r {
	case end:
		l.emit(EndToken)
		return nil
	case '*':
		return l.errorf(`route %q: "*" not allowed after "?"`, l.input)
	default:
		return l.errorf(`route %q: optional "?" must be at the end`, l.input)
	}
}

func lexStar(l *lexer) stateFn {
	// Expect end after
	switch r := l.step(); r {
	case end:
		l.emit(EndToken)
		return nil
	case '?':
		return l.errorf(`route %q: "?" not allowed after "*"`, l.input)
	default:
		return l.errorf(`route %q: wildcard "*" must be at the end`, l.input)
	}
}
