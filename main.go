package main

import (
	"GoSlice/data"
	"fmt"
	"io"
	"os"

	flag "github.com/spf13/pflag"
)

var Version = "unknown development version"

func main() {
	o := data.ParseFlags()

	if o.GoSlice.PrintVersion {
		printVersion(os.Stdout)
		os.Exit(0)
	}

	if o.GoSlice.InputFilePath == "" {
		_, _ = fmt.Fprintf(os.Stderr, "the STL_FILE path has to be specified\n")
		flag.Usage()
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
