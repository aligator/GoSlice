// This file provides renderers for gcode injected at specific layers.

package renderer

import (
	"GoSlice/data"
	"GoSlice/gcode"
)

// PreLayer adds starting gcode, resets the extrude speeds on each layer and enables the fan above a specific layer.
type PreLayer struct{}

func (PreLayer) Init(model data.OptimizedModel) {}

func (PreLayer) Render(b *gcode.Builder, layerNr int, layers []data.PartitionedLayer, z data.Micrometer, options *data.Options) {
	b.AddComment("LAYER:%v", layerNr)
	if layerNr == 0 {
		// starting gcode
		b.AddComment("START_GCODE")
		b.AddComment("Generated with GoSlice")
		b.AddCommand("G1 X0 Y20 Z0.2 F3000 ; get ready to prime")
		b.AddCommand("G92 E0 ; reset extrusion distance")
		b.AddCommand("G1 X200 E20 F600 ; prime nozzle")
		b.AddCommand("G1 Z5 F5000 ; lift nozzle")
		b.AddCommand("G92 E0 ; reset extrusion distance")

		b.SetExtrusion(options.Print.InitialLayerThickness, options.Printer.ExtrusionWidth, options.Filament.FilamentDiameter)

		// set speeds
		b.SetExtrudeSpeed(options.Print.LayerSpeed)
		b.SetMoveSpeed(options.Print.MoveSpeed)

		// force the InitialLayerSpeed for first layer
		b.SetExtrudeSpeedOverride(options.Print.IntialLayerSpeed)
	} else {
		b.DisableExtrudeSpeedOverride()
		b.SetExtrudeSpeed(options.Print.LayerSpeed)
	}

	if layerNr == 2 {
		b.AddCommand("M106 ; enable fan")
	}
}

// PostLayer adds GCode at the last layer.
type PostLayer struct{}

func (PostLayer) Init(model data.OptimizedModel) {}

func (PostLayer) Render(builder *gcode.Builder, layerNr int, layers []data.PartitionedLayer, z data.Micrometer, options *data.Options) {
	// ending gcode
	if layerNr == len(layers)-1 {
		builder.AddComment("END_GCODE")
		builder.SetExtrusion(options.Print.LayerThickness, options.Printer.ExtrusionWidth, options.Filament.FilamentDiameter)
		builder.AddCommand("M107 ; disable fan")
	}
}
