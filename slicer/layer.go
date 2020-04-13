package slicer

import (
	"GoSlicer/model"
	"GoSlicer/util"
)

type Layer interface {
	Segments() []Segment
	addSegment(seg Segment)
	FaceToSegment(index int) (int, bool)
	setFaceToSegment(faceIndex int, segmentIndex int)
	makePolygons(om model.OptimizedModel)
}

type layer struct {
	segments           []Segment
	faceToSegmentIndex map[int]int
	polygons           []Polygon
}

func NewLayer() Layer {
	return &layer{
		faceToSegmentIndex: map[int]int{},
	}
}

func (l *layer) FaceToSegment(index int) (int, bool) {
	val, ok := l.faceToSegmentIndex[index]
	return val, ok
}

func (l *layer) setFaceToSegment(faceIndex int, segmentIndex int) {
	l.faceToSegmentIndex[faceIndex] = segmentIndex
}

func (l *layer) Segments() []Segment {
	return l.segments
}

func (l *layer) addSegment(seg Segment) {
	l.segments = append(l.segments, seg)
}

func (l layer) makePolygons(om model.OptimizedModel) {
	startSegment := 0
	// try for each segment to generate a polygon with other segments
	// if the segment is not already assigned to another polygon
	for _, segment := range l.segments {
		if segment.isAddedToPolygon() {
			continue
		}

		var polygon Polygon = &polygon{}
		polygon.addPoint(l.segments[startSegment].Start())

		segmentIndex := startSegment
		canClose := false

		for {
			canClose = false
			currentSegment := l.segments[segmentIndex]
			currentSegment.setAddedToPolygon(true)
			p0 := currentSegment.End()
			polygon.addPoint(p0)

			nextIndex := -1
			// get the whole face for the index
			face := om.Faces()[currentSegment.FaceIndex()]

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
				segmentIndex, ok := l.FaceToSegment(touchingFaceIndex)
				if touchingFaceIndex > -1 && ok {
					p1 := l.segments[segmentIndex].Start()
					diff := p0.Sub(p1)

					if diff.ShorterThan(30) {
						if segmentIndex == startSegment {
							canClose = true
						}
						if l.segments[segmentIndex].isAddedToPolygon() {
							continue
						}
						nextIndex = segmentIndex
					}
				}
			}

			if nextIndex == -1 {
				break
			}

			segmentIndex = nextIndex
		}

		if canClose {
			polygon.close()
		}

		l.polygons = append(l.polygons, polygon)
	}

	snapDistance := util.Micrometer(100)
	// Connect polygons that are not closed yet.
	// As models are not always perfect manifold we need to join
	// some stuff up to get proper polygons.
	for i, polygon := range l.polygons {
		if polygon == nil || polygon.isClosed() {
			continue
		}

		best := -1
		bestScore := snapDistance + 1
		for j, polygon2 := range l.polygons {
			if polygon == nil || polygon2.isClosed() || i == j {
				continue
			}

			// check the distance of the last point from the first unfinished polygon
			// with the first point of the second unfinished polygon
			diff := polygon.Points()[len(polygon.Points())-1].Sub(polygon2.Points()[0])
			if diff.ShorterThan(snapDistance) {
				score := diff.Size() - util.Micrometer(len(polygon2.Points())*10)
				if score < bestScore {
					best = j
					bestScore = score
				}
			}
		}

		// if a matching polygon was found, connect them
		if best > -1 {
			for _, aPointFromBest := range l.polygons[best].Points() {
				polygon.addPoint(aPointFromBest)
			}

			// close polygon if the start end end now fits inside the snap distance
			if polygon.isAlmostFinished(snapDistance) {
				polygon.removeLastPoint()
				polygon.close()
			}

			// erase the merged polygon
			l.removePolygon(best)
		}
	}

	snapDistance = 1000
	// do not use range to allow modifying i when deleting
	for i := 0; i < len(l.polygons); i++ {
		polygon := l.polygons[i]
		// check if polygon is almost finished
		// if yes just finish it
		if polygon.isAlmostFinished(snapDistance) {
			polygon.removeLastPoint()
			polygon.close()
		}

		// remove tiny polygons or not closed polygons
		length := util.Micrometer(0)
		for n, point := range polygon.Points() {
			// ignore first point
			if n == 0 {
				continue
			}

			length += point.Sub(polygon.Points()[n-1]).Size()
			if polygon.isClosed() && length > snapDistance {
				break
			}

		}

		if length < snapDistance || !polygon.isClosed() {
			l.removePolygon(i)
			i--
		}
	}
}

func (l *layer) removePolygon(i int) {
	l.polygons[i] = l.polygons[len(l.polygons)-1] // Copy last element to index i.
	l.polygons[len(l.polygons)-1] = nil           // Erase last element (write zero value).
	l.polygons = l.polygons[:len(l.polygons)-1]   // Truncate slice.
}
