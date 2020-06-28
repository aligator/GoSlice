package gcode_test

import (
	"GoSlice/data"
	"GoSlice/gcode"
	"GoSlice/util/test"
	"testing"
)

func TestGCodeBuilder(t *testing.T) {
	var tests = map[string]struct {
		exec     func(*gcode.Builder)
		expected string
	}{
		"no commands": {
			exec:     func(b *gcode.Builder) {},
			expected: "",
		},
		"one simple command": {
			exec: func(b *gcode.Builder) {
				b.AddCommand("G1 X%0.2f Y%0.2f", 5.0, 6.9)
			},
			expected: "G1 X5.00 Y6.90\n",
		},
		"several commands": {
			exec: func(b *gcode.Builder) {
				b.AddCommand("G1 X%0.2f Y%0.2f", 5.0, 6.9)
				b.AddCommand("G1 X%0.2f Y%0.2f", 99.0, 6.9)
				b.AddCommand("G1 X%0.2f Y%0.2f", 5.0, 6.88)
			},
			expected: "G1 X5.00 Y6.90\n" +
				"G1 X99.00 Y6.90\n" +
				"G1 X5.00 Y6.88\n",
		},
		"some comments": {
			exec: func(b *gcode.Builder) {
				b.AddComment("This is a comment.")
				b.AddComment("This is another comment.")
				b.AddComment("This is a third comment.")
			},
			expected: ";This is a comment.\n" +
				";This is another comment.\n" +
				";This is a third comment.\n",
		},
		"add polygon": {
			exec: func(b *gcode.Builder) {
				err := b.AddPolygon(nil, data.Path{
					data.NewMicroPoint(0, 0),
					data.NewMicroPoint(100, 0),
					data.NewMicroPoint(100, 100),
					data.NewMicroPoint(0, 100),
				}, 100, true)
				test.Ok(t, err)

				// empty polygon should just be ignored
				err = b.AddPolygon(nil, data.Path{}, 100, false)
				test.Ok(t, err)
				err = b.AddPolygon(nil, data.Path{
					data.NewMicroPoint(0, 0),
					data.NewMicroPoint(50, 0),
					data.NewMicroPoint(50, 50),
					data.NewMicroPoint(0, 50),
				}, 100, false)
				test.Ok(t, err)
			},
			expected: "G0 X0.00 Y0.00 Z0.10\n" +
				"G0 X0.10 Y0.00\n" +
				"G0 X0.10 Y0.10\n" +
				"G0 X0.00 Y0.10\n" +
				"G0 X0.00 Y0.00\n" +
				"G0 X0.05 Y0.00\n" +
				"G0 X0.05 Y0.05\n" +
				"G0 X0.00 Y0.05\n" +
				"G0 X0.00 Y0.00\n",
		},
		"some moves": {
			exec: func(b *gcode.Builder) {
				b.AddMove(data.NewMicroVec3(0, 0, 0), 0)
				b.AddMove(data.NewMicroVec3(10, 0, 0), 0)
				b.AddMove(data.NewMicroVec3(0, 20, 0), 5)
				b.AddMove(data.NewMicroVec3(0, 0, 30), 0)

				b.AddMove(data.NewMicroVec3(0, 0, 0), 2)
				b.AddMove(data.NewMicroVec3(10, 0, 0), 3)
				b.AddMove(data.NewMicroVec3(0, 20, 0), 4)
				b.AddMove(data.NewMicroVec3(0, 0, 30), 5)
			},
			expected: "G0 X0.00 Y0.00\n" +
				"G0 X0.01 Y0.00\n" +
				"G1 X0.00 Y0.02 E5.0000\n" +
				"G0 X0.00 Y0.00 Z0.03\n" +
				"G1 X0.00 Y0.00 Z0.00 E7.0000\n" +
				"G1 X0.01 Y0.00 E10.0000\n" +
				"G1 X0.00 Y0.02 E14.0000\n" +
				"G1 X0.00 Y0.00 Z0.03 E19.0000\n",
		},
		"different speeds": {
			exec: func(b *gcode.Builder) {
				b.SetMoveSpeed(200)
				b.SetExtrudeSpeed(100)

				b.AddMove(data.NewMicroVec3(0, 0, 0), 0)
				b.AddMove(data.NewMicroVec3(10, 0, 0), 5)
				b.AddMove(data.NewMicroVec3(40, 0, 0), 5)
				b.AddMove(data.NewMicroVec3(10, 0, 0), 0)

				b.SetExtrudeSpeedOverride(150)
				b.AddMove(data.NewMicroVec3(0, 0, 0), 0)
				b.AddMove(data.NewMicroVec3(10, 0, 0), 5)
				b.AddMove(data.NewMicroVec3(40, 0, 0), 5)
				b.AddMove(data.NewMicroVec3(10, 0, 0), 0)

				b.DisableExtrudeSpeedOverride()
				b.AddMove(data.NewMicroVec3(0, 0, 0), 0)
				b.AddMove(data.NewMicroVec3(10, 0, 0), 5)
				b.AddMove(data.NewMicroVec3(40, 0, 0), 5)
				b.AddMove(data.NewMicroVec3(10, 0, 0), 0)

				b.SetMoveSpeed(600)
				b.AddMove(data.NewMicroVec3(0, 0, 0), 0)
				b.AddMove(data.NewMicroVec3(10, 0, 0), 5)
				b.AddMove(data.NewMicroVec3(40, 0, 0), 5)
				b.AddMove(data.NewMicroVec3(10, 0, 0), 0)

				b.SetExtrudeSpeed(500)
				b.AddMove(data.NewMicroVec3(0, 0, 0), 0)
				b.AddMove(data.NewMicroVec3(10, 0, 0), 5)
				b.AddMove(data.NewMicroVec3(40, 0, 0), 5)
				b.AddMove(data.NewMicroVec3(10, 0, 0), 0)
			},
			expected: "G0 X0.00 Y0.00 F12000\n" +
				"G1 X0.01 Y0.00 F6000 E5.0000\n" +
				"G1 X0.04 Y0.00 E10.0000\n" +
				"G0 X0.01 Y0.00 F12000\n" +
				"G0 X0.00 Y0.00\n" +
				"G1 X0.01 Y0.00 F9000 E15.0000\n" +
				"G1 X0.04 Y0.00 E20.0000\n" +
				"G0 X0.01 Y0.00 F12000\n" +
				"G0 X0.00 Y0.00\n" +
				"G1 X0.01 Y0.00 F6000 E25.0000\n" +
				"G1 X0.04 Y0.00 E30.0000\n" +
				"G0 X0.01 Y0.00 F12000\n" +
				"G0 X0.00 Y0.00 F36000\n" +
				"G1 X0.01 Y0.00 F6000 E35.0000\n" +
				"G1 X0.04 Y0.00 E40.0000\n" +
				"G0 X0.01 Y0.00 F36000\n" +
				"G0 X0.00 Y0.00\n" +
				"G1 X0.01 Y0.00 F30000 E45.0000\n" +
				"G1 X0.04 Y0.00 E50.0000\n" +
				"G0 X0.01 Y0.00 F36000\n",
		},

		"set extrusion": {
			exec: func(b *gcode.Builder) {
				b.SetExtrusion(200, 400, 175)
				err := b.AddPolygon(nil, []data.MicroPoint{
					data.NewMicroPoint(0, 0),
					data.NewMicroPoint(0, 10000),
				}, 0, true)
				test.Ok(t, err)
			},
			expected: "G0 X0.00 Y0.00\n" +
				"G1 X0.00 Y10.00 E33.2601\n",
		},
	}

	for desc, testCase := range tests {
		t.Log(desc)
		builder := gcode.NewGCodeBuilder()
		testCase.exec(builder)
		test.Equals(t, testCase.expected, builder.String())
	}
}
