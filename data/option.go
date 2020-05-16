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

	// InfillOverlapPercent says how much the infill should overlap into the perimeters.
	InfillOverlapPercent int

	// InfillPercent says how much infill should be generated.
	InfillPercent int
}

// FilamentOptions contains all Filament specific GoSlice options.
type FilamentOptions struct {
	// FilamentDiameter is the filament diameter used by the printer.
	FilamentDiameter Micrometer
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

			InitialLayerThickness: 200,
			LayerThickness:        200,
			InsetCount:            2,
			InfillOverlapPercent:  50,
			InfillPercent:         20,
		},
		Filament: FilamentOptions{
			FilamentDiameter: Millimeter(1.75).ToMicrometer(),
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
		JoinPolygonSnapDistance:   100,
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
	flag.Var(&options.Print.LayerThickness, "layer-thickness", "The layer thickness for the first layer.")
	flag.IntVar(&options.Print.InsetCount, "inset-count", options.Print.InsetCount, "The layer thickness for the first layer.")
	flag.IntVar(&options.Print.InfillOverlapPercent, "infill-overlap-percent", options.Print.InfillOverlapPercent, "The layer thickness for the first layer.")
	flag.IntVar(&options.Print.InfillPercent, "infill-percent", options.Print.InfillPercent, "The layer thickness for the first layer.")

	// filament options
	flag.Var(&options.Filament.FilamentDiameter, "filament-diameter", "The filament diameter used by the printer.")

	// printer options
	flag.Var(&options.Printer.ExtrusionWidth, "extrusion-width", "The diameter of your nozzle.")
	center := microVec3{}
	flag.Var(&center, "center", "The point where the model is finally placed.")

	flag.Parse()

	options.Printer.Center = &center

	if options.InputFilePath == "" {
		panic("you have to pass a filename using the --file flag")
	}

	return options
}
