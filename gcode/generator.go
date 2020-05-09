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
}

func (g *generator) Generate(layers []data.PartitionedLayer) string {
	g.init()

	for layerNr := range layers {
		for _, renderer := range g.renderers {
			z := g.options.Print.InitialLayerThickness + data.Micrometer(layerNr)*g.options.Print.LayerThickness
			renderer.Render(g.builder, layerNr, layers, z, g.options)
		}
	}

	return g.builder.Buffer().String()
}
