// Package data holds basic data structures and interfaces used by GoSlice.
package data

// model.go holds types used for 3d model representation.

// Face represents a triangle face which is defined by three vectors.
type Face interface {
	Points() [3]MicroVec3
}

// Model represents a full model.
type Model interface {
	FaceCount() int
	Face(index int) Face
	Min() MicroVec3
	Max() MicroVec3
}

// OptimizedFace represents a full but optimized face.
// It additionally provides indices of touching faces.
// The corresponding other faces can be found in a matching OptimizedModelInstance.
// So you always need an OptimizedModel instance.
type OptimizedFace interface {
	Face
	TouchingFaceIndices() [3]int
	MinZ() Micrometer
	MaxZ() Micrometer
}

// OptimizedModel represents a full but optimized model.
// Each face contains the the indices of touching faces.
type OptimizedModel interface {
	Model

	// Bounds returns the amount of faces.
	Size() MicroVec3

	OptimizedFace(index int) OptimizedFace
	SaveDebugSTL(filename string) error
}
