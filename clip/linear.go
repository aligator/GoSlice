package clip

import (
	"GoSlice/data"
	"fmt"
	clipper "github.com/aligator/go.clipper"
)

type linear struct {
	paths  clipper.Paths
	paths2 clipper.Paths
}

// verticalLinesByX assumes that each LayerPart contains only a vertical line, specified by two points.
// It can sort them by the x value.
type verticalLinesByX []data.LayerPart

func (a verticalLinesByX) Len() int {
	return len(a)
}

func (a verticalLinesByX) Less(i, j int) bool {
	return a[i].Outline()[0].X() < a[j].Outline()[0].X()
}

func (a verticalLinesByX) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

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

	return linear{paths: verticalLines, paths2: horizontalLines}
}

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

	if layerNr%2 == 0 {
		cl.AddPaths(p.paths, clipper.PtSubject, false)
	} else {
		cl.AddPaths(p.paths2, clipper.PtSubject, false)
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
