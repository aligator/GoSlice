package write

import (
	"GoSlicer/slicer/handle"
	"fmt"
	"os"
)

type writer struct{}

func Writer() handle.GCodeWriter {
	return &writer{}
}

func (w writer) Write(gcode string, filename string) error {
	buf, err := os.Create(filename)
	if err != nil {
		fmt.Println(err)
	}

	defer buf.Close()

	_, err = buf.WriteString(gcode)
	return err
}
