// Package slicer provides an implementation for slicing a model into 2d slices.
//
// How it works:
// For each face (always a triangle) it first fetches the min and max height (z) and then the face is sliced at each
// needed layer height (taking into account the optionally different initial layer thickness).
// For this it first determines which of the three points is below or above the current z height and then based on this
// calls the SliceFace function which simply returns a segment which is one line (2 points) representing the slice of the triangle
// at exactly the current height.
// The segments are saved for each layer.
// At the end it loops through all generated layers and
// - creates polygons out of the segments (see documentation of the makePolygons method to learn how)
// - generates the layer parts out of the polygons. This means it groups them together and calculates which polygons
//   just represents holes of other polygons.

package slicer

import (
	"fmt"
	"github.com/aligator/goslice/clip"
	"github.com/aligator/goslice/data"
	"github.com/aligator/goslice/handler"
)

type slicer struct {
	options *data.Options
}

// NewSlicer provides the built in slicer implementation.
func NewSlicer(options *data.Options) handler.ModelSlicer {
	return &slicer{options: options}
}

func (s slicer) Slice(m data.OptimizedModel) ([]data.PartitionedLayer, error) {
	layerCount := (m.Size().Z()-s.options.Print.InitialLayerThickness)/s.options.Print.LayerThickness + 1
	s.options.GoSlice.Logger.Println("layer count:", layerCount)

	layers := make([]*layer, layerCount)

	for i := 0; i < m.FaceCount(); i++ {
		points := m.Face(i).Points()
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
		for layerNr := int((minZ - s.options.Print.InitialLayerThickness) / s.options.Print.LayerThickness); data.Micrometer(layerNr) <= (maxZ-s.options.Print.InitialLayerThickness)/s.options.Print.LayerThickness; layerNr++ {
			z := data.Micrometer(layerNr)*s.options.Print.LayerThickness + s.options.Print.InitialLayerThickness
			if z < minZ {
				continue
			}
			if layerNr < 0 {
				continue
			}

			if layers[layerNr] == nil {
				layers[layerNr] = newLayer(layerNr, s.options)
			}

			layer := layers[layerNr]

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
				// and if all paths are below or above no face has to be created
				continue
			}

			layer.faceToSegmentIndex[i] = len(layer.segments)
			seg.faceIndex = i
			seg.addedToPolygon = false
			layer.segments = append(layer.segments, seg)
		}
	}

	retLayers := make([]data.PartitionedLayer, len(layers))
	c := clip.NewClipper()

	for i, layer := range layers {
		layer.makePolygons(m, s.options.Slicing.JoinPolygonSnapDistance, s.options.Slicing.FinishPolygonSnapDistance)
		lp, ok := c.GenerateLayerParts(layer)

		if !ok {
			return nil, fmt.Errorf("partitioning failed at layer %v", i)
		}

		retLayers[i] = lp
	}

	return retLayers, nil
}
