// Package data holds basic data structures and interfaces used by GoSlice.
package data

// PrintOptions contains all Print specific GoSlice options.
type PrintOptions struct {
	IntialLayerSpeed    Millimeter
	LayerSpeed          Millimeter
	OuterPerimeterSpeed Millimeter

	InitialLayerThickness Micrometer
	LayerThickness        Micrometer
	InsetCount            int

	InfillOverlapPercent int
}

// FilamentOptions contains all Filament specific GoSlice options.
type FilamentOptions struct {
	FilamentDiameter Micrometer
}

// PrinterOptions contains all Printer specific GoSlice options.
type PrinterOptions struct {
	ExtrusionWidth Micrometer
	Center         MicroVec3
}

// Options contains all GoSlice options.
type Options struct {
	Printer  PrinterOptions
	Filament FilamentOptions
	Print    PrintOptions

	MeldDistance              Micrometer
	JoinPolygonSnapDistance   Micrometer
	FinishPolygonSnapDistance Micrometer
}
