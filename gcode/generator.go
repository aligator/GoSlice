// Package gcode provides a generator for GCode files.
package gcode

import (
	"github.com/aligator/goslice/data"
	"github.com/aligator/goslice/handler"
)

// Renderer can be used to add GCodes based on the current layer and layer data.
// Several renderers can be provided to the generator.
type Renderer interface {
	// Init is called once at the beginning and can be used to set up the renderer.
	// For example the infill patterns can be instantiated int this method.
	Init(model data.OptimizedModel)

	// Render is called for each layer and the provided Builder can be used to add gcode.
	Render(b *Builder, layerNr int, maxLayer int, layer data.PartitionedLayer, z data.Micrometer, options *data.Options) error
}

type generator struct {
	options *data.Options
	gcode   string
	builder *Builder

	renderers []Renderer
}

func (g *generator) Init(model data.OptimizedModel) {
	for _, renderer := range g.renderers {
		renderer.Init(model)
	}
}

type option func(s *generator)

// WithRenderer adds a renderer to the generator.
func WithRenderer(r Renderer) option {
	return func(s *generator) {
		s.renderers = append(s.renderers, r)
	}
}

// NewGenerator returns a new Builder generator which can be customized by adding several renderers using WithRenderer().
func NewGenerator(options *data.Options, generatorOptions ...option) handler.GCodeGenerator {
	g := &generator{
		options: options,
	}

	for _, option := range generatorOptions {
		option(g)
	}

	return g
}

func (g *generator) init() {
	g.builder = NewGCodeBuilder(g.options)
}

// Generate generates the GCode by using the renderers added to the generator.
// The final GCode is just returned as string.
func (g *generator) Generate(layers []data.PartitionedLayer) (string, error) {
	g.init()

	maxLayer := len(layers) - 1

	for layerNr := range layers {
		for _, renderer := range g.renderers {
			z := g.options.Print.InitialLayerThickness + data.Micrometer(layerNr)*g.options.Print.LayerThickness
			err := renderer.Render(g.builder, layerNr, maxLayer, layers[layerNr], z, g.options)
			if err != nil {
				return "", err
			}
		}
	}

	return g.builder.String(), nil
}
