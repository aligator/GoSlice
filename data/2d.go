// This file provides some basic helper functions for 2d vector calculations.

package data

// DotProduct calculates the dot product of two points
func DotProduct(a, b MicroPoint) Micrometer {
	return a.X()*b.X() + a.Y()*b.Y()
}

// PerpendicularDistance2 calculates the (perpendicular Distance)^2 of a point to a line
func PerpendicularDistance2(a, b, point MicroPoint) Micrometer {
	//  x.......a------------b
	//  :
	//  :
	//  p
	// return px_size
	vecAB := b.Sub(a)
	vecAP := point.Sub(a)

	if vecAB.Size2() == 0 {
		return vecAP.Size2() // assume a perpendicular line to p
	}

	dotProduct := DotProduct(vecAB, vecAP)
	axSize2 := dotProduct * dotProduct / vecAB.Size2()
	return Max(0, vecAP.Size2()-axSize2)
}

// douglasPeucker accepts a list of points and epsilon as threshold, simplifies a path by dropping
// points that do not pass threshold values.
func douglasPeucker(points Path, ep Micrometer) Path {
	if len(points) <= 2 {
		return points
	}

	idx, maxDist := seekMostDistantPoint(points[0], points[len(points)-1], points)
	if maxDist >= ep {
		// TODO: check if implementation without recursion would be possible and if it is more performant
		left := douglasPeucker(points[:idx+1], ep)
		right := douglasPeucker(points[idx:], ep)
		return append(left[:len(left)-1], right...)
	}

	// If the most distant point fails to pass the threshold test, then just return the two points
	return Path{points[0], points[len(points)-1]}
}

func seekMostDistantPoint(p1 MicroPoint, p2 MicroPoint, points Path) (idx int, maxDist Micrometer) {
	for i := 0; i < len(points); i++ {

		// TODO: check usage of 'Shortest Distance' from a point to a line segment
		//       suggested here https://karthaus.nl/rdp/ I think slic3r uses that
		d := PerpendicularDistance2(p1, p2, points[i])
		if d > maxDist*maxDist {
			maxDist = d
			idx = i
		}
	}

	return idx, maxDist
}

// DouglasPeucker is an algorithm for simplifying / smoothing polygons by removing some points.
// see https://en.wikipedia.org/wiki/Ramer%E2%80%93Douglas%E2%80%93Peucker_algorithm
func DouglasPeucker(points Path, epsilon Micrometer) Path {
	if epsilon == -1 {
		// TODO: This value may need optimization
		epsilon = 70
	}
	return douglasPeucker(points, epsilon)
}
