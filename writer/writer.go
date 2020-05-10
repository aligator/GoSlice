package writer

import (
	"GoSlice/handler"
	"fmt"
	"os"
)

type writer struct{}

// Writer can write gcode to a file.
func Writer() handler.GCodeWriter {
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
