// Package handler provides interfaces for all steps needed for the whole process of generating GCode out of a model-file.

package handler

import "github.com/aligator/goslice/data"

type Namer interface {
	GetName() string
}

type Named struct {
	Name string
}

func (n Named) GetName() string {
	return n.Name
}

// ModelReader reads a model from a file.
type ModelReader interface {
	Read(filename string) (data.Model, error)
}

// ModelOptimizer can optimize a model and generates an optimized model out of it.
type ModelOptimizer interface {
	Optimize(m data.Model) (data.OptimizedModel, error)
}

// ModelSlicer can slice an optimized model into several layers.
type ModelSlicer interface {
	Slice(m data.OptimizedModel) ([]data.PartitionedLayer, error)
}

// LayerModifier can add new attributes to the layers or even alter the layer directly.
type LayerModifier interface {
	Namer
	Init(m data.OptimizedModel)
	Modify(layers []data.PartitionedLayer) error
}

// GCodeGenerator generates the GCode out of the given layers.
// The layers are already modified by the layer modifiers.
// So the attributes added by them can be used.
type GCodeGenerator interface {
	Init(m data.OptimizedModel)
	Generate(layer []data.PartitionedLayer) (string, error)
}

// GCodeWriter writes the given GCode into the given destination.
type GCodeWriter interface {
	Write(gcode string, destination string) error
}
