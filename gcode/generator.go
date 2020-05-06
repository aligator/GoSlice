package gcode

import (
	"GoSlice/data"
	"GoSlice/gcode/builder"
	"GoSlice/handle"
	"bytes"
)

type Renderer interface {
	Init(model data.OptimizedModel)
	Render(builder builder.Builder, layerNr int, layers []data.PartitionedLayer, z data.Micrometer, options *data.Options)
}

type GCodePaths struct {
	paths data.Paths
	Speed data.Millimeter
}

type LayerMetadata struct {
	Elements map[string]interface{}
}

type RenderStep func(builder *builder.GCode, layerNr int, layers []data.PartitionedLayer, z data.Micrometer, options *data.Options)

type generator struct {
	options *data.Options
	gcode   string
	builder builder.Builder

	renderers []Renderer
}

func (g *generator) Init(model data.OptimizedModel) {
	for _, renderer := range g.renderers {
		renderer.Init(model)
	}
}

type option func(s *generator)

func (s *generator) With(o ...option) {
	for _, option := range o {
		option(s)
	}
}

func WithRenderer(r Renderer) option {
	return func(s *generator) {
		s.renderers = append(s.renderers, r)
	}
}

func NewGenerator(options *data.Options, generatorOptions ...option) handle.GCodeGenerator {
	g := &generator{
		options: options,
	}

	for _, o := range generatorOptions {
		o(g)
	}

	return g
}

func (g *generator) init() {
	var b []byte
	g.builder = builder.NewGCodeBuilder(bytes.NewBuffer(b))

	g.builder.AddComment("Generated with GoSlice")
	g.builder.AddComment("G1 X0 Y20 Z0.2 F3000 ; get ready to prime")
	g.builder.AddComment("G92 E0 ; reset extrusion distance")
	g.builder.AddComment("G1 X200 E20 F600 ; prime nozzle")
	g.builder.AddComment("G1 Z5 F5000 ; lift nozzle")
	g.builder.AddComment("G92 E0 ; reset extrusion distance")

	g.builder.SetExtrusion(g.options.Print.InitialLayerThickness, g.options.Printer.ExtrusionWidth, g.options.Filament.FilamentDiameter)
}

func (g *generator) Generate(layers []data.PartitionedLayer) string {
	g.init()

	for layerNr := range layers {
		for _, renderer := range g.renderers {
			z := g.options.Print.InitialLayerThickness + data.Micrometer(layerNr)*g.options.Print.LayerThickness
			renderer.Render(g.builder, layerNr, layers, z, g.options)
		}
	}

	return g.finish()
}

func (g *generator) finish() string {
	g.builder.SetExtrusion(g.options.Print.LayerThickness, g.options.Printer.ExtrusionWidth, g.options.Filament.FilamentDiameter)
	g.builder.AddCommand("M107 ; enable fan")

	return g.builder.Buffer().String()
}
