package stl

import (
	"GoSlice/data"
	"GoSlice/handle"
	"errors"
	"github.com/hschendel/stl"
	"strings"
)

type face struct {
	vectors [3]data.MicroVec3
}

func (f face) Points() [3]data.MicroVec3 {
	return f.vectors
}

type model struct {
	faces []data.Face
}

func newModel(faces []data.Face) data.Model {
	return &model{
		faces: faces,
	}
}

func (m model) FaceCount() int {
	return len(m.faces)
}

func (m model) Face(index int) data.Face {
	return m.faces[index]
}

func (m model) Min() data.MicroVec3 {
	ret := m.faces[0].Points()[0].Copy()

	for _, face := range m.faces {
		for _, vertice := range face.Points() {
			if ret.X() > vertice.X() {
				ret.SetX(vertice.X())
			}

			if ret.Y() > vertice.Y() {
				ret.SetY(vertice.Y())
			}

			if ret.Z() > vertice.Z() {
				ret.SetZ(vertice.Z())
			}
		}
	}

	return ret
}

func (m model) Max() data.MicroVec3 {
	ret := m.faces[0].Points()[0].Copy()

	for _, face := range m.faces {
		for _, vertice := range face.Points() {
			if ret.X() < vertice.X() {
				ret.SetX(vertice.X())
			}

			if ret.Y() < vertice.Y() {
				ret.SetY(vertice.Y())
			}

			if ret.Z() < vertice.Z() {
				ret.SetZ(vertice.Z())
			}
		}
	}

	return ret
}

type reader struct{}

func Reader(options *data.Options) handle.ModelReader {
	return &reader{}
}

func (r reader) Read(filename string) ([]data.Model, error) {
	var faces = []data.Face{}

	splitted := strings.Split(filename, ".")
	if len(splitted) <= 1 {
		return nil, errors.New("the file has no extension")
	}

	extension := splitted[len(splitted)-1]

	if extension != "stl" {
		return nil, errors.New("the file is not a stl file")
	}

	solid, err := stl.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	for _, triangle := range solid.Triangles {
		faces = append(faces, stlTriangleToFace(triangle))
	}

	return []data.Model{newModel(faces)}, nil
}

func stlTriangleToFace(t stl.Triangle) face {
	return face{vectors: [3]data.MicroVec3{
		data.NewMicroVec3(
			data.Millimeter(t.Vertices[0][0]).ToMicrometer(),
			data.Millimeter(t.Vertices[0][1]).ToMicrometer(),
			data.Millimeter(t.Vertices[0][2]).ToMicrometer()),
		data.NewMicroVec3(
			data.Millimeter(t.Vertices[1][0]).ToMicrometer(),
			data.Millimeter(t.Vertices[1][1]).ToMicrometer(),
			data.Millimeter(t.Vertices[1][2]).ToMicrometer()),
		data.NewMicroVec3(
			data.Millimeter(t.Vertices[2][0]).ToMicrometer(),
			data.Millimeter(t.Vertices[2][1]).ToMicrometer(),
			data.Millimeter(t.Vertices[2][2]).ToMicrometer()),
	}}
}
