package slicer

import "GoSlicer/util"

type Polygon interface {
	addPoint(point util.MicroPoint)
	removeLastPoint()
	Points() []util.MicroPoint
	close()
	isClosed() bool
	isAlmostFinished(snapDistance util.Micrometer) bool
}

type polygon struct {
	points []util.MicroPoint
	closed bool
}

func (p *polygon) addPoint(point util.MicroPoint) {
	p.points = append(p.points, point)
}

func (p *polygon) Points() []util.MicroPoint {
	return p.points
}

func (p *polygon) removeLastPoint() {
	p.points = p.points[:len(p.points)]
}

func (p *polygon) close() {
	p.closed = true
}

func (p *polygon) isClosed() bool {
	return p.closed
}

func (p *polygon) isAlmostFinished(snapDistance util.Micrometer) bool {
	return p.points[0].Sub(p.points[len(p.points)-1]).ShorterThan(snapDistance)
}
