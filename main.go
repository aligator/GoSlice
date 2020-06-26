package main

import (
	"GoSlice/data"
	"fmt"
	"os"
)

func main() {
	o := data.ParseFlags()

	// remove infills to debug support
	o.Print.InfillPercent = 0
	o.Print.NumberBottomLayers = 0
	o.Print.NumberTopLayers = 0

	// enable support to debug it
	o.Print.Support.Enabled = true

	p := NewGoSlice(o)
	err := p.Process()

	if err != nil {
		fmt.Println("error while processing file:", err)
		os.Exit(2)
	}
}
