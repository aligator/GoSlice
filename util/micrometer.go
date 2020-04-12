package util

import "math"

// MicroVec3 represents a point in 3d space
// which is in a Micrometer-grid.
// A value of 1 represents 0.001 mm
type MicroVec3 interface {
	X() Micrometer
	Y() Micrometer
	Z() Micrometer

	SetX(x Micrometer)
	SetY(y Micrometer)
	SetZ(z Micrometer)

	ToMilliVec3() MilliVec3

	Add(vec MicroVec3)
	Sub(vec MicroVec3)
	Mul(value Micrometer)
	Div(value Micrometer)

	TestLength(length Micrometer) bool
	Size2() Micrometer
	Size() Micrometer

	Copy() MicroVec3
}

type microVec3 struct {
	x, y, z Micrometer
}

func NewMicroVec3(x Micrometer, y Micrometer, z Micrometer) MicroVec3 {
	return &microVec3{
		x: x,
		y: y,
		z: z,
	}
}

func (v *microVec3) ToMilliVec3() MilliVec3 {
	return NewMilliVec3(v.x.ToMillimeter(), v.y.ToMillimeter(), v.z.ToMillimeter())
}

func (v *microVec3) X() Micrometer {
	return v.x
}

func (v *microVec3) Y() Micrometer {
	return v.y
}

func (v *microVec3) Z() Micrometer {
	return v.z
}

func (v *microVec3) SetX(x Micrometer) {
	v.x = x
}

func (v *microVec3) SetY(y Micrometer) {
	v.y = y
}

func (v *microVec3) SetZ(z Micrometer) {
	v.z = z
}

func (v *microVec3) Add(vec MicroVec3) {
	v.x += vec.X()
	v.y += vec.Y()
	v.z += vec.Z()
}

func (v *microVec3) Sub(vec MicroVec3) {
	v.x -= vec.X()
	v.y -= vec.Y()
	v.z -= vec.Z()
}

func (v *microVec3) Mul(value Micrometer) {
	v.x *= value
	v.y *= value
	v.z *= value
}

func (v *microVec3) Div(value Micrometer) {
	v.x /= value
	v.y /= value
	v.z /= value
}

func (v *microVec3) TestLength(length Micrometer) bool {
	return v.Size2() <= length*length
}

func (v *microVec3) Size2() Micrometer {
	return v.x*v.x + v.y*v.y + v.z*v.z
}

func (v *microVec3) Size() Micrometer {
	return Micrometer(math.Sqrt(float64(v.Size2())))
}

func (v *microVec3) Copy() MicroVec3 {
	return &microVec3{
		x: v.x,
		y: v.y,
		z: v.z,
	}
}
