package model

import "GoSlicer/util"

type OptimizedFace interface {
	Index() [3]int
	touching() [3]int
}

type optimizedFace struct {
	index    [3]int
	touching [3]int
}

type OptimizedVector interface {
	Vector() util.MicroVec3
	FaceIndices() []int
}

type optimizedVector struct {
	vector      util.MicroVec3
	faceIndices []int
}

func NewOptimizedVector(vec3 util.MicroVec3) OptimizedVector {
	return &optimizedVector{
		vector:      vec3,
		faceIndices: make([]int, 0),
	}
}

func (v *optimizedVector) Vector() util.MicroVec3 {
	return v.vector
}

func (v *optimizedVector) FaceIndices() []int {
	return v.faceIndices
}

type OptimizedModel interface {
	Vectors() []OptimizedVector
	Faces() []OptimizedFace
	Size() util.MicroVec3
}

type optimizedModel struct {
	meldDistance util.Micrometer
	vectors      []OptimizedVector
	faces        []OptimizedFace
	modelSize    util.MicroVec3
}

type vectorHash int

func OptimizeModel(m Model, meldDistance util.Micrometer, center util.MicroVec3) OptimizedModel {
	om := &optimizedModel{
		meldDistance: meldDistance,
		vectors:      make([]OptimizedVector, 0),
		faces:        make([]OptimizedFace, 0),
	}

	//minVector := m.Min()
	//maxVector := m.Max()

	// map of same faces grouped by their calculated hash
	indices := make(map[vectorHash][]int, 0)

	// filter nearly double vectors
	for _, face := range m.Faces() {
		optimizedFace := optimizedFace{
			index:    [3]int{},
			touching: [3]int{},
		}
		for j := 0; j < 3; j++ {
			currentVector := face.Vectors()[j]
			// not exactly sure how this calculation works...
			hash := vectorHash(((currentVector.X() + meldDistance/2) / meldDistance) ^ (((currentVector.Y() + meldDistance/2) / meldDistance) << 10) ^ (((currentVector.Z() + meldDistance/2) / meldDistance) << 20))
			var idx int
			add := true

			// for each vector-index with this hash
			// check if the difference between it and the currentVector
			// is smaller (or same) than the currently tested vector
			for _, index := range indices[hash] {
				differenceVec := om.Vectors()[index].Vector().Copy()
				differenceVec.Sub(currentVector)
				if differenceVec.TestLength(meldDistance) {
					// if true for any of the vectors with the same hash,
					// do not add the current vector to the indices map
					// but save the index of the already existing duplicate
					idx = index
					add = false
					break
				}
			}
			if add {
				// add the new vector-index to the indices
				indices[hash] = append(indices[hash], len(om.vectors))
				idx = len(om.vectors)
				om.vectors = append(om.vectors, NewOptimizedVector(currentVector))
			}

			optimizedFace.index[j] = idx
		}

		if optimizedFace.index[0] != optimizedFace.index[1] &&
			optimizedFace.index[0] != optimizedFace.index[2] &&
			optimizedFace.index[1] != optimizedFace.index[2] {

			// Check if there is a face with the same points
			//duplicate := false

		}
	}

	return om
}

func (m *optimizedModel) Vectors() []OptimizedVector {
	return m.vectors
}

func (m *optimizedModel) Faces() []OptimizedFace {
	return m.faces
}

func (m *optimizedModel) Size() util.MicroVec3 {
	return m.modelSize
}
