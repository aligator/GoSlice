// Package builder provides a GCode builder for creating GCode files

package builder

import (
	"GoSlice/data"
	"bytes"
	"fmt"
	"math"
)

// Builder provides methods for building gcode files.
// The result can be obtained by calling Buffer().
type Builder interface {
	Buffer() *bytes.Buffer

	SetExtrusion(layerThickness, lineWidth, filamentDiameter data.Micrometer)
	SetMoveSpeed(moveSpeed data.Millimeter)
	SetExtrudeSpeed(extrudeSpeed data.Millimeter)
	SetExtrudeSpeedOverride(extrudeSpeed data.Millimeter)
	DisableExtrudeSpeedOverride()
	AddCommand(command string, args ...interface{})
	AddComment(comment string, args ...interface{})
	AddMove(p data.MicroVec3, extrusion data.Millimeter)
	AddPolygon(polygon data.Path, z data.Micrometer)
}

// GCode is an implementation of the Builder interface
type GCode struct {
	buf *bytes.Buffer

	extrusionAmount                                             data.Millimeter
	extrusionPerMM                                              data.Millimeter
	currentPosition                                             data.MicroVec3
	moveSpeed, extrudeSpeed, currentSpeed, extrudeSpeedOverride int
}

func NewGCodeBuilder(buf *bytes.Buffer) Builder {
	g := &GCode{
		currentPosition: data.NewMicroVec3(0, 0, 0),
	}

	g.buf = buf
	return g
}

func (g *GCode) Buffer() *bytes.Buffer {
	return g.buf
}

func (g *GCode) SetExtrusion(layerThickness, lineWidth, filamentDiameter data.Micrometer) {
	filamentArea := data.Millimeter(math.Pi * (filamentDiameter.ToMillimeter() / 2.0) * (filamentDiameter.ToMillimeter() / 2.0))
	g.extrusionPerMM = layerThickness.ToMillimeter() * lineWidth.ToMillimeter() / filamentArea
}

func (g *GCode) SetMoveSpeed(moveSpeed data.Millimeter) {
	g.moveSpeed = int(moveSpeed)
}

func (g *GCode) SetExtrudeSpeed(extrudeSpeed data.Millimeter) {
	g.extrudeSpeed = int(extrudeSpeed)
}

func (g *GCode) SetExtrudeSpeedOverride(extrudeSpeed data.Millimeter) {
	g.extrudeSpeedOverride = int(extrudeSpeed)
}

func (g *GCode) DisableExtrudeSpeedOverride() {
	g.extrudeSpeedOverride = -1
}

func (g *GCode) AddCommand(command string, args ...interface{}) {
	command = command + "\n"
	command = fmt.Sprintf(command, args...)
	g.buf.WriteString(command)
}

func (g *GCode) AddComment(comment string, args ...interface{}) {
	comment = ";" + comment + "\n"
	comment = fmt.Sprintf(comment, args...)
	g.buf.WriteString(comment)
}

func (g *GCode) AddMove(p data.MicroVec3, extrusion data.Millimeter) {
	var speed int
	if extrusion != 0 {
		g.buf.WriteString("G1")

		if g.extrudeSpeedOverride == -1 {
			speed = g.extrudeSpeed
		} else {
			speed = g.extrudeSpeedOverride
		}
	} else {
		g.buf.WriteString("G0")
		speed = g.moveSpeed
	}

	if g.currentSpeed != speed {
		g.buf.WriteString(fmt.Sprintf(" F%v", speed*60))
		g.currentSpeed = speed
	}

	g.buf.WriteString(fmt.Sprintf(" X%0.2f Y%0.2f", p.X().ToMillimeter(), p.Y().ToMillimeter()))
	if p.Z() != g.currentPosition.Z() {
		g.buf.WriteString(fmt.Sprintf(" Z%0.2f Y%0.2f", p.Z().ToMillimeter(), p.Y().ToMillimeter()))
	}

	g.extrusionAmount += extrusion
	if extrusion != 0 {
		g.buf.WriteString(fmt.Sprintf(" E%0.4f", g.extrusionAmount))
	}
	g.buf.WriteString("\n")

	g.currentPosition = p
}

func (g *GCode) AddPolygon(polygon data.Path, z data.Micrometer) {
	if len(polygon) == 0 {
		g.AddComment("ignore Too small polygon")
		return
	}
	for i, p := range polygon {
		if i == 0 {
			g.AddMove(data.NewMicroVec3(
				polygon[0].X(),
				polygon[0].Y(),
				z), 0.0)
			continue
		}

		point := data.NewMicroPoint(p.X(), p.Y())

		prevPoint := data.NewMicroPoint(polygon[i-1].X(), polygon[i-1].Y())

		g.AddMove(
			data.NewMicroVec3(p.X(), p.Y(), z),
			point.Sub(prevPoint).SizeMM()*g.extrusionPerMM,
		)
	}

	point0 := data.NewMicroPoint(polygon[0].X(), polygon[0].Y())

	last := len(polygon) - 1
	pointLast := data.NewMicroPoint(polygon[last].X(), polygon[last].Y())

	g.AddMove(
		data.NewMicroVec3(polygon[0].X(), polygon[0].Y(), z),
		point0.Sub(pointLast).SizeMM()*g.extrusionPerMM,
	)
}
