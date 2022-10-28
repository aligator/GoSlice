// This file provides renderers for gcode injected at specific layers.

package renderer

import (
	"fmt"
	"strings"

	"github.com/aligator/goslice/data"
	"github.com/aligator/goslice/gcode"
)

var (
	HotendTempCodes    = []string{"M104", "M109"}
	BedTempCodes       = []string{"M140", "M190"}
	BedTempTemplate    = "{print_bed_temperature}"
	HotendTempTemplate = "{print_temperature}"
)

// PreLayer adds starting gcode, resets the extrude speeds on each layer and enables the fan above a specific layer.
type PreLayer struct{}

func (PreLayer) Init(model data.OptimizedModel) {}

func (PreLayer) Render(b *gcode.Builder, layerNr int, maxLayer int, layer data.PartitionedLayer, z data.Micrometer, options *data.Options) error {
	if layerNr == 0 {
		b.AddComment("Generated with GoSlice")
		b.AddComment("______________________")

		if options.Printer.ForceSafeStartStopGCode {
			if options.Printer.HasHeatedBed && !options.Printer.StartGCode.DoesInstructionContainCodes(BedTempCodes) {
				b.AddComment("SET BED TEMP")
				b.AddCommand("M190 S%d ; heat and wait for bed", options.Filament.InitialBedTemperature)
			}

			if !options.Printer.StartGCode.DoesInstructionContainCodes(HotendTempCodes) {
				b.AddComment("SET HOTEND TEMP")
				b.AddCommand("M109 S%d ; wait for hot end temperature", options.Filament.InitialHotEndTemperature)
			}

		}
		b.AddComment("START GCODE")
		// starting gcode
		for _, instruction := range options.Printer.StartGCode.GCodeLines {
			if strings.Contains(instruction, BedTempTemplate) {
				instruction = strings.Replace(instruction, BedTempTemplate, fmt.Sprint(options.Filament.InitialBedTemperature), -1)
			}
			if strings.Contains(instruction, HotendTempTemplate) {
				instruction = strings.Replace(instruction, HotendTempTemplate, fmt.Sprint(options.Filament.InitialHotEndTemperature), -1)
			}
			b.AddCommand(instruction)
		}

		b.AddCommand("G92 E0 ; reset extrusion distance")

		b.SetExtrusion(options.Print.InitialLayerThickness, options.Printer.ExtrusionWidth)

		// set speeds
		b.SetExtrudeSpeed(options.Print.LayerSpeed)
		b.SetMoveSpeed(options.Print.MoveSpeed)

		// set retraction
		b.SetRetractionSpeed(options.Filament.RetractionSpeed)
		b.SetRetractionAmount(options.Filament.RetractionLength)
		b.SetRetractionZHop(options.Filament.RetractionZHop)

		// force the InitialLayerSpeed for first layer
		b.SetExtrudeSpeedOverride(options.Print.IntialLayerSpeed)
	} else if layerNr == 1 {
		b.SetExtrusion(options.Print.LayerThickness, options.Printer.ExtrusionWidth)
	}

	b.AddComment("LAYER:%v", layerNr)

	if layerNr > 0 {
		b.DisableExtrudeSpeedOverride()
		b.SetExtrudeSpeed(options.Print.LayerSpeed)
	}

	if fanSpeed, ok := options.Filament.FanSpeed.LayerToSpeedLUT[layerNr]; ok {
		if fanSpeed == 0 {
			b.AddCommand("M107 ; disable fan")
		} else {
			b.AddCommand("M106 S%d; change fan speed", fanSpeed)
		}
	}

	if layerNr == options.Filament.InitialTemperatureLayerCount {
		// set the normal temperature
		// this is done without waiting
		b.AddComment("SET_TEMP")
		b.AddCommand("M140 S%d", options.Filament.BedTemperature)
		b.AddCommand("M104 S%d", options.Filament.HotEndTemperature)
	}

	return nil
}

// PostLayer adds GCode at the last layer.
type PostLayer struct{}

func (PostLayer) Init(model data.OptimizedModel) {}

func (PostLayer) Render(b *gcode.Builder, layerNr int, maxLayer int, layer data.PartitionedLayer, z data.Micrometer, options *data.Options) error {
	// ending gcode
	if layerNr == maxLayer {
		b.AddComment("END_GCODE")
		b.SetExtrusion(options.Print.LayerThickness, options.Printer.ExtrusionWidth)

		if options.Printer.ForceSafeStartStopGCode {
			// disable heaters
			if !options.Printer.EndGCode.DoesInstructionContainCodes(HotendTempCodes) {
				b.AddCommand("M104 S0 ; Set Hot-end to 0C (off)")
			}

			if options.Printer.HasHeatedBed && !options.Printer.EndGCode.DoesInstructionContainCodes(BedTempCodes) {
				b.AddCommand("M140 S0 ; Set bed to 0C (off)")
			}
		}
		for _, instruction := range options.Printer.EndGCode.GCodeLines {
			b.AddCommand(instruction)
		}

	}

	return nil
}
