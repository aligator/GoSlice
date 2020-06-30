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

func TestFanSpeedStringMultipleGood(t *testing.T) {
	fanSpeedText := "1=20,5=100"
	fanSpeedOptions := data.FanSpeedOptions{}
	fanSpeedOptions.Set(fanSpeedText)
	test.Equals(t, fanSpeedOptions.String(), fanSpeedText)
}

func TestFanSpeedStringSingleGood(t *testing.T) {
	fanSpeedText := "1=20"
	fanSpeedOptions := data.FanSpeedOptions{}
	fanSpeedOptions.Set(fanSpeedText)
	test.Equals(t, fanSpeedOptions.String(), fanSpeedText)
}
