package gotemplate_test

import (
	"testing"

	"github.com/livebud/bud/internal/gotemplate"
	"github.com/livebud/bud/internal/is"
)

func TestGenerateGoFile(t *testing.T) {
	is := is.New(t)
	template := `package main

	func main()  {
		  println("{{ .name }}")
}`
	expect := `package main

func main() {
	println("jason")
}
`
	generator := gotemplate.MustParse("test.gotext", template)
	b, err := generator.Generate(map[string]string{"name": "jason"})
	is.NoErr(err)
	is.Equal(string(b), expect)
}

func TestGenerateFreeText(t *testing.T) {
	is := is.New(t)
	template := `Hi {{ .name }}`
	expect := `Hi Kim`
	generator := gotemplate.MustParse("test.gotext", template)
	b, err := generator.Generate(map[string]string{"name": "Kim"})
	is.NoErr(err)
	is.Equal(string(b), expect)
}
