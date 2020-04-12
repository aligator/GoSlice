package util

// type.go holds the basic type definition
// the types behind them may change in the
// future to more precise / larger ones

// Millimeter represents a value in mm
type Millimeter float32

func (m Millimeter) ToMicrometer() Micrometer {
	return Micrometer(m * 1000)
}

type Millimeter2 float32

func (m Millimeter2) ToMicrometer() Micrometer {
	return Micrometer(m * 1000000)
}

type Millimeter3 float32

func (m Millimeter3) ToMicrometer() Micrometer {
	return Micrometer(m * 1000000000)
}

// Micrometer represents a value in 0.001 mm
type Micrometer int

func (m Micrometer) ToMillimeter() Millimeter {
	return Millimeter(m) * 0.001
}

func (m Micrometer) ToMillimeter2() Millimeter2 {
	return Millimeter2(m) * 0.000001
}

func (m Micrometer) ToMillimeter3() Millimeter3 {
	return Millimeter3(m) * 0.000000001
}
