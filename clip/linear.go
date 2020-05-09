// This file implements a basic linear pattern infill.

package clip

import (
	"GoSlice/data"
	"fmt"

	clipper "github.com/aligator/go.clipper"
)

// linear provides an infill which consists of simple parallel lines.
// The direction of the lines is switching for each layer by 90°..
type linear struct {
	verticalPaths   clipper.Paths
	horizontalPaths clipper.Paths
}

// NewLinearPattern provides a simple linear infill pattern consisting of simple parallel lines.
// The direction of the lines is switching for each layer by 90°.
func NewLinearPattern(min data.MicroPoint, max data.MicroPoint, lineWidth data.Micrometer) Pattern {
	verticalLines := clipper.Paths{}
	numLine := 0
	// generate the verticalLines
	for x := min.X(); x <= max.X(); x += lineWidth {
		// switch line direction based on even / odd
		if numLine%2 == 1 {
			verticalLines = append(verticalLines, clipper.Path{
				&clipper.IntPoint{
					X: clipper.CInt(x),
					Y: clipper.CInt(max.Y()),
				},
				&clipper.IntPoint{
					X: clipper.CInt(x),
					Y: clipper.CInt(min.Y()),
				},
			})
		} else {
			verticalLines = append(verticalLines, clipper.Path{
				&clipper.IntPoint{
					X: clipper.CInt(x),
					Y: clipper.CInt(min.Y()),
				},
				&clipper.IntPoint{
					X: clipper.CInt(x),
					Y: clipper.CInt(max.Y()),
				},
			})
		}
		numLine++
	}

	horizontalLines := clipper.Paths{}
	numLine = 0
	// generate the verticalLines
	for y := min.Y(); y <= max.Y(); y += lineWidth {
		// switch line direction based on even / odd
		if numLine%2 == 1 {
			horizontalLines = append(horizontalLines, clipper.Path{
				&clipper.IntPoint{
					X: clipper.CInt(max.X()),
					Y: clipper.CInt(y),
				},
				&clipper.IntPoint{
					X: clipper.CInt(min.X()),
					Y: clipper.CInt(y),
				},
			})
		} else {
			horizontalLines = append(horizontalLines, clipper.Path{
				&clipper.IntPoint{
					X: clipper.CInt(min.X()),
					Y: clipper.CInt(y),
				},
				&clipper.IntPoint{
					X: clipper.CInt(max.X()),
					Y: clipper.CInt(y),
				},
			})
		}
		numLine++
	}

	return linear{verticalPaths: verticalLines, horizontalPaths: horizontalLines}
}

// Fill implements the Pattern interface by using simple linear lines as infill.
func (p linear) Fill(layerNr int, part data.LayerPart) data.Paths {
	resultInfill := p.getInfill(layerNr, clipperPath(part.Outline()), clipperPaths(part.Holes()), 0)
	return microPaths(resultInfill, false)
}

// getInfill fills a polygon (with holes)
func (p linear) getInfill(layerNr int, outline clipper.Path, holes clipper.Paths, overlap float32) clipper.Paths {
	var result clipper.Paths

	// clip the paths with the lines using intersection
	exset := clipper.Paths{outline}

	co := clipper.NewClipperOffset()
	cl := clipper.NewClipper(clipper.IoNone)

	// generate the ex-set for the overlap (only if needed)
	if overlap != 0 {
		co.AddPaths(exset, clipper.JtSquare, clipper.EtClosedPolygon)
		co.MiterLimit = 2
		exset = co.Execute(float64(-overlap))

		co.Clear()
		co.AddPaths(holes, clipper.JtSquare, clipper.EtClosedPolygon)
		co.MiterLimit = 2
		holes = co.Execute(float64(overlap))
	}

	// clip the lines by the outline and holes
	cl.AddPaths(exset, clipper.PtClip, true)
	cl.AddPaths(holes, clipper.PtClip, true)

	if layerNr%2 == 0 {
		cl.AddPaths(p.verticalPaths, clipper.PtSubject, false)
	} else {
		cl.AddPaths(p.horizontalPaths, clipper.PtSubject, false)
	}

	tree, ok := cl.Execute2(clipper.CtIntersection, clipper.PftEvenOdd, clipper.PftEvenOdd)
	if !ok {
		fmt.Println("getLinearFill failed")
		return nil
	}

	for _, c := range tree.Childs() {
		result = append(result, c.Contour())
	}

	return result
}
