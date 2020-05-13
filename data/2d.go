// This file provides some basic helper functions for 2d vector calculations.

package data

// DotProduct calculates the dot product of two points
func DotProduct(a, b MicroPoint) Micrometer {
	return a.X()*b.X() + a.Y()*b.Y()
}

// XDistance2ToLine calculates the X-Distance^2 of a point to a line
func XDistance2ToLine(a, b, point MicroPoint) Micrometer {
	//  x.......a------------b
	//  :
	//  :
	//  p
	// return px_size
	vecAB := b.Sub(a)
	vecAP := point.Sub(a)

	if vecAB.Size2() == 0 {
		return vecAP.Size2() // assume a perpendicular line to p
	}

	dist := DotProduct(vecAB, vecAP)
	axSize2 := dist * dist / vecAB.Size2()
	return Max(0, vecAP.Size2()-axSize2)
}
