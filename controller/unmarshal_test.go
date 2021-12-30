package controller_test

import (
	"bytes"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/matryer/is"
	"gitlab.com/mnm/bud/controller"
)

func TestJSONEmpty(t *testing.T) {
	is := is.New(t)
	type S struct{}
	var s S
	r := httptest.NewRequest("GET", "/", nil)
	r.Header.Add("Content-Type", "application/json")
	err := controller.Unmarshal(r, &s)
	is.NoErr(err)
}

func TestFormEmpty(t *testing.T) {
	is := is.New(t)
	type S struct{}
	var s S
	r := httptest.NewRequest("GET", "/", nil)
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	err := controller.Unmarshal(r, &s)
	is.NoErr(err)
}

func TestStringJSONEmpties(t *testing.T) {
	is := is.New(t)
	type S struct {
		A string
		B string
	}
	var s S
	r := httptest.NewRequest("GET", "/", nil)
	r.Header.Add("Content-Type", "application/json")
	err := controller.Unmarshal(r, &s)
	is.NoErr(err)
	is.Equal("", s.A)
	is.Equal("", s.B)
}

func TestStringFormEmpties(t *testing.T) {
	is := is.New(t)
	type S struct {
		A string
		B string
	}
	var s S
	r := httptest.NewRequest("GET", "/", nil)
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	err := controller.Unmarshal(r, &s)
	is.NoErr(err)
	is.Equal("", s.A)
	is.Equal("", s.B)
}

func TestStringJSONEmptyQuery(t *testing.T) {
	is := is.New(t)
	type S struct {
		A string
		B string
	}
	var s S
	r := httptest.NewRequest("GET", "/?a=a&b=b", nil)
	r.Header.Add("Content-Type", "application/json")
	err := controller.Unmarshal(r, &s)
	is.NoErr(err)
	is.Equal("a", s.A)
	is.Equal("b", s.B)
}

func TestStringFormEmptyQuery(t *testing.T) {
	is := is.New(t)
	type S struct {
		A string
		B string
	}
	var s S
	r := httptest.NewRequest("GET", "/?a=a&b=b", nil)
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	err := controller.Unmarshal(r, &s)
	is.NoErr(err)
	is.Equal("a", s.A)
	is.Equal("b", s.B)
}

func TestStringJSONQuery(t *testing.T) {
	is := is.New(t)
	type S struct {
		A string
		B string
		C string
	}
	var s S
	r := httptest.NewRequest("GET", "/?a=a&b=b", bytes.NewBufferString(`{"c":"c"}`))
	r.Header.Add("Content-Type", "application/json")
	err := controller.Unmarshal(r, &s)
	is.NoErr(err)
	is.Equal("a", s.A)
	is.Equal("b", s.B)
	is.Equal("c", s.C)
}

func TestStringFormQuery(t *testing.T) {
	is := is.New(t)
	type S struct {
		A string
		B string
		C string
	}
	var s S
	r := httptest.NewRequest("GET", "/?a=a&b=b", bytes.NewBufferString(`c=c`))
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	err := controller.Unmarshal(r, &s)
	is.NoErr(err)
	is.Equal("a", s.A)
	is.Equal("b", s.B)
	is.Equal("", s.C)
}

// func TestStringJSONQueryExtra(t *testing.T) {
// is := is.New(t)
// 	type S struct {
// 		A string
// 		B string
// 		C string
// 		D string
// 	}
// 	var s S
// 	r := httptest.NewRequest("GET", "/?a=a&b=b", bytes.NewBufferString(`{"c":"c"}`))
// 	r.Header.Add("Content-Type", "application/json")
// 	err := controller.Unmarshal(r, &s)
// 	is.NoErr(err)
// 	is.Equal("a", s.A)
// 	is.Equal("b", s.B)
// 	is.Equal("", s.C)
// 	is.Equal("", s.D)
// }

// func TestStringFormQueryExtra(t *testing.T) {
// is := is.New(t)
// 	type S struct {
// 		A string
// 		B string
// 		C string
// 		D string
// 	}
// 	var s S
// 	r := httptest.NewRequest("GET", "/?a=a&b=b", bytes.NewBufferString(`c=c`))
// 	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
// 	err := controller.Unmarshal(r, &s)
// 	is.NoErr(err)
// 	is.Equal("a", s.A)
// 	is.Equal("b", s.B)
// 	is.Equal("c", s.C)
// 	is.Equal("", s.D)
// }

func TestStringJSONQueryOverride(t *testing.T) {
	is := is.New(t)
	type S struct {
		A string
		B string
	}
	var s S
	r := httptest.NewRequest("GET", "/?a=a&b=b", bytes.NewBufferString(`{"b":"c"}`))
	r.Header.Add("Content-Type", "application/json")
	err := controller.Unmarshal(r, &s)
	is.NoErr(err)
	is.Equal("a", s.A)
	is.Equal("b", s.B)
}

func TestStringFormQueryOverride(t *testing.T) {
	is := is.New(t)
	type S struct {
		A string
		B string
	}
	var s S
	r := httptest.NewRequest("GET", "/?a=a&b=b", bytes.NewBufferString(`b=c`))
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	err := controller.Unmarshal(r, &s)
	is.NoErr(err)
	is.Equal("a", s.A)
	is.Equal("b", s.B)
}

func TestNumberJSONEmpties(t *testing.T) {
	is := is.New(t)
	type S struct {
		A int
		B float64
	}
	var s S
	r := httptest.NewRequest("GET", "/", nil)
	r.Header.Add("Content-Type", "application/json")
	err := controller.Unmarshal(r, &s)
	is.NoErr(err)
	is.Equal(0, s.A)
	is.Equal(0.0, s.B)
}

func TestNumberFormEmpties(t *testing.T) {
	is := is.New(t)
	type S struct {
		A int
		B float64
	}
	var s S
	r := httptest.NewRequest("GET", "/", nil)
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	err := controller.Unmarshal(r, &s)
	is.NoErr(err)
	is.Equal(0, s.A)
	is.Equal(0.0, s.B)
}

func TestNumberJSONEmptyQuery(t *testing.T) {
	is := is.New(t)
	type S struct {
		A int
		B float64
	}
	var s S
	r := httptest.NewRequest("GET", "/?a=20&b=10.2", nil)
	r.Header.Add("Content-Type", "application/json")
	err := controller.Unmarshal(r, &s)
	is.NoErr(err)
	is.Equal(20, s.A)
	is.Equal(10.2, s.B)
}

func TestNumberFormEmptyQuery(t *testing.T) {
	is := is.New(t)
	type S struct {
		A int
		B float64
	}
	var s S
	r := httptest.NewRequest("GET", "/?a=1&b=2.2", nil)
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	err := controller.Unmarshal(r, &s)
	is.NoErr(err)
	is.Equal(1, s.A)
	is.Equal(2.2, s.B)
}

func TestNumberJSONQuery(t *testing.T) {
	is := is.New(t)
	type S struct {
		A int
		B float64
		C int
		D float64
	}
	var s S
	r := httptest.NewRequest("GET", "/?a=1&b=2.2", bytes.NewBufferString(`{"c":3,"d":4.4}`))
	r.Header.Add("Content-Type", "application/json")
	err := controller.Unmarshal(r, &s)
	is.NoErr(err)
	is.Equal(1, s.A)
	is.Equal(2.2, s.B)
	is.Equal(3, s.C)
	is.Equal(4.4, s.D)
}

// func TestNumberFormQuery(t *testing.T) {
// is := is.New(t)
// 	type S struct {
// 		A int
// 		B float64
// 		C int
// 		D float64
// 	}
// 	var s S
// 	r := httptest.NewRequest("GET", "/?a=1&b=2.2", bytes.NewBufferString(`c=3&d=4.4`))
// 	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
// 	err := controller.Unmarshal(r, &s)
// 	is.NoErr(err)
// 	is.Equal(1, s.A)
// 	is.Equal(2.2, s.B)
// 	is.Equal(3, s.C)
// 	is.Equal(4.4, s.D)
// }

// func TestNumberJSONQueryExtra(t *testing.T) {
// is := is.New(t)
// 	type S struct {
// 		A int
// 		B float64
// 		C int
// 		D float64
// 	}
// 	var s S
// 	r := httptest.NewRequest("GET", "/?a=1&b=2.2", bytes.NewBufferString(`{"c":3}`))
// 	r.Header.Add("Content-Type", "application/json")
// 	err := controller.Unmarshal(r, &s)
// 	is.NoErr(err)
// 	is.Equal(1, s.A)
// 	is.Equal(2.2, s.B)
// 	is.Equal(0, s.C)
// 	is.Equal(0.0, s.D)
// }

// func TestNumberFormQueryExtra(t *testing.T) {
// is := is.New(t)
// 	type S struct {
// 		A int
// 		B float64
// 		C int
// 		D float64
// 	}
// 	var s S
// 	r := httptest.NewRequest("GET", "/?a=1&b=2.2", bytes.NewBufferString(`c=3`))
// 	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
// 	err := controller.Unmarshal(r, &s)
// 	is.NoErr(err)
// 	is.Equal(1, s.A)
// 	is.Equal(2.2, s.B)
// 	is.Equal(3, s.C)
// 	is.Equal(0.0, s.D)
// }

func TestNumberJSONQueryOverride(t *testing.T) {
	is := is.New(t)
	type S struct {
		A int
		B float64
	}
	var s S
	r := httptest.NewRequest("GET", "/?a=1&b=2.2", bytes.NewBufferString(`{"b":3.3}`))
	r.Header.Add("Content-Type", "application/json")
	err := controller.Unmarshal(r, &s)
	is.NoErr(err)
	is.Equal(1, s.A)
	is.Equal(2.2, s.B)
}

func TestNumberFormQueryOverride(t *testing.T) {
	is := is.New(t)
	type S struct {
		A int
		B float64
	}
	var s S
	r := httptest.NewRequest("GET", "/?a=1&b=2.2", bytes.NewBufferString(`b=3.3`))
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	err := controller.Unmarshal(r, &s)
	is.NoErr(err)
	is.Equal(1, s.A)
	is.Equal(2.2, s.B)
}

func TestRequired(t *testing.T) {
	is := is.New(t)
	type S struct {
		A int `validate:"required"`
		B int
	}
	var s S
	r := httptest.NewRequest("GET", "/?b=0", nil)
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	err := controller.Unmarshal(r, &s)
	is.True(err != nil)
	is.True(strings.Contains(err.Error(), "Field validation for 'A' failed on the 'required' tag"))
}
