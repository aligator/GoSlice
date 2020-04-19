package main

import (
	"GoSlicer/go_slicer"
	"fmt"
	"os"
)

func main() {
	args := os.Args[1:]

	if len(args) != 1 {
		fmt.Println("you have to pass a stl file to slice")
		os.Exit(2)
	}

	p := go_slicer.NewGoSlicer()
	err := p.Process(args[0], args[0]+".gcode")

	if err != nil {
		fmt.Println("error while processing file:", err)
		os.Exit(2)
	}
}
