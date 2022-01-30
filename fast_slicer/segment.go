package fast_slicer

import (
	"github.com/aligator/goslice/data"
)

// segment is a line specified by two points.
type segment struct {
	start, end data.MicroPoint

	// faceIndex is the face which this segment belongs to.
	faceIndex      int
	addedToPolygon bool
}

// SliceFace generates a 2d slice out of a triangle at a specific z.
// The triangle is defined by the three points.
func SliceFace(z data.Micrometer, p0, p1, p2 data.MicroVec3) *segment {
	seg := &segment{
		start: data.NewMicroPoint(
			p0.X()+(p1.X()-p0.X())*(z-p0.Z())/(p1.Z()-p0.Z()),
			p0.Y()+(p1.Y()-p0.Y())*(z-p0.Z())/(p1.Z()-p0.Z()),
		),

		end: data.NewMicroPoint(
			p0.X()+(p2.X()-p0.X())*(z-p0.Z())/(p2.Z()-p0.Z()),
			p0.Y()+(p2.Y()-p0.Y())*(z-p0.Z())/(p2.Z()-p0.Z()),
		),
	}
	return seg
}
