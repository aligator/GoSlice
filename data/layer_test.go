package data_test

import (
	"GoSlice/data"
	"GoSlice/util/test"
	"testing"
)

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
