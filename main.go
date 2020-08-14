package main

import (
	"GoSlice/data"
	"fmt"
	"os"
)

func main() {
	o := data.ParseFlags()

	p := NewGoSlice(o)
	err := p.Process()

	if err != nil {
		fmt.Println("error while processing file:", err)
		os.Exit(2)
	}
}
