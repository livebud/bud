// Package ldflag contains flags set at link time via LDFLAGS="...". Runtime
// packages can rely on these variables to access information about the build.
package ldflag

var Env string = "development"

func IsDevelopment() bool {
	return Env == "development"
}

func IsTest() bool {
	return Env == "test"
}

func IsProduction() bool {
	return Env == "production"
}

var Minify bool = false

var Hot bool = true
