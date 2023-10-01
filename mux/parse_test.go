package mux_test

import (
	"testing"

	"github.com/livebud/bud/mux/internal/parser"
	"github.com/matthewmueller/diff"
)

func parseEqual(t *testing.T, input, expected string) {
	t.Helper()
	t.Run(input, func(t *testing.T) {
		t.Helper()
		route, err := parser.Parse(input)
		if err != nil {
			if err.Error() == expected {
				return
			}
			t.Fatal(err)
		}
		actual := route.String()
		diff.TestString(t, expected, actual)
	})
}

func TestParse(t *testing.T) {
	parseEqual(t, `/{name}`, `/{name}`)
	parseEqual(t, `/{na me}`, `invalid character ' ' in slot`)
	parseEqual(t, `/hello/{name}`, `/hello/{name}`)
	parseEqual(t, `hello/{name}`, `path must start with a slash /`)
	parseEqual(t, `/hello/{name}/`, `/hello/{name}/`)
	parseEqual(t, `/hello/{name?}`, `/hello/{name?}`)
	parseEqual(t, `/hello/{name*}`, `/hello/{name*}`)
	parseEqual(t, `/hel_lo/`, `/hel_lo/`)
	parseEqual(t, `/hel lo/`, `unexpected character ' ' in path`)
	parseEqual(t, `/hello/{*name}`, `slot can't start with '*'`)
	parseEqual(t, `/hello/{na*me}`, `expected '}' but got 'me'`)
	parseEqual(t, `/hello/{name}/admin`, `/hello/{name}/admin`)
	parseEqual(t, `/hello/{name}/admin/`, `/hello/{name}/admin/`)
	parseEqual(t, `/hello/{name}/{owner}`, `/hello/{name}/{owner}`)
	parseEqual(t, `/hello/{name`, `unclosed slot`)
	parseEqual(t, `/hello/{name|^[A-Za-z]$}`, `/hello/{name|^[A-Za-z]$}`)
	parseEqual(t, `/hello/{name|^[A-Za-z]{1,3}$}`, `/hello/{name|^[A-Za-z]{1,3}$}`)
}
