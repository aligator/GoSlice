package main

import (
	"GoSlice/data"
	"fmt"
	clipper "github.com/aligator/go.clipper"
	"os"
)

func main() {
	cl := clipper.NewClipper(clipper.IoNone)

	res, ok := cl.Execute2(clipper.CtIntersection, clipper.PftEvenOdd, clipper.PftEvenOdd)

	if ok {
		fmt.Println(res)
	}

	args := os.Args[1:]

	if len(args) != 1 {
		fmt.Println("you have to pass a stl file to slice")
		os.Exit(2)
	}

	p := NewGoSlice(
		Center(data.NewMicroVec3(data.Millimeter(100).ToMicrometer(), data.Millimeter(100).ToMicrometer(), 0)),
		InsetCount(2))
	err := p.Process(args[0], args[0]+".gcode")

	if err != nil {
		fmt.Println("error while processing file:", err)
		os.Exit(2)
	}
}
