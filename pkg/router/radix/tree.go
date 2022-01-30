package radix

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"gitlab.com/mnm/bud/pkg/router/lex"
)

// New radix tree
func New() Tree {
	return &tree{}
}

// Tree interface
type Tree interface {
	Insert(route string, handler http.Handler) error
	Match(path string) (*Match, bool)
	String() string
}

// Slots is a list of key value pairs
type Slots []*Slot

// Slot is a key value pair
type Slot struct {
	Key   string
	Value string
}

// Match struct
type Match struct {
	Handler http.Handler
	Route   string
	Slots   Slots
}

// Match the path to a route
type matchFn func(path string) (index int, slots Slots)

type node struct {
	tokens   lex.Tokens   // tokens that make up part or all of the route
	match    matchFn      // compiled match function
	route    string       // original inserted route
	handler  http.Handler // original handler
	children []*node      // child nodes
	wilds    []*node      // wild children
}

func (n *node) isWild() bool {
	token := n.tokens[0]
	switch token.Type {
	case lex.SlotToken, lex.QuestionToken, lex.StarToken:
		return true
	default:
		return false
	}
}

// Priority of the node
func (n *node) priority() (priority int) {
	for _, token := range n.tokens {
		switch token.Type {
		case lex.SlashToken, lex.PathToken:
			priority++
		}
	}
	return priority
}

type tree struct {
	root *node
}

func (t *tree) Insert(route string, handler http.Handler) error {
	lexer := lex.New(route)
	var tokens lex.Tokens
	for {
		token := lexer.Next()
		switch token.Type {
		case lex.QuestionToken:
			// Each optional tokens insert two routes
			if err := t.insert(stripTokenTrail(tokens), route, handler); err != nil {
				return err
			}
			// Make the optional token required
			tokens = append(tokens, lex.Token{
				Value: strings.TrimRight(token.Value, "?"),
				Type:  lex.SlotToken,
			})
		case lex.StarToken:
			// Each optional tokens insert two routes
			if err := t.insert(stripTokenTrail(tokens), route, handler); err != nil {
				return err
			}
			tokens = append(tokens, token)
		case lex.ErrorToken:
			// Error parsing the route
			return errors.New(token.Value)
		case lex.EndToken:
			// Done parsing the route
			return t.insert(tokens, route, handler)
		default:
			tokens = append(tokens, token)
		}
	}
}

// strip token trail removes path tokens up to either a slot or a slash
// e.g. /:id. => /:id
//      /a/b => /a
func stripTokenTrail(tokens lex.Tokens) lex.Tokens {
	i := len(tokens) - 1
loop:
	for ; i >= 0; i-- {
		switch tokens[i].Type {
		case lex.SlotToken:
			i++ // Include the slot
			break loop
		case lex.SlashToken:
			break loop
		}
	}
	if i == 0 {
		return tokens[:1]
	}
	newTokens := make(lex.Tokens, i)
	copy(newTokens[:], tokens)
	return newTokens
}

func (t *tree) insert(tokens lex.Tokens, route string, handler http.Handler) error {
	if t.root == nil {
		t.root = &node{
			tokens:  tokens,
			match:   matcher(tokens),
			route:   route,
			handler: handler,
		}
		return nil
	}
	return t.insertAt(t.root, tokens, route, handler)
}

func (t *tree) insertAt(parent *node, tokens lex.Tokens, route string, handler http.Handler) error {
	// Compute the longest common prefix between new path and the node's path
	// before slots.
	lcp := longestCommonPrefix(tokens, parent.tokens)
	parts := tokens.Split(lcp)
	inTreeAlready := len(parts) == 1
	// If longest common prefix is not the same length as the node path, We need
	// to split the node path into parent and child to prepare for another child.
	if lcp < parent.tokens.Size() {
		parent = splitAt(parent, lcp)
		// This set of tokens are already in the tree
		// E.g. We've inserted "/a", "/b", then "/". "/" will already be in the tree
		// but not have a handler
		if inTreeAlready {
			parent.handler = handler
			parent.route = route
			return nil
		}
		err := insertChild(parent, &node{
			tokens:  parts[1],
			match:   matcher(parts[1]),
			route:   route,
			handler: handler,
		})
		if err != nil {
			return err
		}
		// Unset the parent data
		parent.handler = nil
		parent.route = ""
		return nil
	}
	// This set of tokens are already in the tree
	// E.g. We've inserted "/a", "/b", then "/". "/" will already be in the tree
	// but not have a handler
	if inTreeAlready {
		// Error out if we have a handler that's exactly the same as another route
		if parent.handler != nil {
			return fmt.Errorf("radix: %q is already in the tree", route)
		}
		parent.handler = handler
		parent.route = route
		return nil
	}
	// For the remaining, non-common part of the path, check if any of the
	// children also start with that non-common part. If so, traverse that child.
	for _, child := range parent.children {
		if child.tokens.At(0) == parts[1].At(0) {
			return t.insertAt(child, parts[1], route, handler)
		}
	}
	// Recurse wild children if the wild child matches exactly
	for _, wild := range parent.wilds {
		if wild != nil && wild.tokens.At(0) == parts[1].At(0) {
			return t.insertAt(wild, parts[1], route, handler)
		}
	}
	// Otherwise, insert a new child on the parent with the remaining non-common
	// part of the path.
	return insertChild(parent, &node{
		tokens:  parts[1],
		match:   matcher(parts[1]),
		route:   route,
		handler: handler,
	})
}

// Split the single node into a parent and child node
func splitAt(parent *node, at int) *node {
	parts := parent.tokens.Split(at)
	if len(parts) == 1 {
		return parent
	}
	child := &node{
		tokens:   parts[1],
		match:    matcher(parts[1]),
		route:    parent.route,
		handler:  parent.handler,
		children: parent.children,
		wilds:    parent.wilds,
	}
	// Add the split child, moving all existing children into the child
	parent.children = []*node{}
	parent.wilds = []*node{}
	insertChild(parent, child)
	// Split the tokens up and recompile the match function
	parent.tokens = parts[0]
	parent.match = matcher(parts[0])
	return parent
}

// Insert a child
func insertChild(parent *node, child *node) error {
	if child.isWild() {
		return insertWild(parent, child)
	}
	parent.children = append(parent.children, child)
	return nil
}

// Insert a wild child
func insertWild(parent *node, child *node) error {
	lwilds := len(parent.wilds)
	childp := child.priority()
	for i := 0; i < lwilds; i++ {
		wild := parent.wilds[i]
		wildp := wild.priority()
		// Prioritize more specific slots over less specific slots.
		if childp > wildp {
			parent.wilds = append(parent.wilds[:i], append([]*node{child}, parent.wilds[i:]...)...)
			return nil
		}
		// Don't allow /:id and /:hi on the same level.
		if childp == wildp {
			return fmt.Errorf("radix: ambiguous routes %q and %q", child.route, wild.route)
		}
	}
	parent.wilds = append(parent.wilds, child)
	return nil
}

// Match the path to a route in the tree
func (t *tree) Match(path string) (*Match, bool) {
	match := t.match(t.root, path, Slots{})
	if match == nil {
		return nil, false
	}
	return match, true
}

// Turn the tokens into a matcher
func matcher(tokens lex.Tokens) matchFn {
	var matchers []matchFn
	for _, token := range tokens {
		switch token.Type {
		case lex.PathToken, lex.SlashToken:
			matchers = append(matchers, matchExact(token))
		case lex.SlotToken:
			matchers = append(matchers, matchSlot(token))
		case lex.StarToken:
			matchers = append(matchers, matchStar(token))
		}
	}
	return compose(matchers)
}

// Compose the match functions into one function
func compose(matchers []matchFn) matchFn {
	return func(path string) (index int, slots Slots) {
		for _, match := range matchers {
			i, matchSlots := match(path)
			if i == -1 {
				return -1, matchSlots
			}
			path = path[i:]
			index += i
			slots = append(slots, matchSlots...)
		}
		return index, slots
	}
}

// Match a slot exactly (/users)
func matchExact(token lex.Token) matchFn {
	route := token.Value
	rlen := len(route)
	return func(path string) (index int, slots Slots) {
		if len(path) < rlen {
			return -1, nil
		}
		for ; index < rlen; index++ {
			if path[index] != route[index] {
				return -1, nil
			}
		}
		return index, nil
	}
}

// Match a slot (/:id)
func matchSlot(token lex.Token) matchFn {
	slotKey := token.Value[1:]
	return func(path string) (index int, slots Slots) {
		lpath := len(path)
		for i := 0; i < lpath; i++ {
			if path[i] == '.' || path[i] == '/' {
				break
			}
			index++
		}
		if index == 0 {
			return -1, nil
		}
		return index, Slots{{
			Key:   slotKey,
			Value: path[:index],
		}}
	}
}

// Match a star (e.g. /:path*)
func matchStar(token lex.Token) matchFn {
	lvalue := len(token.Value)
	slotKey := token.Value[1 : lvalue-1]
	return func(path string) (index int, slots Slots) {
		return len(path), Slots{{
			Key:   slotKey,
			Value: path,
		}}
	}
}

// Match the node
func (t *tree) match(node *node, path string, slots Slots) *Match {
	index, matchSlots := node.match(path)
	if index < 0 {
		return nil
	}
	path = path[index:]
	if matchSlots != nil {
		slots = append(slots, matchSlots...)
	}
	// No more path, we're done!
	if path == "" {
		// At a junction node, but this node isn't a route, so it's not a match
		if node.handler == nil {
			return nil
		}
		return &Match{
			Handler: node.handler,
			Route:   node.route,
			Slots:   slots,
		}
	}
	// First try matching the children
	for _, child := range node.children {
		if match := t.match(child, path, slots); match != nil {
			return match
		}
	}
	// Next try matching the wild children
	for _, wild := range node.wilds {
		if match := t.match(wild, path, slots); match != nil {
			return match
		}
	}
	return nil
}

func (t *tree) String() string {
	return t.string(t.root, "")
}

func (t *tree) string(n *node, indent string) string {
	if n == nil {
		return ""
	}
	route := ""
	for _, token := range n.tokens {
		route += token.Value
	}
	kind := "c"
	if n.isWild() {
		kind = "w"
	}
	out := fmt.Sprintf("%s%s[%d%s] %v\r\n", indent, route, len(n.children)+len(n.wilds), kind, n.handler)
	for l := len(route); l > 0; l-- {
		indent += " "
	}
	for _, child := range n.children {
		out += t.string(child, indent)
	}
	for _, wild := range n.wilds {
		out += t.string(wild, indent)
	}
	return out
}

func min(a, b int) int {
	if a <= b {
		return a
	}
	return b
}

func longestCommonPrefix(a, b lex.Tokens) int {
	i := 0
	max := min(a.Size(), b.Size())
	for i < max && a.At(i) == b.At(i) {
		i++
	}
	return i
}
