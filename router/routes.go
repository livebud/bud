package router

//go:generate peg -switch -inline routes.peg

import (
	"fmt"
	"regexp"
	"strings"
)

// Parse fn
func Parse(input string) (*regexp.Regexp, error) {
	parser := &parser{Buffer: input}
	parser.Init()
	err := parser.Parse()
	if err != nil {
		return nil, err
	}
	parser.Execute()
	// parser.PrintSyntaxTree()
	return parser.route.compile()
}

// Node interface
type Node interface {
	source() string
}

// route struct
type route struct {
	Segments []Segment
}

// Compile a route to regexp
func (r *route) compile() (*regexp.Regexp, error) {
	var sources []string
	for _, segment := range r.Segments {
		sources = append(sources, segment.source())
	}
	source := strings.Join(sources, "")
	return regexp.Compile("^" + source + "$")
}

// Segment interface
type Segment interface {
	Node
	segment()
}

func (*Slash) segment()       {}
func (*OptionalKey) segment() {}
func (*RegexpKey) segment()   {}
func (*BasicKey) segment()    {}
func (*Text) segment()        {}

// Slash struct
type Slash struct {
}

func (*Slash) source() string {
	return `\/`
}

// OptionalKey struct
type OptionalKey struct {
	PrefixSlash *Slash
	PrefixText  *Text
	Key         Key
}

func (o *OptionalKey) source() string {
	prefix := ""
	if o.PrefixSlash != nil {
		prefix += o.PrefixSlash.source()
	}
	if o.PrefixText != nil {
		prefix += o.PrefixText.source()
	}
	return fmt.Sprintf(`(?:%s%s)?`, prefix, o.Key.source())
}

// Key interface
type Key interface {
	Segment
	key()
}

func (*RegexpKey) key() {}
func (*BasicKey) key()  {}

// RegexpKey struct
type RegexpKey struct {
	Name   *Identifier
	Regexp *Regexp
}

func (n *RegexpKey) source() string {
	return fmt.Sprintf(`(?P<%s>%s)`, n.Name.source(), n.Regexp.source())
}

// BasicKey struct
type BasicKey struct {
	Name *Identifier
}

// Matches any ascii letter, number, -, _
func (n *BasicKey) source() string {
	return fmt.Sprintf(`(?P<%s>[A-Za-z0-9-_]+)`, n.Name.source())
}

// Identifier struct
type Identifier struct {
	Value string
}

func (n *Identifier) source() string {
	return n.Value
}

// Regexp struct
type Regexp struct {
	Value string
}

func (n *Regexp) source() string {
	return n.Value
}

// Text struct
type Text struct {
	Value string
}

var rescape = regexp.MustCompile("[.+*?=^!:${}()[\\]|/\\\\]")

func escapeString(str string) string {
	return rescape.ReplaceAllString(str, "\\$0")
}

func (n *Text) source() string {
	return escapeString(n.Value)
}
