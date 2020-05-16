package data_test

import (
	"GoSlice/data"
	"GoSlice/util/test"
	"github.com/google/go-cmp/cmp"
	"testing"
)

// pathComparer returns a cmp.Comparer which can handle Path.
func pathComparer() cmp.Option {
	return cmp.Comparer(func(p1, p2 data.Path) bool {
		for i, path := range p1 {
			if !cmp.Equal(path, p2[i], microPointComparer()) {
				return false
			}
		}

		return true
	})
}

// pathsComparer returns a cmp.Comparer which can handle Paths.
func pathsComparer(forceOrder bool) cmp.Option {
	// implementation which forces the order to be the same
	if forceOrder {
		return cmp.Comparer(func(p1, p2 data.Paths) bool {
			if len(p1) != len(p2) {
				return false
			}

			for i, path := range p1 {
				if !cmp.Equal(path, p2[i], pathComparer()) {
					return false
				}
			}

			return true
		})
	}

	// implementation which also accepts a different order
	return cmp.Comparer(func(p1, p2 data.Paths) bool {
		if len(p1) != len(p2) {
			return false
		}

		found := make([]bool, len(p1))

		for i, path := range p1 {
			if !found[i] && cmp.Equal(path, p2[i], pathComparer()) {
				found[i] = true
			}
		}

		// if any of the paths was not found, return false
		for _, f := range found {
			if !f {
				return false
			}
		}

		return true
	})
}

// layerPartComparer returns a cmp.Comparer which can handle LayerPart.
func layerPartComparer(forceOrder bool) cmp.Option {
	return cmp.Comparer(func(p1, p2 data.LayerPart) bool {
		if !cmp.Equal(p1.Outline(), p2.Outline(), pathComparer()) {
			return false
		}
		if !cmp.Equal(p1.Holes(), p2.Holes(), pathsComparer(true)) {
			return false
		}

		return true
	})
}

// Test Path

func TestPathIsAlmostFinished(t *testing.T) {
	var testCases = []struct {
		toTest   data.Path
		distance data.Micrometer
		expected bool
	}{
		{toTest: data.Path{
			data.NewMicroPoint(0, 0),
			data.NewMicroPoint(100, 0),
			data.NewMicroPoint(100, 100),
			data.NewMicroPoint(0, 100),
			data.NewMicroPoint(0, 0),
		}, distance: 30, expected: true},
		{toTest: data.Path{
			data.NewMicroPoint(0, -31),
			data.NewMicroPoint(100, 0),
			data.NewMicroPoint(100, 100),
			data.NewMicroPoint(0, 100),
			data.NewMicroPoint(0, 0),
		}, distance: 30, expected: false},
		{toTest: data.Path{
			data.NewMicroPoint(0, -30),
			data.NewMicroPoint(100, 0),
			data.NewMicroPoint(100, 100),
			data.NewMicroPoint(0, 100),
			data.NewMicroPoint(0, 0),
		}, distance: 30, expected: true},
		{toTest: data.Path{
			data.NewMicroPoint(-30, 0),
			data.NewMicroPoint(100, 0),
			data.NewMicroPoint(100, 100),
			data.NewMicroPoint(0, 100),
			data.NewMicroPoint(0, 0),
		}, distance: 30, expected: true},
		{toTest: data.Path{
			data.NewMicroPoint(-31, 0),
			data.NewMicroPoint(100, 0),
			data.NewMicroPoint(100, 100),
			data.NewMicroPoint(0, 100),
			data.NewMicroPoint(0, 0),
		}, distance: 30, expected: false},
	}

	for i, testCase := range testCases {
		t.Log("testCase", i)
		test.Equals(t, testCase.expected, testCase.toTest.IsAlmostFinished(testCase.distance))
	}
}

func TestPathSimplify(t *testing.T) {
	// TODO
}

func TestPathBounds(t *testing.T) {
	var testCases = []struct {
		toTest      data.Path
		expectedMin data.MicroPoint
		expectedMax data.MicroPoint
	}{
		{toTest: data.Path{
			data.NewMicroPoint(0, 0),
			data.NewMicroPoint(100, 0),
			data.NewMicroPoint(100, 100),
			data.NewMicroPoint(0, 100),
			data.NewMicroPoint(0, 0),
		}, expectedMin: data.NewMicroPoint(0, 0), expectedMax: data.NewMicroPoint(100, 100)},
		{toTest: data.Path{
			data.NewMicroPoint(0, 50),
			data.NewMicroPoint(50, 100),
			data.NewMicroPoint(100, 50),
			data.NewMicroPoint(50, 0),
			data.NewMicroPoint(0, 50),
		}, expectedMin: data.NewMicroPoint(0, 0), expectedMax: data.NewMicroPoint(100, 100)},
	}

	for i, testCase := range testCases {
		t.Log("testCase", i)
		min, max := testCase.toTest.Bounds()
		test.Equals(t, testCase.expectedMin, min, microPointComparer())
		test.Equals(t, testCase.expectedMax, max, microPointComparer())
	}
}

func TestPathsBounds(t *testing.T) {
	var testCases = []struct {
		toTest      data.Paths
		expectedMin data.MicroPoint
		expectedMax data.MicroPoint
	}{
		{toTest: data.Paths{
			data.Path{
				data.NewMicroPoint(0, 0),
				data.NewMicroPoint(100, 0),
				data.NewMicroPoint(100, 100),
				data.NewMicroPoint(0, 100),
				data.NewMicroPoint(0, 0),
			},
			data.Path{
				data.NewMicroPoint(0, 0),
				data.NewMicroPoint(200, 0),
				data.NewMicroPoint(200, 50),
				data.NewMicroPoint(0, 50),
				data.NewMicroPoint(0, 0),
			},
		}, expectedMin: data.NewMicroPoint(0, 0), expectedMax: data.NewMicroPoint(200, 100)},
		{toTest: data.Paths{
			data.Path{
				data.NewMicroPoint(0, 50),
				data.NewMicroPoint(50, 100),
				data.NewMicroPoint(100, 50),
				data.NewMicroPoint(50, 0),
				data.NewMicroPoint(0, 50),
			},
			data.Path{
				data.NewMicroPoint(-50, 0),
				data.NewMicroPoint(0, 50),
				data.NewMicroPoint(50, 0),
				data.NewMicroPoint(0, -50),
				data.NewMicroPoint(-50, 0),
			},
		}, expectedMin: data.NewMicroPoint(-50, -50), expectedMax: data.NewMicroPoint(100, 100)},
	}

	for i, testCase := range testCases {
		t.Log("testCase", i)
		min, max := testCase.toTest.Bounds()
		test.Equals(t, testCase.expectedMin, min, microPointComparer())
		test.Equals(t, testCase.expectedMax, max, microPointComparer())
	}
}

func TestNewBasicLayerPart(t *testing.T) {
	var testCases = []struct {
		outline data.Path
		holes   data.Paths
	}{
		{
			outline: data.Path{
				data.NewMicroPoint(0, 0),
				data.NewMicroPoint(100, 0),
				data.NewMicroPoint(100, 0),
				data.NewMicroPoint(0, 100),
				data.NewMicroPoint(0, 0),
			},
			holes: data.Paths{
				data.Path{
					data.NewMicroPoint(50, 50),
					data.NewMicroPoint(75, 50),
					data.NewMicroPoint(75, 75),
					data.NewMicroPoint(50, 75),
					data.NewMicroPoint(50, 50),
				},
			},
		},
		{
			outline: data.Path{
				data.NewMicroPoint(0, 0),
				data.NewMicroPoint(100, 0),
				data.NewMicroPoint(100, 0),
				data.NewMicroPoint(0, 100),
				data.NewMicroPoint(0, 0),
			},
			holes: data.Paths{
				data.Path{
					data.NewMicroPoint(50, 50),
					data.NewMicroPoint(75, 50),
					data.NewMicroPoint(75, 75),
					data.NewMicroPoint(50, 75),
					data.NewMicroPoint(50, 50),
				},
			},
		},
	}

	for i, testCase := range testCases {
		t.Log("testCase", i)
		part := data.NewBasicLayerPart(testCase.outline, testCase.holes)
		test.Equals(t, testCase.outline, part.Outline(), pathComparer())
		test.Equals(t, testCase.holes, part.Holes(), pathsComparer(true))
		test.Equals(t, map[string]interface{}(nil), part.Attributes())
	}
}
