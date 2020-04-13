package util

import "math"

// MilliVec3 represents a point in 3d space
// which is in a Millimeter-grid.
// A value of 1 represents 1 mm.
// Millimeter vectors are not used most time,
// because of possible rounding errors.
// They represent millimeters in 3D space.
type MilliVec3 interface {
	X() Millimeter
	Y() Millimeter
	Z() Millimeter

	SetX(x Millimeter)
	SetY(y Millimeter)
	SetZ(z Millimeter)

	ToMicroVec3() MicroVec3

	Add(vec MilliVec3) MilliVec3
	Sub(vec MilliVec3) MilliVec3
	Mul(value Millimeter) MilliVec3
	Div(value Millimeter) MilliVec3

	Max() Millimeter
	TestLength(length Millimeter) bool
	Size2() Millimeter
	Size() Millimeter
	Normalized() MilliVec3
	Cross(p2 MilliVec3) MilliVec3

	Copy() MilliVec3
}

type milliVec3 struct {
	x, y, z Millimeter
}

func NewMilliVec3(x Millimeter, y Millimeter, z Millimeter) MilliVec3 {
	return &milliVec3{
		x: x,
		y: y,
		z: z,
	}
}

func (v milliVec3) X() Millimeter {
	return v.x
}

func (v milliVec3) Y() Millimeter {
	return v.y
}

func (v milliVec3) Z() Millimeter {
	return v.z
}

func (v milliVec3) SetX(x Millimeter) {
	v.x = x
}

func (v milliVec3) SetY(y Millimeter) {
	v.y = y
}

func (v milliVec3) SetZ(z Millimeter) {
	v.z = z
}

func (v *milliVec3) ToMicroVec3() MicroVec3 {
	return NewMicroVec3(v.x.ToMicrometer(), v.y.ToMicrometer(), v.z.ToMicrometer())
}

func (v *milliVec3) Add(vec MilliVec3) MilliVec3 {
	result := v.Copy()
	result.SetX(result.X() + vec.X())
	result.SetY(result.Y() + vec.Y())
	result.SetZ(result.Z() + vec.Z())
	return result
}

func (v *milliVec3) Sub(vec MilliVec3) MilliVec3 {
	result := v.Copy()
	result.SetX(result.X() - vec.X())
	result.SetY(result.Y() - vec.Y())
	result.SetZ(result.Z() - vec.Z())
	return result
}

func (v *milliVec3) Mul(value Millimeter) MilliVec3 {
	result := v.Copy()
	result.SetX(result.X() * value)
	result.SetY(result.Y() * value)
	result.SetZ(result.Z() * value)
	return result
}

func (v *milliVec3) Div(value Millimeter) MilliVec3 {
	result := v.Copy()
	result.SetX(result.X() / value)
	result.SetY(result.Y() / value)
	result.SetZ(result.Z() / value)
	return result
}

func (v *milliVec3) Max() Millimeter {
	if v.x > v.y && v.x > v.z {
		return v.x
	}
	if v.y > v.z {
		return v.y
	}
	return v.z
}

func (v *milliVec3) TestLength(length Millimeter) bool {
	return v.Size2() <= length*length
}

func (v *milliVec3) Size2() Millimeter {
	return v.x*v.x + v.y*v.y + v.z*v.z
}

func (v *milliVec3) Size() Millimeter {
	return Millimeter(math.Sqrt(float64(v.Size2())))
}

func (v *milliVec3) Normalized() MilliVec3 {
	return v.Div(v.Size())
}

func (v *milliVec3) Cross(p2 MilliVec3) MilliVec3 {
	crossVec := NewMilliVec3(
		v.y*p2.Z()-v.z*p2.Y(),
		v.z*p2.X()-v.x*p2.Z(),
		v.x*p2.Y()-v.y*p2.X(),
	)
	return crossVec
}

func (v *milliVec3) Copy() MilliVec3 {
	return &milliVec3{
		x: v.x,
		y: v.y,
		z: v.z,
	}
}

type MilliMatrix3x3 interface {
	Matrix() [][]Millimeter

	ApplyTo(p MilliVec3) MicroVec3
}

type milliMatrix3x3 struct {
	matrix [][]Millimeter
}

func NewFMatrix3x3() MilliMatrix3x3 {
	m := milliMatrix3x3{matrix: [][]Millimeter{
		{1, 0, 0},
		{0, 1, 0},
		{0, 0, 1},
	}}

	return &m
}

func (m *milliMatrix3x3) Matrix() [][]Millimeter {
	return m.matrix
}

func (m *milliMatrix3x3) ApplyTo(p MilliVec3) MicroVec3 {
	return NewMicroVec3(
		(p.X()*m.matrix[0][0] + p.Y()*m.matrix[1][0] + p.Z()*m.matrix[2][0]).ToMicrometer(),
		(p.X()*m.matrix[0][1] + p.Y()*m.matrix[1][1] + p.Z()*m.matrix[2][1]).ToMicrometer(),
		(p.X()*m.matrix[0][2] + p.Y()*m.matrix[1][2] + p.Z()*m.matrix[2][2]).ToMicrometer())
}
