package handle

import "GoSlice/go_slice/data"

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
	Modify(layer data.PartitionedLayer) (data.PartitionedLayer, error)
}

type GCodeGenerator interface {
	Init()
	Generate(layerNum int, layer data.PartitionedLayer)
	Finish() string
}

type GCodeWriter interface {
	Write(gcode string, filename string) error
}
