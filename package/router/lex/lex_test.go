package lex_test

import (
	"testing"

	"github.com/livebud/bud/internal/is"
	"github.com/livebud/bud/package/router/lex"
)

func Test(t *testing.T) {
	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			is := is.New(t)
			l := lex.New(test.input)
			var tokens lex.Tokens
			for {
				token := l.Next()
				switch token.Type {
				case lex.ErrorToken:
					is.Equal(test.err, token.Value)
					return
				case lex.EndToken:
					is.Equal(test.expect, tokens.String())
					return
				default:
					tokens = append(tokens, token)
				}
			}
		})
	}
}

var tests = []struct {
	input  string
	expect string
	err    string
}{
	{input: "a", err: `route "a": must start with a slash "/"`},
	{input: "/a/", err: `route "/a/": remove the slash "/" at the end`},
	{input: "/", expect: `slash:"/"`},
	{input: "/hi", expect: `slash:"/" path:"hi"`},
	{input: "/doc/", err: `route "/doc/": remove the slash "/" at the end`},
	{input: "/doc/go_faq.html", expect: `slash:"/" path:"doc" slash:"/" path:"go_faq.html"`},
	{input: "α", err: `route "α": must start with a slash "/"`},
	{input: "/α", expect: `slash:"/" path:"α"`},
	{input: "/β", expect: `slash:"/" path:"β"`},
	{input: "/ ", err: `route "/ ": invalid character " "`},
	{input: "/café", expect: `slash:"/" path:"café"`},
	{input: "/:", err: `route "/:": missing slot name after ":"`},
	{input: "/:a", expect: `slash:"/" slot:":a"`},
	{input: "/:hi", expect: `slash:"/" slot:":hi"`},
	{input: "/:hi?", expect: `slash:"/" question:":hi?"`},
	{input: "/:hi*", expect: `slash:"/" star:":hi*"`},
	{input: "/:hi*?", err: `route "/:hi*?": "?" not allowed after "*"`},
	{input: "/:hi?*", err: `route "/:hi?*": "*" not allowed after "?"`},
	{input: "/a?*", err: `route "/a?*": unexpected modifier "?"`},
	{input: "/:a/:b", expect: `slash:"/" slot:":a" slash:"/" slot:":b"`},
	{input: "/:a/:b?", expect: `slash:"/" slot:":a" slash:"/" question:":b?"`},
	{input: "/:a/:b*", expect: `slash:"/" slot:":a" slash:"/" star:":b*"`},
	{input: "/users/:id.:format?", expect: `slash:"/" path:"users" slash:"/" slot:":id" path:"." question:":format?"`},
	{input: "/users/:major.:minor", expect: `slash:"/" path:"users" slash:"/" slot:":major" path:"." slot:":minor"`},
	{input: "/users/:major.:minor", expect: `slash:"/" path:"users" slash:"/" slot:":major" path:"." slot:":minor"`},
	{input: "/:-a/:b?", err: `route "/:-a/:b?": first letter after ":" must be a lowercase Latin letter`},
	{input: "/:a-b", err: `route "/:a-b": invalid slot character "-"`},
	{input: ":a", err: `route ":a": must start with a slash "/"`},
	{input: "/:a/", err: `route "/:a/": remove the slash "/" at the end`},
	{input: "/:café", err: `route "/:café": invalid slot character "é"`},
	// From the guide
	{input: "/about", expect: `slash:"/" path:"about"`},
	{input: "/deactivate", expect: `slash:"/" path:"deactivate"`},
	{input: "/archive/:year/:month", expect: `slash:"/" path:"archive" slash:"/" slot:":year" slash:"/" slot:":month"`},
	{input: "/users", expect: `slash:"/" path:"users"`},
	{input: "/users/:id", expect: `slash:"/" path:"users" slash:"/" slot:":id"`},
	{input: "/:id", expect: `slash:"/" slot:":id"`},
	{input: "/v.:version", expect: `slash:"/" path:"v." slot:":version"`},
	{input: "/:post_id.:format", expect: `slash:"/" slot:":post_id" path:"." slot:":format"`},
	{input: "/:from-:to", err: `route "/:from-:to": invalid slot character "-"`},
	{input: "/:key1/:key2", expect: `slash:"/" slot:":key1" slash:"/" slot:":key2"`},
	{input: "/:id?", expect: `slash:"/" question:":id?"`},
	{input: "/v.:version?", expect: `slash:"/" path:"v." question:":version?"`},
	{input: "/:post_id.:format?", expect: `slash:"/" slot:":post_id" path:"." question:":format?"`},
	{input: "/:from-:to?", err: `route "/:from-:to?": invalid slot character "-"`},
	{input: "/:from/:to?", expect: `slash:"/" slot:":from" slash:"/" question:":to?"`},
	{input: "/:id/:path*", expect: `slash:"/" slot:":id" slash:"/" star:":path*"`},
	{input: "/v.:version*", expect: `slash:"/" path:"v." star:":version*"`},
	{input: "/explore", expect: `slash:"/" path:"explore"`},
	// Must be lowercase
	{input: "/Explore", err: `route "/Explore": uppercase letters are not allowed "E"`},
	{input: "/eXPLORE", err: `route "/eXPLORE": uppercase letters are not allowed "X"`},
	{input: "/explorE", err: `route "/explorE": uppercase letters are not allowed "E"`},
	{input: "/explorE/", err: `route "/explorE/": uppercase letters are not allowed "E"`},
	{input: "/:Slot", err: `route "/:Slot": first letter after ":" must be a lowercase Latin letter`},
	{input: "/:sLot", err: `route "/:sLot": uppercase letters are not allowed "L"`},
	{input: "/:sloT", err: `route "/:sloT": uppercase letters are not allowed "T"`},
	{input: "/:sloT/", err: `route "/:sloT/": uppercase letters are not allowed "T"`},
}
