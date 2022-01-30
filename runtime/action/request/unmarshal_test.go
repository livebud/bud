package request_test

import (
	"bytes"
	"net/http/httptest"
	"testing"

	"github.com/tj/assert"
	. "gitlab.com/mnm/bud/runtime/action/request"
)

func TestJSONEmpty(t *testing.T) {
	type S struct{}
	s := S{}
	r := httptest.NewRequest("GET", "/", nil)
	r.Header.Add("Content-Type", "application/json")
	err := Unmarshal(r, &s)
	assert.NoError(t, err)
}

func TestFormEmpty(t *testing.T) {
	type S struct{}
	s := S{}
	r := httptest.NewRequest("GET", "/", nil)
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	err := Unmarshal(r, &s)
	assert.NoError(t, err)
}

func TestStringJSONEmpties(t *testing.T) {
	type S struct {
		A string
		B string
	}
	s := S{}
	r := httptest.NewRequest("GET", "/", nil)
	r.Header.Add("Content-Type", "application/json")
	err := Unmarshal(r, &s)
	assert.NoError(t, err)
	assert.Equal(t, "", s.A)
	assert.Equal(t, "", s.B)
}

func TestStringFormEmpties(t *testing.T) {
	type S struct {
		A string
		B string
	}
	s := S{}
	r := httptest.NewRequest("GET", "/", nil)
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	err := Unmarshal(r, &s)
	assert.NoError(t, err)
	assert.Equal(t, "", s.A)
	assert.Equal(t, "", s.B)
}

func TestStringJSONEmptyQuery(t *testing.T) {
	type S struct {
		A string
		B string
	}
	s := S{}
	r := httptest.NewRequest("GET", "/?a=a&b=b", nil)
	r.Header.Add("Content-Type", "application/json")
	err := Unmarshal(r, &s)
	assert.NoError(t, err)
	assert.Equal(t, "a", s.A)
	assert.Equal(t, "b", s.B)
}

func TestStringFormEmptyQuery(t *testing.T) {
	type S struct {
		A string
		B string
	}
	s := S{}
	r := httptest.NewRequest("GET", "/?a=a&b=b", nil)
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	err := Unmarshal(r, &s)
	assert.NoError(t, err)
	assert.Equal(t, "a", s.A)
	assert.Equal(t, "b", s.B)
}

func TestStringJSONQuery(t *testing.T) {
	type S struct {
		A string
		B string
		C string
	}
	s := S{}
	r := httptest.NewRequest("GET", "/?a=a&b=b", bytes.NewBufferString(`{"c":"c"}`))
	r.Header.Add("Content-Type", "application/json")
	err := Unmarshal(r, &s)
	assert.NoError(t, err)
	assert.Equal(t, "a", s.A)
	assert.Equal(t, "b", s.B)
	assert.Equal(t, "c", s.C)
}

func TestStringFormQuery(t *testing.T) {
	type S struct {
		A string
		B string
		C string
	}
	s := S{}
	r := httptest.NewRequest("GET", "/?a=a&b=b", bytes.NewBufferString(`c=c`))
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	err := Unmarshal(r, &s)
	assert.NoError(t, err)
	assert.Equal(t, "a", s.A)
	assert.Equal(t, "b", s.B)
	assert.Equal(t, "", s.C)
}

// func TestStringJSONQueryExtra(t *testing.T) {
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
// 	assert.NoError(t, err)
// 	assert.Equal(t, "a", s.A)
// 	assert.Equal(t, "b", s.B)
// 	assert.Equal(t, "", s.C)
// 	assert.Equal(t, "", s.D)
// }

// func TestStringFormQueryExtra(t *testing.T) {
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
// 	assert.NoError(t, err)
// 	assert.Equal(t, "a", s.A)
// 	assert.Equal(t, "b", s.B)
// 	assert.Equal(t, "c", s.C)
// 	assert.Equal(t, "", s.D)
// }

func TestStringJSONQueryOverride(t *testing.T) {
	type S struct {
		A string
		B string
	}
	s := S{}
	r := httptest.NewRequest("GET", "/?a=a&b=b", bytes.NewBufferString(`{"b":"c"}`))
	r.Header.Add("Content-Type", "application/json")
	err := Unmarshal(r, &s)
	assert.NoError(t, err)
	assert.Equal(t, "a", s.A)
	assert.Equal(t, "b", s.B)
}

func TestStringFormQueryOverride(t *testing.T) {
	type S struct {
		A string
		B string
	}
	s := S{}
	r := httptest.NewRequest("GET", "/?a=a&b=b", bytes.NewBufferString(`b=c`))
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	err := Unmarshal(r, &s)
	assert.NoError(t, err)
	assert.Equal(t, "a", s.A)
	assert.Equal(t, "b", s.B)
}

func TestNumberJSONEmpties(t *testing.T) {
	type S struct {
		A int
		B float64
	}
	s := S{}
	r := httptest.NewRequest("GET", "/", nil)
	r.Header.Add("Content-Type", "application/json")
	err := Unmarshal(r, &s)
	assert.NoError(t, err)
	assert.Equal(t, 0, s.A)
	assert.Equal(t, 0.0, s.B)
}

func TestNumberFormEmpties(t *testing.T) {
	type S struct {
		A int
		B float64
	}
	s := S{}
	r := httptest.NewRequest("GET", "/", nil)
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	err := Unmarshal(r, &s)
	assert.NoError(t, err)
	assert.Equal(t, 0, s.A)
	assert.Equal(t, 0.0, s.B)
}

func TestNumberJSONEmptyQuery(t *testing.T) {
	type S struct {
		A int
		B float64
	}
	s := S{}
	r := httptest.NewRequest("GET", "/?a=20&b=10.2", nil)
	r.Header.Add("Content-Type", "application/json")
	err := Unmarshal(r, &s)
	assert.NoError(t, err)
	assert.Equal(t, 20, s.A)
	assert.Equal(t, 10.2, s.B)
}

func TestNumberFormEmptyQuery(t *testing.T) {
	type S struct {
		A int
		B float64
	}
	s := S{}
	r := httptest.NewRequest("GET", "/?a=1&b=2.2", nil)
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	err := Unmarshal(r, &s)
	assert.NoError(t, err)
	assert.Equal(t, 1, s.A)
	assert.Equal(t, 2.2, s.B)
}

func TestNumberJSONQuery(t *testing.T) {
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
	assert.NoError(t, err)
	assert.Equal(t, 1, s.A)
	assert.Equal(t, 2.2, s.B)
	assert.Equal(t, 3, s.C)
	assert.Equal(t, 4.4, s.D)
}

// func TestNumberFormQuery(t *testing.T) {
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
// 	assert.NoError(t, err)
// 	assert.Equal(t, 1, s.A)
// 	assert.Equal(t, 2.2, s.B)
// 	assert.Equal(t, 3, s.C)
// 	assert.Equal(t, 4.4, s.D)
// }

// func TestNumberJSONQueryExtra(t *testing.T) {
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
// 	assert.NoError(t, err)
// 	assert.Equal(t, 1, s.A)
// 	assert.Equal(t, 2.2, s.B)
// 	assert.Equal(t, 0, s.C)
// 	assert.Equal(t, 0.0, s.D)
// }

// func TestNumberFormQueryExtra(t *testing.T) {
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
// 	assert.NoError(t, err)
// 	assert.Equal(t, 1, s.A)
// 	assert.Equal(t, 2.2, s.B)
// 	assert.Equal(t, 3, s.C)
// 	assert.Equal(t, 0.0, s.D)
// }

func TestNumberJSONQueryOverride(t *testing.T) {
	type S struct {
		A int
		B float64
	}
	s := S{}
	r := httptest.NewRequest("GET", "/?a=1&b=2.2", bytes.NewBufferString(`{"b":3.3}`))
	r.Header.Add("Content-Type", "application/json")
	err := Unmarshal(r, &s)
	assert.NoError(t, err)
	assert.Equal(t, 1, s.A)
	assert.Equal(t, 2.2, s.B)
}

func TestNumberFormQueryOverride(t *testing.T) {
	type S struct {
		A int
		B float64
	}
	s := S{}
	r := httptest.NewRequest("GET", "/?a=1&b=2.2", bytes.NewBufferString(`b=3.3`))
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	err := Unmarshal(r, &s)
	assert.NoError(t, err)
	assert.Equal(t, 1, s.A)
	assert.Equal(t, 2.2, s.B)
}

func TestJSONKey(t *testing.T) {
	type S struct {
		PostID int `json:"post_id"`
	}
	s := S{}
	r := httptest.NewRequest("GET", "/?post_id=1", nil)
	err := Unmarshal(r, &s)
	assert.NoError(t, err)
	assert.Equal(t, 1, s.PostID)
}

func TestFormKey(t *testing.T) {
	type S struct {
		PostID int `form:"post_id" json:"post_ids"`
	}
	s := S{}
	r := httptest.NewRequest("GET", "/?post_id=1", nil)
	err := Unmarshal(r, &s)
	assert.NoError(t, err)
	assert.Equal(t, 1, s.PostID)
}
