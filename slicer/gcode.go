package slicer

import (
	"GoSlicer/util"
	"fmt"
	clipper "github.com/ctessum/go.clipper"
	"math"
	"os"
)

type gcodeBuilder struct {
	buf *os.File

	extrusionAmount                       util.Millimeter
	extrusionPerMM                        util.Millimeter
	currentPosition                       util.MicroVec3
	moveSpeed, extrudeSpeed, currentSpeed int
}

func newGCodeBuilder(filename string) *gcodeBuilder {
	g := &gcodeBuilder{
		moveSpeed:       150,
		extrudeSpeed:    50,
		currentPosition: util.NewMicroVec3(0, 0, 0),
	}
	buf, err := os.Create(filename)
	if err != nil {
		fmt.Println(err)
	}

	g.buf = buf

	return g
}

func (g *gcodeBuilder) setExtrusion(layerThickness, lineWidth, filamentDiameter util.Micrometer) {
	filamentArea := util.Millimeter(math.Pi * (filamentDiameter.ToMillimeter() / 2.0) * (filamentDiameter.ToMillimeter() / 2.0))
	g.extrusionPerMM = layerThickness.ToMillimeter() * lineWidth.ToMillimeter() / filamentArea
}

func (g *gcodeBuilder) setSpeeds(moveSpeed, extrudeSpeed int) {
	g.moveSpeed = moveSpeed
	g.extrudeSpeed = extrudeSpeed
}

func (g *gcodeBuilder) addComment(comment string, args ...interface{}) {
	comment = ";" + comment + "\n"
	comment = fmt.Sprintf(comment, args...)
	g.buf.WriteString(comment)
}

func (g *gcodeBuilder) addMove(p util.MicroVec3, extrusion util.Millimeter) {
	var speed int
	if extrusion != 0 {
		g.buf.WriteString("G1")
		speed = g.extrudeSpeed
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

func (g *gcodeBuilder) addPolygon(polygon clipper.Path, z util.Micrometer) {
	if len(polygon) == 0 {
		g.addComment("ignore Too small polygon")
		return
	}
	for i, p := range polygon {
		if i == 0 {
			g.addMove(util.NewMicroVec3(
				util.Micrometer(polygon[0].X),
				util.Micrometer(polygon[0].Y),
				z), 0.0)
			continue
		}

		point := util.NewMicroPoint(
			util.Micrometer(p.X),
			util.Micrometer(p.Y))

		prevPoint := util.NewMicroPoint(
			util.Micrometer(polygon[i-1].X),
			util.Micrometer(polygon[i-1].Y))

		g.addMove(util.NewMicroVec3(
			util.Micrometer(p.X),
			util.Micrometer(p.Y),
			z),
			point.Sub(prevPoint).SizeMM()*g.extrusionPerMM,
		)
	}

	point0 := util.NewMicroPoint(
		util.Micrometer(polygon[0].X),
		util.Micrometer(polygon[0].Y))

	last := len(polygon) - 1
	pointLast := util.NewMicroPoint(
		util.Micrometer(polygon[last].X),
		util.Micrometer(polygon[last].Y))

	g.addMove(util.NewMicroVec3(
		util.Micrometer(polygon[0].X),
		util.Micrometer(polygon[0].Y),
		z),
		point0.Sub(pointLast).SizeMM()*g.extrusionPerMM,
	)
}

func (g *gcodeBuilder) Close() error {
	return g.buf.Close()
}
