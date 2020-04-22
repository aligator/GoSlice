package gcode

import (
	"GoSlice/data"
	"bytes"
	"fmt"
	"math"
)

type gcodeBuilder struct {
	buf *bytes.Buffer

	extrusionAmount                                             data.Millimeter
	extrusionPerMM                                              data.Millimeter
	currentPosition                                             data.MicroVec3
	moveSpeed, extrudeSpeed, currentSpeed, extrudeSpeedOverride int
}

func newGCodeBuilder(buf *bytes.Buffer) *gcodeBuilder {
	g := &gcodeBuilder{
		moveSpeed:       150,
		extrudeSpeed:    50,
		currentPosition: data.NewMicroVec3(0, 0, 0),
	}

	g.buf = buf
	return g
}

func (g *gcodeBuilder) setExtrusion(layerThickness, lineWidth, filamentDiameter data.Micrometer) {
	filamentArea := data.Millimeter(math.Pi * (filamentDiameter.ToMillimeter() / 2.0) * (filamentDiameter.ToMillimeter() / 2.0))
	g.extrusionPerMM = layerThickness.ToMillimeter() * lineWidth.ToMillimeter() / filamentArea
}

func (g *gcodeBuilder) setMoveSpeed(moveSpeed data.Millimeter) {
	g.moveSpeed = int(moveSpeed)
}

func (g *gcodeBuilder) setExtrudeSpeed(extrudeSpeed data.Millimeter) {
	g.extrudeSpeed = int(extrudeSpeed)
}

func (g *gcodeBuilder) setExtrudeSpeedOverride(extrudeSpeed data.Millimeter) {
	g.extrudeSpeedOverride = int(extrudeSpeed)
}

func (g *gcodeBuilder) disableExtrudeSpeedOverride() {
	g.extrudeSpeedOverride = -1
}

func (g *gcodeBuilder) addCommand(command string, args ...interface{}) {
	command = command + "\n"
	command = fmt.Sprintf(command, args...)
	g.buf.WriteString(command)
}

func (g *gcodeBuilder) addComment(comment string, args ...interface{}) {
	comment = ";" + comment + "\n"
	comment = fmt.Sprintf(comment, args...)
	g.buf.WriteString(comment)
}

func (g *gcodeBuilder) addMove(p data.MicroVec3, extrusion data.Millimeter) {
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

func (g *gcodeBuilder) addPolygon(polygon data.Path, z data.Micrometer) {
	if len(polygon) == 0 {
		g.addComment("ignore Too small polygon")
		return
	}
	for i, p := range polygon {
		if i == 0 {
			g.addMove(data.NewMicroVec3(
				polygon[0].X(),
				polygon[0].Y(),
				z), 0.0)
			continue
		}

		point := data.NewMicroPoint(p.X(), p.Y())

		prevPoint := data.NewMicroPoint(polygon[i-1].X(), polygon[i-1].Y())

		g.addMove(
			data.NewMicroVec3(p.X(), p.Y(), z),
			point.Sub(prevPoint).SizeMM()*g.extrusionPerMM,
		)
	}

	point0 := data.NewMicroPoint(polygon[0].X(), polygon[0].Y())

	last := len(polygon) - 1
	pointLast := data.NewMicroPoint(polygon[last].X(), polygon[last].Y())

	g.addMove(
		data.NewMicroVec3(polygon[0].X(), polygon[0].Y(), z),
		point0.Sub(pointLast).SizeMM()*g.extrusionPerMM,
	)
}
