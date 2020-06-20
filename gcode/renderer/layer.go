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
		b.AddComment("Generated with GoSlice")
		b.AddComment("______________________")

		b.AddCommand("M107 ; disable fan")

		// set and wait for the initial temperature
		b.AddComment("SET_INITIAL_TEMP")
		b.AddCommand("M104 S%d ; start heating hot end", options.Filament.InitialHotEndTemperature)
		b.AddCommand("M190 S%d ; heat and wait for bed", options.Filament.InitialBedTemperature)
		b.AddCommand("M109 S%d ; wait for hot end temperature", options.Filament.InitialHotEndTemperature)

		// starting gcode
		b.AddComment("START_GCODE")
		b.AddCommand("G1 X0 Y20 Z0.2 F3000 ; get ready to prime")
		b.AddCommand("G92 E0 ; reset extrusion distance")
		b.AddCommand("G1 X200 E20 F600 ; prime nozzle")
		b.AddCommand("G1 Z5 F5000 ; lift nozzle")
		b.AddCommand("G92 E0 ; reset extrusion distance")

		b.SetExtrusion(options.Print.InitialLayerThickness, options.Printer.ExtrusionWidth, options.Filament.FilamentDiameter)

		// set speeds
		b.SetExtrudeSpeed(options.Print.LayerSpeed)
		b.SetMoveSpeed(options.Print.MoveSpeed)

		// set retraction
		b.SetRetractionSpeed(options.Filament.RetractionSpeed)
		b.SetRetractionAmount(options.Filament.RetractionLength)

		// force the InitialLayerSpeed for first layer
		b.SetExtrudeSpeedOverride(options.Print.IntialLayerSpeed)
	} else {
		b.DisableExtrudeSpeedOverride()
		b.SetExtrudeSpeed(options.Print.LayerSpeed)
	}

	if layerNr == 2 {
		b.AddCommand("M106 ; enable fan")
	}

	if layerNr == options.Filament.InitialTemeratureLayerCount {
		// set the normal temperature
		// this is done without waiting
		b.AddComment("SET_TEMP")
		b.AddCommand("M140 S%d", options.Filament.BedTemperature)
		b.AddCommand("M104 S%d", options.Filament.HotEndTemperature)
	}
}

// PostLayer adds GCode at the last layer.
type PostLayer struct{}

func (PostLayer) Init(model data.OptimizedModel) {}

func (PostLayer) Render(b *gcode.Builder, layerNr int, layers []data.PartitionedLayer, z data.Micrometer, options *data.Options) {
	// ending gcode
	if layerNr == len(layers)-1 {
		b.AddComment("END_GCODE")
		b.SetExtrusion(options.Print.LayerThickness, options.Printer.ExtrusionWidth, options.Filament.FilamentDiameter)
		b.AddCommand("M107 ; disable fan")

		// disable heaters
		b.AddCommand("M104 S0 ; Set Hot-end to 0C (off)")
		b.AddCommand("M140 S0 ; Set bed to 0C (off)")
	}
}
