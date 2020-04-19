package main

import (
	"GoSlice/util"
	"fmt"
	"os"
)

func main() {
	args := os.Args[1:]

	if len(args) != 1 {
		fmt.Println("you have to pass a stl file to slice")
		os.Exit(2)
	}

	p := NewGoSlice(
		Center(util.NewMicroVec3(util.Millimeter(100).ToMicrometer(), util.Millimeter(100).ToMicrometer(), 0)),
		InsetCount(3))
	err := p.Process(args[0], args[0]+".gcode")

	if err != nil {
		fmt.Println("error while processing file:", err)
		os.Exit(2)
	}
}
