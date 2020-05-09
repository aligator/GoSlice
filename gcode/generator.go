// Package gcode provides a generator for GCode files.
package gcode

import (
	"GoSlice/data"
	"GoSlice/gcode/builder"
	"GoSlice/handle"
	"bytes"
)

// Renderer can be used to add GCodes based on the current layer and layer data.
// Several renderers can be provided to the generator.
type Renderer interface {
	// Init is called once at the beginning and can be used to set up the renderer.
	// For example the infill patterns can be instanciated int this method.
	Init(model data.OptimizedModel)

	// Render is called for each layer and the provided builder can be used to add gcode.
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

// WithRenderer adds a renderer to the generator.
func WithRenderer(r Renderer) option {
	return func(s *generator) {
		s.renderers = append(s.renderers, r)
	}
}

// NewGenerator returns a new GCode generator which can be customized by adding several renderers using WithRenderer().
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

// Generate generates the GCode by using the renderers added to the generator.
// The final GCode is just returned.
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
