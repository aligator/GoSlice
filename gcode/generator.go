package gcode

import (
	"GoSlice/clip"
	"GoSlice/data"
	"GoSlice/handle"
	"GoSlice/util"
	"bytes"
	"fmt"
)

type GCodePaths struct {
	paths data.Paths
	Speed util.Millimeter
}

type LayerMetadata struct {
	Elements map[string]interface{}
}

type GeneratorStep func(layerNr int, layers []data.PartitionedLayer, meta []LayerMetadata, options *data.Options) LayerMetadata
type RenderStep func(builder *gcodeBuilder, layerNr int, meta []LayerMetadata, z util.Micrometer, options *data.Options)

type generator struct {
	options   *data.Options
	gcode     string
	builder   *gcodeBuilder
	steps     []GeneratorStep
	renderers []RenderStep
}

func NewGenerator(options *data.Options) handle.GCodeGenerator {
	return &generator{
		options: options,
		steps: []GeneratorStep{
			// perimeters
			func(layerNr int, layers []data.PartitionedLayer, meta []LayerMetadata, options *data.Options) LayerMetadata {
				// perimeters per object
				innerPerimeters := []GCodePaths{}
				outerPerimeters := []GCodePaths{}
				middlePerimeters := []GCodePaths{}

				// generate perimeters
				c := clip.NewClip()
				insetParts := c.InsetLayer(layers[layerNr], options.Printer.ExtrusionWidth, options.Print.InsetCount)

				// iterate over all generated perimeters
				for _, part := range insetParts {
					for _, wall := range part {
						for insetNum, wallInset := range wall {
							var speed util.Millimeter
							// set the speed based on the current perimeter
							if insetNum == 0 {
								if layerNr > 0 {
									speed = options.Print.OuterPerimeterSpeed
								}

								outerPerimeters = append(outerPerimeters, GCodePaths{
									paths: wallInset,
									Speed: speed,
								})
							} else {
								if layerNr > 0 {
									speed = options.Print.LayerSpeed
								}
							}

							if insetNum > 0 && insetNum < len(wall)-1 {
								middlePerimeters = append(middlePerimeters, GCodePaths{
									paths: wallInset,
									Speed: speed,
								})
							} else {
								innerPerimeters = append(innerPerimeters, GCodePaths{
									paths: wallInset,
									Speed: speed,
								})
							}
						}
					}
				}

				meta[layerNr].Elements["perimeter"] = [3][]GCodePaths{
					innerPerimeters,
					middlePerimeters,
					outerPerimeters,
				}

				return meta[layerNr]
			},
		},
		renderers: []RenderStep{
			// pre layer
			func(builder *gcodeBuilder, layerNr int, meta []LayerMetadata, z util.Micrometer, options *data.Options) {
				if layerNr == 0 {
					builder.setExtrudeSpeed(options.Print.IntialLayerSpeed)
				} else {
					builder.setExtrudeSpeed(options.Print.LayerSpeed)
				}
			},

			// fan control
			func(builder *gcodeBuilder, layerNr int, meta []LayerMetadata, z util.Micrometer, options *data.Options) {
				if layerNr == 2 {
					builder.addComment("\nM106 ; enable fan")
				}
			},

			// perimeters
			func(builder *gcodeBuilder, layerNr int, meta []LayerMetadata, z util.Micrometer, options *data.Options) {
				p, ok := meta[layerNr].Elements["perimeter"].([3][]GCodePaths)
				if !ok {
					fmt.Print("wrong type for perimeter elements")
					return
				}

				for i, perimeter := range p {
					if i == 0 {
						builder.addComment("TYPE:WALL-OUTER")
					} else {
						builder.addComment("TYPE:WALL-INNER")
					}

					for _, paths := range perimeter {
						for _, path := range paths.paths {
							builder.setExtrudeSpeed(paths.Speed)
							builder.addPolygon(path, z, false)
						}
					}
				}
			},
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
	meta := []LayerMetadata{}

	for _, step := range g.steps {
		for layerNr := range layers {
			if len(meta) <= layerNr {
				meta = append(meta, LayerMetadata{
					Elements: map[string]interface{}{},
				})
			}
			meta[layerNr] = step(layerNr, layers, meta, g.options)
		}
	}

	g.init()

	for layerNr := range layers {
		for _, renderer := range g.renderers {
			z := g.options.Print.InitialLayerThickness + util.Micrometer(layerNr)*g.options.Print.LayerThickness
			renderer(g.builder, layerNr, meta, z, g.options)
		}
	}

	return g.finish()
}

func (g *generator) finish() string {
	g.builder.setExtrusion(g.options.Print.LayerThickness, g.options.Printer.ExtrusionWidth, g.options.Filament.FilamentDiameter)
	g.builder.addComment("\nM107 ; enable fan")

	return g.builder.buf.String()
}
