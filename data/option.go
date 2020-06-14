// This file provides all options of GoSlice.

package data

import (
	"errors"
	"strconv"
	"strings"

	flag "github.com/spf13/pflag"
)

// implement the Value interface for all types which can occur in the options

func (m Micrometer) String() string {
	return strconv.FormatInt(int64(m), 10)
}

func (m *Micrometer) Set(s string) error {
	v, err := strconv.ParseInt(s, 0, 64)
	*m = Micrometer(v)
	return err
}

func (m Micrometer) Type() string {
	return "Micrometer"
}

func (m Millimeter) String() string {
	return strconv.FormatInt(int64(m), 10)
}

func (m *Millimeter) Set(s string) error {
	v, err := strconv.ParseInt(s, 0, 32)
	*m = Millimeter(v)
	return err
}

func (m Millimeter) Type() string {
	return "Millimeter"
}

func (v microVec3) String() string {
	return v.X().String() + "_" + v.Y().String() + "_" + v.Z().String()
}

func (v *microVec3) Set(s string) error {
	const errorMsg = "the string should contain three integers separated by _"
	parts := strings.Split(s, "_")
	if len(parts) != 3 {
		return errors.New(errorMsg)
	}

	result := microVec3{}

	if err := result.x.Set(parts[0]); err != nil {
		return errors.New(errorMsg)
	}

	if err := result.y.Set(parts[1]); err != nil {
		return errors.New(errorMsg)
	}

	if err := result.z.Set(parts[2]); err != nil {
		return errors.New(errorMsg)
	}

	v.x = result.x
	v.y = result.y
	v.z = result.z

	return nil
}

func (v microVec3) Type() string {
	return "Micrometer"
}

// PrintOptions contains all Print specific GoSlice options.
type PrintOptions struct {
	// InitialLayerSpeed is the speed only for the first layer in mm per second.
	IntialLayerSpeed Millimeter

	// LayerSpeed is the speed for all but the first layer in mm per second.
	LayerSpeed Millimeter

	// OuterPerimeterSpeed is the speed only for outer perimeters.
	OuterPerimeterSpeed Millimeter

	// MoveSpeed is the speed for all non printing moves.
	MoveSpeed Millimeter

	// InitialLayerThickness is the layer thickness for the first layer.
	InitialLayerThickness Micrometer

	// LayerThickness is the thickness for all but the first layer.
	LayerThickness Micrometer

	// InsetCount is the number of perimeters.
	InsetCount int

	// InfillOverlapPercent is the percentage of overlap into the perimeters.
	InfillOverlapPercent int

	// AdditionalInternalInfillOverlapPercent is the percentage used to make the internal
	// infill (infill not blocked by the perimeters) even bigger so that it grows a bit into the model.
	AdditionalInternalInfillOverlapPercent int

	// InfillPercent is the amount of infill which should be generated.
	InfillPercent int

	// InfillRotationDegree is the rotation used for the infill.
	InfillRotationDegree int

	// NumberBottomLayers is the amount of layers the bottom layers should grow into the model.
	NumberBottomLayers int

	// NumberBottomLayers is the amount of layers the bottom layers should grow into the model.
	NumberTopLayers int
}

// FilamentOptions contains all Filament specific GoSlice options.
type FilamentOptions struct {
	// FilamentDiameter is the filament diameter used by the printer.
	FilamentDiameter Micrometer

	// InitialBedTemperature is the temperature for the heated bed for the first layers.
	InitialBedTemperature int

	// InitialHotendTemperature is the temperature for the hot end for the first layers.
	InitialHotEndTemperature int

	// BedTemperature is the temperature for the heated bed after the first layers.
	BedTemperature int

	// HotEndTemperature is the temperature for the hot end after the first layers.
	HotEndTemperature int

	// InitialTemeratureLayerCount is the number of layers which use the initial temperatures.
	// After this amount of layers, the normal temperatures are used.
	InitialTemeratureLayerCount int
}

// PrinterOptions contains all Printer specific GoSlice options.
type PrinterOptions struct {
	// ExtrusionWidth is the diameter of your nozzle.
	ExtrusionWidth Micrometer

	// Center is the point where the model is finally placed.
	Center MicroVec3
}

// Options contains all GoSlice options.
type Options struct {
	Printer  PrinterOptions
	Filament FilamentOptions
	Print    PrintOptions

	// MeldDistance is the distance which two points have to be
	// within to count them as one point.
	MeldDistance Micrometer

	// JoinPolygonSnapDistance is the distance used to check if two open
	// polygons can be snapped together to one bigger polygon.
	// Checked by the start and endpoints of the polygons.
	JoinPolygonSnapDistance Micrometer

	// FinishPolygonSnapDistance is the max distance between start end endpoint of
	// a polygon used to check if a open polygon can be closed.
	FinishPolygonSnapDistance Micrometer

	// InputFilePath specifies the path to the input stl file.
	InputFilePath string
}

func DefaultOptions() Options {
	return Options{
		Print: PrintOptions{
			IntialLayerSpeed:    30,
			LayerSpeed:          60,
			OuterPerimeterSpeed: 40,
			MoveSpeed:           150,

			InitialLayerThickness:                  200,
			LayerThickness:                         200,
			InsetCount:                             2,
			InfillOverlapPercent:                   50,
			AdditionalInternalInfillOverlapPercent: 400,
			InfillPercent:                          20,
			InfillRotationDegree:                   45,
			NumberBottomLayers:                     3,
			NumberTopLayers:                        4,
		},
		Filament: FilamentOptions{
			FilamentDiameter:            Millimeter(1.75).ToMicrometer(),
			InitialBedTemperature:       60,
			InitialHotEndTemperature:    205,
			BedTemperature:              55,
			HotEndTemperature:           200,
			InitialTemeratureLayerCount: 3,
		},
		Printer: PrinterOptions{
			ExtrusionWidth: 400,
			Center: NewMicroVec3(
				Millimeter(100).ToMicrometer(),
				Millimeter(100).ToMicrometer(),
				0,
			),
		},

		MeldDistance:              30,
		JoinPolygonSnapDistance:   160,
		FinishPolygonSnapDistance: 1000,
	}
}

// ParseFlags parses the command line flags.
// It returns the default options but sets all passed options.
func ParseFlags() Options {
	options := DefaultOptions()

	flag.StringVar(&options.InputFilePath, "file", "", "The path to the input stl file.")

	flag.Var(&options.MeldDistance, "meld-distance", "The distance which two points have to be within to count them as one point.")
	flag.Var(&options.JoinPolygonSnapDistance, "join-polygon-snap-distance", "The distance used to check if two open polygons can be snapped together to one bigger polygon. Checked by the start and endpoints of the polygons.")
	flag.Var(&options.FinishPolygonSnapDistance, "finish-polygon-snap-distance", "The max distance between start end endpoint of a polygon used to check if a open polygon can be closed.")

	// print options
	flag.Var(&options.Print.IntialLayerSpeed, "initial-layer-speed", "The speed only for the first layer in mm per second.")
	flag.Var(&options.Print.LayerSpeed, "layer-speed", "The speed for all but the first layer in mm per second.")
	flag.Var(&options.Print.OuterPerimeterSpeed, "outer-perimeter-speed", "The speed only for outer perimeters.")
	flag.Var(&options.Print.MoveSpeed, "move-speed", "The speed for all non printing moves.")
	flag.Var(&options.Print.InitialLayerThickness, "initial-layer-thickness", "The layer thickness for the first layer.")
	flag.Var(&options.Print.LayerThickness, "layer-thickness", "The thickness for all but the first layer.")
	flag.IntVar(&options.Print.InsetCount, "inset-count", options.Print.InsetCount, "The number of perimeters.")
	flag.IntVar(&options.Print.InfillOverlapPercent, "infill-overlap-percent", options.Print.InfillOverlapPercent, "The percentage of overlap into the perimeters.")
	flag.IntVar(&options.Print.AdditionalInternalInfillOverlapPercent, "additional-internal-infill-overlap-percent", options.Print.AdditionalInternalInfillOverlapPercent, "The percentage used to make the internal infill (infill not blocked by the perimeters) even bigger so that it grows a bit into the model.")
	flag.IntVar(&options.Print.InfillPercent, "infill-percent", options.Print.InfillPercent, "The amount of infill which should be generated.")
	flag.IntVar(&options.Print.InfillRotationDegree, "infill-rotation-degree", options.Print.InfillRotationDegree, "The rotation used for the infill.")
	flag.IntVar(&options.Print.NumberBottomLayers, "number-bottom-layers", options.Print.NumberBottomLayers, "The amount of layers the bottom layers should grow into the model.")
	flag.IntVar(&options.Print.NumberTopLayers, "number-top-layers", options.Print.NumberTopLayers, "The amount of layers the bottom layers should grow into the model.")

	// filament options
	flag.Var(&options.Filament.FilamentDiameter, "filament-diameter", "The filament diameter used by the printer.")
	flag.IntVar(&options.Filament.InitialBedTemperature, "initial-bed-temperature", options.Filament.InitialBedTemperature, "The temperature for the heated bed for the first layers.")
	flag.IntVar(&options.Filament.InitialHotEndTemperature, "initial-hot-end-temperature", options.Filament.InitialHotEndTemperature, "The filament diameter used by the printer.")
	flag.IntVar(&options.Filament.BedTemperature, "bed-temperature", options.Filament.BedTemperature, "The temperature for the heated bed after the first layers.")
	flag.IntVar(&options.Filament.HotEndTemperature, "hot-end-temperature", options.Filament.HotEndTemperature, "The temperature for the hot end after the first layers.")
	flag.IntVar(&options.Filament.InitialTemeratureLayerCount, "initial-temperature-layer-count", options.Filament.InitialTemeratureLayerCount, "The number of layers which use the initial temperatures. After this amount of layers, the normal temperatures are used.")

	// printer options
	flag.Var(&options.Printer.ExtrusionWidth, "extrusion-width", "The diameter of your nozzle.")
	center := microVec3{
		options.Printer.Center.X(),
		options.Printer.Center.Y(),
		options.Printer.Center.Z(),
	}
	flag.Var(&center, "center", "The point where the model is finally placed.")

	flag.Parse()

	options.Printer.Center = &center

	if options.InputFilePath == "" {
		panic("you have to pass a filename using the --file flag")
	}

	return options
}
