package main

import (
	"GoSlice/data"
	"GoSlice/util/test"
	"testing"
)

const (
	folder = "../../test_stl/"

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
		s.options.GoSlice.InputFilePath = folder + testCase.path
		err := s.Process()
		test.Ok(t, err)
	}
}
