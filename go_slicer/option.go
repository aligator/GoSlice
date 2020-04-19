package go_slicer

import (
	"GoSlicer/util"
)

type option func(s *GoSlicer)

func (s *GoSlicer) With(o ...option) {
	for _, option := range o {
		option(s)
	}
}

// MeldDistance is the distance which two points have to be within to
// count them as one point.
func ḾeldDistance(m util.Micrometer) option {
	return func(s *GoSlicer) {
		s.o.MeldDistance = m
	}
}

// Center is the point where the model is finally placed
func Center(p util.MicroVec3) option {
	return func(s *GoSlicer) {
		s.o.Center = p
	}
}

// InitialLayerThickness is the layer thickness for the first layer
func InitialLayerThickness(m util.Micrometer) option {
	return func(s *GoSlicer) {
		s.o.Print.InitialLayerThickness = m
	}
}

// LayerThickness is the thickness for all but the first layer
func LayerThickness(m util.Micrometer) option {
	return func(s *GoSlicer) {
		s.o.Print.LayerThickness = m
	}
}

// JoinPolygonSnapDistance is the distance used to check if two open
// polygons can be snapped together to one bigger polygon.
// Checked by the start and endpoints of the polygons.
func JoinPolygonSnapDistance(m util.Micrometer) option {
	return func(s *GoSlicer) {
		s.o.JoinPolygonSnapDistance = m
	}
}

// FinishPolygonSnapDistance is the max distance between start end endpoint of
// a polygon used to check if a open polygon can be closed.
func FinishPolygonSnapDistance(m util.Micrometer) option {
	return func(s *GoSlicer) {
		s.o.FinishPolygonSnapDistance = m
	}
}

// FilamentDiameter is the filament diameter used by the printer
func FilamentDiameter(m util.Millimeter) option {
	return func(s *GoSlicer) {
		s.o.Filament.FilamentDiameter = m.ToMicrometer()
	}
}

// ExtrusionWidth is the diameter of your nozzle
func ExtrusionWidth(m util.Micrometer) option {
	return func(s *GoSlicer) {
		s.o.Printer.ExtrusionWidth = m
	}
}

// InsetCount is the number of perimeters
func InsetCount(n int) option {
	return func(s *GoSlicer) {
		s.o.Print.InsetCount = n
	}
}

// InitialLayerSpeed is speed only for the first layer in mm per second
func InitialLayerSpeed(mmPerS util.Millimeter) option {
	return func(s *GoSlicer) {
		s.o.Print.IntialLayerSpeed = mmPerS
	}
}

// LayerSpeed is speed only for all but the first layer in mm per second
func LayerSpeed(mmPerS util.Millimeter) option {
	return func(s *GoSlicer) {
		s.o.Print.LayerSpeed = mmPerS
	}
}
