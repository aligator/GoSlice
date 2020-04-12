package main

import (
	"GoSlicer/slicer"
	"fmt"
	"os"
)

func main() {
	args := os.Args[1:]

	if len(args) != 1 {
		fmt.Println("you have to pass a stl file to slice")
		os.Exit(2)
	}

	s := slicer.Slicer{Path: args[0]}
	err := s.Process()

	if err != nil {
		fmt.Println("error while processing file")
		os.Exit(2)
	}
}
