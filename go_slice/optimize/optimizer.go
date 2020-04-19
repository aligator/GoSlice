package optimize

import (
	"GoSlice/go_slice/data"
	"GoSlice/go_slice/handle"
	"GoSlice/util"
	"fmt"
)

type optimizer struct {
	options *data.Options
}

func NewOptimizer(options *data.Options) handle.ModelOptimizer {
	return &optimizer{
		options: options,
	}
}

type pointHash uint

func (o optimizer) Optimize(m data.Model) (data.OptimizedModel, error) {
	om := &optimizedModel{}

	// map of same faces grouped by their calculated hash
	indices := make(map[pointHash][]int, 0)

FacesLoop:
	for i := 0; i < m.FaceCount(); i++ {
		face := m.Face(i)

		optimizedFace := optimizedFace{
			indices:  [3]int{},
			touching: [3]int{},
			model:    om,
		}
		for j := 0; j < 3; j++ {
			currentPoint := face.Points()[j]
			// create hash for the pos
			// points which are within the meldDistance fall into the same category of the indices map
			meldDistanceHash := pointHash(o.options.MeldDistance)
			hash := ((pointHash(currentPoint.X()) + meldDistanceHash/2) / meldDistanceHash) ^
				(((pointHash(currentPoint.Y()) + meldDistanceHash/2) / meldDistanceHash) << 10) ^
				(((pointHash(currentPoint.Z()) + meldDistanceHash/2) / meldDistanceHash) << 20)
			var idx int
			add := true

			// for each pos-indices with this hash
			// check if the difference between it and the currentPoint
			// is smaller (or same) than the currently tested pos
			for _, index := range indices[hash] {
				differenceVec := om.points[index].pos.Sub(currentPoint)
				if differenceVec.TestLength(o.options.MeldDistance) {
					// if true for any of the points with the same hash,
					// do not add the current pos to the indices map
					// but save the indices of the already existing duplicate
					idx = index
					add = false
					break
				}
			}
			if add {
				// add the new pos-indices to the indices
				indices[hash] = append(indices[hash], len(om.points))
				idx = len(om.points)
				om.points = append(om.points, point{
					pos: currentPoint,
				})
			}

			optimizedFace.indices[j] = idx
		}

		// ignore duplicate search for faces
		// which have two vertices with the same location
		if optimizedFace.indices[0] == optimizedFace.indices[1] ||
			optimizedFace.indices[0] == optimizedFace.indices[2] ||
			optimizedFace.indices[1] == optimizedFace.indices[2] {
			continue
		}

		// check if there is a face with the exact same points
		for _, faceIndex0 := range om.points[optimizedFace.indices[0]].faceIndices {
			for _, faceIndex1 := range om.points[optimizedFace.indices[1]].faceIndices {
				for _, faceIndex2 := range om.points[optimizedFace.indices[2]].faceIndices {
					if faceIndex0 == faceIndex1 &&
						faceIndex0 == faceIndex2 {
						// no need to go further
						continue FacesLoop
					}
				}
			}
		}

		// if it reaches this code, no duplicate was detected
		om.points[optimizedFace.indices[0]].faceIndices = append(om.points[optimizedFace.indices[0]].faceIndices, len(om.faces))
		om.points[optimizedFace.indices[1]].faceIndices = append(om.points[optimizedFace.indices[1]].faceIndices, len(om.faces))
		om.points[optimizedFace.indices[2]].faceIndices = append(om.points[optimizedFace.indices[2]].faceIndices, len(om.faces))

		optimizedFace.index = len(om.faces)
		om.faces = append(om.faces, optimizedFace)
	}

	// count open faces
	openFaces := 0
	for i, face := range om.faces {
		face.touching = [3]int{
			om.getFaceIdxWithPoints(face.indices[0], face.indices[1], i),
			om.getFaceIdxWithPoints(face.indices[1], face.indices[2], i),
			om.getFaceIdxWithPoints(face.indices[2], face.indices[0], i),
		}

		if face.touching[0] == -1 {
			openFaces++
		}

		if face.touching[1] == -1 {
			openFaces++
		}

		if face.touching[2] == -1 {
			openFaces++
		}

		om.faces[i] = face
	}

	fmt.Printf("Number of open faces: %v\n", openFaces)

	min := m.Min()
	max := m.Max()
	// move points according to the center value
	vectorOffset := util.NewMicroVec3((min.X()+max.X())/2, (min.Y()+max.Y())/2, min.Z())
	vectorOffset = vectorOffset.Sub(o.options.Center)
	for i, point := range om.points {
		om.points[i].pos = point.pos.Sub(vectorOffset)
	}

	om.modelSize = max.Sub(min)

	return om, nil
}
