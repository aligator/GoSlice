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
	lineDistance data.Micrometer
	lineWidth    data.Micrometer
	degree       int
	min, max     data.MicroPoint
	rectlinear   bool
	zigZag       bool
}

// NewLinearPattern provides a simple linear infill pattern consisting of simple parallel lines.
// The direction of the lines is switching for each layer by 90°.
func NewLinearPattern(lineWidth data.Micrometer, lineDistance data.Micrometer, min data.MicroPoint, max data.MicroPoint, degree int, zigZag bool) Pattern {
	return linear{
		lineDistance: lineDistance,
		lineWidth:    lineWidth,
		degree:       degree,
		min:          min,
		max:          max,
		zigZag:       zigZag,
	}
}

// Fill implements the Pattern interface by using simple linear lines as infill.
func (p linear) Fill(layerNr int, part data.LayerPart) data.Paths {
	rotation := float64(p.degree)

	if layerNr%2 == 0 {
		rotation += 90
	}

	// copy holes and outline as the original layer part should not be modified by the rotation (slices are passed by reference)
	var holes = data.Paths{}
	for _, points := range part.Holes() {
		var copied = make(data.Path, len(points))
		copy(copied, points)
		holes = append(holes, copied)
	}
	var outline = make(data.Path, len(part.Outline()))
	copy(outline, part.Outline())

	// rotate them
	outline.Rotate(rotation)
	holes.Rotate(rotation)

	// create rectangle for the max bounding box and rotate it,
	// then get the min and max from the rotated bounding rectangle.
	bounds := data.Path{
		p.min,
		data.NewMicroPoint(p.max.X(), p.min.Y()),
		p.max,
		data.NewMicroPoint(p.min.X(), p.max.Y()),
	}
	bounds.Rotate(rotation)
	min, max := bounds.Bounds()

	smallerLines := data.Micrometer(0)
	if p.zigZag {
		smallerLines = p.lineWidth
	}

	resultInfill := p.getInfill(min, max, clipperPath(outline), clipperPaths(holes), 0, smallerLines)

	result := p.sortInfill(microPaths(resultInfill, false), p.zigZag, data.NewBasicLayerPart(outline, holes))

	result.Rotate(-rotation)

	return result
}

// sortInfill optimizes the order of the infill lines.
func (p linear) sortInfill(unsorted data.Paths, zigZag bool, part data.LayerPart) data.Paths {
	if len(unsorted) == 0 {
		return unsorted
	}

	cl := NewClipper()

	// Save all sorted paths here.
	sorted := data.Paths{unsorted[0]}

	// save the amount of lines already saved without the extra zigZag lines
	savedPointsNum := 1

	// Saves already used indices.
	isUsed := make([]bool, len(unsorted))
	isUsed[0] = true

	// Saves the last path to know where to continue.
	lastIndex := 0

	// Save if the first(0) or second(1) point from the lastPath was the last point.
	lastPoint := 0

	for savedPointsNum < len(unsorted) {
		point := unsorted[lastIndex][lastPoint]

		bestIndex := -1
		bestDiff := data.Micrometer(-1)

		// get the line with the nearest point (of the same side)
		for i, line := range unsorted {
			if isUsed[i] {
				continue
			}

			point2 := line[lastPoint]

			differenceVec := point.Sub(point2)
			if bestDiff == -1 || differenceVec.ShorterThanOrEqual(bestDiff) {
				bestIndex = i
				bestDiff = differenceVec.Size()
				continue
			}
		}

		if bestIndex > -1 {
			lastIndex = bestIndex

			if zigZag {
				p1 := sorted[len(sorted)-1][1]
				p2 := unsorted[lastIndex][1-lastPoint]

				if p1.Sub(p2).ShorterThanOrEqual(p.lineWidth + p.lineDistance*2) {

					connectionLine := []data.MicroPoint{p1, p2}

					isCrossing, ok := cl.IsCrossingPerimeter([]data.LayerPart{part}, connectionLine)

					if !ok {
						// TODO: return error
						panic("could not calculate the difference between the current layer and the non-extrusion-move")
					}

					if !isCrossing {
						sorted = append(sorted, connectionLine)
					}
				}
			}

			sorted = append(sorted, unsorted[lastIndex])
			savedPointsNum++
			isUsed[bestIndex] = true
			lastPoint = 1 - lastPoint
		} else {
			// should never go here
			panic("there should always be a bestIndex > -1")
		}

		if lastPoint == 1 {
			sorted[len(sorted)-1] = []data.MicroPoint{
				sorted[len(sorted)-1][1],
				sorted[len(sorted)-1][0],
			}
		}
	}

	return sorted
}

// getInfill fills a polygon (with holes)
func (p linear) getInfill(min data.MicroPoint, max data.MicroPoint, outline clipper.Path, holes clipper.Paths, overlap float32, smallerLines data.Micrometer) clipper.Paths {
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

	verticalLines := clipper.Paths{}
	numLine := 0
	// generate the verticalLines
	for x := min.X(); x <= max.X(); x += p.lineDistance {
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
		numLine++
	}

	cl.AddPaths(verticalLines, clipper.PtSubject, false)

	tree, ok := cl.Execute2(clipper.CtIntersection, clipper.PftEvenOdd, clipper.PftEvenOdd)
	if !ok {
		fmt.Println("getLinearFill failed")
		return nil
	}

	for _, c := range tree.Childs() {
		if smallerLines != 0 {
			// shorten the lines if smallerLines is set
			p1 := c.Contour()[0]
			p2 := c.Contour()[1]

			// shorten them by the half value on each side
			// only do this if the line is bigger than the smallerLines value
			if p1.Y-clipper.CInt(smallerLines)/2 > p2.Y {
				p1.Y = p1.Y - clipper.CInt(smallerLines)/2
				p2.Y = p2.Y + clipper.CInt(smallerLines)/2
			}

			result = append(result, []*clipper.IntPoint{p1, p2})
		} else {
			result = append(result, c.Contour())
		}
	}

	return result
}
