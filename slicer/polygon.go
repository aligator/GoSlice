package slicer

import (
	"GoSlicer/util"
)

// slicePolygon is a internal polygon used while slicing
type slicePolygon struct {
	points []util.MicroPoint
	closed bool
}

func (p *slicePolygon) removeLastPoint() {
	p.points = p.points[:len(p.points)]
}

func (p *slicePolygon) isAlmostFinished(snapDistance util.Micrometer) bool {
	return p.points[0].Sub(p.points[len(p.points)-1]).ShorterThan(snapDistance)
}
