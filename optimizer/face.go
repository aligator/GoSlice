package optimizer

import (
	"github.com/aligator/goslice/data"
)

type optimizedFace struct {
	model    *optimizedModel
	indices  [3]int
	touching [3]int
	index    int
}

func (o optimizedFace) Points() [3]data.MicroVec3 {
	return [3]data.MicroVec3{
		o.model.points[o.indices[0]].pos,
		o.model.points[o.indices[1]].pos,
		o.model.points[o.indices[2]].pos,
	}
}

func (o optimizedFace) TouchingFaceIndices() [3]int {
	return o.touching
}

func (o optimizedFace) MinZ() data.Micrometer {
	points := o.Points()
	minZ := points[0].Z()

	if points[1].Z() < minZ {
		minZ = points[1].Z()
	}
	if points[2].Z() < minZ {
		minZ = points[2].Z()
	}

	return minZ
}

func (o optimizedFace) MaxZ() data.Micrometer {
	points := o.Points()
	maxZ := points[0].Z()

	if points[1].Z() > maxZ {
		maxZ = points[1].Z()
	}
	if points[2].Z() > maxZ {
		maxZ = points[2].Z()
	}

	return maxZ
}
