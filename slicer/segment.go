package slicer

import (
	"GoSlicer/util"
)

type Segment interface {
	setFaceIndex(index int)
	FaceIndex() int
	setAddedToPolygon(value bool)
	isAddedToPolygon() bool
	Start() util.MicroPoint
	End() util.MicroPoint
}

type segment struct {
	start, end     util.MicroPoint
	faceIndex      int
	addedToPolygon bool
}

func SliceFace(z util.Micrometer, p0, p1, p2 util.MicroVec3) Segment {
	seg := &segment{
		start: util.NewMicroPoint(
			p0.X()+(p1.X()-p0.X())*(z-p0.Z()/(p1.Z()-p0.Z())),
			p0.Y()+(p1.Y()-p0.Y())*(z-p0.Z()/(p1.Z()-p0.Z()))),
		end: util.NewMicroPoint(
			p0.X()+(p2.X()-p0.X())*(z-p0.Z()/(p2.Z()-p0.Z())),
			p0.Y()+(p2.Y()-p0.Y())*(z-p0.Z()/(p2.Z()-p0.Z()))),
	}
	return seg
}

func (s *segment) setFaceIndex(index int) {
	s.faceIndex = index
}

func (s *segment) FaceIndex() int {
	return s.faceIndex
}

func (s *segment) setAddedToPolygon(value bool) {
	s.addedToPolygon = value
}

func (s *segment) isAddedToPolygon() bool {
	return s.addedToPolygon
}

func (s *segment) Start() util.MicroPoint {
	return s.start
}

func (s *segment) End() util.MicroPoint {
	return s.end
}
