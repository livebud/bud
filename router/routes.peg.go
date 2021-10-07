package router

import (
	"fmt"
	"math"
	"sort"
	"strconv"
)

const endSymbol rune = 1114112

/* The rule types inferred from the grammar are below. */
type pegRule uint8

const (
	ruleUnknown pegRule = iota
	ruleRoute
	ruleFirstSegment
	ruleSegment
	ruleOptionalKey
	ruleKey
	ruleRegexpKey
	ruleBasicKey
	ruleRegexp
	ruleIdentifier
	ruleEscaped
	ruleText
	ruleSlash
	ruleEnd
	ruleAction0
	ruleAction1
	ruleAction2
	ruleAction3
	ruleAction4
	ruleAction5
	ruleAction6
	ruleAction7
	ruleAction8
	ruleAction9
	ruleAction10
	ruleAction11
	rulePegText
	ruleAction12
	ruleAction13
	ruleAction14
	ruleAction15
	ruleAction16

	rulePre
	ruleIn
	ruleSuf
)

var rul3s = [...]string{
	"Unknown",
	"Route",
	"FirstSegment",
	"Segment",
	"OptionalKey",
	"Key",
	"RegexpKey",
	"BasicKey",
	"Regexp",
	"Identifier",
	"Escaped",
	"Text",
	"Slash",
	"End",
	"Action0",
	"Action1",
	"Action2",
	"Action3",
	"Action4",
	"Action5",
	"Action6",
	"Action7",
	"Action8",
	"Action9",
	"Action10",
	"Action11",
	"PegText",
	"Action12",
	"Action13",
	"Action14",
	"Action15",
	"Action16",

	"Pre_",
	"_In_",
	"_Suf",
}

type node32 struct {
	token32
	up, next *node32
}

func (node *node32) print(depth int, buffer string) {
	for node != nil {
		for c := 0; c < depth; c++ {
			fmt.Printf(" ")
		}
		fmt.Printf("\x1B[34m%v\x1B[m %v\n", rul3s[node.pegRule], strconv.Quote(string(([]rune(buffer)[node.begin:node.end]))))
		if node.up != nil {
			node.up.print(depth+1, buffer)
		}
		node = node.next
	}
}

func (node *node32) Print(buffer string) {
	node.print(0, buffer)
}

type element struct {
	node *node32
	down *element
}

/* ${@} bit structure for abstract syntax tree */
type token32 struct {
	pegRule
	begin, end, next uint32
}

func (t *token32) isZero() bool {
	return t.pegRule == ruleUnknown && t.begin == 0 && t.end == 0 && t.next == 0
}

func (t *token32) isParentOf(u token32) bool {
	return t.begin <= u.begin && t.end >= u.end && t.next > u.next
}

func (t *token32) getToken32() token32 {
	return token32{pegRule: t.pegRule, begin: uint32(t.begin), end: uint32(t.end), next: uint32(t.next)}
}

func (t *token32) String() string {
	return fmt.Sprintf("\x1B[34m%v\x1B[m %v %v %v", rul3s[t.pegRule], t.begin, t.end, t.next)
}

type tokens32 struct {
	tree    []token32
	ordered [][]token32
}

func (t *tokens32) trim(length int) {
	t.tree = t.tree[0:length]
}

func (t *tokens32) Print() {
	for _, token := range t.tree {
		fmt.Println(token.String())
	}
}

func (t *tokens32) Order() [][]token32 {
	if t.ordered != nil {
		return t.ordered
	}

	depths := make([]int32, 1, math.MaxInt16)
	for i, token := range t.tree {
		if token.pegRule == ruleUnknown {
			t.tree = t.tree[:i]
			break
		}
		depth := int(token.next)
		if length := len(depths); depth >= length {
			depths = depths[:depth+1]
		}
		depths[depth]++
	}
	depths = append(depths, 0)

	ordered, pool := make([][]token32, len(depths)), make([]token32, len(t.tree)+len(depths))
	for i, depth := range depths {
		depth++
		ordered[i], pool, depths[i] = pool[:depth], pool[depth:], 0
	}

	for i, token := range t.tree {
		depth := token.next
		token.next = uint32(i)
		ordered[depth][depths[depth]] = token
		depths[depth]++
	}
	t.ordered = ordered
	return ordered
}

type state32 struct {
	token32
	depths []int32
	leaf   bool
}

func (t *tokens32) AST() *node32 {
	tokens := t.Tokens()
	stack := &element{node: &node32{token32: <-tokens}}
	for token := range tokens {
		if token.begin == token.end {
			continue
		}
		node := &node32{token32: token}
		for stack != nil && stack.node.begin >= token.begin && stack.node.end <= token.end {
			stack.node.next = node.up
			node.up = stack.node
			stack = stack.down
		}
		stack = &element{node: node, down: stack}
	}
	return stack.node
}

func (t *tokens32) PreOrder() (<-chan state32, [][]token32) {
	s, ordered := make(chan state32, 6), t.Order()
	go func() {
		var states [8]state32
		for i := range states {
			states[i].depths = make([]int32, len(ordered))
		}
		depths, state, depth := make([]int32, len(ordered)), 0, 1
		write := func(t token32, leaf bool) {
			S := states[state]
			state, S.pegRule, S.begin, S.end, S.next, S.leaf = (state+1)%8, t.pegRule, t.begin, t.end, uint32(depth), leaf
			copy(S.depths, depths)
			s <- S
		}

		states[state].token32 = ordered[0][0]
		depths[0]++
		state++
		a, b := ordered[depth-1][depths[depth-1]-1], ordered[depth][depths[depth]]
	depthFirstSearch:
		for {
			for {
				if i := depths[depth]; i > 0 {
					if c, j := ordered[depth][i-1], depths[depth-1]; a.isParentOf(c) &&
						(j < 2 || !ordered[depth-1][j-2].isParentOf(c)) {
						if c.end != b.begin {
							write(token32{pegRule: ruleIn, begin: c.end, end: b.begin}, true)
						}
						break
					}
				}

				if a.begin < b.begin {
					write(token32{pegRule: rulePre, begin: a.begin, end: b.begin}, true)
				}
				break
			}

			next := depth + 1
			if c := ordered[next][depths[next]]; c.pegRule != ruleUnknown && b.isParentOf(c) {
				write(b, false)
				depths[depth]++
				depth, a, b = next, b, c
				continue
			}

			write(b, true)
			depths[depth]++
			c, parent := ordered[depth][depths[depth]], true
			for {
				if c.pegRule != ruleUnknown && a.isParentOf(c) {
					b = c
					continue depthFirstSearch
				} else if parent && b.end != a.end {
					write(token32{pegRule: ruleSuf, begin: b.end, end: a.end}, true)
				}

				depth--
				if depth > 0 {
					a, b, c = ordered[depth-1][depths[depth-1]-1], a, ordered[depth][depths[depth]]
					parent = a.isParentOf(b)
					continue
				}

				break depthFirstSearch
			}
		}

		close(s)
	}()
	return s, ordered
}

func (t *tokens32) PrintSyntax() {
	tokens, ordered := t.PreOrder()
	max := -1
	for token := range tokens {
		if !token.leaf {
			fmt.Printf("%v", token.begin)
			for i, leaf, depths := 0, int(token.next), token.depths; i < leaf; i++ {
				fmt.Printf(" \x1B[36m%v\x1B[m", rul3s[ordered[i][depths[i]-1].pegRule])
			}
			fmt.Printf(" \x1B[36m%v\x1B[m\n", rul3s[token.pegRule])
		} else if token.begin == token.end {
			fmt.Printf("%v", token.begin)
			for i, leaf, depths := 0, int(token.next), token.depths; i < leaf; i++ {
				fmt.Printf(" \x1B[31m%v\x1B[m", rul3s[ordered[i][depths[i]-1].pegRule])
			}
			fmt.Printf(" \x1B[31m%v\x1B[m\n", rul3s[token.pegRule])
		} else {
			for c, end := token.begin, token.end; c < end; c++ {
				if i := int(c); max+1 < i {
					for j := max; j < i; j++ {
						fmt.Printf("skip %v %v\n", j, token.String())
					}
					max = i
				} else if i := int(c); i <= max {
					for j := i; j <= max; j++ {
						fmt.Printf("dupe %v %v\n", j, token.String())
					}
				} else {
					max = int(c)
				}
				fmt.Printf("%v", c)
				for i, leaf, depths := 0, int(token.next), token.depths; i < leaf; i++ {
					fmt.Printf(" \x1B[34m%v\x1B[m", rul3s[ordered[i][depths[i]-1].pegRule])
				}
				fmt.Printf(" \x1B[34m%v\x1B[m\n", rul3s[token.pegRule])
			}
			fmt.Printf("\n")
		}
	}
}

func (t *tokens32) PrintSyntaxTree(buffer string) {
	tokens, _ := t.PreOrder()
	for token := range tokens {
		for c := 0; c < int(token.next); c++ {
			fmt.Printf(" ")
		}
		fmt.Printf("\x1B[34m%v\x1B[m %v\n", rul3s[token.pegRule], strconv.Quote(string(([]rune(buffer)[token.begin:token.end]))))
	}
}

func (t *tokens32) Add(rule pegRule, begin, end, depth uint32, index int) {
	t.tree[index] = token32{pegRule: rule, begin: uint32(begin), end: uint32(end), next: uint32(depth)}
}

func (t *tokens32) Tokens() <-chan token32 {
	s := make(chan token32, 16)
	go func() {
		for _, v := range t.tree {
			s <- v.getToken32()
		}
		close(s)
	}()
	return s
}

func (t *tokens32) Error() []token32 {
	ordered := t.Order()
	length := len(ordered)
	tokens, length := make([]token32, length), length-1
	for i := range tokens {
		o := ordered[length-i]
		if len(o) > 1 {
			tokens[i] = o[len(o)-2].getToken32()
		}
	}
	return tokens
}

func (t *tokens32) Expand(index int) {
	tree := t.tree
	if index >= len(tree) {
		expanded := make([]token32, 2*len(tree))
		copy(expanded, tree)
		t.tree = expanded
	}
}

type parser struct {
	route       *route
	segments    []Segment
	key         Key
	regexpKey   *RegexpKey
	basicKey    *BasicKey
	optionalKey *OptionalKey
	slash       *Slash
	regexp      *Regexp
	identifier  *Identifier
	text        *Text

	Buffer string
	buffer []rune
	rules  [32]func() bool
	Parse  func(rule ...int) error
	Reset  func()
	Pretty bool
	tokens32
}

type textPosition struct {
	line, symbol int
}

type textPositionMap map[int]textPosition

func translatePositions(buffer []rune, positions []int) textPositionMap {
	length, translations, j, line, symbol := len(positions), make(textPositionMap, len(positions)), 0, 1, 0
	sort.Ints(positions)

search:
	for i, c := range buffer {
		if c == '\n' {
			line, symbol = line+1, 0
		} else {
			symbol++
		}
		if i == positions[j] {
			translations[positions[j]] = textPosition{line, symbol}
			for j++; j < length; j++ {
				if i != positions[j] {
					continue search
				}
			}
			break search
		}
	}

	return translations
}

type parseError struct {
	p   *parser
	max token32
}

func (e *parseError) Error() string {
	tokens, error := []token32{e.max}, "\n"
	positions, p := make([]int, 2*len(tokens)), 0
	for _, token := range tokens {
		positions[p], p = int(token.begin), p+1
		positions[p], p = int(token.end), p+1
	}
	translations := translatePositions(e.p.buffer, positions)
	format := "parse error near %v (line %v symbol %v - line %v symbol %v):\n%v\n"
	if e.p.Pretty {
		format = "parse error near \x1B[34m%v\x1B[m (line %v symbol %v - line %v symbol %v):\n%v\n"
	}
	for _, token := range tokens {
		begin, end := int(token.begin), int(token.end)
		error += fmt.Sprintf(format,
			rul3s[token.pegRule],
			translations[begin].line, translations[begin].symbol,
			translations[end].line, translations[end].symbol,
			strconv.Quote(string(e.p.buffer[begin:end])))
	}

	return error
}

func (p *parser) PrintSyntaxTree() {
	p.tokens32.PrintSyntaxTree(p.Buffer)
}

func (p *parser) Highlighter() {
	p.PrintSyntax()
}

func (p *parser) Execute() {
	buffer, _buffer, text, begin, end := p.Buffer, p.buffer, "", 0, 0
	for token := range p.Tokens() {
		switch token.pegRule {

		case rulePegText:
			begin, end = int(token.begin), int(token.end)
			text = string(_buffer[begin:end])

		case ruleAction0:

			p.route = &route{
				Segments: p.segments,
			}
			p.segments = nil

		case ruleAction1:

			p.segments = append(p.segments, p.slash)

		case ruleAction2:
			p.segments = append(p.segments, p.optionalKey)
			p.optionalKey = nil
		case ruleAction3:
			p.segments = append(p.segments, p.slash)
			p.slash = nil
		case ruleAction4:
			p.segments = append(p.segments, p.text)
			p.text = nil
		case ruleAction5:
			p.segments = append(p.segments, p.key)
			p.key = nil
		case ruleAction6:
			p.segments = append(p.segments, p.text)
			p.text = nil
		case ruleAction7:

			p.optionalKey = &OptionalKey{
				PrefixSlash: p.slash,
				PrefixText:  p.text,
				Key:         p.key,
			}
			p.slash = nil
			p.text = nil
			p.key = nil

		case ruleAction8:
			p.key = p.regexpKey
			p.regexpKey = nil
		case ruleAction9:
			p.key = p.basicKey
			p.basicKey = nil
		case ruleAction10:

			p.regexpKey = &RegexpKey{
				Name:   p.identifier,
				Regexp: p.regexp,
			}
			p.identifier = nil
			p.regexp = nil

		case ruleAction11:

			p.basicKey = &BasicKey{Name: p.identifier}
			p.identifier = nil

		case ruleAction12:

			p.regexp = &Regexp{Value: text}

		case ruleAction13:

			p.identifier = &Identifier{Value: text}

		case ruleAction14:

			p.text = &Text{Value: text}

		case ruleAction15:

			p.text = &Text{Value: text}

		case ruleAction16:

			p.slash = &Slash{}

		}
	}
	_, _, _, _, _ = buffer, _buffer, text, begin, end
}

func (p *parser) Init() {
	p.buffer = []rune(p.Buffer)
	if len(p.buffer) == 0 || p.buffer[len(p.buffer)-1] != endSymbol {
		p.buffer = append(p.buffer, endSymbol)
	}

	tree := tokens32{tree: make([]token32, math.MaxInt16)}
	var max token32
	position, depth, tokenIndex, buffer, _rules := uint32(0), uint32(0), 0, p.buffer, p.rules

	p.Parse = func(rule ...int) error {
		r := 1
		if len(rule) > 0 {
			r = rule[0]
		}
		matches := p.rules[r]()
		p.tokens32 = tree
		if matches {
			p.trim(tokenIndex)
			return nil
		}
		return &parseError{p, max}
	}

	p.Reset = func() {
		position, tokenIndex, depth = 0, 0, 0
	}

	add := func(rule pegRule, begin uint32) {
		tree.Expand(tokenIndex)
		tree.Add(rule, begin, position, depth, tokenIndex)
		tokenIndex++
		if begin != position && position > max.end {
			max = token32{rule, begin, position, depth}
		}
	}

	matchDot := func() bool {
		if buffer[position] != endSymbol {
			position++
			return true
		}
		return false
	}

	/*matchChar := func(c byte) bool {
		if buffer[position] == c {
			position++
			return true
		}
		return false
	}*/

	/*matchRange := func(lower byte, upper byte) bool {
		if c := buffer[position]; c >= lower && c <= upper {
			position++
			return true
		}
		return false
	}*/

	_rules = [...]func() bool{
		nil,
		/* 0 Route <- <(FirstSegment Segment* End Action0)> */
		func() bool {
			position0, tokenIndex0, depth0 := position, tokenIndex, depth
			{
				position1 := position
				depth++
				{
					position2 := position
					depth++
					if !_rules[ruleSlash]() {
						goto l0
					}
					{
						add(ruleAction1, position)
					}
					depth--
					add(ruleFirstSegment, position2)
				}
			l4:
				{
					position5, tokenIndex5, depth5 := position, tokenIndex, depth
					{
						position6 := position
						depth++
						{
							position7, tokenIndex7, depth7 := position, tokenIndex, depth
							{
								position9 := position
								depth++
								{
									position10, tokenIndex10, depth10 := position, tokenIndex, depth
									if !_rules[ruleSlash]() {
										goto l10
									}
									goto l11
								l10:
									position, tokenIndex, depth = position10, tokenIndex10, depth10
								}
							l11:
								{
									position12, tokenIndex12, depth12 := position, tokenIndex, depth
									if !_rules[ruleText]() {
										goto l12
									}
									goto l13
								l12:
									position, tokenIndex, depth = position12, tokenIndex12, depth12
								}
							l13:
								if !_rules[ruleKey]() {
									goto l8
								}
								if buffer[position] != rune('?') {
									goto l8
								}
								position++
								{
									add(ruleAction7, position)
								}
								depth--
								add(ruleOptionalKey, position9)
							}
							{
								add(ruleAction2, position)
							}
							goto l7
						l8:
							position, tokenIndex, depth = position7, tokenIndex7, depth7
							if !_rules[ruleSlash]() {
								goto l16
							}
							{
								add(ruleAction3, position)
							}
							goto l7
						l16:
							position, tokenIndex, depth = position7, tokenIndex7, depth7
							{
								position19 := position
								depth++
								if buffer[position] != rune('\\') {
									goto l18
								}
								position++
								{
									position20 := position
									depth++
									{
										position21, tokenIndex21, depth21 := position, tokenIndex, depth
										if buffer[position] != rune(':') {
											goto l22
										}
										position++
										goto l21
									l22:
										position, tokenIndex, depth = position21, tokenIndex21, depth21
										if buffer[position] != rune('/') {
											goto l18
										}
										position++
									}
								l21:
									depth--
									add(rulePegText, position20)
								}
								{
									add(ruleAction14, position)
								}
								depth--
								add(ruleEscaped, position19)
							}
							{
								add(ruleAction4, position)
							}
							goto l7
						l18:
							position, tokenIndex, depth = position7, tokenIndex7, depth7
							if !_rules[ruleKey]() {
								goto l25
							}
							{
								add(ruleAction5, position)
							}
							goto l7
						l25:
							position, tokenIndex, depth = position7, tokenIndex7, depth7
							if !_rules[ruleText]() {
								goto l5
							}
							{
								add(ruleAction6, position)
							}
						}
					l7:
						depth--
						add(ruleSegment, position6)
					}
					goto l4
				l5:
					position, tokenIndex, depth = position5, tokenIndex5, depth5
				}
				{
					position28 := position
					depth++
					{
						position29, tokenIndex29, depth29 := position, tokenIndex, depth
						if !matchDot() {
							goto l29
						}
						goto l0
					l29:
						position, tokenIndex, depth = position29, tokenIndex29, depth29
					}
					depth--
					add(ruleEnd, position28)
				}
				{
					add(ruleAction0, position)
				}
				depth--
				add(ruleRoute, position1)
			}
			return true
		l0:
			position, tokenIndex, depth = position0, tokenIndex0, depth0
			return false
		},
		/* 1 FirstSegment <- <(Slash Action1)> */
		nil,
		/* 2 Segment <- <((OptionalKey Action2) / (Slash Action3) / (Escaped Action4) / (Key Action5) / (Text Action6))> */
		nil,
		/* 3 OptionalKey <- <(Slash? Text? Key '?' Action7)> */
		nil,
		/* 4 Key <- <((RegexpKey Action8) / (BasicKey Action9))> */
		func() bool {
			position34, tokenIndex34, depth34 := position, tokenIndex, depth
			{
				position35 := position
				depth++
				{
					position36, tokenIndex36, depth36 := position, tokenIndex, depth
					{
						position38 := position
						depth++
						if buffer[position] != rune(':') {
							goto l37
						}
						position++
						if !_rules[ruleIdentifier]() {
							goto l37
						}
						{
							position39 := position
							depth++
							if buffer[position] != rune('(') {
								goto l37
							}
							position++
							{
								position40 := position
								depth++
								{
									position43, tokenIndex43, depth43 := position, tokenIndex, depth
									if buffer[position] != rune(')') {
										goto l43
									}
									position++
									goto l37
								l43:
									position, tokenIndex, depth = position43, tokenIndex43, depth43
								}
								if !matchDot() {
									goto l37
								}
							l41:
								{
									position42, tokenIndex42, depth42 := position, tokenIndex, depth
									{
										position44, tokenIndex44, depth44 := position, tokenIndex, depth
										if buffer[position] != rune(')') {
											goto l44
										}
										position++
										goto l42
									l44:
										position, tokenIndex, depth = position44, tokenIndex44, depth44
									}
									if !matchDot() {
										goto l42
									}
									goto l41
								l42:
									position, tokenIndex, depth = position42, tokenIndex42, depth42
								}
								depth--
								add(rulePegText, position40)
							}
							if buffer[position] != rune(')') {
								goto l37
							}
							position++
							{
								add(ruleAction12, position)
							}
							depth--
							add(ruleRegexp, position39)
						}
						{
							add(ruleAction10, position)
						}
						depth--
						add(ruleRegexpKey, position38)
					}
					{
						add(ruleAction8, position)
					}
					goto l36
				l37:
					position, tokenIndex, depth = position36, tokenIndex36, depth36
					{
						position48 := position
						depth++
						if buffer[position] != rune(':') {
							goto l34
						}
						position++
						if !_rules[ruleIdentifier]() {
							goto l34
						}
						{
							add(ruleAction11, position)
						}
						depth--
						add(ruleBasicKey, position48)
					}
					{
						add(ruleAction9, position)
					}
				}
			l36:
				depth--
				add(ruleKey, position35)
			}
			return true
		l34:
			position, tokenIndex, depth = position34, tokenIndex34, depth34
			return false
		},
		/* 5 RegexpKey <- <(':' Identifier Regexp Action10)> */
		nil,
		/* 6 BasicKey <- <(':' Identifier Action11)> */
		nil,
		/* 7 Regexp <- <('(' <(!')' .)+> ')' Action12)> */
		nil,
		/* 8 Identifier <- <(<(([A-Z] / [a-z]) ((&('0' | '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9') [0-9]) | (&('a' | 'b' | 'c' | 'd' | 'e' | 'f' | 'g' | 'h' | 'i' | 'j' | 'k' | 'l' | 'm' | 'n' | 'o' | 'p' | 'q' | 'r' | 's' | 't' | 'u' | 'v' | 'w' | 'x' | 'y' | 'z') [a-z]) | (&('A' | 'B' | 'C' | 'D' | 'E' | 'F' | 'G' | 'H' | 'I' | 'J' | 'K' | 'L' | 'M' | 'N' | 'O' | 'P' | 'Q' | 'R' | 'S' | 'T' | 'U' | 'V' | 'W' | 'X' | 'Y' | 'Z') [A-Z]))*)> Action13)> */
		func() bool {
			position54, tokenIndex54, depth54 := position, tokenIndex, depth
			{
				position55 := position
				depth++
				{
					position56 := position
					depth++
					{
						position57, tokenIndex57, depth57 := position, tokenIndex, depth
						if c := buffer[position]; c < rune('A') || c > rune('Z') {
							goto l58
						}
						position++
						goto l57
					l58:
						position, tokenIndex, depth = position57, tokenIndex57, depth57
						if c := buffer[position]; c < rune('a') || c > rune('z') {
							goto l54
						}
						position++
					}
				l57:
				l59:
					{
						position60, tokenIndex60, depth60 := position, tokenIndex, depth
						{
							switch buffer[position] {
							case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
								if c := buffer[position]; c < rune('0') || c > rune('9') {
									goto l60
								}
								position++
								break
							case 'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z':
								if c := buffer[position]; c < rune('a') || c > rune('z') {
									goto l60
								}
								position++
								break
							default:
								if c := buffer[position]; c < rune('A') || c > rune('Z') {
									goto l60
								}
								position++
								break
							}
						}

						goto l59
					l60:
						position, tokenIndex, depth = position60, tokenIndex60, depth60
					}
					depth--
					add(rulePegText, position56)
				}
				{
					add(ruleAction13, position)
				}
				depth--
				add(ruleIdentifier, position55)
			}
			return true
		l54:
			position, tokenIndex, depth = position54, tokenIndex54, depth54
			return false
		},
		/* 9 Escaped <- <('\\' <(':' / '/')> Action14)> */
		nil,
		/* 10 Text <- <(<(!(':' / '/') .)+> Action15)> */
		func() bool {
			position64, tokenIndex64, depth64 := position, tokenIndex, depth
			{
				position65 := position
				depth++
				{
					position66 := position
					depth++
					{
						position69, tokenIndex69, depth69 := position, tokenIndex, depth
						{
							position70, tokenIndex70, depth70 := position, tokenIndex, depth
							if buffer[position] != rune(':') {
								goto l71
							}
							position++
							goto l70
						l71:
							position, tokenIndex, depth = position70, tokenIndex70, depth70
							if buffer[position] != rune('/') {
								goto l69
							}
							position++
						}
					l70:
						goto l64
					l69:
						position, tokenIndex, depth = position69, tokenIndex69, depth69
					}
					if !matchDot() {
						goto l64
					}
				l67:
					{
						position68, tokenIndex68, depth68 := position, tokenIndex, depth
						{
							position72, tokenIndex72, depth72 := position, tokenIndex, depth
							{
								position73, tokenIndex73, depth73 := position, tokenIndex, depth
								if buffer[position] != rune(':') {
									goto l74
								}
								position++
								goto l73
							l74:
								position, tokenIndex, depth = position73, tokenIndex73, depth73
								if buffer[position] != rune('/') {
									goto l72
								}
								position++
							}
						l73:
							goto l68
						l72:
							position, tokenIndex, depth = position72, tokenIndex72, depth72
						}
						if !matchDot() {
							goto l68
						}
						goto l67
					l68:
						position, tokenIndex, depth = position68, tokenIndex68, depth68
					}
					depth--
					add(rulePegText, position66)
				}
				{
					add(ruleAction15, position)
				}
				depth--
				add(ruleText, position65)
			}
			return true
		l64:
			position, tokenIndex, depth = position64, tokenIndex64, depth64
			return false
		},
		/* 11 Slash <- <('/' Action16)> */
		func() bool {
			position76, tokenIndex76, depth76 := position, tokenIndex, depth
			{
				position77 := position
				depth++
				if buffer[position] != rune('/') {
					goto l76
				}
				position++
				{
					add(ruleAction16, position)
				}
				depth--
				add(ruleSlash, position77)
			}
			return true
		l76:
			position, tokenIndex, depth = position76, tokenIndex76, depth76
			return false
		},
		/* 12 End <- <!.> */
		nil,
		/* 14 Action0 <- <{
		  p.route = &route{
		    Segments: p.segments,
		  }
		  p.segments = nil
		}> */
		nil,
		/* 15 Action1 <- <{
		  p.segments = append(p.segments, p.slash)
		}> */
		nil,
		/* 16 Action2 <- <{ p.segments = append(p.segments, p.optionalKey); p.optionalKey = nil }> */
		nil,
		/* 17 Action3 <- <{ p.segments = append(p.segments, p.slash); p.slash = nil }> */
		nil,
		/* 18 Action4 <- <{ p.segments = append(p.segments, p.text); p.text = nil }> */
		nil,
		/* 19 Action5 <- <{ p.segments = append(p.segments, p.key); p.key = nil }> */
		nil,
		/* 20 Action6 <- <{ p.segments = append(p.segments, p.text); p.text = nil }> */
		nil,
		/* 21 Action7 <- <{
		  p.optionalKey = &OptionalKey{
		    PrefixSlash: p.slash,
		    PrefixText: p.text,
		    Key: p.key,
		  }
		  p.slash = nil
		  p.text = nil
		  p.key = nil
		}> */
		nil,
		/* 22 Action8 <- <{ p.key = p.regexpKey; p.regexpKey = nil }> */
		nil,
		/* 23 Action9 <- <{ p.key = p.basicKey; p.basicKey = nil }> */
		nil,
		/* 24 Action10 <- <{
		  p.regexpKey = &RegexpKey{
		    Name: p.identifier,
		    Regexp: p.regexp,
		  }
		  p.identifier = nil
		  p.regexp = nil
		}> */
		nil,
		/* 25 Action11 <- <{
		  p.basicKey = &BasicKey{Name: p.identifier}
		  p.identifier = nil
		}> */
		nil,
		nil,
		/* 27 Action12 <- <{
		  p.regexp = &Regexp{Value: text}
		}> */
		nil,
		/* 28 Action13 <- <{
		  p.identifier = &Identifier{Value: text}
		}> */
		nil,
		/* 29 Action14 <- <{
		  p.text = &Text{Value: text}
		}> */
		nil,
		/* 30 Action15 <- <{
		  p.text = &Text{Value: text}
		}> */
		nil,
		/* 31 Action16 <- <{
		  p.slash = &Slash{}
		}> */
		nil,
	}
	p.rules = _rules
}
