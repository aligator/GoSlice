package data

import "GoSlicer/util"

type Face interface {
	Points() [3]util.MicroVec3
}

type Model interface {
	FaceCount() int
	Face(index int) Face
	Min() util.MicroVec3
	Max() util.MicroVec3
}

type OptimizedFace interface {
	Face
	TouchingFaceIndices() [3]int
	MinZ() util.Micrometer
	MaxZ() util.Micrometer
}

type OptimizedModel interface {
	Model
	Size() util.MicroVec3
	OptimizedFace(index int) OptimizedFace
	SaveDebugSTL(filename string) error
}
