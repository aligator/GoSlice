package data_test

import (
	"GoSlice/data"
	"GoSlice/util/test"
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

	test.Equals(t, data.Micrometer(400), data.XDistance2ToLine(vec1, vec2, point))
}
