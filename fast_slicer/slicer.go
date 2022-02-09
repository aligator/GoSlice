package fast_slicer

import (
	"fmt"
	"sort"

	"github.com/aligator/goslice/clip"
	"github.com/aligator/goslice/data"
	"github.com/aligator/goslice/handler"
)

type Edge struct {
	Point1 data.MicroVec3
	Point2 data.MicroVec3
}

type IntersectionEdges struct {
	ForwardEdge  *Edge
	BackwardEdge *Edge
}

// Intersection is used to store the information related
// to the intersection between the triangle and plane.
// (Section 4.1.1)
type Intersection struct {
	Pre *Intersection

	// IntersectionVertice between the forwardEdge and the current
	// slicing Plane.
	IntersectionVertice data.MicroVec3
	Edge                *IntersectionEdges

	Next *Intersection
}

type slicer struct {
	options *data.Options
}

// NewSlicer provides the built in slicer implementation.
func NewSlicer(options *data.Options) handler.ModelSlicer {
	return &slicer{options: options}
}

func (s slicer) Slice(m data.OptimizedModel) ([]data.PartitionedLayer, error) {
	layerCount := (m.Size().Z()-s.options.Print.InitialLayerThickness)/s.options.Print.LayerThickness + 1

	layers := make([]*layer, layerCount)

	for i := 0; i < m.FaceCount(); i++ {
		currentFace := m.Face(i)
		currentOptimizedFace := m.OptimizedFace(i)

		minZ := currentOptimizedFace.MinZ()
		maxZ := currentOptimizedFace.MaxZ()

		if minZ == maxZ {
			continue
		}

		points := currentFace.Points()

		// vertices of the face sorted by min, med, max Z
		vertices := points[:]
		sort.Slice(vertices, func(i, j int) bool {
			return vertices[i].Z() < vertices[j].Z()
		})

		slices := make([]int, 3)
		for i, v := range vertices {
			// Get the number of the slice where the vertice lies in.
			// Take into account the thickness of the first layer and count +1 at the end for the ignored layer.
			if v.Z() <= s.options.Print.InitialLayerThickness {
				continue // leave the slice number = 0
			}

			// (Section 4.2)
			// sj = (vj−szmin) / t
			//
			// where sj stands for s_min,s_med,s_max; vj stands for vzmin,
			// vzmed,vzmax; szmin denotes the z coordinate of the 0th slice
			// which can be designed to coincide to the x-y plane of the
			// system; and t denotes the slice thickness.

			// Assume layer 0 to be on the bottom plane. -> can be ignored
			slices[i] = int((v.Z()-s.options.Print.InitialLayerThickness)/s.options.Print.LayerThickness + 1)
		}

		// (Section 4.3)
		// Judge forward and backwardEdge for {vecZMin, vecZMax} and {vecZMin, vecZMed}

		// Judge forward and backwardEdge for {vecZMin, vecZMax} and {vecZMed, vecZMax}
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
