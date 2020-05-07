// Package data holds basic data structures and interfaces used by GoSlice.
package data

import (
	"math"
)

// MicroVec3 represents a point in 3d space
// which is in a Micrometer-grid.
// A value of 1 represents 0.001 mm.
// Micro vectors are used as soon as possible to avoid rounding errors
// because the Micrometer datatype uses integers.
type MicroVec3 interface {
	X() Micrometer
	Y() Micrometer
	Z() Micrometer

	SetX(x Micrometer)
	SetY(y Micrometer)
	SetZ(z Micrometer)

	// PointXY returns a new point only with the x and y coordinates of the vector.
	PointXY() MicroPoint

	// Add returns a new vector which is the sum of the vectors. (this + vec)
	//
	// By convention it should never mutate the instance and instead return a new copy.
	Add(vec MicroVec3) MicroVec3

	// Sub returns a new vector which is the difference of the vectors. (this - vec)
	//
	// By convention it should never mutate the instance and instead return a new copy.
	Sub(vec MicroVec3) MicroVec3

	// Mul returns a new vector which is the multiplication by the given value. (this * value)
	//
	// By convention it should never mutate the instance and instead return a new copy.
	Mul(value Micrometer) MicroVec3

	// Div returns a new vector which is the division by the given value. (this / value)
	//
	// By convention it should never mutate the instance and instead return a new copy.
	Div(value Micrometer) MicroVec3

	Max() Micrometer

	// TestLength checks if the length of the vector fits inside the given length.
	// Returns true if the vector length is <= the given length.
	TestLength(length Micrometer) bool

	// Size2 returns the length of the vector^2.
	//
	// Use this whenever possible as it may be faster than Bounds().
	Size2() Micrometer

	// Size2 returns the length of the vector.
	//
	// Use Size2() this whenever possible as it may be faster than Bounds().
	Size() Micrometer

	Normalized() MicroVec3
	Cross(p2 MicroVec3) MicroVec3

	// Copy returns a completely new copy of the vector.
	Copy() MicroVec3

	// String implements the value interface needed for the options.
	String() string
	// Set implements the value interface needed for the options.
	Set(s string) error
	// Type implements the value interface needed for the options.
	Type() string
}

// microVec implements MicroVec3
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

func Max(a, b Micrometer) Micrometer {
	if a > b {
		return a
	}
	return b
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

func (v *microVec3) PointXY() MicroPoint {
	return &microPoint{
		x: v.x,
		y: v.y,
	}
}

func (v *microVec3) Add(vec MicroVec3) MicroVec3 {
	result := v.Copy()
	result.SetX(result.X() + vec.X())
	result.SetY(result.Y() + vec.Y())
	result.SetZ(result.Z() + vec.Z())
	return result
}

func (v *microVec3) Sub(vec MicroVec3) MicroVec3 {
	result := v.Copy()
	result.SetX(result.X() - vec.X())
	result.SetY(result.Y() - vec.Y())
	result.SetZ(result.Z() - vec.Z())
	return result
}

func (v *microVec3) Mul(value Micrometer) MicroVec3 {
	result := v.Copy()
	result.SetX(result.X() * value)
	result.SetY(result.Y() * value)
	result.SetZ(result.Z() * value)
	return result
}

func (v *microVec3) Div(value Micrometer) MicroVec3 {
	result := v.Copy()
	result.SetX(result.X() / value)
	result.SetY(result.Y() / value)
	result.SetZ(result.Z() / value)
	return result
}

func (v *microVec3) Max() Micrometer {
	if v.x > v.y && v.x > v.z {
		return v.x
	}
	if v.y > v.z {
		return v.y
	}
	return v.z
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

func (v *microVec3) Normalized() MicroVec3 {
	return v.Div(v.Size())
}

func (v *microVec3) Cross(p2 MicroVec3) MicroVec3 {
	crossVec := NewMicroVec3(
		v.y*p2.Z()-v.z*p2.Y(),
		v.z*p2.X()-v.x*p2.Z(),
		v.x*p2.Y()-v.y*p2.X(),
	)
	return crossVec
}

func (v *microVec3) Copy() MicroVec3 {
	return &microVec3{
		x: v.x,
		y: v.y,
		z: v.z,
	}
}

// MicroPoint represents a point in 2d space
// which is in a Micrometer-grid.
// A value of 1 represents 0.001 mm.
type MicroPoint interface {
	X() Micrometer
	Y() Micrometer

	SetX(x Micrometer)
	SetY(y Micrometer)

	// Add returns a new vector which is the sum of the vectors. (this + vec)
	//
	// By convention it should never mutate the instance and instead return a new copy.
	Add(vec MicroPoint) MicroPoint

	// Sub returns a new vector which is the difference of the vectors. (this - vec)
	//
	// By convention it should never mutate the instance and instead return a new copy.
	Sub(vec MicroPoint) MicroPoint

	// Mul returns a new vector which is the multiplication by the given value. (this * value)
	//
	// By convention it should never mutate the instance and instead return a new copy.
	Mul(value Micrometer) MicroPoint

	// Div returns a new vector which is the division by the given value. (this / value)
	//
	// By convention it should never mutate the instance and instead return a new copy.
	Div(value Micrometer) MicroPoint

	// ShorterThanOrEqual checks if the length of the vector fits inside the given length.
	// Returns true if the vector length is <= the given length.
	ShorterThanOrEqual(length Micrometer) bool

	// Size2 returns the length of the vector^2.
	//
	// Use this whenever possible as it may be faster than Bounds().
	Size2() Micrometer

	// Size2 returns the length of the vector.
	//
	// Use Size2() this whenever possible as it may be faster than Bounds().
	Size() Micrometer

	// SizeMM returns the length of the vector in mm.
	SizeMM() Millimeter

	// Copy returns a completely new copy of the vector.
	Copy() MicroPoint
}

// microPoint implements MicroPoint
type microPoint struct {
	x, y Micrometer
}

func NewMicroPoint(x, y Micrometer) MicroPoint {
	return &microPoint{
		x, y,
	}
}

func (p *microPoint) X() Micrometer {
	return p.x
}

func (p *microPoint) Y() Micrometer {
	return p.y
}

func (p *microPoint) SetX(x Micrometer) {
	p.x = x
}

func (p *microPoint) SetY(y Micrometer) {
	p.y = y
}

func (p *microPoint) Add(p2 MicroPoint) MicroPoint {
	result := p.Copy()
	result.SetX(result.X() + p2.X())
	result.SetY(result.Y() + p2.Y())
	return result
}

func (p *microPoint) Sub(p2 MicroPoint) MicroPoint {
	result := p.Copy()
	result.SetX(result.X() - p2.X())
	result.SetY(result.Y() - p2.Y())
	return result
}

func (p *microPoint) Mul(value Micrometer) MicroPoint {
	result := p.Copy()
	result.SetX(result.X() * value)
	result.SetY(result.Y() * value)
	return result
}

func (p *microPoint) Div(value Micrometer) MicroPoint {
	result := p.Copy()
	result.SetX(result.X() / value)
	result.SetY(result.Y() / value)
	return result
}

// ShorterThanOrEqual just checks if the given length is smaller than the
// length of the vector to this point.
// This implementation first tries a more performant way before actually calculating the suize
func (p *microPoint) ShorterThanOrEqual(length Micrometer) bool {
	if p.x > length || p.x < -length ||
		p.y > length || p.y < -length {
		return false
	}

	return p.Size2() <= length*length
}

func (p *microPoint) Size2() Micrometer {
	return p.x*p.x + p.y*p.y
}

func (p *microPoint) Size() Micrometer {
	return Micrometer(math.Sqrt(float64(p.Size2())))
}

func (p *microPoint) SizeMM() Millimeter {
	x := p.x.ToMillimeter()
	y := p.y.ToMillimeter()
	return Millimeter(math.Sqrt(float64(x*x + y*y)))
}

func (p *microPoint) Copy() MicroPoint {
	return &microPoint{
		p.x, p.y,
	}
}
