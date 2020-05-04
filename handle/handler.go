package handle

import "GoSlice/data"

type ModelReader interface {
	Read(filename string) ([]data.Model, error)
}

type ModelOptimizer interface {
	Optimize(m data.Model) (data.OptimizedModel, error)
}

type ModelSlicer interface {
	Slice(m data.OptimizedModel) ([]data.PartitionedLayer, error)
}

type LayerModifier interface {
	Init(m data.OptimizedModel)
	Modify(layerNr int, layers []data.PartitionedLayer) ([]data.PartitionedLayer, error)
}

type GCodeGenerator interface {
	Init(m data.OptimizedModel)
	Generate(layer []data.PartitionedLayer) string
}

type GCodeWriter interface {
	Write(gcode string, filename string) error
}
