package gcode

import (
	"GoSlice/data"
	"GoSlice/handle"
	"bytes"
)

type GCodePaths struct {
	paths data.Paths
	Speed data.Millimeter
}

type LayerMetadata struct {
	Elements map[string]interface{}
}

type RenderStep func(builder *gcodeBuilder, layerNr int, layers []data.PartitionedLayer, z data.Micrometer, options *data.Options)

type generator struct {
	options   *data.Options
	gcode     string
	builder   *gcodeBuilder
	renderers []RenderStep
}

func NewGenerator(options *data.Options) handle.GCodeGenerator {
	// The following steps and renderers are the builtin ones.
	// Later it will be possible to add custom ones to extend the functionality.

	return &generator{
		options: options,
		renderers: []RenderStep{
			// pre layer
			func(builder *gcodeBuilder, layerNr int, layers []data.PartitionedLayer, z data.Micrometer, options *data.Options) {
				builder.addComment("LAYER:%v", layerNr)
				if layerNr == 0 {
					// force the InitialLayerSpeed for first layer
					builder.setExtrudeSpeedOverride(options.Print.IntialLayerSpeed)
				} else {
					builder.disableExtrudeSpeedOverride()
					builder.setExtrudeSpeed(options.Print.LayerSpeed)
				}
			},

			// fan control
			func(builder *gcodeBuilder, layerNr int, layers []data.PartitionedLayer, z data.Micrometer, options *data.Options) {
				if layerNr == 2 {
					builder.addCommand("M106 ; enable fan")
				}
			},

			// perimeters
			func(builder *gcodeBuilder, layerNr int, layers []data.PartitionedLayer, z data.Micrometer, options *data.Options) {
				for _, part := range layers[layerNr].LayerParts() {
					perimeters, ok := part.Attributes()["perimeters"].([][]data.Paths)

					if !ok {
						continue
					}

					// perimeters contains them as [wallNr][insetNr]data.Paths
					for _, wall := range perimeters {
						for insetNr, insets := range wall {
							// set the speed based on outer or inner layer
							if insetNr == 0 {
								builder.addComment("TYPE:WALL-OUTER")
								builder.setExtrudeSpeed(options.Print.OuterPerimeterSpeed)
							} else {
								builder.addComment("TYPE:WALL-INNER")
								builder.setExtrudeSpeed(options.Print.LayerSpeed)
							}

							for _, inset := range insets {
								builder.addPolygon(inset, z)
							}
						}
					}
				}
			},
			/*
				// bottom layer TODO: bottom and top layers
				func(builder *gcodeBuilder, layerNr int, layers []data.PartitionedLayer, z data.Micrometer, options *data.Options) {
					if meta[layerNr].Elements["bottomLayer"] == nil {
						return
					}

					layer, ok := meta[layerNr].Elements["bottomLayer"].([]data.Paths)
					if !ok {
						fmt.Println("wrong type for bottomLayer elements")
						return
					}
					builder.addComment("bottomLayer")

					for _, paths := range layer {
						for _, path := range paths {
							builder.addPolygon(path, z)
						}

					}
				},
			*/
			// TODO: infill, support, bridges,...
		},
	}
}

func (g *generator) init() {
	b := []byte{}
	g.builder = newGCodeBuilder(bytes.NewBuffer(b))

	g.builder.addComment("Generated with GoSlice")
	g.builder.addCommand("G1 X0 Y20 Z0.2 F3000 ; get ready to prime")
	g.builder.addCommand("G92 E0 ; reset extrusion distance")
	g.builder.addCommand("G1 X200 E20 F600 ; prime nozzle")
	g.builder.addCommand("G1 Z5 F5000 ; lift nozzle")
	g.builder.addCommand("G92 E0 ; reset extrusion distance")

	g.builder.setExtrusion(g.options.Print.InitialLayerThickness, g.options.Printer.ExtrusionWidth, g.options.Filament.FilamentDiameter)
}

func (g *generator) Generate(layers []data.PartitionedLayer) string {
	g.init()

	for layerNr := range layers {
		for _, renderer := range g.renderers {
			z := g.options.Print.InitialLayerThickness + data.Micrometer(layerNr)*g.options.Print.LayerThickness
			renderer(g.builder, layerNr, layers, z, g.options)
		}
	}

	return g.finish()
}

func (g *generator) finish() string {
	g.builder.setExtrusion(g.options.Print.LayerThickness, g.options.Printer.ExtrusionWidth, g.options.Filament.FilamentDiameter)
	g.builder.addCommand("M107 ; enable fan")

	return g.builder.buf.String()
}
