package clip

import (
	"GoSlice/data"
	"fmt"
	clipper "github.com/aligator/go.clipper"
)

type linear struct {
	paths clipper.Paths
}

func NewLinearPattern(min data.MicroPoint, max data.MicroPoint, lineWidth data.Micrometer) Pattern {
	lines := clipper.Paths{}
	numLine := 0
	// generate the lines
	for x := min.X(); x <= max.X(); x += lineWidth {
		// switch line direction based on even / odd
		if numLine%2 == 1 {
			lines = append(lines, clipper.Path{
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
			lines = append(lines, clipper.Path{
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

	return linear{lines}
}

func (p linear) Fill(layerNr int, paths data.LayerPart, outline data.LayerPart, lineWidth data.Micrometer, overlapPercentage int, additionalInternalOverlap int) data.Paths {
	cPath := clipperPath(paths.Outline())
	cHoles := clipperPaths(paths.Holes())

	// The inside overlap is for parts which are smaller than the outline.
	// These parts are overlapped a bit more to avoid linos which are printed only in the air.
	insideOverlap := float32(lineWidth) * (100.0 - float32(overlapPercentage+additionalInternalOverlap)) / 100.0

	// The perimeter overlap is the overlap into the outline.
	perimeterOverlap := float32(lineWidth) * (100.0 - float32(overlapPercentage)) / 100.0

	// generate infill with the full inside overlap
	var result = p.getInfill(cPath, cHoles, insideOverlap)

	// then clip the result by the outline, so that the big overlap from the inside is cut at the outline
	if outline == nil {
		outline = paths
	}

	cl := clipper.NewClipper(clipper.IoNone)

	// generate the exset for the overlap (only if needed)
	if perimeterOverlap != 0 {
		co := clipper.NewClipperOffset()
		co.AddPath(clipperPath(outline.Outline()), clipper.JtSquare, clipper.EtClosedPolygon)
		co.AddPaths(clipperPaths(outline.Holes()), clipper.JtSquare, clipper.EtClosedPolygon)
		co.MiterLimit = 2
		maxOutline := co.Execute(float64(-perimeterOverlap))
		cl.AddPaths(maxOutline, clipper.PtClip, true)
	} else {
		cl.AddPath(clipperPath(outline.Outline()), clipper.PtClip, true)
		cl.AddPaths(clipperPaths(outline.Holes()), clipper.PtClip, true)
	}

	cl.AddPaths(result, clipper.PtSubject, false)

	res, ok := cl.Execute2(clipper.CtIntersection, clipper.PftEvenOdd, clipper.PftEvenOdd)
	if !ok {
		return nil
	}

	// TODO: finish infill by connecting the individual line-ends:
	// For this the upper and lower points have to be offset by the line width/2.
	//   _    _
	// |  | |  |
	// |  | |  |
	// |  | |  |
	// |  |_|  |

	// convert result to data.Paths
	var resultInfill data.Paths
	parts, _ := polyTreeToLayerParts(res)
	for _, part := range parts {
		resultInfill = append(resultInfill, part.Outline())

		for _, path := range part.Holes() {
			resultInfill = append(resultInfill, path)
		}
	}
	return resultInfill
}

// getInfill fills a polygon (with holes)
func (p linear) getInfill(outline clipper.Path, holes clipper.Paths, overlap float32) clipper.Paths {
	var result clipper.Paths

	// clip the paths with the lines using intersection
	exset := clipper.Paths{outline}

	co := clipper.NewClipperOffset()
	cl := clipper.NewClipper(clipper.IoNone)

	// generate the exset for the overlap (only if needed)
	if overlap != 0 {
		co.AddPaths(exset, clipper.JtSquare, clipper.EtClosedPolygon)
		co.MiterLimit = 2
		exset = co.Execute(float64(-overlap))

		co.Clear()
		co.AddPaths(holes, clipper.JtSquare, clipper.EtClosedPolygon)
		co.MiterLimit = 2
		holes = co.Execute(float64(overlap))
	}

	// clip the lines by the resulting inset
	cl.AddPaths(exset, clipper.PtClip, true)
	cl.AddPaths(holes, clipper.PtClip, true)
	cl.AddPaths(p.paths, clipper.PtSubject, false)

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
