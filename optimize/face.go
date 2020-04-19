package optimize

import (
	"GoSlice/util"
)

type optimizedFace struct {
	model    *optimizedModel
	indices  [3]int
	touching [3]int
	index    int
}

func (o optimizedFace) Points() [3]util.MicroVec3 {
	return [3]util.MicroVec3{
		o.model.points[o.indices[0]].pos,
		o.model.points[o.indices[1]].pos,
		o.model.points[o.indices[2]].pos,
	}
}

func (o optimizedFace) TouchingFaceIndices() [3]int {
	return o.touching
}

func (o optimizedFace) MinZ() util.Micrometer {
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

func (o optimizedFace) MaxZ() util.Micrometer {
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
