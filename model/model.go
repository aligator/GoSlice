package model

import (
	"GoSlicer/util"
	"errors"
	"github.com/hschendel/stl"
	"strings"
)

type Face interface {
	Vectors() [3]util.MicroVec3
}

type face struct {
	vectors [3]util.MicroVec3
}

func (f *face) Vectors() [3]util.MicroVec3 {
	return f.vectors
}

func NewFace(vec0, vec1, vec2 util.MicroVec3) Face {
	return &face{
		vectors: [3]util.MicroVec3{
			vec0, vec1, vec2,
		}}
}

type Model interface {
	Min() util.MicroVec3
	Max() util.MicroVec3
	AddFace(face Face)
	Faces() []Face
}

type model struct {
	faces []Face
}

func stlTriangleToFace(t stl.Triangle) Face {
	return &face{vectors: [3]util.MicroVec3{
		util.NewMicroVec3(
			util.Millimeter(t.Vertices[0][0]).ToMicrometer(),
			util.Millimeter(t.Vertices[0][1]).ToMicrometer(),
			util.Millimeter(t.Vertices[0][2]).ToMicrometer()),
		util.NewMicroVec3(
			util.Millimeter(t.Vertices[1][0]).ToMicrometer(),
			util.Millimeter(t.Vertices[1][1]).ToMicrometer(),
			util.Millimeter(t.Vertices[1][2]).ToMicrometer()),
		util.NewMicroVec3(
			util.Millimeter(t.Vertices[2][0]).ToMicrometer(),
			util.Millimeter(t.Vertices[2][1]).ToMicrometer(),
			util.Millimeter(t.Vertices[2][2]).ToMicrometer()),
	}}
}

// LoadSTL loads a model from a stl file and converts it to a Model
func LoadSTL(filename string) (Model, error) {
	m := model{}

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
		m.AddFace(stlTriangleToFace(triangle))
	}

	return &m, nil
}

func (m *model) AddFace(face Face) {
	m.faces = append(m.faces, face)
}

func (m *model) Min() util.MicroVec3 {
	ret := m.faces[0].Vectors()[0].Copy()

	for _, face := range m.faces {
		for _, vertice := range face.Vectors() {
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

func (m *model) Max() util.MicroVec3 {
	ret := m.faces[0].Vectors()[0].Copy()

	for _, face := range m.faces {
		for _, vertice := range face.Vectors() {
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

func (m *model) Faces() []Face {
	return m.faces
}
