// Package ldflag contains flags set at link time via LDFLAGS="...". Runtime
// packages can rely on these variables to access information about the build.
//
// e.g.
// go run -ldflags="-X 'github.com/go-duo/ldflag/ldflag.hot=true'" main.go
//
package ldflag

import "strconv"

var env string = "development"

func Env() string {
	return env
}

func IsDevelopment() bool {
	return env == "development"
}

func IsTest() bool {
	return env == "test"
}

func IsProduction() bool {
	return env == "production"
}

var minify string = "false"

func Minify() bool {
	value, err := strconv.ParseBool(minify)
	if err != nil {
		panic("ldflag: minify flag is invalid")
	}
	return value
}

var hot string = "true"

func Hot() bool {
	value, err := strconv.ParseBool(hot)
	if err != nil {
		panic("ldflag: hot flag is invalid")
	}
	return value
}
