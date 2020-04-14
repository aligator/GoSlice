package slicer

import (
	"GoSlicer/model"
	"GoSlicer/util"
)

type layer struct {
	segments           []*segment
	faceToSegmentIndex map[int]int
	polygons           []*polygon
}

func NewLayer() *layer {
	return &layer{
		faceToSegmentIndex: map[int]int{},
	}
}

func (l *layer) makePolygons(om model.OptimizedModel) {
	// try for each segment to generate a polygon with other segments
	// if the segment is not already assigned to another polygon
	for startSegmentIndex, segment := range l.segments {
		if segment.addedToPolygon {
			continue
		}

		var polygon = &polygon{}
		polygon.points = append(polygon.points, l.segments[startSegmentIndex].start)

		currentSegmentIndex := startSegmentIndex
		canClose := false

		for {
			canClose = false
			currentSegment := l.segments[currentSegmentIndex]
			currentSegment.addedToPolygon = true
			p0 := currentSegment.end
			polygon.points = append(polygon.points, p0)

			nextIndex := -1
			// get the whole face for the index
			face := om.Faces()[currentSegment.faceIndex]

			// For each touching face of the current face
			// check if touching face is in this layer.
			// Then calculate the difference between the current end-point (p0)
			// and the starting points of the segments of the touching faces.
			// if it is below the threshold
			// * check if it is the same segment as the starting segment of this round
			//   -> close it as a polygon is finished
			// * if the segment is already added just continue
			// then set the next index to the touching segment
			for _, touchingFaceIndex := range face.Touching() {
				touchingSegmentIndex, ok := l.faceToSegmentIndex[touchingFaceIndex]
				if touchingFaceIndex > -1 && ok {
					p1 := l.segments[touchingSegmentIndex].start
					diff := p0.Sub(p1)

					if diff.ShorterThan(30) {
						if touchingSegmentIndex == startSegmentIndex {
							canClose = true
						}
						if l.segments[touchingSegmentIndex].addedToPolygon {
							continue
						}
						nextIndex = touchingSegmentIndex
					}
				}
			}

			if nextIndex == -1 {
				break
			}

			currentSegmentIndex = nextIndex
		}

		if canClose {
			polygon.closed = true
		}

		l.polygons = append(l.polygons, polygon)
	}

	snapDistance := util.Micrometer(100)
	// Connect polygons that are not closed yet.
	// As models are not always perfect manifold we need to join
	// some stuff up to get proper polygons.
	for i, polygon := range l.polygons {
		if polygon == nil || polygon.closed {
			continue
		}

		best := -1
		bestScore := snapDistance + 1
		for j, polygon2 := range l.polygons {
			if polygon2 == nil || polygon2.closed || i == j {
				continue
			}

			// check the distance of the last point from the first unfinished polygon
			// with the first point of the second unfinished polygon
			diff := polygon.points[len(polygon.points)-1].Sub(polygon2.points[0])
			if diff.ShorterThan(snapDistance) {
				score := diff.Size() - util.Micrometer(len(polygon2.points)*10)
				if score < bestScore {
					best = j
					bestScore = score
				}
			}
		}

		// if a matching polygon was found, connect them
		if best > -1 {
			for _, aPointFromBest := range l.polygons[best].points {
				polygon.points = append(polygon.points, aPointFromBest)
			}

			// close polygon if the start end end now fits inside the snap distance
			if polygon.isAlmostFinished(snapDistance) {
				polygon.removeLastPoint()
				polygon.closed = true
			}

			// erase the merged polygon
			l.polygons[best] = nil
		}
	}

	var clearedPolygons []*polygon
	snapDistance = util.Micrometer(1000)
	// do not use range to allow modifying i when deleting
	for i, poly := range l.polygons {
		if poly == nil {
			continue
		}

		// check if polygon is almost finished
		// if yes just finish it
		if poly.isAlmostFinished(snapDistance) {
			poly.removeLastPoint()
			poly.closed = true
		}

		// remove tiny polygons or not closed polygons
		length := util.Micrometer(0)
		for n, point := range poly.points {
			// ignore first point
			if n == 0 {
				continue
			}

			length += point.Sub(poly.points[n-1]).Size()
			if poly.closed && length > snapDistance {
				break
			}

		}

		// remove already cleared polygons and filter also not closed / too small ones
		if l.polygons[i] != nil && length > snapDistance && poly.closed {
			clearedPolygons = append(clearedPolygons, l.polygons[i])
		}
	}

	l.polygons = clearedPolygons
}
