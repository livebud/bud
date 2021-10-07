package router_test

import (
	"testing"

	"github.com/go-duo/bud/router"
	"github.com/matryer/is"
	"github.com/matthewmueller/diff"
)

func TestRoute(t *testing.T) {
	for _, test := range tests {
		t.Run(test.route, func(t *testing.T) {
			is := is.New(t)
			re, err := router.Parse(test.route)
			if err != nil {
				if test.err != "" {
					diff.TestString(t, test.err, err.Error())
					return
				}
				is.NoErr(err)
			}
			match := re.FindStringSubmatch(test.path)
			is.Equal(test.match, match != nil)
			results := map[string]string{}
			if match != nil {
				for i, name := range re.SubexpNames() {
					if i != 0 && name != "" {
						results[name] = match[i]
					}
				}
			}
			for key, param := range test.expect {
				diff.TestString(t, param, results[key])
			}
		})
	}
}

var tests = []struct {
	route  string
	path   string
	match  bool
	expect map[string]string
	err    string
}{
	{
		route: "/",
		path:  "",
		match: false,
	},
	{
		route: "/",
		path:  "/",
		match: true,
	},
	{
		route: "/",
		path:  "/hi",
		match: false,
	},
	{
		route: "/hi",
		path:  "/hi",
		match: true,
	},
	{
		route: "/:hi",
		path:  "/hi",
		match: true,
		expect: map[string]string{
			"hi": "hi",
		},
	},
	{
		route:  "/:hi?",
		path:   "/",
		match:  true,
		expect: map[string]string{},
	},
	{
		route:  "/:hi?",
		path:   "",
		match:  false,
		expect: map[string]string{},
	},
	{
		route:  "/:hi?",
		path:   "",
		match:  false,
		expect: map[string]string{},
	},
	{
		route:  "/:a/:b",
		path:   "/a",
		match:  false,
		expect: map[string]string{},
	},
	{
		route: "/:a/:b",
		path:  "/a1/b1",
		match: true,
		expect: map[string]string{
			"a": "a1",
			"b": "b1",
		},
	},
	{
		route: "/:a/:b?",
		path:  "/a1/b1",
		match: true,
		expect: map[string]string{
			"a": "a1",
			"b": "b1",
		},
	},
	{
		route: "/:a/:b?",
		path:  "/a1",
		match: true,
		expect: map[string]string{
			"a": "a1",
		},
	},
	{
		route: "/:a/:b?",
		path:  "/10",
		match: true,
		expect: map[string]string{
			"a": "10",
		},
	},
	{
		route: "/:a/:b?",
		path:  "/10/20",
		match: true,
		expect: map[string]string{
			"a": "10",
			"b": "20",
		},
	},
	{
		route: "/users/:id.:format?",
		path:  "/users/10.json",
		match: true,
		expect: map[string]string{
			"id":     "10",
			"format": "json",
		},
	},
	{
		route: "/users/:id.:format?",
		path:  "/users/10",
		match: true,
		expect: map[string]string{
			"id":     "10",
			"format": "",
		},
	},
	{
		route: "/users/:major.:minor",
		path:  "/users/10.1",
		match: true,
		expect: map[string]string{
			"major": "10",
			"minor": "1",
		},
	},
	{
		route:  "/\\:a",
		path:   "/:a",
		match:  true,
		expect: map[string]string{},
	},
	{
		route: "/:-a/:b?",
		path:  "/a1",
		err:   "\nparse error near Slash (line 1 symbol 1 - line 1 symbol 2):\n\"/\"\n",
	},
	{
		route: ":a",
		path:  "/a1",
		err:   "\nparse error near Unknown (line 1 symbol 1 - line 1 symbol 1):\n\"\"\n",
	},
	{
		route: "/:a",
		path:  "/unique_id",
		match: true,
		expect: map[string]string{
			"a": "unique_id",
		},
	},
}
