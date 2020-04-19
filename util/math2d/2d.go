package math2d

import (
	"GoSlice/util"
)

func PointsDistance(a, b util.MicroPoint) util.Micrometer {
	return a.X()*b.X() + a.Y()*b.Y()
}

func XDistance2ToLine(a, b, point util.MicroPoint) util.Micrometer {
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

	dist := PointsDistance(vecAB, vecAP)
	axSize2 := dist * dist / vecAB.Size2()
	return util.Max(0, vecAP.Size2()-axSize2)
}
