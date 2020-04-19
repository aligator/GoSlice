package gcode

import (
	"GoSlice/go_slice/clip"
	"GoSlice/go_slice/data"
	"GoSlice/go_slice/handle"
	"GoSlice/util"
	"bytes"
	"fmt"
)

type generator struct {
	options *data.Options
	gcode   string
	builder *gcodeBuilder
}

func NewGenerator(options *data.Options) handle.GCodeGenerator {
	return &generator{
		options: options,
	}
}

func (g *generator) Init() {
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

func (g *generator) Generate(layerNum int, layer data.PartitionedLayer) {
	if layerNum == 0 {
		g.builder.setExtrudeSpeed(g.options.Print.IntialLayerSpeed)
	} else {
		g.builder.setExtrudeSpeed(g.options.Print.LayerSpeed)
	}

	z := g.options.Print.InitialLayerThickness + util.Micrometer(layerNum)*g.options.Print.LayerThickness

	if layerNum == 2 {
		g.builder.addComment("\nM106 ; enable fan")
	}

	// generate perimeters
	c := clip.NewClip()
	insetParts := c.InsetLayer(layer, g.options.Printer.ExtrusionWidth, g.options.Print.InsetCount)
	fmt.Printf("Processing layer %v...\n", layerNum)
	g.builder.addComment("LAYER:%v", layerNum)

	// iterate over all generated perimeters
	for _, part := range insetParts {
		for _, wall := range part {
			for insetNum, wallInset := range wall {
				for _, inset := range wallInset {
					// set the speed based on the current perimeter
					if insetNum == 0 {
						g.builder.addComment("TYPE:WALL-OUTER")
						g.builder.setExtrudeSpeed(g.options.Print.OuterPerimeterSpeed)
					} else {
						g.builder.addComment("TYPE:WALL-INNER")
						g.builder.setExtrudeSpeed(g.options.Print.LayerSpeed)
					}

					// add the perimeter and check if it is a bottom layer
					// -> generate one if it is
					fill := insetNum == len(wall)-1 && layerNum == 0
					g.builder.addPolygon(inset, z, fill)
				}
			}
		}
	}
}

func (g *generator) Finish() string {
	g.builder.setExtrusion(g.options.Print.LayerThickness, g.options.Printer.ExtrusionWidth, g.options.Filament.FilamentDiameter)
	g.builder.addComment("\nM107 ; enable fan")

	return g.builder.buf.String()
}
