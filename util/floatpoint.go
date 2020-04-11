package util

import "math"

type Point3 struct {
	X, Y, Z coord
}

func NewPoint3(x int, y int, z int) *Point3 {
	return &Point3{
		X: coord(x),
		Y: coord(y),
		Z: coord(z),
	}
}

func (p *Point3) ToFPoint3() *FPoint3 {
	return NewFPoint3(float64(p.X)*0.001, float64(p.Y)*0.001, float64(p.Z)*0.001)
}

type FPoint3 struct {
	X, Y, Z float64
}

func NewFPoint3(x float64, y float64, z float64) *FPoint3 {
	return &FPoint3{
		X: x,
		Y: y,
		Z: z,
	}
}

func AddFPoint3(p1, p2 FPoint3) FPoint3 {
	return *p1.Add(p2)
}

func SubFPoint3(p1, p2 FPoint3) FPoint3 {
	return *p1.Sub(p2)
}

func MulFPoint3(p FPoint3, f float64) FPoint3 {
	return *p.Mul(f)
}

func CrossFPoint3(p1 FPoint3, p2 FPoint3) FPoint3 {
	return p1.Cross(p2)
}

func DivFPoint3(p FPoint3, f float64) FPoint3 {
	p.X /= f
	p.Y /= f
	p.Z /= f
	return p
}

func (p *FPoint3) Add(p2 FPoint3) *FPoint3 {
	p.X += p2.X
	p.Y += p2.Y
	p.Z += p2.Z
	return p
}

func (p *FPoint3) Sub(p2 FPoint3) *FPoint3 {
	p.X -= p2.X
	p.Y -= p2.Y
	p.Z -= p2.Z
	return p
}

func (p *FPoint3) Mul(f float64) *FPoint3 {
	p.X *= f
	p.Y *= f
	p.Z *= f
	return p
}

func (p *FPoint3) Equals(p2 FPoint3) bool {
	return *p == p2
}

func (p *FPoint3) Max() float64 {
	if p.X > p.Y && p.X > p.Z {
		return p.X
	}
	if p.Y > p.Z {
		return p.Y
	}
	return p.Z
}

func (p *FPoint3) TestLength(length float64) bool {
	return p.VSize2() <= length*length
}

func (p *FPoint3) VSize2() float64 {
	return p.X*p.X + p.Y*p.Y + p.Z*p.Z
}

func (p *FPoint3) VSize() float64 {
	return math.Sqrt(p.VSize2())
}

func (p *FPoint3) Normalized() FPoint3 {
	return DivFPoint3(*p, p.VSize())
}

func (p *FPoint3) Cross(p2 FPoint3) FPoint3 {
	crossPoint := NewFPoint3(
		p.Y*p2.Z-p.Z*p2.Y,
		p.Z*p2.X-p.X*p2.Z,
		p.X*p2.Y-p.Y*p2.X,
	)
	return *crossPoint
}

func (p *FPoint3) Point3() Point3 {
	return *NewPoint3(int(p.X)*1000, int(p.Y)*1000, int(p.Z)*1000)
}

type FMatrix3x3 struct {
	matrix [][]float64
}

func NewFMatrix3x3() FMatrix3x3 {
	m := FMatrix3x3{matrix: [][]float64{
		{1, 0, 0},
		{0, 1, 0},
		{0, 0, 1},
	}}

	return m
}

func (m *FMatrix3x3) apply(p FPoint3) Point3 {
	return Point3{
		X: mmToInt(p.X*m.matrix[0][0] + p.Y*m.matrix[1][0] + p.Z*m.matrix[2][0]),
		Y: mmToInt(p.X*m.matrix[0][1] + p.Y*m.matrix[1][1] + p.Z*m.matrix[2][1]),
		Z: mmToInt(p.X*m.matrix[0][2] + p.Y*m.matrix[1][2] + p.Z*m.matrix[2][2]),
	}
}
