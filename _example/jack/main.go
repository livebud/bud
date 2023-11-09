package main

import (
	"fmt"

	"github.com/livebud/bud/pkg/mux"
)

func main() {
	router := mux.New()
	preact := preact.New()
	// router.Get("/", http.HandlerFunc())
	fmt.Println("OK!!!")
}
