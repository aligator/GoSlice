package reader

import (
	"errors"
	"os"

	"github.com/aligator/goslice/data"
	"github.com/aligator/goslice/handler"
	"github.com/hschendel/stl"
)

// face is a 3d triangle face defined by three 3d vectors.
type face struct {
	vectors [3]data.MicroVec3
}

func (f face) Points() [3]data.MicroVec3 {
	return f.vectors
}

type Model struct {
	faces []data.Face
}

func (m *Model) SetName(name string) {
	// not used yet
	return
}

func (m *Model) SetBinaryHeader(header []byte) {
	// not used yet
	return
}

func (m *Model) SetASCII(isASCII bool) {
	// not used yet
	return
}

func (m *Model) SetTriangleCount(n uint32) {
	// not used yet
	return
}

func (m *Model) AppendTriangle(t stl.Triangle) {
	m.faces = append(m.faces, stlTriangleToFace(t))
}

func newModel(faces []data.Face) data.Model {
	return &Model{
		faces: faces,
	}
}

func (m Model) FaceCount() int {
	return len(m.faces)
}

func (m Model) Face(index int) data.Face {
	return m.faces[index]
}

func (m Model) Min() data.MicroVec3 {
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

func (m Model) Max() data.MicroVec3 {
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

// Reader returns a stl model reader.
func Reader(options *data.Options) handler.ModelReader {
	return &reader{}
}

func (r reader) Read(filename string) (data.Model, error) {
	model := &Model{}
	if _, err := os.Stat(filename); errors.Is(err, os.ErrNotExist) {
		return model, os.ErrNotExist
	}
	stl.CopyFile(filename, model)
	return model, nil
}

// stlTriangleToFace converts a triangle from the stl package
// into a face.
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
