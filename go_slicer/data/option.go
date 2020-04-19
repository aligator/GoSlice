package data

import "GoSlicer/util"

type Options struct {
	MeldDistance              util.Micrometer
	Center                    util.MicroVec3
	JoinPolygonSnapDistance   util.Micrometer
	FinishPolygonSnapDistance util.Micrometer

	InitialLayerThickness util.Micrometer
	LayerThickness        util.Micrometer
	FilamentDiameter      util.Micrometer
	ExtrusionWidth        util.Micrometer
	InsetCount            int
}
