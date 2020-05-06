package renderer

import (
	"GoSlice/data"
	"GoSlice/gcode/builder"
)

type PreLayer struct{}

func (p PreLayer) Init(model data.OptimizedModel) {}

func (p PreLayer) Render(builder builder.Builder, layerNr int, layers []data.PartitionedLayer, z data.Micrometer, options *data.Options) {
	builder.AddComment("LAYER:%v", layerNr)
	if layerNr == 0 {
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
