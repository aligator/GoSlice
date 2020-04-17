package data

import (
	"GoSlicer/util"
)

type Path []util.MicroPoint

func (p Path) IsAlmostFinished(distance util.Micrometer) bool {
	return p[0].Sub(p[len(p)-1]).ShorterThan(distance)
}

type Paths []Path

type ComplexPolygon interface {
	Paths() Paths
}

type LayerPart interface {
	Polygons() Paths
}

type Layer interface {
	Polygons() Paths
}

type PartitionedLayer interface {
	LayerParts() []LayerPart
}
