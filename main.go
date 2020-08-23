package main

import (
	"GoSlice/data"
	"fmt"
	"io"
	"os"
)

var Version = "UNDEFINED"

func main() {
	o := data.ParseFlags()

	if o.GoSlice.PrintVersion {
		printVersion(os.Stdout)
		os.Exit(0)
	}

	p := NewGoSlice(o)
	err := p.Process()

	if err != nil {
		fmt.Println("error while processing file:", err)
		os.Exit(2)
	}
}

func printVersion(w io.Writer) {
	str := fmt.Sprintf("GoSlice %s", Version)
	_, _ = w.Write([]byte(str))
}
