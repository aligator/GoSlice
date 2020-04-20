package data

import "GoSlice/util"

type PrintOptions struct {
	IntialLayerSpeed    util.Millimeter
	LayerSpeed          util.Millimeter
	OuterPerimeterSpeed util.Millimeter

	InitialLayerThickness util.Micrometer
	LayerThickness        util.Micrometer
	InsetCount            int

	InfillOverlapPercent int
}

type FilamentOptions struct {
	FilamentDiameter util.Micrometer
}

type PrinterOptions struct {
	ExtrusionWidth util.Micrometer
}

type Options struct {
	Printer  PrinterOptions
	Filament FilamentOptions
	Print    PrintOptions

	MeldDistance              util.Micrometer
	Center                    util.MicroVec3
	JoinPolygonSnapDistance   util.Micrometer
	FinishPolygonSnapDistance util.Micrometer
}
