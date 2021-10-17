package gotext_test

import (
	"testing"

	"gitlab.com/mnm/bud/internal/gotext"
	"github.com/matryer/is"
)

type fntest struct {
	fn  func(...string) string
	in  string
	out string
}

func TestSlim(t *testing.T) {
	is := is.New(t)
	tests := []fntest{
		{gotext.Slim, "", ""},
		{gotext.Slim, "hi world", "hiworld"},
		{gotext.Slim, "user id", "userid"},
		{gotext.Slim, "id user", "iduser"},
		{gotext.Slim, "id dns http", "iddnshttp"},
		{gotext.Slim, "http user", "httpuser"},
		{gotext.Slim, "string", "string_"},
		{gotext.Slim, "newIn", "newin"},
		{gotext.Slim, "new", "new_"},
	}
	for _, test := range tests {
		is.Equal(test.out, test.fn(test.in))
	}
}

func TestPascal(t *testing.T) {
	is := is.New(t)
	tests := []fntest{
		{gotext.Pascal, "", ""},
		{gotext.Pascal, "hi world", "HiWorld"},
		{gotext.Pascal, "user id", "UserID"},
		{gotext.Pascal, "id user", "IDUser"},
		{gotext.Pascal, "id dns http", "IDDNSHTTP"},
		{gotext.Pascal, "http user", "HTTPUser"},
		{gotext.Pascal, "string", "String"},
		{gotext.Pascal, "newIn", "NewIn"},
		{gotext.Pascal, "new", "New"},
	}
	for _, test := range tests {
		is.Equal(test.out, test.fn(test.in))
	}
}

func TestCamel(t *testing.T) {
	is := is.New(t)
	tests := []fntest{
		{gotext.Camel, "", ""},
		{gotext.Camel, "hi world", "hiWorld"},
		{gotext.Camel, "user id", "userID"},
		{gotext.Camel, "id user", "idUser"},
		{gotext.Camel, "id dns http", "idDNSHTTP"},
		{gotext.Camel, "http user", "httpUser"},
		{gotext.Camel, "string", "string_"},
		{gotext.Camel, "newIn", "newIn"},
		{gotext.Camel, "new", "new_"},
	}
	for _, test := range tests {
		is.Equal(test.out, test.fn(test.in))
	}
}

func TestShort(t *testing.T) {
	is := is.New(t)
	tests := []fntest{
		{gotext.Short, "", ""},
		{gotext.Short, "hi world", "hw"},
		{gotext.Short, "user id", "ui"},
		{gotext.Short, "id user", "iu"},
		{gotext.Short, "id dns http", "idh"},
		{gotext.Short, "http user", "hu"},
		{gotext.Short, "string", "s"},
		{gotext.Short, "newIn", "ni"},
		{gotext.Short, "new", "n"},
		{gotext.Short, "newEwW", "new_"},
	}
	for _, test := range tests {
		is.Equal(test.out, test.fn(test.in))
	}
}

func TestSnake(t *testing.T) {
	is := is.New(t)
	tests := []fntest{
		{gotext.Snake, "", ""},
		{gotext.Snake, "hi world", "hi_world"},
		{gotext.Snake, "user id", "user_id"},
		{gotext.Snake, "id user", "id_user"},
		{gotext.Snake, "id dns http", "id_dns_http"},
		{gotext.Snake, "http user", "http_user"},
		{gotext.Snake, "string", "string"},
		{gotext.Snake, "newIn", "new_in"},
		{gotext.Snake, "new", "new"},
	}
	for _, test := range tests {
		is.Equal(test.out, test.fn(test.in))
	}
}

func TestUpper(t *testing.T) {
	is := is.New(t)
	tests := []fntest{
		{gotext.Upper, "", ""},
		{gotext.Upper, "hi world", "HI_WORLD"},
		{gotext.Upper, "user id", "USER_ID"},
		{gotext.Upper, "id user", "ID_USER"},
		{gotext.Upper, "id dns http", "ID_DNS_HTTP"},
		{gotext.Upper, "http user", "HTTP_USER"},
		{gotext.Upper, "string", "STRING"},
		{gotext.Upper, "newIn", "NEW_IN"},
		{gotext.Upper, "new", "NEW"},
	}
	for _, test := range tests {
		is.Equal(test.out, test.fn(test.in))
	}
}

func TestPath(t *testing.T) {
	is := is.New(t)
	tests := []fntest{
		{gotext.Path, "", ""},
		{gotext.Path, "hi world", "hi/world"},
		{gotext.Path, "user id", "user/id"},
		{gotext.Path, "id user", "id/user"},
		{gotext.Path, "id dns http", "id/dns/http"},
		{gotext.Path, "http user", "http/user"},
		{gotext.Path, "string", "string"},
		{gotext.Path, "newIn", "new/in"},
		{gotext.Path, "new", "new"},
	}
	for _, test := range tests {
		is.Equal(test.out, test.fn(test.in))
	}
}

func TestPlural(t *testing.T) {
	is := is.New(t)
	tests := []fntest{
		{gotext.Plural, "", ""},
		{gotext.Plural, "hi world", "hi worlds"},
		{gotext.Plural, "user id", "user ids"},
		{gotext.Plural, "id user", "id users"},
		{gotext.Plural, "id dns http", "id dns https"},
		{gotext.Plural, "http user", "http users"},
		{gotext.Plural, "string", "strings"},
		{gotext.Plural, "newIn", "newIns"},
		{gotext.Plural, "news", "news"},
		{gotext.Plural, "new", "news"},
		{gotext.Plural, "hi worlds", "hi worlds"},
		{gotext.Plural, "Hi World", "Hi Worlds"},
		{gotext.Plural, "Hi$World$$$", "Hi$Worlds$$$"},
		{gotext.Plural, "hi_world", "hi_worlds"},
		{gotext.Plural, "his_world", "his_worlds"},
		{gotext.Plural, "his", "his"},
		{gotext.Plural, "hi", "his"},
	}
	for _, test := range tests {
		is.Equal(test.out, test.fn(test.in))
	}
}

func TestSingular(t *testing.T) {
	is := is.New(t)
	tests := []fntest{
		{gotext.Singular, "", ""},
		{gotext.Singular, "hi worlds", "hi world"},
		{gotext.Singular, "user ids", "user id"},
		{gotext.Singular, "users ids", "users id"},
		{gotext.Singular, "id users", "id user"},
		{gotext.Singular, "dns", "dn"},
		{gotext.Singular, "id dns https", "id dns http"},
		{gotext.Singular, "id dns http", "id dns http"},
		{gotext.Singular, "http users", "http user"},
		{gotext.Singular, "strings", "string"},
		{gotext.Singular, "newIns", "newIn"},
		{gotext.Singular, "news", "news"},
		{gotext.Singular, "new", "new"},
		{gotext.Singular, "hi world", "hi world"},
		{gotext.Singular, "Hi Worlds", "Hi World"},
		{gotext.Singular, "Hi$Worlds$$$", "Hi$World$$$"},
		{gotext.Singular, "hi_worlds", "hi_world"},
		{gotext.Singular, "his_worlds", "his_world"},
		{gotext.Singular, "his", "hi"},
	}
	for _, test := range tests {
		is.Equal(test.out, test.fn(test.in))
	}
}
