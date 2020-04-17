package optimize

import (
	"GoSlicer/slicer/data"
	"GoSlicer/util"
)

type optimizedModel struct {
	points    []point
	faces     []optimizedFace
	modelSize util.MicroVec3
}

func (o optimizedModel) FaceCount() int {
	return len(o.faces)
}

func (o optimizedModel) Face(index int) data.Face {
	return o.faces[index]
}

func (o optimizedModel) OptimizedFace(index int) data.OptimizedFace {
	return o.faces[index]
}

func (o optimizedModel) Size() util.MicroVec3 {
	return o.modelSize
}

func (o optimizedModel) Min() util.MicroVec3 {
	panic("implement me")
}

func (o optimizedModel) Max() util.MicroVec3 {
	panic("implement me")
}

func (o optimizedModel) getFaceIdxWithPoints(idx0, idx1, notFaceIdx int) int {
	for _, faceIndex0 := range o.points[idx0].faceIndices {
		if faceIndex0 == notFaceIdx {
			continue
		}
		for _, faceIndex1 := range o.points[idx1].faceIndices {
			if faceIndex1 == notFaceIdx {
				continue
			}
			if faceIndex0 == faceIndex1 {
				return faceIndex0
			}
		}
	}
	return -1
}
