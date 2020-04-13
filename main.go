package main

import (
	"fmt"
	"os"
)

func main() {
	args := os.Args[1:]

	if len(args) != 1 {
		fmt.Println("you have to pass a stl file to slice")
		os.Exit(2)
	}

	p := Process{Path: args[0]}
	err := p.Process()

	if err != nil {
		fmt.Println("error while processing file")
		os.Exit(2)
	}
}
