package parser_test

import (
	"testing"

	"github.com/matryer/is"
	"gitlab.com/mnm/bud/internal/modtest"
	"gitlab.com/mnm/bud/internal/parser"
)

// TODO: test the parser types

func isBuiltin(f *parser.Field) bool {
	if f == nil {
		return false
	}
	return parser.IsBuiltin(f.Type())
}

func isNotBuiltin(f *parser.Field) bool {
	return !isBuiltin(f)
}

func TestRequalify(t *testing.T) {
	is := is.New(t)
	module := modtest.Make(t, modtest.Module{
		Files: map[string][]byte{
			"go.mod": []byte(`module app.com/app`),
			"js/v8/v8.go": []byte(`
				package v8

				type VM struct {}
			`),
			"app.go": []byte(`
				package app

				import v8 "app.com/app/js/v8"

				type Test struct {
					vms []*[]v8.VM
					none []*[]string
					same []*[]v8.VM
				}
			`),
		},
	})
	p := parser.New(module)
	pkg, err := p.Parse(".")
	is.NoErr(err)
	stct := pkg.Struct("Test")
	is.True(stct != nil)
	// vms
	field := stct.Field("vms")
	is.True(field != nil)
	typ := parser.Requalify(field.Type(), "js")
	is.Equal(typ.String(), "[]*[]js.VM")
	// none
	field = stct.Field("none")
	is.True(field != nil)
	typ = parser.Requalify(field.Type(), "js")
	is.Equal(typ.String(), "[]*[]string")
	// same
	field = stct.Field("same")
	is.True(field != nil)
	typ = parser.Requalify(field.Type(), "v8")
	is.Equal(typ.String(), "[]*[]v8.VM")
}

func TestBuiltins(t *testing.T) {
	is := is.New(t)
	module := modtest.Make(t, modtest.Module{
		Files: map[string][]byte{
			"go.mod": []byte(`module app.com/app`),
			"app.go": []byte(`
				package app

				type Builtin struct {
					String string
					StringPtr *string
					Strings []string
					StringPtrs []*string
					StringsPtr *[]string
					String2D [][]string
					StringArr [2][3]string

					Bool bool
					BoolPtr *bool
					Bools []bool
					BoolPtrs []*bool
					BoolsPtr *[]bool
					Bool2D [][]bool
					BoolArr [2][3]bool

					Int int
					IntPtr *int
					Ints []int
					IntPtrs []*int
					IntsPtr *[]int
					Int2D [][]int
					IntArr [2][3]int

					Int32 int32
					Int32Ptr *int32
					Int32s []int32
					Int32Ptrs []*int32
					Int32sPtr *[]int32
					Int322D [][]int32
					Int32Arr [2][3]int32

					Struct struct{}
				}
			`),
		},
	})
	p := parser.New(module)
	pkg, err := p.Parse(".")
	is.NoErr(err)
	stct := pkg.Struct("Builtin")
	is.True(stct != nil)
	// Strings
	is.True(isBuiltin(stct.Field("String")))
	is.True(isBuiltin(stct.Field("StringPtr")))
	is.True(isBuiltin(stct.Field("Strings")))
	is.True(isBuiltin(stct.Field("StringPtrs")))
	is.True(isBuiltin(stct.Field("StringsPtr")))
	is.True(isBuiltin(stct.Field("String2D")))
	is.True(isBuiltin(stct.Field("StringArr")))
	// Bools
	is.True(isBuiltin(stct.Field("Bool")))
	is.True(isBuiltin(stct.Field("BoolPtr")))
	is.True(isBuiltin(stct.Field("Bools")))
	is.True(isBuiltin(stct.Field("BoolPtrs")))
	is.True(isBuiltin(stct.Field("BoolsPtr")))
	is.True(isBuiltin(stct.Field("Bool2D")))
	is.True(isBuiltin(stct.Field("BoolArr")))
	// Int
	is.True(isBuiltin(stct.Field("Int")))
	is.True(isBuiltin(stct.Field("IntPtr")))
	is.True(isBuiltin(stct.Field("Ints")))
	is.True(isBuiltin(stct.Field("IntPtrs")))
	is.True(isBuiltin(stct.Field("IntsPtr")))
	is.True(isBuiltin(stct.Field("Int2D")))
	is.True(isBuiltin(stct.Field("IntArr")))
	// Int32
	is.True(isBuiltin(stct.Field("Int32")))
	is.True(isBuiltin(stct.Field("Int32Ptr")))
	is.True(isBuiltin(stct.Field("Int32s")))
	is.True(isBuiltin(stct.Field("Int32Ptrs")))
	is.True(isBuiltin(stct.Field("Int32sPtr")))
	is.True(isBuiltin(stct.Field("Int322D")))
	is.True(isBuiltin(stct.Field("Int32Arr")))
	// TODO: add remaining types
	// Structs
	is.True(isNotBuiltin(stct.Field("Struct")))
}
