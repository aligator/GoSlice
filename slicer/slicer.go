package slicer

import (
	"GoSlicer/model"
	"GoSlicer/util"
	"fmt"
)

type Slicer interface {
	LayerParts()
	GenerateGCode()
}

type slicer struct {
	modelSize util.MicroVec3
	layers    []*layer
}

func NewSlicer(om model.OptimizedModel, initialThickness util.Micrometer, layerThickness util.Micrometer) Slicer {
	s := &slicer{}

	s.modelSize = om.Size()
	layerCount := (s.modelSize.Z()-initialThickness)/layerThickness + 1

	fmt.Println("layer count:", layerCount, s.modelSize.Z(), initialThickness, layerThickness)

	s.layers = make([]*layer, layerCount)

	for i, _ := range om.Faces() {
		points := om.FacePoints(i)
		minZ := points[0].Z()
		maxZ := points[0].Z()

		if points[1].Z() < minZ {
			minZ = points[1].Z()
		}
		if points[2].Z() < minZ {
			minZ = points[2].Z()
		}

		if points[1].Z() > maxZ {
			maxZ = points[1].Z()
		}
		if points[2].Z() > maxZ {
			maxZ = points[2].Z()
		}

		// for each layerNr
		for layerNr := int((minZ - initialThickness) / layerThickness); util.Micrometer(layerNr) <= (maxZ-initialThickness)/layerThickness; layerNr++ {
			z := util.Micrometer(layerNr)*layerThickness + initialThickness
			if z < minZ {
				continue
			}
			if layerNr < 0 {
				continue
			}

			if s.layers[layerNr] == nil {
				s.layers[layerNr] = NewLayer()
			}

			layer := s.layers[layerNr]

			var seg *segment
			switch {
			// only p0 is below z
			case points[0].Z() < z && points[1].Z() >= z && points[2].Z() >= z:
				seg = SliceFace(z, points[0], points[2], points[1])
			// only p1 and p2 are below z
			case points[0].Z() > z && points[1].Z() < z && points[2].Z() < z:
				seg = SliceFace(z, points[0], points[1], points[2])

			// only p1 is below z
			case points[1].Z() < z && points[0].Z() >= z && points[2].Z() >= z:
				seg = SliceFace(z, points[1], points[0], points[2])
			// only p0 and p2 are below z
			case points[1].Z() > z && points[0].Z() < z && points[2].Z() < z:
				seg = SliceFace(z, points[1], points[2], points[0])

			// only p2 is below z
			case points[2].Z() < z && points[1].Z() >= z && points[0].Z() >= z:
				seg = SliceFace(z, points[2], points[1], points[0])

			// only p1 and p0 are below z
			case points[2].Z() > z && points[1].Z() < z && points[0].Z() < z:
				seg = SliceFace(z, points[2], points[0], points[1])
			default:
				// not all cases create a segment, because
				// a point of a face could create just a dot
				// and if all points are below or above no face has to be created
				continue
			}

			layer.faceToSegmentIndex[i] = len(layer.segments)
			seg.faceIndex = i
			seg.addedToPolygon = false
			layer.segments = append(layer.segments, seg)
		}
	}

	for _, layer := range s.layers {
		layer.makePolygons(om)
	}
	return s
}

func (s *slicer) LayerParts() {
	panic("implement me")
}

func (s *slicer) GenerateGCode() {
	panic("implement me")
}
