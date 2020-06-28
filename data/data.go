// Package data holds basic data structures and interfaces used by GoSlice.
package data

import (
	"math"
)

// Millimeter represents a value in mm
// It should not be used for calculations, convert to micrometer.
// using ToMicrometer() before calculating to prevent rounding
// errors because of the float-type.
type Millimeter float32

func (m Millimeter) ToMicrometer() Micrometer {
	return Micrometer(math.RoundToEven(float64(m * 1000)))
}

// Micrometer represents a value in 0.001 mm
type Micrometer int64

const MaxMicrometer = Micrometer(math.MaxInt64)
const MinMicrometer = Micrometer(math.MinInt64)

func (m Micrometer) ToMillimeter() Millimeter {
	return Millimeter(float64(m) / 1000)
}
