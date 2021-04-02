package data_test

import (
	"github.com/aligator/goslice/data"
	"github.com/aligator/goslice/util/test"
	"testing"
)

func TestMillimeter(t *testing.T) {
	var tests = []struct {
		val      data.Millimeter
		expected data.Micrometer
	}{
		{val: 10, expected: 10000},
		{val: 1.123, expected: 1123},
		{val: 1.1235, expected: 1124},
		{val: 1.1234, expected: 1123},
		{val: 0, expected: 0},
		{val: -1.1234, expected: -1123},
	}

	for _, testCase := range tests {
		test.Equals(t, testCase.expected, testCase.val.ToMicrometer())
	}
}

func TestMicrometer(t *testing.T) {
	var tests = []struct {
		val      data.Micrometer
		expected data.Millimeter
	}{
		{val: 10, expected: 0.01},
		{val: 1123, expected: 1.123},
		{val: 1124, expected: 1.124},
		{val: 0, expected: 0},
		{val: -1124, expected: -1.124},
	}

	for _, testCase := range tests {
		test.Equals(t, testCase.expected, testCase.val.ToMillimeter())
	}
}
