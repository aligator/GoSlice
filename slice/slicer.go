package slice

import (
	"GoSlice/data"
	"GoSlice/handle"
	"errors"
	"fmt"
)

type slicer struct {
	options *data.Options
}

func NewSlicer(options *data.Options) handle.ModelSlicer {
	return &slicer{options: options}
}

func (s slicer) Slice(m data.OptimizedModel) ([]data.PartitionedLayer, error) {
	layerCount := (m.Size().Z()-s.options.Print.InitialLayerThickness)/s.options.Print.LayerThickness + 1
	fmt.Println("layer count:", layerCount)

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
	for i, layer := range layers {
		layer.makePolygons(m, s.options.JoinPolygonSnapDistance, s.options.FinishPolygonSnapDistance)
		lp, ok := layer.generateLayerParts()

		if !ok {
			return nil, errors.New(fmt.Sprintf("partitioning failed at layer %v", i))
		}

		retLayers[i] = lp
	}

	return retLayers, nil
}
