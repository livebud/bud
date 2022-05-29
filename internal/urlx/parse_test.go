package urlx_test

import (
	"testing"

	"github.com/livebud/bud/internal/is"
	urlx "github.com/livebud/bud/internal/urlx"
)

func Test(t *testing.T) {
	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			is := is.New(t)
			u, err := urlx.Parse(test.input)
			if err != nil {
				is.Equal(err.Error(), test.expect)
				return
			}
			is.Equal(u.String(), test.expect)
		})
	}
}

var tests = []struct {
	input  string
	expect string
}{
	{
		input:  "5000",
		expect: "http://127.0.0.1:5000",
	},
	{
		input:  ":5000",
		expect: "http://127.0.0.1:5000",
	},
	{
		input:  "0",
		expect: "http://127.0.0.1:0",
	},
	{
		input:  "0.0.0.0",
		expect: "http://0.0.0.0:3000",
	},
	{
		input:  "127.0.0.1",
		expect: "http://127.0.0.1:3000",
	},
	{
		input:  "127.0.0.1:5000",
		expect: "http://127.0.0.1:5000",
	},
	{
		input:  "localhost",
		expect: "http://localhost:3000",
	},
	{
		input:  "otherhost",
		expect: "http://otherhost:3000",
	},
	{
		input:  "/tmp.sock",
		expect: "http:///tmp.sock",
	},
	{
		input:  "/whatever/tmp.sock",
		expect: "http:///whatever/tmp.sock",
	},
	{
		input:  "https:",
		expect: "https://127.0.0.1:443",
	},
	{
		input:  "https://localhost:8000/a/b/c",
		expect: "https://localhost:8000/a/b/c",
	},
	{
		input:  "80.ab",
		expect: `urlx: unable to parse "80.ab"`,
	},
	{
		input:  "http://127.0.0.1:49341",
		expect: "http://127.0.0.1:49341",
	},
	{
		input:  "[::]:50516",
		expect: "http://[::]:50516",
	},
	{
		input:  "[::]:443",
		expect: "https://[::]:443",
	},
}
