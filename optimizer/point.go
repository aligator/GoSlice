package optimizer

import (
	"GoSlice/data"
)

// point is a simple point together
// with all indices of the faces where it belongs to
type point struct {
	pos         data.MicroVec3
	faceIndices []int
}
