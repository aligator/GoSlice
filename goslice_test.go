package goslice

import (
	"strings"
	"testing"

	"github.com/aligator/goslice/data"
	"github.com/aligator/goslice/util/test"
)

const (
	folder = "./test_stl/"

	// The models are copied to the project just to avoid downloading them for each test.

	// 3DBenchy is the unmodified model from here:
	// https://www.thingiverse.com/thing:763622
	// using the following license
	// https://creativecommons.org/licenses/by-nd/4.0/
	benchy = "3DBenchy.stl"

	// Go Gopher mascot is the unmodified model from here:
	// https://www.thingiverse.com/thing:3413597
	// using the following license
	// https://creativecommons.org/licenses/by/4.0/
	gopher = "gopher_union.stl"
)

func TestWholeSlicer(t *testing.T) {
	o := data.DefaultOptions()
	// enable support so that it is tested also
	o.Print.Support.Enabled = true
	o.Print.BrimSkirt.BrimCount = 3
	s := NewGoSlice(o)

	var tests = []struct {
		path string
	}{
		{
			path: benchy,
		},
		{
			path: gopher,
		},
	}

	for _, testCase := range tests {
		t.Log("slice " + testCase.path)
		s.Options.InputFilePath = folder + testCase.path
		err := s.Process()
		test.Ok(t, err)
	}
}

func TestStartGCode(t *testing.T) {
	var tests = []struct {
		Name       string
		StartGCode data.GCodeHunk
		expected   string
	}{
		{
			Name:       "basic",
			StartGCode: data.NewGCodeHunk([]string{";Set Hotend", "M104 S0"}),
			expected:   ";Set Hotend\nM104 S0",
		},
		{
			Name:       "no startcode supplied",
			StartGCode: data.DefaultOptions().Printer.StartGCode,
			expected:   ";SET BED TEMP\nM190 S60 ; heat and wait for bed\n;SET HOTEND TEMP\nM109 S205 ; wait for hot end temperature",
		},
		{
			Name:       "no temp setting in start code",
			StartGCode: data.NewGCodeHunk([]string{"printer_start"}),
			expected:   ";SET BED TEMP\nM190 S60 ; heat and wait for bed\n;SET HOTEND TEMP\nM109 S205 ; wait for hot end temperature\n;START GCODE\nprinter_start",
		},
	}

	for _, testCase := range tests {
		o := data.DefaultOptions()
		t.Log(testCase.Name)
		o.Printer.StartGCode = testCase.StartGCode
		s := NewGoSlice(o)
		s.Options.InputFilePath = folder + benchy
		gcode, err := s.GetGCode()
		test.Ok(t, err)
		test.Assert(t, strings.Contains(strings.Join(gcode, "\n"), testCase.expected),
			"final gcode does not contain the expected gcode.\nexpected:"+
				testCase.expected+
				"\nActual: "+
				strings.Join(gcode[0:10], "\n"))
	}
}

func TestEndGCode(t *testing.T) {
	var tests = []struct {
		Name     string
		EndGCode data.GCodeHunk
		Options  data.Options
		expected string
	}{
		{
			Name:     "basic",
			EndGCode: data.NewGCodeHunk([]string{";Set Hotend", "M104 S0"}),
			Options:  data.DefaultOptions(),
			expected: ";Set Hotend\nM104 S0",
		},
		{
			Name:     "no endcode supplied",
			EndGCode: data.DefaultOptions().Printer.EndGCode,
			Options:  data.DefaultOptions(),
			expected: ";END_GCODE\nM104 S0 ; Set Hot-end to 0C (off)\nM140 S0 ; Set bed to 0C (off)",
		},
		{
			Name:     "no temp setting in end code",
			EndGCode: data.NewGCodeHunk([]string{"printer_stop"}),
			Options:  data.DefaultOptions(),
			expected: ";END_GCODE\nM104 S0 ; Set Hot-end to 0C (off)\nM140 S0 ; Set bed to 0C (off)\nprinter_stop",
		},
		{
			Name:     "no heated bed",
			EndGCode: data.NewGCodeHunk([]string{"printer_stop"}),
			Options:  data.DefaultOptions().SetHasHeatedBed(false),
			expected: ";END_GCODE\nM104 S0 ; Set Hot-end to 0C (off)\nprinter_stop",
		},
	}

	for _, testCase := range tests {
		testCase.Options.Printer.EndGCode = testCase.EndGCode
		t.Log(testCase.Name)
		s := NewGoSlice(testCase.Options)
		s.Options.InputFilePath = folder + benchy
		gcode, err := s.GetGCode()
		test.Ok(t, err)
		test.Assert(t, strings.Contains(strings.Join(gcode, "\n"), testCase.expected),
			"final gcode does not contain the expected gcode.\nexpected:"+
				testCase.expected+
				"\nActual: "+
				strings.Join(gcode[len(gcode)-10:], "\n"))
	}
}
