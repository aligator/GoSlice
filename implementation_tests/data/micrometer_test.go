package data

import (
	"GoSlice/data"
	"testing"
)

const (
	x = 10
	y = 20
	z = 30
)

func setupMicroVec3() data.MicroVec3 {
	return data.NewMicroVec3(x, y, z)
}

func TestNewMicroVec3(t *testing.T) {
	vec := data.NewMicroVec3(x, y, z)
	if vec == nil {
		t.Error("vec should not be nil")
	}
}

func TestMicroVec3Values(t *testing.T) {
	vec := setupMicroVec3()

	if vec.X() != x {
		t.Errorf("X() should be %v but it is %v", x, vec.X())
	}
	if vec.Y() != y {
		t.Errorf("Y() should be %v but it is %v", y, vec.Y())
	}
	if vec.Z() != z {
		t.Errorf("Z() should be %v but it is %v", z, vec.Z())
	}
}

func TestMicroVec3Add(t *testing.T) {
	var expected = [3]data.Micrometer{20, 40, 60}

	vec := setupMicroVec3()
	vec2 := setupMicroVec3()
	actual := vec.Add(vec2)

	new := setupMicroVec3()
	if vec == new {
		t.Errorf("it should not modify the instance directly")
	}

	if actual.X() != expected[0] {
		t.Errorf("X() should be %v but it is %v", expected[0], actual.X())
	}
	if actual.Y() != expected[1] {
		t.Errorf("Y() should be %v but it is %v", expected[1], actual.Y())
	}
	if actual.Z() != expected[2] {
		t.Errorf("Z() should be %v but it is %v", expected[2], actual.Z())
	}
}

func TestMicroVec3Sub(t *testing.T) {
	var expected = [3]data.Micrometer{0, 0, 0}

	vec := setupMicroVec3()
	vec2 := setupMicroVec3()
	actual := vec.Sub(vec2)

	new := setupMicroVec3()
	if vec == new {
		t.Errorf("it should not modify the instance directly")
	}

	if actual.X() != expected[0] {
		t.Errorf("X() should be %v but it is %v", expected[0], actual.X())
	}
	if actual.Y() != expected[1] {
		t.Errorf("Y() should be %v but it is %v", expected[1], actual.Y())
	}
	if actual.Z() != expected[2] {
		t.Errorf("Z() should be %v but it is %v", expected[2], actual.Z())
	}
}

func TestMicroVec3Mul(t *testing.T) {
	var expected = [3]data.Micrometer{30, 60, 90}

	vec := setupMicroVec3()
	actual := vec.Mul(3)

	new := setupMicroVec3()
	if vec == new {
		t.Errorf("it should not modify the instance directly")
	}

	if actual.X() != expected[0] {
		t.Errorf("X() should be %v but it is %v", expected[0], actual.X())
	}
	if actual.Y() != expected[1] {
		t.Errorf("Y() should be %v but it is %v", expected[1], actual.Y())
	}
	if actual.Z() != expected[2] {
		t.Errorf("Z() should be %v but it is %v", expected[2], actual.Z())
	}
}

func TestMicroVec3Div(t *testing.T) {
	var expected = [3]data.Micrometer{5, 10, 15}

	vec := setupMicroVec3()
	actual := vec.Div(2)

	new := setupMicroVec3()
	if vec == new {
		t.Errorf("it should not modify the instance directly")
	}

	if actual.X() != expected[0] {
		t.Errorf("X() should be %v but it is %v", expected[0], actual.X())
	}
	if actual.Y() != expected[1] {
		t.Errorf("Y() should be %v but it is %v", expected[1], actual.Y())
	}
	if actual.Z() != expected[2] {
		t.Errorf("Z() should be %v but it is %v", expected[2], actual.Z())
	}
}

func TestMicroVec3Max(t *testing.T) {
	var expected = data.Micrometer(30)

	vec := setupMicroVec3()
	actual := vec.Max()

	if actual != expected {
		t.Errorf("the maximum should be %v but it is %v", expected, actual)
	}
}

func TestMicroVec3PointXY(t *testing.T) {
	var expected = [2]data.Micrometer{10, 20}

	vec := setupMicroVec3()
	result := vec.PointXY()

	if result.X() != expected[0] {
		t.Errorf("Y() should be %v but it is %v", expected, result)
	}
	if result.Y() != expected[1] {
		t.Errorf("Y() should be %v but it is %v", expected, result)
	}
}

func TestMicroVec3TestLength(t *testing.T) {
	vec := setupMicroVec3()

	var tests = []struct {
		expected bool
		length   data.Micrometer
	}{
		{true, 100},
		{true, 38},
		{false, 37},
		{false, 36},
		{false, 0},
	}

	for _, test := range tests {
		actual := vec.TestLength(test.length)
		if test.expected != actual {
			t.Errorf("the length %v should return %v but it returns %v", test.length, test.expected, actual)
		}
	}
}

func TestMicroVec3TestSize2(t *testing.T) {
	var expected = data.Micrometer(1400)
	vec := setupMicroVec3()

	actual := vec.Size2()
	if expected != actual {
		t.Errorf("Size2() should return %v but it returns %v", expected, actual)
	}
}

func TestMicroVec3TestSize(t *testing.T) {
	var expected = data.Micrometer(37)
	vec := setupMicroVec3()

	actual := vec.Size()
	if expected != actual {
		t.Errorf("Size() should return %v but it returns %v", expected, actual)
	}
}

func TestMicroVec3TestCopy(t *testing.T) {
	vec := setupMicroVec3()

	actual := vec.Copy()

	if &actual == &vec {
		t.Errorf("Copy should create a new instance")
	}
	if actual.X() != vec.X() {
		t.Errorf("X() should be %v but it is %v", vec.X(), actual.X())
	}
	if actual.Y() != vec.Y() {
		t.Errorf("Y() should be %v but it is %v", vec.Y(), actual.Y())
	}
	if actual.Z() != vec.Z() {
		t.Errorf("Z() should be %v but it is %v", vec.Z(), actual.Z())
	}
}

func TestMicroVec3TestString(t *testing.T) {
	var expected = "10_20_30"
	vec := setupMicroVec3()

	actual := vec.String()
	if expected != actual {
		t.Errorf("String() should return %v but it returns %v", expected, actual)
	}
}

func TestMicroVec3TestSet(t *testing.T) {
	var expected = [3]data.Micrometer{40, 60, 200}

	actual := setupMicroVec3()
	err := actual.Set("40_60_200")

	if err != nil {
		t.Error(err)
	}

	if actual.X() != expected[0] {
		t.Errorf("X() should be %v but it is %v", expected[0], actual.X())
	}
	if actual.Y() != expected[1] {
		t.Errorf("Y() should be %v but it is %v", expected[1], actual.Y())
	}
	if actual.Z() != expected[2] {
		t.Errorf("Z() should be %v but it is %v", expected[2], actual.Z())
	}
}

func TestMicroVec3TestType(t *testing.T) {
	var expected = "Micrometer"
	vec := setupMicroVec3()

	actual := vec.Type()
	if expected != actual {
		t.Errorf("Type() should return %v but it returns %v", expected, actual)
	}
}
