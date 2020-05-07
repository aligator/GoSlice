package main

import (
	"GoSlice/data"
	"fmt"
	"os"
)

func main() {

	args := os.Args[1:]

	if len(args) != 1 {
		fmt.Println("you have to pass a stl file to slice")
		os.Exit(2)
	}

	o := DefaultOptions()
	o.Printer.Center = data.NewMicroVec3(data.Millimeter(100).ToMicrometer(), data.Millimeter(100).ToMicrometer(), 0)
	o.Print.InsetCount = 2

	p := NewGoSlice(o)

	err := p.Process(args[0], args[0]+".gcode")

	if err != nil {
		fmt.Println("error while processing file:", err)
		os.Exit(2)
	}
}
