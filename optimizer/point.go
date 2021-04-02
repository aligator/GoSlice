package optimizer

import (
	"github.com/aligator/goslice/data"
)

// point is a simple point together
// with all indices of the faces where it belongs to
type point struct {
	pos         data.MicroVec3
	faceIndices []int
}
