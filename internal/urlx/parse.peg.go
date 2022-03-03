package urlx

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
	ruleURL
	ruleURI
	ruleScheme
	ruleHost
	ruleIPPort
	ruleHostNamePort
	ruleBracketsPort
	ruleIP
	ruleIPV4
	ruleHostName
	ruleOnlyPort
	rulePort
	rulePath
	ruleBrackets
	ruleEnd
	rulePegText
	ruleAction0
	ruleAction1
	ruleAction2
	ruleAction3
	ruleAction4
	ruleAction5
	ruleAction6

	rulePre
	ruleIn
	ruleSuf
)

var rul3s = [...]string{
	"Unknown",
	"URL",
	"URI",
	"Scheme",
	"Host",
	"IPPort",
	"HostNamePort",
	"BracketsPort",
	"IP",
	"IPV4",
	"HostName",
	"OnlyPort",
	"Port",
	"Path",
	"Brackets",
	"End",
	"PegText",
	"Action0",
	"Action1",
	"Action2",
	"Action3",
	"Action4",
	"Action5",
	"Action6",

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
	url uri

	Buffer string
	buffer []rune
	rules  [24]func() bool
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

			p.url.uri = text

		case ruleAction1:

			p.url.scheme = text[:len(text)-1]

		case ruleAction2:

			p.url.host = text

		case ruleAction3:

			p.url.host = text

		case ruleAction4:

			p.url.port = text

		case ruleAction5:

			p.url.path = text

		case ruleAction6:

			p.url.host = "[::]"

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
		/* 0 URL <- <(URI / Path / Scheme / Host / (OnlyPort End))> */
		func() bool {
			position0, tokenIndex0, depth0 := position, tokenIndex, depth
			{
				position1 := position
				depth++
				{
					position2, tokenIndex2, depth2 := position, tokenIndex, depth
					{
						position4 := position
						depth++
						{
							position5 := position
							depth++
							if !_rules[ruleScheme]() {
								goto l3
							}
							if buffer[position] != rune('/') {
								goto l3
							}
							position++
							if buffer[position] != rune('/') {
								goto l3
							}
							position++
							if !_rules[ruleHost]() {
								goto l3
							}
							{
								position6, tokenIndex6, depth6 := position, tokenIndex, depth
								if !_rules[rulePath]() {
									goto l6
								}
								goto l7
							l6:
								position, tokenIndex, depth = position6, tokenIndex6, depth6
							}
						l7:
							depth--
							add(rulePegText, position5)
						}
						{
							add(ruleAction0, position)
						}
						depth--
						add(ruleURI, position4)
					}
					goto l2
				l3:
					position, tokenIndex, depth = position2, tokenIndex2, depth2
					if !_rules[rulePath]() {
						goto l9
					}
					goto l2
				l9:
					position, tokenIndex, depth = position2, tokenIndex2, depth2
					if !_rules[ruleScheme]() {
						goto l10
					}
					goto l2
				l10:
					position, tokenIndex, depth = position2, tokenIndex2, depth2
					if !_rules[ruleHost]() {
						goto l11
					}
					goto l2
				l11:
					position, tokenIndex, depth = position2, tokenIndex2, depth2
					{
						position12 := position
						depth++
						{
							position13, tokenIndex13, depth13 := position, tokenIndex, depth
							if buffer[position] != rune(':') {
								goto l14
							}
							position++
							if !_rules[rulePort]() {
								goto l14
							}
							goto l13
						l14:
							position, tokenIndex, depth = position13, tokenIndex13, depth13
							if !_rules[rulePort]() {
								goto l0
							}
						}
					l13:
						depth--
						add(ruleOnlyPort, position12)
					}
					{
						position15 := position
						depth++
						{
							position16, tokenIndex16, depth16 := position, tokenIndex, depth
							if !matchDot() {
								goto l16
							}
							goto l0
						l16:
							position, tokenIndex, depth = position16, tokenIndex16, depth16
						}
						depth--
						add(ruleEnd, position15)
					}
				}
			l2:
				depth--
				add(ruleURL, position1)
			}
			return true
		l0:
			position, tokenIndex, depth = position0, tokenIndex0, depth0
			return false
		},
		/* 1 URI <- <(<(Scheme ('/' '/') Host Path?)> Action0)> */
		nil,
		/* 2 Scheme <- <(<(([a-z] / [A-Z]) ((&('0' | '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9') [0-9]) | (&('+') '+') | (&('A' | 'B' | 'C' | 'D' | 'E' | 'F' | 'G' | 'H' | 'I' | 'J' | 'K' | 'L' | 'M' | 'N' | 'O' | 'P' | 'Q' | 'R' | 'S' | 'T' | 'U' | 'V' | 'W' | 'X' | 'Y' | 'Z') [A-Z]) | (&('a' | 'b' | 'c' | 'd' | 'e' | 'f' | 'g' | 'h' | 'i' | 'j' | 'k' | 'l' | 'm' | 'n' | 'o' | 'p' | 'q' | 'r' | 's' | 't' | 'u' | 'v' | 'w' | 'x' | 'y' | 'z') [a-z]))* ':')> Action1)> */
		func() bool {
			position18, tokenIndex18, depth18 := position, tokenIndex, depth
			{
				position19 := position
				depth++
				{
					position20 := position
					depth++
					{
						position21, tokenIndex21, depth21 := position, tokenIndex, depth
						if c := buffer[position]; c < rune('a') || c > rune('z') {
							goto l22
						}
						position++
						goto l21
					l22:
						position, tokenIndex, depth = position21, tokenIndex21, depth21
						if c := buffer[position]; c < rune('A') || c > rune('Z') {
							goto l18
						}
						position++
					}
				l21:
				l23:
					{
						position24, tokenIndex24, depth24 := position, tokenIndex, depth
						{
							switch buffer[position] {
							case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
								if c := buffer[position]; c < rune('0') || c > rune('9') {
									goto l24
								}
								position++
								break
							case '+':
								if buffer[position] != rune('+') {
									goto l24
								}
								position++
								break
							case 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z':
								if c := buffer[position]; c < rune('A') || c > rune('Z') {
									goto l24
								}
								position++
								break
							default:
								if c := buffer[position]; c < rune('a') || c > rune('z') {
									goto l24
								}
								position++
								break
							}
						}

						goto l23
					l24:
						position, tokenIndex, depth = position24, tokenIndex24, depth24
					}
					if buffer[position] != rune(':') {
						goto l18
					}
					position++
					depth--
					add(rulePegText, position20)
				}
				{
					add(ruleAction1, position)
				}
				depth--
				add(ruleScheme, position19)
			}
			return true
		l18:
			position, tokenIndex, depth = position18, tokenIndex18, depth18
			return false
		},
		/* 3 Host <- <(IPPort / HostNamePort / BracketsPort / ((&('[') Brackets) | (&('0' | '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9') IPV4) | (&('A' | 'B' | 'C' | 'D' | 'E' | 'F' | 'G' | 'H' | 'I' | 'J' | 'K' | 'L' | 'M' | 'N' | 'O' | 'P' | 'Q' | 'R' | 'S' | 'T' | 'U' | 'V' | 'W' | 'X' | 'Y' | 'Z' | 'a' | 'b' | 'c' | 'd' | 'e' | 'f' | 'g' | 'h' | 'i' | 'j' | 'k' | 'l' | 'm' | 'n' | 'o' | 'p' | 'q' | 'r' | 's' | 't' | 'u' | 'v' | 'w' | 'x' | 'y' | 'z') HostName)))> */
		func() bool {
			position27, tokenIndex27, depth27 := position, tokenIndex, depth
			{
				position28 := position
				depth++
				{
					position29, tokenIndex29, depth29 := position, tokenIndex, depth
					{
						position31 := position
						depth++
						{
							position32 := position
							depth++
							if !_rules[ruleIPV4]() {
								goto l30
							}
							depth--
							add(ruleIP, position32)
						}
						if buffer[position] != rune(':') {
							goto l30
						}
						position++
						if !_rules[rulePort]() {
							goto l30
						}
						depth--
						add(ruleIPPort, position31)
					}
					goto l29
				l30:
					position, tokenIndex, depth = position29, tokenIndex29, depth29
					{
						position34 := position
						depth++
						if !_rules[ruleHostName]() {
							goto l33
						}
						if buffer[position] != rune(':') {
							goto l33
						}
						position++
						if !_rules[rulePort]() {
							goto l33
						}
						depth--
						add(ruleHostNamePort, position34)
					}
					goto l29
				l33:
					position, tokenIndex, depth = position29, tokenIndex29, depth29
					{
						position36 := position
						depth++
						if !_rules[ruleBrackets]() {
							goto l35
						}
						if buffer[position] != rune(':') {
							goto l35
						}
						position++
						if !_rules[rulePort]() {
							goto l35
						}
						depth--
						add(ruleBracketsPort, position36)
					}
					goto l29
				l35:
					position, tokenIndex, depth = position29, tokenIndex29, depth29
					{
						switch buffer[position] {
						case '[':
							if !_rules[ruleBrackets]() {
								goto l27
							}
							break
						case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
							if !_rules[ruleIPV4]() {
								goto l27
							}
							break
						default:
							if !_rules[ruleHostName]() {
								goto l27
							}
							break
						}
					}

				}
			l29:
				depth--
				add(ruleHost, position28)
			}
			return true
		l27:
			position, tokenIndex, depth = position27, tokenIndex27, depth27
			return false
		},
		/* 4 IPPort <- <(IP ':' Port)> */
		nil,
		/* 5 HostNamePort <- <(HostName ':' Port)> */
		nil,
		/* 6 BracketsPort <- <(Brackets ':' Port)> */
		nil,
		/* 7 IP <- <IPV4> */
		nil,
		/* 8 IPV4 <- <(<([0-9]+ '.' [0-9]+ '.' [0-9]+ '.' [0-9]+)> Action2)> */
		func() bool {
			position42, tokenIndex42, depth42 := position, tokenIndex, depth
			{
				position43 := position
				depth++
				{
					position44 := position
					depth++
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l42
					}
					position++
				l45:
					{
						position46, tokenIndex46, depth46 := position, tokenIndex, depth
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l46
						}
						position++
						goto l45
					l46:
						position, tokenIndex, depth = position46, tokenIndex46, depth46
					}
					if buffer[position] != rune('.') {
						goto l42
					}
					position++
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l42
					}
					position++
				l47:
					{
						position48, tokenIndex48, depth48 := position, tokenIndex, depth
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l48
						}
						position++
						goto l47
					l48:
						position, tokenIndex, depth = position48, tokenIndex48, depth48
					}
					if buffer[position] != rune('.') {
						goto l42
					}
					position++
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l42
					}
					position++
				l49:
					{
						position50, tokenIndex50, depth50 := position, tokenIndex, depth
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l50
						}
						position++
						goto l49
					l50:
						position, tokenIndex, depth = position50, tokenIndex50, depth50
					}
					if buffer[position] != rune('.') {
						goto l42
					}
					position++
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l42
					}
					position++
				l51:
					{
						position52, tokenIndex52, depth52 := position, tokenIndex, depth
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l52
						}
						position++
						goto l51
					l52:
						position, tokenIndex, depth = position52, tokenIndex52, depth52
					}
					depth--
					add(rulePegText, position44)
				}
				{
					add(ruleAction2, position)
				}
				depth--
				add(ruleIPV4, position43)
			}
			return true
		l42:
			position, tokenIndex, depth = position42, tokenIndex42, depth42
			return false
		},
		/* 9 HostName <- <(<(([a-z] / [A-Z]) ((&('0' | '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9') [0-9]) | (&('A' | 'B' | 'C' | 'D' | 'E' | 'F' | 'G' | 'H' | 'I' | 'J' | 'K' | 'L' | 'M' | 'N' | 'O' | 'P' | 'Q' | 'R' | 'S' | 'T' | 'U' | 'V' | 'W' | 'X' | 'Y' | 'Z') [A-Z]) | (&('a' | 'b' | 'c' | 'd' | 'e' | 'f' | 'g' | 'h' | 'i' | 'j' | 'k' | 'l' | 'm' | 'n' | 'o' | 'p' | 'q' | 'r' | 's' | 't' | 'u' | 'v' | 'w' | 'x' | 'y' | 'z') [a-z]))*)> Action3)> */
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
						if c := buffer[position]; c < rune('a') || c > rune('z') {
							goto l58
						}
						position++
						goto l57
					l58:
						position, tokenIndex, depth = position57, tokenIndex57, depth57
						if c := buffer[position]; c < rune('A') || c > rune('Z') {
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
							case 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z':
								if c := buffer[position]; c < rune('A') || c > rune('Z') {
									goto l60
								}
								position++
								break
							default:
								if c := buffer[position]; c < rune('a') || c > rune('z') {
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
					add(ruleAction3, position)
				}
				depth--
				add(ruleHostName, position55)
			}
			return true
		l54:
			position, tokenIndex, depth = position54, tokenIndex54, depth54
			return false
		},
		/* 10 OnlyPort <- <((':' Port) / Port)> */
		nil,
		/* 11 Port <- <(<[0-9]+> Action4)> */
		func() bool {
			position64, tokenIndex64, depth64 := position, tokenIndex, depth
			{
				position65 := position
				depth++
				{
					position66 := position
					depth++
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l64
					}
					position++
				l67:
					{
						position68, tokenIndex68, depth68 := position, tokenIndex, depth
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l68
						}
						position++
						goto l67
					l68:
						position, tokenIndex, depth = position68, tokenIndex68, depth68
					}
					depth--
					add(rulePegText, position66)
				}
				{
					add(ruleAction4, position)
				}
				depth--
				add(rulePort, position65)
			}
			return true
		l64:
			position, tokenIndex, depth = position64, tokenIndex64, depth64
			return false
		},
		/* 12 Path <- <(<('/' .*)> Action5)> */
		func() bool {
			position70, tokenIndex70, depth70 := position, tokenIndex, depth
			{
				position71 := position
				depth++
				{
					position72 := position
					depth++
					if buffer[position] != rune('/') {
						goto l70
					}
					position++
				l73:
					{
						position74, tokenIndex74, depth74 := position, tokenIndex, depth
						if !matchDot() {
							goto l74
						}
						goto l73
					l74:
						position, tokenIndex, depth = position74, tokenIndex74, depth74
					}
					depth--
					add(rulePegText, position72)
				}
				{
					add(ruleAction5, position)
				}
				depth--
				add(rulePath, position71)
			}
			return true
		l70:
			position, tokenIndex, depth = position70, tokenIndex70, depth70
			return false
		},
		/* 13 Brackets <- <('[' ':' ':' ']' Action6)> */
		func() bool {
			position76, tokenIndex76, depth76 := position, tokenIndex, depth
			{
				position77 := position
				depth++
				if buffer[position] != rune('[') {
					goto l76
				}
				position++
				if buffer[position] != rune(':') {
					goto l76
				}
				position++
				if buffer[position] != rune(':') {
					goto l76
				}
				position++
				if buffer[position] != rune(']') {
					goto l76
				}
				position++
				{
					add(ruleAction6, position)
				}
				depth--
				add(ruleBrackets, position77)
			}
			return true
		l76:
			position, tokenIndex, depth = position76, tokenIndex76, depth76
			return false
		},
		/* 14 End <- <!.> */
		nil,
		nil,
		/* 17 Action0 <- <{
		  p.url.uri = text
		}> */
		nil,
		/* 18 Action1 <- <{
		  p.url.scheme = text[:len(text)-1]
		}> */
		nil,
		/* 19 Action2 <- <{
		  p.url.host = text
		}> */
		nil,
		/* 20 Action3 <- <{
		  p.url.host = text
		}> */
		nil,
		/* 21 Action4 <- <{
		  p.url.port = text
		}> */
		nil,
		/* 22 Action5 <- <{
		  p.url.path = text
		}> */
		nil,
		/* 23 Action6 <- <{
		  p.url.host = "[::]"
		}> */
		nil,
	}
	p.rules = _rules
}
