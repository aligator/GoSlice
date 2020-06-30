package data_test

import (
	"GoSlice/data"
	"GoSlice/util/test"
	"github.com/google/go-cmp/cmp"
	"strings"
	"testing"
)

// fanSpeedOptionsComparer returns a cmp.Comparer which can handle data.FanSpeedOptions.
func fanSpeedOptionsComparer() cmp.Option {
	return cmp.Comparer(func(p1, p2 data.FanSpeedOptions) bool {
		for layer, speed := range p1.LayerToSpeedLUT {
			if p2.LayerToSpeedLUT[layer] != speed {
				return false
			}
		}
		return true
	})
}

func TestSetFanSpeed(t *testing.T) {
	var testCases = map[string]struct {
		optionString  string
		expectedError string
		expected      *data.FanSpeedOptions
	}{
		"FanSpeedNegative": {
			optionString:  "1=-5",
			expectedError: "fan control needs to be in format",
			expected:      nil,
		},
		"FanSpeedOneSuccessful": {
			optionString:  "1=10",
			expectedError: "",
			expected: &data.FanSpeedOptions{LayerToSpeedLUT: map[int]int{
				1: 10,
			}},
		},
		"TestFanSpeedLayerNegative": {
			optionString:  "-1=5",
			expectedError: "fan control needs to be in format",
			expected:      nil,
		},
		"TestFanSpeedGreaterThan255": {
			optionString:  "1=256",
			expectedError: "fan control needs to be in format",
			expected:      nil,
		},
		"TestFanSpeedSingleGood": {
			optionString:  "1=20",
			expectedError: "",
			expected: &data.FanSpeedOptions{LayerToSpeedLUT: map[int]int{
				1: 20,
			}},
		},
		"TestFanSpeedMultipleGood": {
			optionString:  "1=20,5=100",
			expectedError: "",
			expected: &data.FanSpeedOptions{LayerToSpeedLUT: map[int]int{
				1: 20,
				5: 100,
			}},
		},
		"TestFanSpeedMultipleOneBadOneGood": {
			optionString:  "1=-20,5=100",
			expectedError: "fan control needs to be in format",
			expected:      nil,
		},
	}

	for testName, testCase := range testCases {
		t.Log("testCase:", testName)
		actual := data.FanSpeedOptions{}
		err := actual.Set(testCase.optionString)

		// negative cases
		if testCase.expectedError != "" {
			test.Assert(t, strings.Contains(err.Error(), testCase.expectedError), "error expected")
		} else {
			// positive cases
			test.Equals(t, testCase.expected, &actual, fanSpeedOptionsComparer())
		}
	}
}

// fanSpeedOptionsStringComparer returns a cmp.Comparer which can handle generated strings.
func fanSpeedOptionsStringComparer() cmp.Option {
	return cmp.Comparer(func(p1, p2 string) bool {
		return p1 == p2
	})
}

func TestSetFanSpeedString(t *testing.T) {
	var testCases = map[string]struct {
		optionString  string
		expectedError string
		expected      string
	}{
		"TestFanSpeedStringMultipleGood": {
			optionString: "1=20,5=100",
			expected:     "1=20,5=100",
		},
		"TestFanSpeedStringSingleGood": {
			optionString: "1=20",
			expected:     "1=20",
		},
		"TestFanSpeedStringSingleBad": {
			optionString: "1=-20",
			expected:     "", // expect default value for string to be returned in bad case.
		},
	}

	for testName, testCase := range testCases {
		t.Log("testCase:", testName)
		actual := data.FanSpeedOptions{}
		actual.Set(testCase.optionString)
		test.Equals(t, actual.String(), testCase.expected, fanSpeedOptionsStringComparer())
	}
}
