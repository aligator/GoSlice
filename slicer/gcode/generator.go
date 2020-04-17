package gcode

import (
	"GoSlicer/slicer/clip"
	"GoSlicer/slicer/data"
	"GoSlicer/slicer/handle"
	"GoSlicer/util"
	"bytes"
	"fmt"
)

// TODO use interface with functional options?
type GeneratorOptions struct {
	LayerThickness        util.Micrometer
	InitialLayerThickness util.Micrometer
	FilamentDiameter      util.Micrometer
	ExtrusionWidth        util.Micrometer
	InsetCount            int
}

type generator struct {
	options GeneratorOptions
	gcode   string
	builder *gcodeBuilder
}

func NewGenerator(options GeneratorOptions) handle.GCodeGenerator {
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

	g.builder.setExtrusion(g.options.InitialLayerThickness, g.options.ExtrusionWidth, g.options.FilamentDiameter)
}

func (g *generator) Generate(layerNum int, layer data.PartitionedLayer) {
	if layerNum == 2 {
		g.builder.addComment("\nM106 ; enable fan")
	}

	c := clip.NewClip()
	insetParts := c.InsetLayer(layer, g.options.ExtrusionWidth, g.options.InsetCount)
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
				g.builder.addPolygon(poly, g.options.InitialLayerThickness+util.Micrometer(layerNum)*g.options.LayerThickness)
			}
		}
	}
}

func (g *generator) Finish() string {
	g.builder.setExtrusion(g.options.LayerThickness, g.options.ExtrusionWidth, g.options.FilamentDiameter)
	g.builder.addComment("\nM107 ; enable fan")

	return g.builder.buf.String()
}
