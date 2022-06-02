package main

import "fmt"

func main() {
	files := []*File{{"A"}, {"B"}}
	add(&files, &File{"C"})
	fmt.Println(files)
}

type File struct {
	Name string
}

func add(list *[]*File, files ...*File) {
	*list = append(*list, files...)
}
