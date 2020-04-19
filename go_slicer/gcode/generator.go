package gcode

import (
	"GoSlicer/go_slicer/clip"
	"GoSlicer/go_slicer/data"
	"GoSlicer/go_slicer/handle"
	"GoSlicer/util"
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

	g.builder.addComment("Generated with GoSlicer")
	g.builder.addComment("\nG1 X0 Y20 Z0.2 F3000 ; get ready to prime")
	g.builder.addComment("\nG92 E0 ; reset extrusion distance")
	g.builder.addComment("\nG1 X200 E20 F600 ; prime nozzle")
	g.builder.addComment("\nG1 Z5 F5000 ; lift nozzle")

	g.builder.setExtrusion(g.options.Print.InitialLayerThickness, g.options.Printer.ExtrusionWidth, g.options.Filament.FilamentDiameter)
}

func (g *generator) Generate(layerNum int, layer data.PartitionedLayer) {
	if layerNum == 0 {
		g.builder.setExtrudeSpeed(g.options.Print.IntialLayerSpeed)
	} else {
		g.builder.setExtrudeSpeed(g.options.Print.LayerSpeed)
	}

	if layerNum == 2 {
		g.builder.addComment("\nM106 ; enable fan")
	}

	c := clip.NewClip()
	insetParts := c.InsetLayer(layer, g.options.Printer.ExtrusionWidth, g.options.Print.InsetCount)
	fmt.Printf("Processing layer %v...\n", layerNum)
	g.builder.addComment("LAYER:%v", layerNum)

	for _, part := range insetParts {
		for insetNr := len(part) - 1; insetNr > -1; insetNr-- {
			if insetNr == 0 {
				g.builder.addComment("TYPE:WALL-OUTER")
			} else {
				g.builder.addComment("TYPE:WALL-INNER")
			}

			for _, poly := range part[insetNr] {
				g.builder.addPolygon(poly, g.options.Print.InitialLayerThickness+util.Micrometer(layerNum)*g.options.Print.LayerThickness)
			}
		}
	}
}

func (g *generator) Finish() string {
	g.builder.setExtrusion(g.options.Print.LayerThickness, g.options.Printer.ExtrusionWidth, g.options.Filament.FilamentDiameter)
	g.builder.addComment("\nM107 ; enable fan")

	return g.builder.buf.String()
}
