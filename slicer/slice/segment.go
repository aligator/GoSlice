package slice

import (
	"GoSlicer/util"
)

type segment struct {
	start, end     util.MicroPoint
	faceIndex      int
	addedToPolygon bool
}

func SliceFace(z util.Micrometer, p0, p1, p2 util.MicroVec3) *segment {
	seg := &segment{
		start: util.NewMicroPoint(
			p0.X()+(p1.X()-p0.X())*(z-p0.Z())/(p1.Z()-p0.Z()),
			p0.Y()+(p1.Y()-p0.Y())*(z-p0.Z())/(p1.Z()-p0.Z()),
		),

		end: util.NewMicroPoint(
			p0.X()+(p2.X()-p0.X())*(z-p0.Z())/(p2.Z()-p0.Z()),
			p0.Y()+(p2.Y()-p0.Y())*(z-p0.Z())/(p2.Z()-p0.Z()),
		),
	}
	return seg
}
