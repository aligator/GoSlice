package gcode

import (
	"GoSlice/data"
	"bytes"
	"fmt"
	"math"
)

// Builder creates GCode by combining several commands.
type Builder struct {
	buf *bytes.Buffer

	extrusionAmount                                             data.Millimeter
	extrusionPerMM                                              data.Millimeter
	currentPosition                                             data.MicroVec3
	moveSpeed, extrudeSpeed, currentSpeed, extrudeSpeedOverride int
}

func NewGCodeBuilder(buf *bytes.Buffer) *Builder {
	g := &Builder{
		currentPosition: data.NewMicroVec3(0, 0, 0),
	}

	g.buf = buf
	return g
}

func (g *Builder) Buffer() *bytes.Buffer {
	return g.buf
}

func (g *Builder) SetExtrusion(layerThickness, lineWidth, filamentDiameter data.Micrometer) {
	filamentArea := data.Millimeter(math.Pi * (filamentDiameter.ToMillimeter() / 2.0) * (filamentDiameter.ToMillimeter() / 2.0))
	g.extrusionPerMM = layerThickness.ToMillimeter() * lineWidth.ToMillimeter() / filamentArea
}

func (g *Builder) SetMoveSpeed(moveSpeed data.Millimeter) {
	g.moveSpeed = int(moveSpeed)
}

func (g *Builder) SetExtrudeSpeed(extrudeSpeed data.Millimeter) {
	g.extrudeSpeed = int(extrudeSpeed)
}

func (g *Builder) SetExtrudeSpeedOverride(extrudeSpeed data.Millimeter) {
	g.extrudeSpeedOverride = int(extrudeSpeed)
}

func (g *Builder) DisableExtrudeSpeedOverride() {
	g.extrudeSpeedOverride = -1
}

func (g *Builder) AddCommand(command string, args ...interface{}) {
	command = command + "\n"
	command = fmt.Sprintf(command, args...)
	g.buf.WriteString(command)
}

func (g *Builder) AddComment(comment string, args ...interface{}) {
	comment = ";" + comment + "\n"
	comment = fmt.Sprintf(comment, args...)
	g.buf.WriteString(comment)
}

func (g *Builder) AddMove(p data.MicroVec3, extrusion data.Millimeter) {
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

func (g *Builder) AddPolygon(polygon data.Path, z data.Micrometer, open bool) {
	if len(polygon) == 0 {
		g.AddComment("ignore Too small polygon")
		return
	}

	// smooth the polygon
	polygon = data.DouglasPeucker(polygon, -1)

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

	// add the move from the last point to the first point only if the path is closed
	if open {
		return
	}

	point0 := data.NewMicroPoint(polygon[0].X(), polygon[0].Y())

	last := len(polygon) - 1
	pointLast := data.NewMicroPoint(polygon[last].X(), polygon[last].Y())

	g.AddMove(
		data.NewMicroVec3(polygon[0].X(), polygon[0].Y(), z),
		point0.Sub(pointLast).SizeMM()*g.extrusionPerMM,
	)
}
