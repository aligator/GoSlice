package slicer

import (
	"GoSlicer/slicer/data"
	"GoSlicer/slicer/gcode"
	"GoSlicer/slicer/handle"
	"GoSlicer/slicer/optimize"
	"GoSlicer/slicer/slice"
	"GoSlicer/slicer/stl"
	"GoSlicer/slicer/write"
	"GoSlicer/util"
	"fmt"
	"time"
)

type GoSlicer struct {
	reader    handle.ModelReader
	optimizer handle.ModelOptimizer
	slicer    handle.ModelSlicer
	modifiers []handle.LayerModifier
	generator handle.GCodeGenerator
	writer    handle.GCodeWriter
}

func NewGoSlicer() *GoSlicer {
	s := &GoSlicer{
		reader: stl.Reader(),
		optimizer: optimize.NewOptimizer(optimize.OptimizerOptions{
			MeldDistance: util.Micrometer(30),
			Center:       util.NewMicroVec3(102500, 102500, 0),
		}),
		slicer: slice.NewSlicer(slice.SlicerOptions{
			InitialThickness: 200,
			LayerThickness:   200,
		}),
		modifiers: nil,
		generator: gcode.NewGenerator(gcode.GeneratorOptions{
			LayerThickness:        200,
			InitialLayerThickness: 200,
			FilamentDiameter:      1500,
			ExtrusionWidth:        400,
			InsetCount:            5,
		}),
		writer: write.Writer(),
	}
	return s
}

func (s GoSlicer) Process(filename string, outFilename string) error {
	startTime := time.Now()
	models, err := s.reader.Read(filename)
	if err != nil {
		return err
	}

	var optimizedModel data.OptimizedModel

	// TODO: support several model processing
	//for i, model := range models {
	optimizedModel, err = s.optimizer.Optimize(models[0])
	if err != nil {
		return err
	}
	//}

	optimizedModel.SaveDebugSTL("test.stl")

	layers, err := s.slicer.Slice(optimizedModel)
	if err != nil {
		return err
	}

	for _, m := range s.modifiers {
		for i, layer := range layers {
			layers[i], err = m.Modify(layer)
			if err != nil {
				return err
			}
		}
	}

	s.generator.Init()
	for i, layer := range layers {
		s.generator.Generate(i, layer)
	}

	gcode := s.generator.Finish()

	err = s.writer.Write(gcode, outFilename)

	fmt.Println("full processing time:", time.Now().Sub(startTime))

	return nil
}
