package main

import (
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/aligator/goslice"
	"github.com/aligator/goslice/data"

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
		os.Exit(1)
	}

	if _, err := os.Stat(o.GoSlice.InputFilePath); errors.Is(err, os.ErrNotExist) {
		_, _ = fmt.Fprintf(os.Stderr, "the file doesn't exist\n")
		os.Exit(2)
	}

	p := goslice.NewGoSlice(o)
	err := p.Process()

	if err != nil {
		fmt.Println("error while processing file:", err)
		os.Exit(3)
	}
}

func printVersion(w io.Writer) {
	str := fmt.Sprintf("GoSlice %s", Version)
	_, _ = w.Write([]byte(str))
}
