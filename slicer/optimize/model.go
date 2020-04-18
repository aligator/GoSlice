package optimize

import (
	"GoSlicer/slicer/data"
	"GoSlicer/util"
	"github.com/hschendel/stl"
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

func (o optimizedModel) SaveDebugSTL(filename string) error {
	triangles := make([]stl.Triangle, 0)

	for _, face := range o.faces {
		triangles = append(triangles, stl.Triangle{
			Normal: [3]float32{
				0, 0, 0,
			},
			Vertices: [3]stl.Vec3{
				[3]float32{
					float32(o.points[face.indices[0]].pos.X().ToMillimeter()),
					float32(o.points[face.indices[0]].pos.Y().ToMillimeter()),
					float32(o.points[face.indices[0]].pos.Z().ToMillimeter()),
				},
				[3]float32{
					float32(o.points[face.indices[1]].pos.X().ToMillimeter()),
					float32(o.points[face.indices[1]].pos.Y().ToMillimeter()),
					float32(o.points[face.indices[1]].pos.Z().ToMillimeter()),
				},
				[3]float32{
					float32(o.points[face.indices[2]].pos.X().ToMillimeter()),
					float32(o.points[face.indices[2]].pos.Y().ToMillimeter()),
					float32(o.points[face.indices[2]].pos.Z().ToMillimeter()),
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
