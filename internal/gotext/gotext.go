// Package text is a wrapper over matthewmueller/text for generating
// Go-compatible code. This means replacing keywords with valid alternatives
// and handling common acronyms that golint warns about.
//
// One of the goals for this package is to have orthogonal functions. For
// example, text.Snake("Hi World") should return `Hi_World`, not `hi_word`. This
// can be challenging to pull off at times and I don't think we're quite there
// yet. Once we feel like we've reached the correct balance, I think it makes
// sense to pull out this package and make it available for everyone.
package gotext

import (
	"strings"

	"github.com/matthewmueller/text"
)

// Slim case. Often used for package names.
func Slim(s ...string) string {
	return unreserve(text.Slim(text.Lower(strings.Join(s, " "))))
}

// Pascal case. Often used for declaration names.
func Pascal(s ...string) string {
	a := strings.Split(text.Space(strings.Join(s, " ")), " ")
	for i, w := range a {
		a[i] = acronym(text.Title(w))
	}
	return unreserve(strings.Join(a, ""))
}

// Camel case. Often used for variable names.
func Camel(s ...string) string {
	a := strings.Split(text.Space(strings.Join(s, " ")), " ")
	for i, w := range a {
		if i == 0 {
			a[i] = text.Lower(w)
			continue
		}
		a[i] = acronym(text.Title(w))
	}
	return unreserve(strings.Join(a, ""))
}

// Short case. Often used for receiver names.
func Short(s ...string) string {
	return unreserve(text.Short(text.Lower(text.Space(strings.Join(s, " ")))))
}

// Snake case. Often used for keys and json struct tags.
// Missing unreserve is intentional, the result isn't meant to be an identifier.
func Snake(s ...string) string {
	return text.Snake(text.Lower(text.Space(strings.Join(s, " "))))
}

// Upper case. Often used for environment struct tags.
// Missing unreserve is intentional, the result isn't meant to be an identifier.
func Upper(s ...string) string {
	return text.Upper(text.Snake(text.Space(strings.Join(s, " "))))
}

// Path case. Often used for HTTP paths.
// Missing unreserve is intentional, the result isn't meant to be an identifier.
func Path(s ...string) string {
	return text.Path(text.Snake(text.Lower(text.Space(strings.Join(s, " ")))))
}

// Plural case. Often used in combination with other functions.
// Missing unreserve is intentional, result meant to be used inside others.
func Plural(s ...string) string {
	return text.Plural(strings.Join(s, " "))
}

// Singular case. Often used in combination with other functions.
// Missing unreserve is intentional, result meant to be used inside others.
func Singular(s ...string) string {
	return text.Singular(strings.Join(s, " "))
}

// Slug case. Often used for file names.
// Missing unreserve is intentional, file names don't have reserved names.
func Slug(s ...string) string {
	return text.Slug(text.Lower(text.Space(strings.Join(s, " "))))
}

// acronym case may turn s into a proper acronym
func acronym(s string) string {
	u := strings.ToUpper(s)
	if initialisms[u] {
		return u
	}
	return s
}

// unreserve turns a reserved keyword into a valid identifier
func unreserve(s string) string {
	if v := builtins[s]; v != "" {
		return v
	}
	return s
}
