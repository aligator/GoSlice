package model

import (
	"GoSlicer/util"
	"fmt"
	"github.com/hschendel/stl"
)

type OptimizedFace interface {
	Indices() [3]int
	Touching() [3]int
	SetTouching([3]int)
}

type optimizedFace struct {
	indices  [3]int
	touching [3]int
}

func (o *optimizedFace) Indices() [3]int {
	return o.indices
}

func (o *optimizedFace) Touching() [3]int {
	return o.touching
}

func (o *optimizedFace) SetTouching(touching [3]int) {
	o.touching = touching
}

type OptimizedPoint interface {
	Point() util.MicroVec3
	SetPoint(point util.MicroVec3)
	FaceIndices() []int
	AddFaceIndex(index int)
}

type optimizedPoint struct {
	point       util.MicroVec3
	faceIndices []int
}

func newOptimizedPoint(vec util.MicroVec3) OptimizedPoint {
	return &optimizedPoint{
		point: vec,
	}
}

func (p *optimizedPoint) Point() util.MicroVec3 {
	return p.point
}

func (p *optimizedPoint) SetPoint(point util.MicroVec3) {
	p.point = point
}

func (p *optimizedPoint) FaceIndices() []int {
	return p.faceIndices
}

func (p *optimizedPoint) AddFaceIndex(index int) {
	p.faceIndices = append(p.faceIndices, index)
}

type OptimizedModel interface {
	Points() []OptimizedPoint
	Faces() []OptimizedFace
	Size() util.MicroVec3
	FacePoints(i int) [3]util.MicroVec3

	SaveDebugSTL(filename string) error
}

type optimizedModel struct {
	points    []OptimizedPoint
	faces     []OptimizedFace
	modelSize util.MicroVec3
}

type pointHash uint

func OptimizeModel(m Model, meldDistance util.Micrometer, center util.MicroVec3) OptimizedModel {
	om := &optimizedModel{}

	minVector := m.Min()
	maxVector := m.Max()

	// map of same faces grouped by their calculated hash
	indices := make(map[pointHash][]int, 0)

FacesLoop:
	for _, face := range m.Faces() {
		optimizedFace := &optimizedFace{
			indices:  [3]int{},
			touching: [3]int{},
		}
		for j := 0; j < 3; j++ {
			currentPoint := face.Vectors()[j]
			// create hash for the point
			// points which are within the meldDistance fall into the same category of the indices map
			meldDistanceHash := pointHash(meldDistance)
			hash := ((pointHash(currentPoint.X()) + meldDistanceHash/2) / meldDistanceHash) ^
				(((pointHash(currentPoint.Y()) + meldDistanceHash/2) / meldDistanceHash) << 10) ^
				(((pointHash(currentPoint.Z()) + meldDistanceHash/2) / meldDistanceHash) << 20)
			var idx int
			add := true

			// for each point-indices with this hash
			// check if the difference between it and the currentPoint
			// is smaller (or same) than the currently tested point
			for _, index := range indices[hash] {
				differenceVec := om.Points()[index].Point().Sub(currentPoint)
				if differenceVec.TestLength(meldDistance) {
					// if true for any of the points with the same hash,
					// do not add the current point to the indices map
					// but save the indices of the already existing duplicate
					idx = index
					add = false
					break
				}
			}
			if add {
				// add the new point-indices to the indices
				indices[hash] = append(indices[hash], len(om.points))
				idx = len(om.points)
				om.points = append(om.points, newOptimizedPoint(currentPoint))
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

		// check if there is a face with the same points
		for _, faceIndex0 := range om.points[optimizedFace.indices[0]].FaceIndices() {
			for _, faceIndex1 := range om.points[optimizedFace.indices[1]].FaceIndices() {
				for _, faceIndex2 := range om.points[optimizedFace.indices[2]].FaceIndices() {
					if faceIndex0 == faceIndex1 &&
						faceIndex0 == faceIndex2 {
						// no need to go further
						continue FacesLoop
					}
				}
			}
		}

		// if it comes here, no duplicate was detected
		om.points[optimizedFace.indices[0]].AddFaceIndex(len(om.faces))
		om.points[optimizedFace.indices[1]].AddFaceIndex(len(om.faces))
		om.points[optimizedFace.indices[2]].AddFaceIndex(len(om.faces))
		om.faces = append(om.faces, optimizedFace)
	}

	// count open faces
	openFaces := 0
	for i, face := range om.faces {
		touching := [3]int{
			om.getFaceIdxWithPoints(face.Indices()[0], face.Indices()[1], i),
			om.getFaceIdxWithPoints(face.Indices()[1], face.Indices()[2], i),
			om.getFaceIdxWithPoints(face.Indices()[2], face.Indices()[0], i),
		}
		face.SetTouching(touching)

		if face.Touching()[0] == -1 {
			openFaces++
		}

		if face.Touching()[1] == -1 {
			openFaces++
		}

		if face.Touching()[2] == -1 {
			openFaces++
		}
	}

	fmt.Printf("Number of open faces: %v\n", openFaces)

	// move points according to the center value
	vectorOffset := util.NewMicroVec3((minVector.X()+maxVector.X())/2, (minVector.Y()+maxVector.Y())/2, minVector.Z())
	vectorOffset = vectorOffset.Sub(center)
	for _, point := range om.points {
		point.SetPoint(point.Point().Sub(vectorOffset))
	}

	om.modelSize = maxVector.Sub(minVector)

	return om
}

func (m *optimizedModel) Points() []OptimizedPoint {
	return m.points
}

func (m *optimizedModel) Faces() []OptimizedFace {
	return m.faces
}

func (m *optimizedModel) Size() util.MicroVec3 {
	return m.modelSize
}

func (m *optimizedModel) getFaceIdxWithPoints(idx0, idx1, notFaceIdx int) int {
	for _, faceIndex0 := range m.points[idx0].FaceIndices() {
		if faceIndex0 == notFaceIdx {
			continue
		}
		for _, faceIndex1 := range m.points[idx1].FaceIndices() {
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

func (m *optimizedModel) FacePoints(i int) [3]util.MicroVec3 {
	return [3]util.MicroVec3{
		m.points[m.faces[i].Indices()[0]].Point(),
		m.points[m.faces[i].Indices()[1]].Point(),
		m.points[m.faces[i].Indices()[2]].Point(),
	}
}

func (m *optimizedModel) SaveDebugSTL(filename string) error {
	triangles := make([]stl.Triangle, 0)

	for _, face := range m.Faces() {
		triangles = append(triangles, stl.Triangle{
			Normal: [3]float32{
				0, 0, 0,
			},
			Vertices: [3]stl.Vec3{
				[3]float32{
					float32(m.points[face.Indices()[0]].Point().X().ToMillimeter()),
					float32(m.points[face.Indices()[0]].Point().Y().ToMillimeter()),
					float32(m.points[face.Indices()[0]].Point().Z().ToMillimeter()),
				},
				[3]float32{
					float32(m.points[face.Indices()[1]].Point().X().ToMillimeter()),
					float32(m.points[face.Indices()[1]].Point().Y().ToMillimeter()),
					float32(m.points[face.Indices()[1]].Point().Z().ToMillimeter()),
				},
				[3]float32{
					float32(m.points[face.Indices()[2]].Point().X().ToMillimeter()),
					float32(m.points[face.Indices()[2]].Point().Y().ToMillimeter()),
					float32(m.points[face.Indices()[2]].Point().Z().ToMillimeter()),
				},
			},
			Attributes: 0,
		})
	}

	solid := stl.Solid{
		BinaryHeader: nil,
		Name:         "GoSlice_STL_export",
		Triangles:    triangles,
		IsAscii:      false,
	}

	err := solid.WriteFile(filename)
	if err != nil {
		return err
	}

	return nil
}
