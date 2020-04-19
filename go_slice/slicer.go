package go_slice

import (
	"GoSlice/go_slice/data"
	"GoSlice/go_slice/gcode"
	"GoSlice/go_slice/handle"
	"GoSlice/go_slice/optimize"
	"GoSlice/go_slice/slice"
	"GoSlice/go_slice/stl"
	"GoSlice/go_slice/write"
	"GoSlice/util"
	"fmt"
	"time"
)

type GoSlice struct {
	o         *data.Options
	reader    handle.ModelReader
	optimizer handle.ModelOptimizer
	slicer    handle.ModelSlicer
	modifiers []handle.LayerModifier
	generator handle.GCodeGenerator
	writer    handle.GCodeWriter
}

func NewGoSlice(o ...option) *GoSlice {
	options := data.Options{
		Print: data.PrintOptions{
			IntialLayerSpeed:    30,
			LayerSpeed:          60,
			OuterPerimeterSpeed: 40,

			InitialLayerThickness: 200,
			LayerThickness:        200,
			InsetCount:            2,
		},
		Filament: data.FilamentOptions{
			FilamentDiameter: util.Millimeter(1.75).ToMicrometer(),
		},
		Printer: data.PrinterOptions{
			ExtrusionWidth: 400,
		},

		MeldDistance: 30,
		Center: util.NewMicroVec3(
			util.Millimeter(100).ToMicrometer(),
			util.Millimeter(100).ToMicrometer(),
			0,
		),
		JoinPolygonSnapDistance:   100,
		FinishPolygonSnapDistance: 1000,
	}

	s := &GoSlice{
		o:         &options,
		reader:    stl.Reader(),
		optimizer: optimize.NewOptimizer(&options),
		slicer:    slice.NewSlicer(&options),
		modifiers: nil,
		generator: gcode.NewGenerator(&options),
		writer:    write.Writer(),
	}

	s.With(o...)

	return s
}

func (s *GoSlice) Process(filename string, outFilename string) error {
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
