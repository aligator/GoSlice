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
			test.Assert(t, err != nil && strings.Contains(err.Error(), testCase.expectedError), "error expected")
		} else {
			// positive cases
			test.Equals(t, testCase.expected, &actual, fanSpeedOptionsComparer())
		}
	}
}

func TestSetFanSpeedString(t *testing.T) {
	var testCases = map[string]struct {
		optionString string
		expected     []string // use slice as the order is not guaranteed in maps
	}{
		"TestFanSpeedStringMultipleGood": {
			optionString: "1=20,5=100",
			expected: []string{
				"1=20",
				"5=100",
			},
		},
		"TestFanSpeedStringSingleGood": {
			optionString: "1=20",
			expected: []string{
				"1=20",
			},
		},
		"TestFanSpeedStringSingleBad": {
			optionString: "1=-20",
			expected:     []string{}, // expect empty value for string to be returned in bad case.
		},
	}

	for testName, testCase := range testCases {
		t.Log("testCase:", testName)

		option := data.FanSpeedOptions{}
		_ = option.Set(testCase.optionString) // ignore error
		actual := option.String()

		// It is a bit tricky to check the string as the order is not guaranteed in maps.
		// So this checks first if the string split by ',' is the same length as the expected values,
		// or if the expected length is 0 and the actual string is "" (as split of "" results in length 1).
		test.Assert(t, len(strings.Split(actual, ",")) == len(testCase.expected) ||
			len(testCase.expected) == 0 && actual == "", "actual: %v, expected: %d", actual, testCase.expected)

		for _, expected := range testCase.expected {
			test.Assert(t, strings.Contains(actual, expected), "actual '%s' should contain '%s'", actual, expected)
		}
	}
}
