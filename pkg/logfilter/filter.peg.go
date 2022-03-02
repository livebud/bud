package logfilter

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
	ruleFilters
	ruleFilter
	ruleInnerFilter
	ruleLevelPackages
	rulePackages
	rulePackage
	ruleExclude
	ruleInclude
	ruleLevel
	ruleDebug
	ruleInfo
	ruleNotice
	ruleWarn
	ruleError
	rule_
	ruleEnd
	ruleAction0
	ruleAction1
	ruleAction2
	ruleAction3
	ruleAction4
	ruleAction5
	ruleAction6
	ruleAction7
	rulePegText
	ruleAction8
	ruleAction9
	ruleAction10
	ruleAction11
	ruleAction12
	ruleAction13
	ruleAction14
	ruleAction15
	ruleAction16
	ruleAction17
	ruleAction18
	ruleAction19

	rulePre
	ruleIn
	ruleSuf
)

var rul3s = [...]string{
	"Unknown",
	"Filters",
	"Filter",
	"InnerFilter",
	"LevelPackages",
	"Packages",
	"Package",
	"Exclude",
	"Include",
	"Level",
	"Debug",
	"Info",
	"Notice",
	"Warn",
	"Error",
	"_",
	"End",
	"Action0",
	"Action1",
	"Action2",
	"Action3",
	"Action4",
	"Action5",
	"Action6",
	"Action7",
	"PegText",
	"Action8",
	"Action9",
	"Action10",
	"Action11",
	"Action12",
	"Action13",
	"Action14",
	"Action15",
	"Action16",
	"Action17",
	"Action18",
	"Action19",

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
	filters filters

	// scratch
	filter filter
	pkgs   []string
	pkg    string
	level  string

	Buffer string
	buffer []rune
	rules  [38]func() bool
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

			p.filters = append(p.filters, p.filter)

		case ruleAction1:

			p.filter = filter{
				level:    p.level,
				packages: p.pkgs,
			}

		case ruleAction2:

		case ruleAction3:

		case ruleAction4:

		case ruleAction5:

		case ruleAction6:
			p.pkgs = append(p.pkgs, p.pkg)
		case ruleAction7:
			p.pkgs = append(p.pkgs, p.pkg)
		case ruleAction8:

			p.pkg = text

		case ruleAction9:

			p.pkg = text

		case ruleAction10:

		case ruleAction11:

		case ruleAction12:

		case ruleAction13:

		case ruleAction14:

		case ruleAction15:

			p.level = "debug"

		case ruleAction16:

			p.level = "info"

		case ruleAction17:

			p.level = "notice"

		case ruleAction18:

			p.level = "warn"

		case ruleAction19:

			p.level = "error"

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
		/* 0 Filters <- <(Filter* End Action0)> */
		func() bool {
			position0, tokenIndex0, depth0 := position, tokenIndex, depth
			{
				position1 := position
				depth++
			l2:
				{
					position3, tokenIndex3, depth3 := position, tokenIndex, depth
					{
						position4 := position
						depth++
						if !_rules[rule_]() {
							goto l3
						}
						{
							position5 := position
							depth++
							{
								position6, tokenIndex6, depth6 := position, tokenIndex, depth
								{
									position8 := position
									depth++
									if !_rules[ruleLevel]() {
										goto l7
									}
									if buffer[position] != rune(':') {
										goto l7
									}
									position++
									{
										position9 := position
										depth++
										if !_rules[rulePackage]() {
											goto l7
										}
									l10:
										{
											position11, tokenIndex11, depth11 := position, tokenIndex, depth
											if buffer[position] != rune(',') {
												goto l11
											}
											position++
											if !_rules[rulePackage]() {
												goto l11
											}
											goto l10
										l11:
											position, tokenIndex, depth = position11, tokenIndex11, depth11
										}
										{
											add(ruleAction5, position)
										}
										depth--
										add(rulePackages, position9)
									}
									{
										add(ruleAction4, position)
									}
									depth--
									add(ruleLevelPackages, position8)
								}
								{
									add(ruleAction2, position)
								}
								goto l6
							l7:
								position, tokenIndex, depth = position6, tokenIndex6, depth6
								if !_rules[ruleLevel]() {
									goto l3
								}
								{
									add(ruleAction3, position)
								}
							}
						l6:
							depth--
							add(ruleInnerFilter, position5)
						}
						if !_rules[rule_]() {
							goto l3
						}
						{
							add(ruleAction1, position)
						}
						depth--
						add(ruleFilter, position4)
					}
					goto l2
				l3:
					position, tokenIndex, depth = position3, tokenIndex3, depth3
				}
				{
					position17 := position
					depth++
					{
						position18, tokenIndex18, depth18 := position, tokenIndex, depth
						if !matchDot() {
							goto l18
						}
						goto l0
					l18:
						position, tokenIndex, depth = position18, tokenIndex18, depth18
					}
					depth--
					add(ruleEnd, position17)
				}
				{
					add(ruleAction0, position)
				}
				depth--
				add(ruleFilters, position1)
			}
			return true
		l0:
			position, tokenIndex, depth = position0, tokenIndex0, depth0
			return false
		},
		/* 1 Filter <- <(_ InnerFilter _ Action1)> */
		nil,
		/* 2 InnerFilter <- <((LevelPackages Action2) / (Level Action3))> */
		nil,
		/* 3 LevelPackages <- <(Level ':' Packages Action4)> */
		nil,
		/* 4 Packages <- <(Package (',' Package)* Action5)> */
		nil,
		/* 5 Package <- <((Include Action6) / (Exclude Action7))> */
		func() bool {
			position24, tokenIndex24, depth24 := position, tokenIndex, depth
			{
				position25 := position
				depth++
				{
					position26, tokenIndex26, depth26 := position, tokenIndex, depth
					if !_rules[ruleInclude]() {
						goto l27
					}
					{
						add(ruleAction6, position)
					}
					goto l26
				l27:
					position, tokenIndex, depth = position26, tokenIndex26, depth26
					{
						position29 := position
						depth++
						{
							position30 := position
							depth++
							if buffer[position] != rune('-') {
								goto l24
							}
							position++
							if !_rules[ruleInclude]() {
								goto l24
							}
							depth--
							add(rulePegText, position30)
						}
						{
							add(ruleAction8, position)
						}
						depth--
						add(ruleExclude, position29)
					}
					{
						add(ruleAction7, position)
					}
				}
			l26:
				depth--
				add(rulePackage, position25)
			}
			return true
		l24:
			position, tokenIndex, depth = position24, tokenIndex24, depth24
			return false
		},
		/* 6 Exclude <- <(<('-' Include)> Action8)> */
		nil,
		/* 7 Include <- <(<(([a-z] / [A-Z]) ((&('*') '*') | (&('.') '.') | (&('/') '/') | (&('0' | '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9') [0-9]) | (&('A' | 'B' | 'C' | 'D' | 'E' | 'F' | 'G' | 'H' | 'I' | 'J' | 'K' | 'L' | 'M' | 'N' | 'O' | 'P' | 'Q' | 'R' | 'S' | 'T' | 'U' | 'V' | 'W' | 'X' | 'Y' | 'Z') [A-Z]) | (&('a' | 'b' | 'c' | 'd' | 'e' | 'f' | 'g' | 'h' | 'i' | 'j' | 'k' | 'l' | 'm' | 'n' | 'o' | 'p' | 'q' | 'r' | 's' | 't' | 'u' | 'v' | 'w' | 'x' | 'y' | 'z') [a-z]))*)> Action9)> */
		func() bool {
			position34, tokenIndex34, depth34 := position, tokenIndex, depth
			{
				position35 := position
				depth++
				{
					position36 := position
					depth++
					{
						position37, tokenIndex37, depth37 := position, tokenIndex, depth
						if c := buffer[position]; c < rune('a') || c > rune('z') {
							goto l38
						}
						position++
						goto l37
					l38:
						position, tokenIndex, depth = position37, tokenIndex37, depth37
						if c := buffer[position]; c < rune('A') || c > rune('Z') {
							goto l34
						}
						position++
					}
				l37:
				l39:
					{
						position40, tokenIndex40, depth40 := position, tokenIndex, depth
						{
							switch buffer[position] {
							case '*':
								if buffer[position] != rune('*') {
									goto l40
								}
								position++
								break
							case '.':
								if buffer[position] != rune('.') {
									goto l40
								}
								position++
								break
							case '/':
								if buffer[position] != rune('/') {
									goto l40
								}
								position++
								break
							case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
								if c := buffer[position]; c < rune('0') || c > rune('9') {
									goto l40
								}
								position++
								break
							case 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z':
								if c := buffer[position]; c < rune('A') || c > rune('Z') {
									goto l40
								}
								position++
								break
							default:
								if c := buffer[position]; c < rune('a') || c > rune('z') {
									goto l40
								}
								position++
								break
							}
						}

						goto l39
					l40:
						position, tokenIndex, depth = position40, tokenIndex40, depth40
					}
					depth--
					add(rulePegText, position36)
				}
				{
					add(ruleAction9, position)
				}
				depth--
				add(ruleInclude, position35)
			}
			return true
		l34:
			position, tokenIndex, depth = position34, tokenIndex34, depth34
			return false
		},
		/* 8 Level <- <((&('E' | 'e') (Error Action14)) | (&('W' | 'w') (Warn Action13)) | (&('N' | 'n') (Notice Action12)) | (&('I' | 'i') (Info Action11)) | (&('D' | 'd') (Debug Action10)))> */
		func() bool {
			position43, tokenIndex43, depth43 := position, tokenIndex, depth
			{
				position44 := position
				depth++
				{
					switch buffer[position] {
					case 'E', 'e':
						{
							position46 := position
							depth++
							{
								position47 := position
								depth++
								{
									position48, tokenIndex48, depth48 := position, tokenIndex, depth
									if buffer[position] != rune('e') {
										goto l49
									}
									position++
									goto l48
								l49:
									position, tokenIndex, depth = position48, tokenIndex48, depth48
									if buffer[position] != rune('E') {
										goto l43
									}
									position++
								}
							l48:
								{
									position50, tokenIndex50, depth50 := position, tokenIndex, depth
									if buffer[position] != rune('r') {
										goto l51
									}
									position++
									goto l50
								l51:
									position, tokenIndex, depth = position50, tokenIndex50, depth50
									if buffer[position] != rune('R') {
										goto l43
									}
									position++
								}
							l50:
								{
									position52, tokenIndex52, depth52 := position, tokenIndex, depth
									if buffer[position] != rune('r') {
										goto l53
									}
									position++
									goto l52
								l53:
									position, tokenIndex, depth = position52, tokenIndex52, depth52
									if buffer[position] != rune('R') {
										goto l43
									}
									position++
								}
							l52:
								{
									position54, tokenIndex54, depth54 := position, tokenIndex, depth
									if buffer[position] != rune('o') {
										goto l55
									}
									position++
									goto l54
								l55:
									position, tokenIndex, depth = position54, tokenIndex54, depth54
									if buffer[position] != rune('O') {
										goto l43
									}
									position++
								}
							l54:
								{
									position56, tokenIndex56, depth56 := position, tokenIndex, depth
									if buffer[position] != rune('r') {
										goto l57
									}
									position++
									goto l56
								l57:
									position, tokenIndex, depth = position56, tokenIndex56, depth56
									if buffer[position] != rune('R') {
										goto l43
									}
									position++
								}
							l56:
								depth--
								add(rulePegText, position47)
							}
							{
								add(ruleAction19, position)
							}
							depth--
							add(ruleError, position46)
						}
						{
							add(ruleAction14, position)
						}
						break
					case 'W', 'w':
						{
							position60 := position
							depth++
							{
								position61 := position
								depth++
								{
									position62, tokenIndex62, depth62 := position, tokenIndex, depth
									if buffer[position] != rune('w') {
										goto l63
									}
									position++
									goto l62
								l63:
									position, tokenIndex, depth = position62, tokenIndex62, depth62
									if buffer[position] != rune('W') {
										goto l43
									}
									position++
								}
							l62:
								{
									position64, tokenIndex64, depth64 := position, tokenIndex, depth
									if buffer[position] != rune('a') {
										goto l65
									}
									position++
									goto l64
								l65:
									position, tokenIndex, depth = position64, tokenIndex64, depth64
									if buffer[position] != rune('A') {
										goto l43
									}
									position++
								}
							l64:
								{
									position66, tokenIndex66, depth66 := position, tokenIndex, depth
									if buffer[position] != rune('r') {
										goto l67
									}
									position++
									goto l66
								l67:
									position, tokenIndex, depth = position66, tokenIndex66, depth66
									if buffer[position] != rune('R') {
										goto l43
									}
									position++
								}
							l66:
								{
									position68, tokenIndex68, depth68 := position, tokenIndex, depth
									if buffer[position] != rune('n') {
										goto l69
									}
									position++
									goto l68
								l69:
									position, tokenIndex, depth = position68, tokenIndex68, depth68
									if buffer[position] != rune('N') {
										goto l43
									}
									position++
								}
							l68:
								depth--
								add(rulePegText, position61)
							}
							{
								add(ruleAction18, position)
							}
							depth--
							add(ruleWarn, position60)
						}
						{
							add(ruleAction13, position)
						}
						break
					case 'N', 'n':
						{
							position72 := position
							depth++
							{
								position73 := position
								depth++
								{
									position74, tokenIndex74, depth74 := position, tokenIndex, depth
									if buffer[position] != rune('n') {
										goto l75
									}
									position++
									goto l74
								l75:
									position, tokenIndex, depth = position74, tokenIndex74, depth74
									if buffer[position] != rune('N') {
										goto l43
									}
									position++
								}
							l74:
								{
									position76, tokenIndex76, depth76 := position, tokenIndex, depth
									if buffer[position] != rune('o') {
										goto l77
									}
									position++
									goto l76
								l77:
									position, tokenIndex, depth = position76, tokenIndex76, depth76
									if buffer[position] != rune('O') {
										goto l43
									}
									position++
								}
							l76:
								{
									position78, tokenIndex78, depth78 := position, tokenIndex, depth
									if buffer[position] != rune('t') {
										goto l79
									}
									position++
									goto l78
								l79:
									position, tokenIndex, depth = position78, tokenIndex78, depth78
									if buffer[position] != rune('T') {
										goto l43
									}
									position++
								}
							l78:
								{
									position80, tokenIndex80, depth80 := position, tokenIndex, depth
									if buffer[position] != rune('i') {
										goto l81
									}
									position++
									goto l80
								l81:
									position, tokenIndex, depth = position80, tokenIndex80, depth80
									if buffer[position] != rune('I') {
										goto l43
									}
									position++
								}
							l80:
								{
									position82, tokenIndex82, depth82 := position, tokenIndex, depth
									if buffer[position] != rune('c') {
										goto l83
									}
									position++
									goto l82
								l83:
									position, tokenIndex, depth = position82, tokenIndex82, depth82
									if buffer[position] != rune('C') {
										goto l43
									}
									position++
								}
							l82:
								{
									position84, tokenIndex84, depth84 := position, tokenIndex, depth
									if buffer[position] != rune('e') {
										goto l85
									}
									position++
									goto l84
								l85:
									position, tokenIndex, depth = position84, tokenIndex84, depth84
									if buffer[position] != rune('E') {
										goto l43
									}
									position++
								}
							l84:
								depth--
								add(rulePegText, position73)
							}
							{
								add(ruleAction17, position)
							}
							depth--
							add(ruleNotice, position72)
						}
						{
							add(ruleAction12, position)
						}
						break
					case 'I', 'i':
						{
							position88 := position
							depth++
							{
								position89 := position
								depth++
								{
									position90, tokenIndex90, depth90 := position, tokenIndex, depth
									if buffer[position] != rune('i') {
										goto l91
									}
									position++
									goto l90
								l91:
									position, tokenIndex, depth = position90, tokenIndex90, depth90
									if buffer[position] != rune('I') {
										goto l43
									}
									position++
								}
							l90:
								{
									position92, tokenIndex92, depth92 := position, tokenIndex, depth
									if buffer[position] != rune('n') {
										goto l93
									}
									position++
									goto l92
								l93:
									position, tokenIndex, depth = position92, tokenIndex92, depth92
									if buffer[position] != rune('N') {
										goto l43
									}
									position++
								}
							l92:
								{
									position94, tokenIndex94, depth94 := position, tokenIndex, depth
									if buffer[position] != rune('f') {
										goto l95
									}
									position++
									goto l94
								l95:
									position, tokenIndex, depth = position94, tokenIndex94, depth94
									if buffer[position] != rune('F') {
										goto l43
									}
									position++
								}
							l94:
								{
									position96, tokenIndex96, depth96 := position, tokenIndex, depth
									if buffer[position] != rune('o') {
										goto l97
									}
									position++
									goto l96
								l97:
									position, tokenIndex, depth = position96, tokenIndex96, depth96
									if buffer[position] != rune('O') {
										goto l43
									}
									position++
								}
							l96:
								depth--
								add(rulePegText, position89)
							}
							{
								add(ruleAction16, position)
							}
							depth--
							add(ruleInfo, position88)
						}
						{
							add(ruleAction11, position)
						}
						break
					default:
						{
							position100 := position
							depth++
							{
								position101 := position
								depth++
								{
									position102, tokenIndex102, depth102 := position, tokenIndex, depth
									if buffer[position] != rune('d') {
										goto l103
									}
									position++
									goto l102
								l103:
									position, tokenIndex, depth = position102, tokenIndex102, depth102
									if buffer[position] != rune('D') {
										goto l43
									}
									position++
								}
							l102:
								{
									position104, tokenIndex104, depth104 := position, tokenIndex, depth
									if buffer[position] != rune('e') {
										goto l105
									}
									position++
									goto l104
								l105:
									position, tokenIndex, depth = position104, tokenIndex104, depth104
									if buffer[position] != rune('E') {
										goto l43
									}
									position++
								}
							l104:
								{
									position106, tokenIndex106, depth106 := position, tokenIndex, depth
									if buffer[position] != rune('b') {
										goto l107
									}
									position++
									goto l106
								l107:
									position, tokenIndex, depth = position106, tokenIndex106, depth106
									if buffer[position] != rune('B') {
										goto l43
									}
									position++
								}
							l106:
								{
									position108, tokenIndex108, depth108 := position, tokenIndex, depth
									if buffer[position] != rune('u') {
										goto l109
									}
									position++
									goto l108
								l109:
									position, tokenIndex, depth = position108, tokenIndex108, depth108
									if buffer[position] != rune('U') {
										goto l43
									}
									position++
								}
							l108:
								{
									position110, tokenIndex110, depth110 := position, tokenIndex, depth
									if buffer[position] != rune('g') {
										goto l111
									}
									position++
									goto l110
								l111:
									position, tokenIndex, depth = position110, tokenIndex110, depth110
									if buffer[position] != rune('G') {
										goto l43
									}
									position++
								}
							l110:
								depth--
								add(rulePegText, position101)
							}
							{
								add(ruleAction15, position)
							}
							depth--
							add(ruleDebug, position100)
						}
						{
							add(ruleAction10, position)
						}
						break
					}
				}

				depth--
				add(ruleLevel, position44)
			}
			return true
		l43:
			position, tokenIndex, depth = position43, tokenIndex43, depth43
			return false
		},
		/* 9 Debug <- <(<(('d' / 'D') ('e' / 'E') ('b' / 'B') ('u' / 'U') ('g' / 'G'))> Action15)> */
		nil,
		/* 10 Info <- <(<(('i' / 'I') ('n' / 'N') ('f' / 'F') ('o' / 'O'))> Action16)> */
		nil,
		/* 11 Notice <- <(<(('n' / 'N') ('o' / 'O') ('t' / 'T') ('i' / 'I') ('c' / 'C') ('e' / 'E'))> Action17)> */
		nil,
		/* 12 Warn <- <(<(('w' / 'W') ('a' / 'A') ('r' / 'R') ('n' / 'N'))> Action18)> */
		nil,
		/* 13 Error <- <(<(('e' / 'E') ('r' / 'R') ('r' / 'R') ('o' / 'O') ('r' / 'R'))> Action19)> */
		nil,
		/* 14 _ <- <((&('\r') '\r') | (&('\n') '\n') | (&('\t') '\t') | (&(' ') ' '))*> */
		func() bool {
			{
				position120 := position
				depth++
			l121:
				{
					position122, tokenIndex122, depth122 := position, tokenIndex, depth
					{
						switch buffer[position] {
						case '\r':
							if buffer[position] != rune('\r') {
								goto l122
							}
							position++
							break
						case '\n':
							if buffer[position] != rune('\n') {
								goto l122
							}
							position++
							break
						case '\t':
							if buffer[position] != rune('\t') {
								goto l122
							}
							position++
							break
						default:
							if buffer[position] != rune(' ') {
								goto l122
							}
							position++
							break
						}
					}

					goto l121
				l122:
					position, tokenIndex, depth = position122, tokenIndex122, depth122
				}
				depth--
				add(rule_, position120)
			}
			return true
		},
		/* 15 End <- <!.> */
		nil,
		/* 17 Action0 <- <{
		  p.filters = append(p.filters, p.filter)
		}> */
		nil,
		/* 18 Action1 <- <{
		  p.filter = filter{
		    level: p.level,
		    packages: p.pkgs,
		  }
		}> */
		nil,
		/* 19 Action2 <- <{}> */
		nil,
		/* 20 Action3 <- <{}> */
		nil,
		/* 21 Action4 <- <{}> */
		nil,
		/* 22 Action5 <- <{
		}> */
		nil,
		/* 23 Action6 <- <{ p.pkgs = append(p.pkgs, p.pkg) }> */
		nil,
		/* 24 Action7 <- <{ p.pkgs = append(p.pkgs, p.pkg) }> */
		nil,
		nil,
		/* 26 Action8 <- <{
		  p.pkg = text
		}> */
		nil,
		/* 27 Action9 <- <{
		  p.pkg = text
		}> */
		nil,
		/* 28 Action10 <- <{}> */
		nil,
		/* 29 Action11 <- <{}> */
		nil,
		/* 30 Action12 <- <{}> */
		nil,
		/* 31 Action13 <- <{}> */
		nil,
		/* 32 Action14 <- <{}> */
		nil,
		/* 33 Action15 <- <{
		  p.level = "debug"
		}> */
		nil,
		/* 34 Action16 <- <{
		  p.level = "info"
		}> */
		nil,
		/* 35 Action17 <- <{
		  p.level = "notice"
		}> */
		nil,
		/* 36 Action18 <- <{
		  p.level = "warn"
		}> */
		nil,
		/* 37 Action19 <- <{
		  p.level = "error"
		}> */
		nil,
	}
	p.rules = _rules
}
