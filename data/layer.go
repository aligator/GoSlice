// This file provides some types which are needed to represent layers in the different slicing steps.

package data

// Path is a simple list of points.
// It can be used to represent polygons (if they are closed) or just lines.
type Path []MicroPoint

// IsAlmostFinished returns true if the path represents an almost closed polygon.
// It checks if the distance between the first and last point is smaller
// than the given threshold distance.
func (p Path) IsAlmostFinished(distance Micrometer) bool {
	return p[0].Sub(p[len(p)-1]).ShorterThanOrEqual(distance)
}

// Simplify removes consecutive line segments with same orientation and changes this polygon.
// If a parameter is -1 a default value is used.
//
// Removes verts which are connected to line segments which are both too small.
// Removes verts which detour from a direct line from the previous and next vert by a too small amount.
//
// Criteria:
// 1. Never remove a vertex if either of the connected segments is larger than \p smallest_line_segment
// 2. Never remove a vertex if the distance between that vertex and the final resulting polygon would be higher than \p allowed_error_distance
// 3. Simplify uses a heuristic and doesn't necessarily remove all removable vertices under the above criteria.
// 4. But simplify may never violate these criteria.
// 5. Unless the segments or the distance is smaller than the rounding error of 5 micron
//
// smallestLineSegmentSquared is the maximal squared length of removed line segments
// allowedErrorDistanceSquared is the square of the distance of the middle point to the line segment of the consecutive and previous point for which the middle point is removed
//
// Note: this is directly ported from a newer CuraEngine version.
func (p Path) Simplify(smallestLineSegmentSquared, allowedErrorDistanceSquared Micrometer) Path {
	if smallestLineSegmentSquared == -1 {
		smallestLineSegmentSquared = 100
	}

	if allowedErrorDistanceSquared == -1 {
		allowedErrorDistanceSquared = 25
	}

	if len(p) <= 2 {
		return Path{}
	}
	if len(p) == 3 {
		return p
	}

	newPath := Path{}
	previous := p[len(p)-1]
	current := p[0]

	/*
		When removing a vertex, we check the height of the triangle of the area
		being removed from the original polygon by the simplification. However,
		when consecutively removing multiple vertices the height of the previously
		removed vertices w.r.t. the shortcut path changes.
		In order to not recompute the new height value of previously removed
		vertices we compute the height of a representative triangle, which covers
		the same amount of area as the area being cut off. We use the Shoelace
		formula to accumulate the area under the removed segments. This works by
		computing the area in a 'fan' where each of the blades of the fan go from
		the origin to one of the segments. While removing vertices the area in
		this fan accumulates. By subtracting the area of the blade connected to
		the shortcutting segment we obtain the total area of the cutoff region.
		From this area we compute the height of the represenatative triangle
		using the standard formula for a triangle area: A = .5*b*h
	*/

	// Twice the Shoelace formula for area of polygon per line segment.
	areaRemoved := previous.X()*current.Y() - previous.Y()*current.X()

	for i := 0; i < len(p); i++ {
		current = p[i%len(p)]

		// Check if the accumulated area doesn't exceed the maximum.
		var next MicroPoint

		switch {
		case i+1 < len(p):
			next = p[i+1]

		// don't spill over if the [next] vertex will then be equal to [previous]
		case i+1 == len(p) && len(newPath) > 1:
			next = newPath[0] // Spill over to new polygon for checking removed area.

		default:
			next = p[(i+1)%len(p)]
		}

		// twice the Shoelace formula for area of polygon per line segment.
		areaRemoveNext := current.X()*next.Y() - current.Y()*next.X()

		// area between the origin and the shortcutting segment
		negativeAreaClosing := next.X()*previous.Y() - next.Y()*previous.X()

		areaRemoved += areaRemoveNext

		length2 := current.Sub(previous).Size2()
		nextLength2 := current.Sub(next).Size2()

		// close the shortcut area polygon
		areaRemovedSoFar := areaRemoved + negativeAreaClosing

		baseLength2 := next.Sub(previous).Size2()

		// Two line segments form a line back and forth with no area.
		if baseLength2 == 0 {
			continue // Remove the vertex.
		}

		// We want to check if the height of the triangle formed by previous, current and next vertices is less than allowedErrorDistanceSquared.
		// 1/2 L = A           [actual area is half of the computed shoelace value] // Shoelace formula is .5*(...) , but we simplify the computation and take out the .5
		// A = 1/2 * b * h     [triangle area formula]
		// L = b * h           [apply above two and take out the 1/2]
		// h = L / b           [divide by b]
		// h^2 = (L / b)^2     [square it]
		// h^2 = L^2 / b^2     [factor the divisor]
		height2 := areaRemovedSoFar * areaRemovedSoFar / baseLength2
		if (height2 <= 25 && //Almost exactly collinear (barring rounding errors).
			PerpendicularDistance2(current, previous, next) <= 25) ||
			(length2 < smallestLineSegmentSquared &&
				nextLength2 < smallestLineSegmentSquared && // segments are small
				height2 <= allowedErrorDistanceSquared) { // removing the vertex doesn't introduce too much error.

			continue // remove the vertex
		}

		// don't remove vertex

		// so that in the next iteration it's the area between the origin, [previous] and [current]
		areaRemoved = areaRemoveNext
		previous = current // Note that "previous" is only updated if we don't remove the vertex.
		newPath = append(newPath, current)
	}
	return newPath
}

// Bounds calculates the bounding box of the Path
// The returned points are the min-X-Y-Point and the max-X-Y-Point.
func (p Path) Bounds() (MicroPoint, MicroPoint) {
	if len(p) == 0 {
		return NewMicroPoint(0, 0), NewMicroPoint(0, 0)
	}

	minX := MaxMicrometer
	minY := MaxMicrometer

	maxX := MinMicrometer
	maxY := MinMicrometer

	// return 0, 0, 0, 0 if everything is empty
	any := false
	for _, point := range p {
		any = true
		if point.X() < minX {
			minX = point.X()
		}
		if point.X() > maxX {
			maxX = point.X()
		}

		if point.Y() < minY {
			minY = point.Y()
		}
		if point.Y() > maxY {
			maxY = point.Y()
		}
	}

	if !any {
		return NewMicroPoint(0, 0), NewMicroPoint(0, 0)
	}

	return NewMicroPoint(minX, minY), NewMicroPoint(maxX, maxY)
}

// Rotate rotates all points around (0|0) by the given degree.
func (p Path) Rotate(degree float64) {
	for i, point := range p {
		p[i] = point.Rotate(degree)
	}
}

// Paths represents a group of Paths.
type Paths []Path

// Bounds calculates the bounding box of all Paths
// The returned points are the min-X-Y-Point and the max-X-Y-Point.
func (p Paths) Bounds() (MicroPoint, MicroPoint) {
	if len(p) == 0 {
		return NewMicroPoint(0, 0), NewMicroPoint(0, 0)
	}

	minX := MaxMicrometer
	minY := MaxMicrometer

	maxX := MinMicrometer
	maxY := MinMicrometer

	// return 0, 0, 0, 0 if everything is empty
	any := false
	for _, path := range p {
		for _, point := range path {
			any = true
			if point.X() < minX {
				minX = point.X()
			}
			if point.X() > maxX {
				maxX = point.X()
			}

			if point.Y() < minY {
				minY = point.Y()
			}
			if point.Y() > maxY {
				maxY = point.Y()
			}
		}
	}

	if !any {
		return NewMicroPoint(0, 0), NewMicroPoint(0, 0)
	}

	return NewMicroPoint(minX, minY), NewMicroPoint(maxX, maxY)
}

// Rotate rotates all points around (0|0) by the given degree.
func (p Paths) Rotate(degree float64) {
	for _, path := range p {
		path.Rotate(degree)
	}
}

// LayerPart represents one part of a layer.
// It consists of an outline and may have several holes
// Some implementations may also provide Attributes for it.
type LayerPart interface {
	Outline() Path
	Holes() Paths

	// Attributes can be any additional data, referenced by a key.
	// Note that you have to know what type the attribute has to
	// use proper type assertion.
	//
	// If the implementation does not support attributes, it should return nil.
	// If the implementation supports attributes but doesn't have any, it should return an empty map.
	Attributes() map[string]interface{}
}

// Layer represents one layer which can consist of several polygons.
// These polygons can consist of several paths, with some of them just representing holes.
// Holes have to be clockwise and outlines counter clockwise.
type Layer interface {
	Polygons() Paths
}

// PartitionedLayer represents one layer with separated layer parts.
// In contrast to the interface Layer this one contains already processed
// polygons in the form of LayerParts.
type PartitionedLayer interface {
	LayerParts() []LayerPart

	// Attributes can be any additional data, referenced by a key.
	// Note that you have to know what type the attribute has to
	// use proper type assertion.
	//
	// If the implementation does not support attributes, it should return nil.
	// If the implementation supports attributes but doesn't have ane, it should return an empty map.
	Attributes() map[string]interface{}

	// Bounds returns the min and max Points which specify the bounding box.
	Bounds() (MicroPoint, MicroPoint)
}

// basicLayerPart is the simplest implementation of LayerPart.
// It holds one outline and several hole-paths.
// You can assume that all paths are closed polygons.
// (If the instance is created by GoSlice...)
type basicLayerPart struct {
	outline Path
	holes   Paths
}

// NewBasicLayerPart returns a new, simple LayerPart.
func NewBasicLayerPart(outline Path, holes Paths) LayerPart {
	return basicLayerPart{
		outline: outline,
		holes:   holes,
	}
}

func (l basicLayerPart) Outline() Path {
	return l.outline
}

func (l basicLayerPart) Holes() Paths {
	return l.holes
}

func (l basicLayerPart) Attributes() map[string]interface{} {
	return nil
}

type partitionedLayer struct {
	parts []LayerPart
}

// NewPartitionedLayer returns a new simple PartitionedLayer which just contains several LayerParts.
func NewPartitionedLayer(parts []LayerPart) PartitionedLayer {
	return partitionedLayer{
		parts: parts,
	}
}

func (p partitionedLayer) LayerParts() []LayerPart {
	return p.parts
}

func (p partitionedLayer) Attributes() map[string]interface{} {
	return nil
}

func (p partitionedLayer) Bounds() (MicroPoint, MicroPoint) {
	paths := Paths{}
	for _, p := range p.LayerParts() {
		paths = append(paths, p.Outline())
	}

	return paths.Bounds()
}
