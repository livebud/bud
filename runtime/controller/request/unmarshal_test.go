package request_test

import (
	"bytes"
	"net/http/httptest"
	"testing"

	"github.com/livebud/bud/internal/is"
	. "github.com/livebud/bud/runtime/controller/request"
)

func TestJSONEmpty(t *testing.T) {
	is := is.New(t)
	type S struct{}
	s := S{}
	r := httptest.NewRequest("GET", "/", nil)
	r.Header.Add("Content-Type", "application/json")
	err := Unmarshal(r, &s)
	is.NoErr(err)
}

func TestFormEmpty(t *testing.T) {
	is := is.New(t)
	type S struct{}
	s := S{}
	r := httptest.NewRequest("GET", "/", nil)
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	err := Unmarshal(r, &s)
	is.NoErr(err)
}

func TestStringJSONEmpties(t *testing.T) {
	is := is.New(t)
	type S struct {
		A string
		B string
	}
	s := S{}
	r := httptest.NewRequest("GET", "/", nil)
	r.Header.Add("Content-Type", "application/json")
	err := Unmarshal(r, &s)
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
	s := S{}
	r := httptest.NewRequest("GET", "/", nil)
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	err := Unmarshal(r, &s)
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
	s := S{}
	r := httptest.NewRequest("GET", "/?a=a&b=b", nil)
	r.Header.Add("Content-Type", "application/json")
	err := Unmarshal(r, &s)
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
	s := S{}
	r := httptest.NewRequest("GET", "/?a=a&b=b", nil)
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	err := Unmarshal(r, &s)
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
	s := S{}
	r := httptest.NewRequest("GET", "/?a=a&b=b", bytes.NewBufferString(`{"c":"c"}`))
	r.Header.Add("Content-Type", "application/json")
	err := Unmarshal(r, &s)
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
	s := S{}
	r := httptest.NewRequest("GET", "/?a=a&b=b", bytes.NewBufferString(`c=c`))
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	err := Unmarshal(r, &s)
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
// 	s := S{}
// 	r := httptest.NewRequest("GET", "/?a=a&b=b", bytes.NewBufferString(`{"c":"c"}`))
// 	r.Header.Add("Content-Type", "application/json")
// 	err := Unmarshal(r, &s)
// 	is.NoErr(err)
// 	is.Equal( "a", s.A)
// 	is.Equal( "b", s.B)
// 	is.Equal( "", s.C)
// 	is.Equal( "", s.D)
// }

// func TestStringFormQueryExtra(t *testing.T) {
// is := is.New(t)
// 	type S struct {
// 		A string
// 		B string
// 		C string
// 		D string
// 	}
// 	s := S{}
// 	r := httptest.NewRequest("GET", "/?a=a&b=b", bytes.NewBufferString(`c=c`))
// 	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
// 	err := Unmarshal(r, &s)
// 	is.NoErr(err)
// 	is.Equal( "a", s.A)
// 	is.Equal( "b", s.B)
// 	is.Equal( "c", s.C)
// 	is.Equal( "", s.D)
// }

func TestStringJSONQueryOverride(t *testing.T) {
	is := is.New(t)
	type S struct {
		A string
		B string
	}
	s := S{}
	r := httptest.NewRequest("GET", "/?a=a&b=b", bytes.NewBufferString(`{"b":"c"}`))
	r.Header.Add("Content-Type", "application/json")
	err := Unmarshal(r, &s)
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
	s := S{}
	r := httptest.NewRequest("GET", "/?a=a&b=b", bytes.NewBufferString(`b=c`))
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	err := Unmarshal(r, &s)
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
	s := S{}
	r := httptest.NewRequest("GET", "/", nil)
	r.Header.Add("Content-Type", "application/json")
	err := Unmarshal(r, &s)
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
	s := S{}
	r := httptest.NewRequest("GET", "/", nil)
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	err := Unmarshal(r, &s)
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
	s := S{}
	r := httptest.NewRequest("GET", "/?a=20&b=10.2", nil)
	r.Header.Add("Content-Type", "application/json")
	err := Unmarshal(r, &s)
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
	s := S{}
	r := httptest.NewRequest("GET", "/?a=1&b=2.2", nil)
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	err := Unmarshal(r, &s)
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
	s := S{}
	r := httptest.NewRequest("GET", "/?a=1&b=2.2", bytes.NewBufferString(`{"c":3,"d":4.4}`))
	r.Header.Add("Content-Type", "application/json")
	err := Unmarshal(r, &s)
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
// 	s := S{}
// 	r := httptest.NewRequest("GET", "/?a=1&b=2.2", bytes.NewBufferString(`c=3&d=4.4`))
// 	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
// 	err := Unmarshal(r, &s)
// 	is.NoErr(err)
// 	is.Equal( 1, s.A)
// 	is.Equal( 2.2, s.B)
// 	is.Equal( 3, s.C)
// 	is.Equal( 4.4, s.D)
// }

// func TestNumberJSONQueryExtra(t *testing.T) {
// is := is.New(t)
// 	type S struct {
// 		A int
// 		B float64
// 		C int
// 		D float64
// 	}
// 	s := S{}
// 	r := httptest.NewRequest("GET", "/?a=1&b=2.2", bytes.NewBufferString(`{"c":3}`))
// 	r.Header.Add("Content-Type", "application/json")
// 	err := Unmarshal(r, &s)
// 	is.NoErr(err)
// 	is.Equal( 1, s.A)
// 	is.Equal( 2.2, s.B)
// 	is.Equal( 0, s.C)
// 	is.Equal( 0.0, s.D)
// }

// func TestNumberFormQueryExtra(t *testing.T) {
// is := is.New(t)
// 	type S struct {
// 		A int
// 		B float64
// 		C int
// 		D float64
// 	}
// 	s := S{}
// 	r := httptest.NewRequest("GET", "/?a=1&b=2.2", bytes.NewBufferString(`c=3`))
// 	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
// 	err := Unmarshal(r, &s)
// 	is.NoErr(err)
// 	is.Equal( 1, s.A)
// 	is.Equal( 2.2, s.B)
// 	is.Equal( 3, s.C)
// 	is.Equal( 0.0, s.D)
// }

func TestNumberJSONQueryOverride(t *testing.T) {
	is := is.New(t)
	type S struct {
		A int
		B float64
	}
	s := S{}
	r := httptest.NewRequest("GET", "/?a=1&b=2.2", bytes.NewBufferString(`{"b":3.3}`))
	r.Header.Add("Content-Type", "application/json")
	err := Unmarshal(r, &s)
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
	s := S{}
	r := httptest.NewRequest("GET", "/?a=1&b=2.2", bytes.NewBufferString(`b=3.3`))
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	err := Unmarshal(r, &s)
	is.NoErr(err)
	is.Equal(1, s.A)
	is.Equal(2.2, s.B)
}

func TestJSONKey(t *testing.T) {
	is := is.New(t)
	type S struct {
		PostID int `json:"post_id"`
	}
	s := S{}
	r := httptest.NewRequest("GET", "/?post_id=1", nil)
	err := Unmarshal(r, &s)
	is.NoErr(err)
	is.Equal(1, s.PostID)
}

func TestFormKey(t *testing.T) {
	is := is.New(t)
	type S struct {
		PostID int `form:"post_id" json:"post_ids"`
	}
	s := S{}
	r := httptest.NewRequest("GET", "/?post_id=1", nil)
	err := Unmarshal(r, &s)
	is.NoErr(err)
	is.Equal(1, s.PostID)
}
