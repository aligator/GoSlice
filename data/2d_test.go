package data_test

import (
	"GoSlice/data"
	"GoSlice/util/test"
	"math"
	"testing"
)

// Test MicroVec3 implementation

func TestDotProduct(t *testing.T) {
	vec1 := data.NewMicroPoint(10, 20)
	vec2 := data.NewMicroPoint(110, 120)

	test.Equals(t, data.Micrometer(3500), data.DotProduct(vec1, vec2))
}

func TestXDistance2ToLine(t *testing.T) {
	vec1 := data.NewMicroPoint(0, 20)
	vec2 := data.NewMicroPoint(50, 20)

	point := data.NewMicroPoint(0, 40)

	test.Equals(t, data.Micrometer(400), data.PerpendicularDistance2(vec1, vec2, point))
}

func TestToRadians(t *testing.T) {
	var testCases = []struct {
		expected float64
		degree   float64
	}{
		{
			expected: 1,
			degree:   57.29577951308232,
		},
		{
			expected: math.Pi,
			degree:   180,
		},
		{
			expected: 0.5235987755982988, // math.Pi / 6
			degree:   30,
		},
	}

	for i, testCase := range testCases {
		t.Log("testCase", i)
		test.Equals(t, testCase.expected, data.ToRadians(testCase.degree))
	}
}
