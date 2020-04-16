package slicer

import (
	"GoSlicer/util"
	clipper "github.com/ctessum/go.clipper"
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
	}

	return in
}
