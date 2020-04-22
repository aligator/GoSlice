package gcode

import (
	"GoSlice/clip"
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
				perimeters, ok := layers[layerNr].Attributes()["perimeters"].([3][]data.Paths)
				if !ok {
					return
				}

				var perimeterPaths [][][3]data.Paths

				// reorder the perimeters to [wallNr][insetNr][typeNr]data.Paths
				for typeNum, _ := range perimeters {
					for wallNr, wall := range perimeters[typeNum] {
						if len(perimeterPaths) <= wallNr {
							perimeterPaths = append(perimeterPaths, [][3]data.Paths{})
						}

						for insetNr, inset := range wall {
							if len(perimeterPaths[wallNr]) <= insetNr {
								perimeterPaths[wallNr] = append(perimeterPaths[wallNr], [3]data.Paths{})
							}

							if len(perimeterPaths[wallNr][insetNr]) <= typeNum {
								perimeterPaths[wallNr][insetNr][typeNum] = append(perimeterPaths[wallNr][insetNr][typeNum], data.Path{})
							}

							perimeterPaths[wallNr][insetNr][typeNum] = append(perimeterPaths[wallNr][insetNr][typeNum], inset)
						}
					}
				}

				// perimeters contains them as [wallNr][insetNr][typeNr]data.Paths
				for _, wall := range perimeterPaths {
					for insetNr, insets := range wall {
						for typeNr, inset := range insets {
							// set the speed based on outer or inner layer
							if typeNr == 0 || insetNr == len(insets)-1 {
								builder.addComment("TYPE:WALL-OUTER")
								builder.setExtrudeSpeed(options.Print.OuterPerimeterSpeed)
							} else {
								builder.addComment("TYPE:WALL-INNER")
								builder.setExtrudeSpeed(options.Print.LayerSpeed)
							}

							for _, path := range inset {
								builder.addPolygon(path, z)
							}
						}
					}
				}
			},

			// bottom layer TODO: bottom and top layers
			func(builder *gcodeBuilder, layerNr int, layers []data.PartitionedLayer, z data.Micrometer, options *data.Options) {
				bottom, ok := layers[layerNr].Attributes()["bottom"].([]data.Paths)
				if !ok {
					return
				}

				c := clip.NewClipper()

				for _, paths := range bottom {
					infill := c.Fill(paths, options.Printer.ExtrusionWidth, options.Print.InfillOverlapPercent)
					if infill != nil {
						builder.addComment("bottomLayer")
						for _, path := range infill {
							builder.addPolygon(path, z)
						}
					}
				}
			},
			// TODO: infill, support, bridges,...
		},
	}
}

func (g *generator) init() {
	var b []byte
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
