package glob

import (
	"fmt"

	"github.com/gobwas/glob/syntax/ast"
	"github.com/gobwas/glob/syntax/lexer"
	"github.com/livebud/bud/internal/orderedset"
)

func Expand(str string) ([]string, error) {
	lex := lexer.NewLexer(str)
	node, err := ast.Parse(lex)
	if err != nil {
		return nil, err
	}
	patterns, err := expand(node)
	if err != nil {
		return nil, err
	}
	return orderedset.Strings(patterns...), nil
}

func expand(node *ast.Node) (patterns []string, err error) {
	prefix := ""
	write := func(value string) {
		prefix += value
		for i := range patterns {
			patterns[i] += value
		}
	}
	for _, child := range node.Children {
		switch child.Kind {
		case ast.KindText:
			text, ok := child.Value.(ast.Text)
			if !ok {
				return nil, fmt.Errorf("expected text value, got %T", child.Value)
			}
			write(text.Text)
		case ast.KindList:
			list, ok := child.Value.(ast.List)
			if !ok {
				return nil, fmt.Errorf("expected list value, got %T", child.Value)
			}
			write("[")
			if list.Not {
				write("^")
			}
			write(list.Chars)
			write("]")
		case ast.KindRange:
			rng, ok := child.Value.(ast.Range)
			if !ok {
				return nil, fmt.Errorf("expected rng value, got %T", child.Value)
			}
			write("[")
			if rng.Not {
				write("^")
			}
			write(string(rng.Lo))
			write("-")
			write(string(rng.Hi))
			write("]")
		case ast.KindAny:
			write("*")
		case ast.KindSuper:
			write("**")
		case ast.KindSingle:
			write("?")
		case ast.KindAnyOf:
			for _, child := range child.Children {
				results, err := expand(child)
				if err != nil {
					return nil, err
				}
				for _, result := range results {
					patterns = append(patterns, prefix+result)
				}
			}
		default:
			return nil, fmt.Errorf("unknown node kind: %v", child.Kind)
		}
	}
	if len(patterns) == 0 {
		patterns = append(patterns, prefix)
	}
	return patterns, nil
}
