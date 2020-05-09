package renderer

import (
	"GoSlice/data"
	"GoSlice/gcode/builder"
)

type PreLayer struct{}

func (PreLayer) Init(model data.OptimizedModel) {}

func (PreLayer) Render(builder builder.Builder, layerNr int, layers []data.PartitionedLayer, z data.Micrometer, options *data.Options) {
	builder.AddComment("LAYER:%v", layerNr)
	if layerNr == 0 {
		// starting gcode
		builder.AddComment("START_GCODE")
		builder.AddComment("Generated with GoSlice")
		builder.AddCommand("G1 X0 Y20 Z0.2 F3000 ; get ready to prime")
		builder.AddCommand("G92 E0 ; reset extrusion distance")
		builder.AddCommand("G1 X200 E20 F600 ; prime nozzle")
		builder.AddCommand("G1 Z5 F5000 ; lift nozzle")
		builder.AddCommand("G92 E0 ; reset extrusion distance")

		builder.SetExtrusion(options.Print.InitialLayerThickness, options.Printer.ExtrusionWidth, options.Filament.FilamentDiameter)

		// force the InitialLayerSpeed for first layer
		builder.SetExtrudeSpeedOverride(options.Print.IntialLayerSpeed)
	} else {
		builder.DisableExtrudeSpeedOverride()
		builder.SetExtrudeSpeed(options.Print.LayerSpeed)
	}

	if layerNr == 2 {
		builder.AddCommand("M106 ; enable fan")
	}
}

type PostLayer struct{}

func (PostLayer) Init(model data.OptimizedModel) {}

func (PostLayer) Render(builder builder.Builder, layerNr int, layers []data.PartitionedLayer, z data.Micrometer, options *data.Options) {
	// ending gcode
	if layerNr == len(layers)-1 {
		builder.AddComment("END_GCODE")
		builder.SetExtrusion(options.Print.LayerThickness, options.Printer.ExtrusionWidth, options.Filament.FilamentDiameter)
		builder.AddCommand("M107 ; disable fan")
	}
}
