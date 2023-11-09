package js

import (
	"fmt"
	"os"

	"github.com/dop251/goja"
)

func New() *goja.Runtime {
	rt := goja.New()
	// Setup the console
	rt.GlobalObject().Set("console", map[string]interface{}{
		"log": func(args ...goja.Value) {
			var params []interface{}
			for _, arg := range args {
				params = append(params, arg.String())
			}
			os.Stdout.Write([]byte(fmt.Sprintln(params...)))
		},
		"warn": func(args ...goja.Value) {
			var params []interface{}
			for _, arg := range args {
				params = append(params, arg.String())
			}
			os.Stderr.Write([]byte(fmt.Sprintln(params...)))
		},
		"error": func(args ...goja.Value) {
			var params []interface{}
			for _, arg := range args {
				params = append(params, arg.String())
			}
			os.Stderr.Write([]byte(fmt.Sprintln(params...)))
		},
	})
	return rt
}
