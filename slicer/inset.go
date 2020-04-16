package slicer

import (
	"GoSlicer/util"
	clipper "github.com/ctessum/go.clipper"
	"math"
)

type insetPart struct {
	inset []clipper.Paths
}

func newInsetPart(part *layerPart, offset util.Micrometer, insetCount int) *insetPart {
	in := &insetPart{}

	for i := 0; i < insetCount; i++ {
		in.inset = append(in.inset, clipper.NewPaths())
		o := clipper.NewClipperOffset()
		o.AddPaths(part.polygons, clipper.JtRound, clipper.EtClosedPolygon)
		o.MiterLimit = 2
		in.inset[i] = o.Execute(float64(-int(offset)*i) - float64(offset/2))
		if len(in.inset[i]) < 1 {
			break
		}

		for j, path := range in.inset[i] {
			var filteredPath []*clipper.IntPoint

			for k, point := range path {
				if k == 0 || k == len(path)-1 {
					filteredPath = append(filteredPath, point)
					continue
				}

				// TODO: Experimental: remove too small segments from path
				prev := path[k-1]
				if math.Sqrt(float64((prev.Y-point.Y)*(prev.Y-point.Y)+(point.X-prev.X)*(point.X-prev.X))) > 30.0 {
					filteredPath = append(filteredPath, point)
				}
			}
			in.inset[i][j] = filteredPath
		}
	}

	return in
}
