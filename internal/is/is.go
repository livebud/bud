// Package is has a modified superset of the excellent github.com/matryer/is.
//
// It adds new methods commonly needed in Bud, removes indentation and replaces
// displaying comments with optional messages that can be formatted.
//
// Learn more in this discussion: https://github.com/livebud/bud/discussions/86
package is

import (
	"reflect"
)

// Helper marks the calling function as a test helper function.
// When printing file and line information, that function will be skipped.
//
// Available with Go 1.7 and later.
func (is *I) Helper() {
	is.helpers[callerName(1)] = struct{}{}
}

// Equal asserts that a and b are equal.
//
//	func Test(t *testing.T) {
//		is := is.New(t)
//		a := greet("Mat")
//		is.Equal(a, "Hi Mat") // greeting
//	}
//
// Will output:
//
//	your_test.go:123: Hey Mat != Hi Mat // greeting
func (is *I) Equal(a interface{}, b interface{}, args ...interface{}) {
	if areEqual(a, b) {
		return
	}
	if isNil(a) || isNil(b) {
		is.logf("%s != %s%s", is.valWithType(a), is.valWithType(b), is.formatArgs(args))
	} else if reflect.ValueOf(a).Type() == reflect.ValueOf(b).Type() {
		is.logf("%v != %v%s", a, b, is.formatArgs(args))
	} else {
		is.logf("%s != %s%s", is.valWithType(a), is.valWithType(b), is.formatArgs(args))
	}
}

// NoErr asserts that err is nil.
//
//	func Test(t *testing.T) {
//		is := is.New(t)
//		val, err := getVal()
//		is.NoErr(err)        // getVal error
//		is.True(len(val) > 10) // val cannot be short
//	}
//
// Will output:
//
//	your_test.go:123: err: not found // getVal error
func (is *I) NoErr(err error, args ...interface{}) {
	if err != nil {
		is.logf("err: %s%s", err.Error(), is.formatArgs(args))
	}
}

// New is a method wrapper around the New function.
// It allows you to write subtests using a similar
// pattern:
//
//	func Test(t *testing.T) {
//		is := is.New(t)
//		t.Run("sub", func(t *testing.T) {
//			is := is.New(t)
//			// TODO: test
//		})
//	}
func (is *I) New(t T) *I {
	return New(t)
}

// Fail immediately fails the test.
//
//	func Test(t *testing.T) {
//		is := is.New(t)
//		is.Fail() // TODO: write this test
//	}
func (is *I) Fail(args ...interface{}) {
	is.logf("failed%s", is.formatArgs(args))
}

// True asserts that the expression is true. The expression
// code itself will be reported if the assertion fails.
//
//	func Test(t *testing.T) {
//		is := is.New(t)
//		val := method()
//		is.True(val != nil) // val should never be nil
//	}
//
// Will output:
//
//	your_test.go:123: not true: val != nil
func (is *I) True(expr bool, args ...interface{}) {
	if !expr {
		is.logf("not true: ${1}%s", is.formatArgs(args))
	}
}

// In asserts that item is contained within list.
func (is *I) In(list interface{}, item interface{}, args ...interface{}) {
	ok, found := containsElement(list, item)
	if !ok {
		is.logf("%s is not a list%s", is.valWithType(list), is.formatArgs(args))
	} else if !found {
		is.logf("%v not in %+v%s", item, list, is.formatArgs(args))
	}
}

// NotIn asserts that item is not contained within list
func (is *I) NotIn(list interface{}, item interface{}, args ...interface{}) {
	panic("TODO: implement")
}

// ErrIs asserts that err is the same type of error as target
func (is *I) ErrIs(err, target error, args ...interface{}) {
	panic("TODO: implement")
}
