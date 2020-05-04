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
	options *data.Options
	gcode   string
	builder *gcodeBuilder

	topBottomInfill clip.Pattern
	internalInfill  clip.Pattern

	renderers []RenderStep
}

func (g *generator) Init(model data.OptimizedModel) {
	g.topBottomInfill = clip.NewLinearPattern(model.Min().PointXY(), model.Max().PointXY(), g.options.Printer.ExtrusionWidth)

	// TODO: the calculation of the percentage is currently very basic and may not be correct.
	// It needs improvement.

	if g.options.Print.InfillPercent != 0 {
		mm10 := data.Millimeter(10).ToMicrometer()
		linesPer10mmFor100Percent := mm10 / g.options.Printer.ExtrusionWidth
		linesPerArea10x10ForInfillPercent := float64(linesPer10mmFor100Percent) * float64(g.options.Print.InfillPercent) / 100.0

		lineWidth := data.Micrometer(float64(mm10) / linesPerArea10x10ForInfillPercent)

		g.internalInfill = clip.NewLinearPattern(model.Min().PointXY(), model.Max().PointXY(), lineWidth)
	}
}

func NewGenerator(options *data.Options) handle.GCodeGenerator {
	g := &generator{
		options: options,
	}

	// The following steps and renderers are the builtin ones.
	// Later it will be possible to add custom ones to extend the functionality.

	renderers := []RenderStep{
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
			perimeters, ok := layers[layerNr].Attributes()["perimeters"].([][][]data.LayerPart)
			if !ok {
				return
			}

			// perimeters contains them as [part][insetNr][insetParts]
			for _, part := range perimeters {
				for insetNr := range part {
					// print the outer perimeter as last perimeter
					if insetNr >= len(part)-1 {
						insetNr = 0
					} else {
						insetNr++
					}

					for _, insetParts := range part[insetNr] {
						if insetNr == 0 {
							builder.addComment("TYPE:WALL-OUTER")
							builder.setExtrudeSpeed(options.Print.OuterPerimeterSpeed)
						} else {
							builder.addComment("TYPE:WALL-INNER")
							builder.setExtrudeSpeed(options.Print.LayerSpeed)
						}

						for _, hole := range insetParts.Holes() {
							builder.addPolygon(hole, z)
						}

						builder.addPolygon(insetParts.Outline(), z)
					}
				}
			}
		},

		// bottom and top layers
		func(builder *gcodeBuilder, layerNr int, layers []data.PartitionedLayer, z data.Micrometer, options *data.Options) {
			bottom, ok := layers[layerNr].Attributes()["bottom"].([]data.LayerPart)
			if !ok {
				return
			}

			for _, part := range bottom {
				builder.addComment("TYPE:FILL")
				builder.addComment("BOTTOM-FILL")
				for _, path := range g.topBottomInfill.Fill(layerNr, part) {
					builder.addPolygon(path, z)
				}
			}
		},
		func(builder *gcodeBuilder, layerNr int, layers []data.PartitionedLayer, z data.Micrometer, options *data.Options) {
			top, ok := layers[layerNr].Attributes()["top"].([]data.LayerPart)
			if !ok {
				return
			}

			for _, part := range top {
				builder.addComment("TYPE:FILL")
				builder.addComment("TOP-FILL")
				for _, path := range g.topBottomInfill.Fill(layerNr, part) {
					builder.addPolygon(path, z)
				}
			}
		},
		func(builder *gcodeBuilder, layerNr int, layers []data.PartitionedLayer, z data.Micrometer, options *data.Options) {
			infill, ok := layers[layerNr].Attributes()["infill"].([]data.LayerPart)
			if !ok {
				return
			}

			for _, part := range infill {
				builder.addComment("TYPE:FILL")
				builder.addComment("INTERNAL-FILL")
				for _, path := range g.internalInfill.Fill(layerNr, part) {
					builder.addPolygon(path, z)
				}
			}
		},
		// TODO: infill, support, bridges,...
	}

	g.renderers = renderers

	return g
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
