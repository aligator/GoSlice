package clip

import (
	"GoSlicer/slicer/data"
	"GoSlicer/util"
	clipper "github.com/ctessum/go.clipper"
	"math"
)

type Clip interface {
	// GenerateLayerParts partitions the whole layer into several partition parts
	GenerateLayerParts(l data.Layer) (data.PartitionedLayer, bool)
	InsetLayer(layer data.PartitionedLayer, offset util.Micrometer, insetCount int) [][]data.Paths
	Inset(part data.LayerPart, offset util.Micrometer, insetCount int) []data.Paths
}

// clipperClip implements Clip using the external clipper library
type clipperClip struct{}

func NewClip() Clip {
	return clipperClip{}
}

type layerPart struct {
	// clipperPolygons are the polys set by clipper
	clipperPolygons clipper.Paths

	// clipPolygons holds the lazy converted polygons.
	// When calling Polygons the first time all clipperPolygons
	// are converted to MicroPoint polygons.
	clipPolygons data.Paths
}

func (l *layerPart) Polygons() data.Paths {
	if l.clipPolygons != nil {
		return l.clipPolygons
	}

	result := data.Paths{}

	for _, poly := range l.clipperPolygons {
		newPath := data.Path{}
		for _, point := range poly {
			newPath = append(newPath, microPoint(point))
		}
		result = append(result, newPath)
	}

	l.clipPolygons = result
	l.clipperPolygons = nil
	return l.clipPolygons
}

type partitionedLayer struct {
	parts    []data.LayerPart
	children []data.PartitionedLayer
}

func (p partitionedLayer) LayerParts() []data.LayerPart {
	return p.parts
}

func clipperPoint(p util.MicroPoint) *clipper.IntPoint {
	return &clipper.IntPoint{
		X: clipper.CInt(p.X()),
		Y: clipper.CInt(p.Y()),
	}
}

func clipperPaths(p data.Paths) clipper.Paths {
	var result clipper.Paths
	for _, path := range p {
		var newPath clipper.Path
		for _, point := range path {
			newPath = append(newPath, clipperPoint(point))
		}
		result = append(result, newPath)
	}

	return result
}

func microPoint(p *clipper.IntPoint) util.MicroPoint {
	return util.NewMicroPoint(util.Micrometer(p.X), util.Micrometer(p.Y))
}

func (c clipperClip) GenerateLayerParts(l data.Layer) (data.PartitionedLayer, bool) {
	polyList := clipper.Paths{}
	// convert all polygons to clipper polygons
	for _, layerPolygon := range l.Polygons() {
		var path = clipper.Path{}

		prev := 0
		// convert all points of this polygons
		for j, layerPoint := range layerPolygon {
			// ignore first as the next check would fail otherwise
			if j == 1 {
				path = append(path, clipperPoint(layerPolygon[0]))
				continue
			}

			// filter too near points
			// check this always with the previous point
			if layerPoint.Sub(layerPolygon[prev]).ShorterThan(200) {
				continue
			}

			path = append(path, clipperPoint(layerPoint))
			prev = j
		}

		polyList = append(polyList, path)
	}

	layer := partitionedLayer{}

	clip := clipper.NewClipper(clipper.IoNone)
	clip.AddPaths(polyList, clipper.PtSubject, true)
	resultPolys, ok := clip.Execute2(clipper.CtUnion, clipper.PftEvenOdd, clipper.PftEvenOdd)
	if !ok {
		return nil, false
	}

	for _, p := range resultPolys.Childs() {
		part := layerPart{}
		part.clipperPolygons = append(part.clipperPolygons, p.Contour())
		for _, child := range p.Childs() {
			part.clipperPolygons = append(part.clipperPolygons, child.Contour())
		}
		layer.parts = append(layer.parts, &part)
	}

	return layer, true
}

func (c clipperClip) InsetLayer(layer data.PartitionedLayer, offset util.Micrometer, insetCount int) [][]data.Paths {
	var result [][]data.Paths
	for _, part := range layer.LayerParts() {
		result = append(result, c.Inset(part, offset, insetCount))
	}

	return result
}

func (c clipperClip) Inset(part data.LayerPart, offset util.Micrometer, insetCount int) []data.Paths {
	var clipperInsets []clipper.Paths
	var insets []data.Paths

	for i := 0; i < insetCount; i++ {
		clipperInsets = append(clipperInsets, clipper.NewPaths())
		o := clipper.NewClipperOffset()
		o.AddPaths(clipperPaths(part.Polygons()), clipper.JtRound, clipper.EtClosedPolygon)
		o.MiterLimit = 2
		clipperInsets[i] = o.Execute(float64(-int(offset)*i) - float64(offset/2))
		if len(clipperInsets[i]) < 1 {
			break
		}

		insets = append(insets, data.Paths{})

		for _, path := range clipperInsets[i] {
			var filteredPath []util.MicroPoint

			for k, point := range path {
				if k == 0 || k == len(path)-1 {
					filteredPath = append(filteredPath, microPoint(point))
					continue
				}

				// TODO: Experimental: remove too small segments from path
				prev := filteredPath[len(filteredPath)-1]
				if math.Sqrt(float64((clipper.CInt(prev.Y())-point.Y)*(clipper.CInt(prev.Y())-point.Y)+(point.X-clipper.CInt(prev.X()))*(point.X-clipper.CInt(prev.X())))) > 100.0 {
					filteredPath = append(filteredPath, microPoint(point))
				}
			}
			insets[i] = append(insets[i], filteredPath)
		}
	}

	return insets
}
