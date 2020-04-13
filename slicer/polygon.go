package slicer

import "GoSlicer/util"

type polygon struct {
	points []util.MicroPoint
	closed bool
}

func (p *polygon) removeLastPoint() {
	p.points = p.points[:len(p.points)]
}

func (p *polygon) isAlmostFinished(snapDistance util.Micrometer) bool {
	return p.points[0].Sub(p.points[len(p.points)-1]).ShorterThan(snapDistance)
}
