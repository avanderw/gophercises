package main

import (
	"fmt"
	"os"
)

func main() {
	runFile("ex1.html")
	runFile("ex2.html")
	runFile("ex3.html")
	runFile("ex4.html")
}

func runFile(path string) {
	f, err := os.Open(path)
	if err != nil {
		panic(err)
	}

	links, err := Parse(f)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v\n", links)
}
