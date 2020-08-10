package main

import (
	"GoSlice/clip"
	"GoSlice/data"
	"GoSlice/gcode"
	"GoSlice/gcode/renderer"
	"GoSlice/handler"
	"GoSlice/modifier"
	"GoSlice/optimizer"
	"GoSlice/reader"
	"GoSlice/slicer"
	"GoSlice/writer"
	"fmt"
	"time"
)

// GoSlice combines all logic  needed to slice
// a model and generate a GCode file.
type GoSlice struct {
	options   *data.Options
	reader    handler.ModelReader
	optimizer handler.ModelOptimizer
	slicer    handler.ModelSlicer
	modifiers []handler.LayerModifier
	generator handler.GCodeGenerator
	writer    handler.GCodeWriter
}

// NewGoSlice provides a GoSlice with all built in implementations.
func NewGoSlice(options data.Options) *GoSlice {
	s := &GoSlice{
		options: &options,
	}

	// create handlers
	topBottomPatternFactory := func(min data.MicroPoint, max data.MicroPoint) clip.Pattern {
		return clip.NewLinearPattern(options.Printer.ExtrusionWidth, options.Printer.ExtrusionWidth, min, max, options.Print.InfillRotationDegree, false)
	}

	s.reader = reader.Reader(&options)
	s.optimizer = optimizer.NewOptimizer(&options)
	s.slicer = slicer.NewSlicer(&options)
	s.modifiers = []handler.LayerModifier{
		modifier.NewPerimeterModifier(&options),
		modifier.NewInfillModifier(&options),
		modifier.NewInternalInfillModifier(&options),
	}
	s.generator = gcode.NewGenerator(
		&options,
		gcode.WithRenderer(renderer.PreLayer{}),
		gcode.WithRenderer(renderer.Perimeter{}),
		gcode.WithRenderer(&renderer.Infill{
			PatternSetup: topBottomPatternFactory,
			AttrName:     "bottom",
			Comments:     []string{"TYPE:FILL", "BOTTOM-FILL"},
		}),
		gcode.WithRenderer(&renderer.Infill{
			PatternSetup: topBottomPatternFactory,
			AttrName:     "top",
			Comments:     []string{"TYPE:FILL", "TOP-FILL"},
		}),
		gcode.WithRenderer(&renderer.Infill{
			PatternSetup: func(min data.MicroPoint, max data.MicroPoint) clip.Pattern {
				// TODO: the calculation of the percentage is currently very basic and may not be correct.

				if options.Print.InfillPercent != 0 {
					mm10 := data.Millimeter(10).ToMicrometer()
					linesPer10mmFor100Percent := mm10 / options.Printer.ExtrusionWidth
					linesPer10mmForInfillPercent := float64(linesPer10mmFor100Percent) * float64(options.Print.InfillPercent) / 100.0

					lineWidth := data.Micrometer(float64(mm10) / linesPer10mmForInfillPercent)

					return clip.NewLinearPattern(options.Printer.ExtrusionWidth, lineWidth, min, max, options.Print.InfillRotationDegree, options.Print.InfillZigZag)
				}

				return nil
			},
			AttrName: "infill",
			Comments: []string{"TYPE:FILL", "INTERNAL-FILL"},
		}),
		gcode.WithRenderer(renderer.PostLayer{}),
	)
	s.writer = writer.Writer()

	return s
}

func (s *GoSlice) Process() error {
	startTime := time.Now()

	// 1. Load model
	models, err := s.reader.Read(s.options.InputFilePath)
	if err != nil {
		return err
	}

	// 2. Optimize model
	var optimizedModel data.OptimizedModel

	optimizedModel, err = s.optimizer.Optimize(models)
	if err != nil {
		return err
	}

	err = optimizedModel.SaveDebugSTL("test.stl")
	if err != nil {
		return err
	}

	// 3. Slice model into layers
	layers, err := s.slicer.Slice(optimizedModel)
	if err != nil {
		return err
	}

	// 4. Modify the layers
	// e.g. generate perimeter paths,
	// generate the parts which should be filled in, ...
	for _, m := range s.modifiers {
		m.Init(optimizedModel)
		for layerNr := range layers {
			err = m.Modify(layerNr, layers)
			if err != nil {
				return err
			}
		}
	}

	// 5. generate gcode from the layers
	s.generator.Init(optimizedModel)
	finalGcode, err := s.generator.Generate(layers)
	if err != nil {
		return err
	}

	err = s.writer.Write(finalGcode, s.options.InputFilePath+".gcode")
	fmt.Println("full processing time:", time.Now().Sub(startTime))

	return err
}
