package gcode

import (
	"bytes"
	"errors"
	"fmt"
	"math"

	"github.com/aligator/goslice/clip"
	"github.com/aligator/goslice/data"
)

// Builder creates GCode by combining several commands.
// Use  NewGCodeBuilder to create a new builder.
type Builder struct {
	buf *bytes.Buffer

	extrusionAmount                                             data.Millimeter
	extrusionPerMM                                              data.Millimeter
	currentPosition                                             data.MicroVec3
	notFirstMove                                                bool
	moveSpeed, extrudeSpeed, currentSpeed, extrudeSpeedOverride int

	retractionSpeed  int
	retractionAmount data.Millimeter
	zHopOnRetract    data.Millimeter

	filamentDiameter    data.Micrometer
	extrusionMultiplier int
}

func NewGCodeBuilder(options *data.Options) *Builder {
	g := &Builder{
		currentPosition:     data.NewMicroVec3(0, 0, 0),
		filamentDiameter:    options.Filament.FilamentDiameter,
		extrusionMultiplier: options.Filament.ExtrusionMultiplier,
	}
	g.buf = bytes.NewBuffer([]byte{})
	return g
}

func (g *Builder) String() string {
	return g.buf.String()
}

func (g *Builder) SetExtrusion(layerThickness, lineWidth data.Micrometer) {
	filamentArea := math.Pi * (g.filamentDiameter.ToMillimeter() / 2.0) * (g.filamentDiameter.ToMillimeter() / 2.0)
	g.extrusionPerMM = (layerThickness.ToMillimeter() * lineWidth.ToMillimeter() / filamentArea) * (data.Millimeter(g.extrusionMultiplier) / 100)
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
	g.extrudeSpeedOverride = 0
}

func (g *Builder) SetRetractionSpeed(retractionSpeed data.Millimeter) {
	g.retractionSpeed = int(retractionSpeed)
}

func (g *Builder) SetRetractionAmount(retractionAmount data.Millimeter) {
	g.retractionAmount = retractionAmount
}

func (g *Builder) SetRetractionZHop(zHopOnRetract data.Millimeter) {
	g.zHopOnRetract = zHopOnRetract
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

func (g *Builder) AddMoveSpeed(p data.MicroVec3, extrusion data.Millimeter, speed int) {
	// Ignore moves which are of zero length.
	if g.notFirstMove && g.currentPosition.X() == p.X() && g.currentPosition.Y() == p.Y() && g.currentPosition.Z() == p.Z() && extrusion == 0 {
		return
	}
	g.notFirstMove = true

	if extrusion != 0 {
		g.buf.WriteString("G1")
	} else {
		g.buf.WriteString("G0")
	}

	if p.X() != g.currentPosition.X() {
		g.buf.WriteString(fmt.Sprintf(" X%0.2f", p.X().ToMillimeter()))
	}
	if p.Y() != g.currentPosition.Y() {
		g.buf.WriteString(fmt.Sprintf(" Y%0.2f", p.Y().ToMillimeter()))
	}
	if p.Z() != g.currentPosition.Z() {
		g.buf.WriteString(fmt.Sprintf(" Z%0.2f", p.Z().ToMillimeter()))
	}

	if g.currentSpeed != speed {
		g.buf.WriteString(fmt.Sprintf(" F%v", speed*60))
		g.currentSpeed = speed
	}

	g.extrusionAmount += extrusion
	if extrusion != 0 {
		g.buf.WriteString(fmt.Sprintf(" E%0.4f", g.extrusionAmount))
	}
	g.buf.WriteString("\n")

	g.currentPosition = p

}

func (g *Builder) AddMove(p data.MicroVec3, extrusion data.Millimeter) {

	var speed int
	if extrusion != 0 {
		if g.extrudeSpeedOverride <= 0 {
			speed = g.extrudeSpeed
		} else {
			speed = g.extrudeSpeedOverride
		}
	} else {
		speed = g.moveSpeed
	}

	g.AddMoveSpeed(p, extrusion, speed)
}

func (g *Builder) AddPolygon(currentLayer data.PartitionedLayer, polygon data.Path, z data.Micrometer, open bool) error {
	if len(polygon) == 0 {
		return nil
	}

	// smooth the polygon
	polygon = data.DouglasPeucker(polygon, -1)

	for i, p := range polygon {
		if i == 0 {
			// for the move to the polygon: detect move through perimeters and add retraction if needed
			// TODO: this is very ineffective, as it has to clip for every first move of every polygon with the whole layer...
			move := data.Path{
				g.currentPosition.PointXY(),
				polygon[0],
			}

			isCrossing := false
			if currentLayer != nil && g.retractionSpeed != 0 && g.retractionAmount != 0 {
				c := clip.NewClipper()
				var ok bool
				isCrossing, ok = c.IsCrossingPerimeter(currentLayer.LayerParts(), move)

				if !ok {
					return errors.New("could not calculate the difference between the current layer and the non-extrusion-move")
				}
			}

			zMove := z

			if isCrossing {
				g.AddMoveSpeed(g.currentPosition, -g.retractionAmount, g.retractionSpeed)

				if g.zHopOnRetract > 0 {
					zMove = z + g.zHopOnRetract.ToMicrometer()

					g.AddMove(data.NewMicroVec3(
						g.currentPosition.X(),
						g.currentPosition.Y(),
						zMove,
					), 0.0)
				}
			}

			g.AddMove(data.NewMicroVec3(
				polygon[i].X(),
				polygon[i].Y(),
				zMove), 0.0)

			if isCrossing {
				if g.zHopOnRetract > 0 {
					g.AddMove(data.NewMicroVec3(
						polygon[i].X(),
						polygon[i].Y(),
						z,
					), 0.0)
				}

				g.AddMoveSpeed(g.currentPosition, g.retractionAmount, g.retractionSpeed)
			}

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
		return nil
	}

	point0 := data.NewMicroPoint(polygon[0].X(), polygon[0].Y())

	last := len(polygon) - 1
	pointLast := data.NewMicroPoint(polygon[last].X(), polygon[last].Y())

	g.AddMove(
		data.NewMicroVec3(polygon[0].X(), polygon[0].Y(), z),
		point0.Sub(pointLast).SizeMM()*g.extrusionPerMM,
	)

	return nil
}
