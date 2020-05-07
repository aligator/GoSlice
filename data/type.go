// Package data holds basic data structures and interfaces used by GoSlice.
package data

import (
	"math"
)

// infill.go holds the basic type definition.
// The types behind them may change in the
// future to more precise / larger ones.

// Millimeter represents a value in mm
type Millimeter float32

func (m Millimeter) ToMicrometer() Micrometer {
	return Micrometer(m * 1000)
}

// Micrometer represents a value in 0.001 mm
type Micrometer int64

const MaxMicrometer = Micrometer(math.MaxInt64)
const MinMicrometer = Micrometer(math.MinInt64)

func (m Micrometer) ToMillimeter() Millimeter {
	return Millimeter(m) * 0.001
}
