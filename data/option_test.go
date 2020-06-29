package data_test

import (
	"GoSlice/data"
	"GoSlice/util/test"
	"strings"
	"testing"
)

func TestFanSpeedDefaults(t *testing.T) {
	opts := data.DefaultOptions()
	test.Equals(t, opts.Filament.FanSpeed.LayerToSpeedLUT[2], 255)
}

func TestFanSpeedNegative(t *testing.T) {
	fanSpeedOptions := data.FanSpeedOptions{}
	err := fanSpeedOptions.Set("1=-5")
	test.Assert(t, strings.Contains(err.Error(), "fan control needs to be in format"), "fan speed error expected")
}

func TestFanSpeedLayerNegative(t *testing.T) {
	fanSpeedOptions := data.FanSpeedOptions{}
	err := fanSpeedOptions.Set("-1=5")
	test.Assert(t, strings.Contains(err.Error(), "fan control needs to be in format"), "fan speed error expected")
}

func TestFanSpeedGreaterThan255(t *testing.T) {
	fanSpeedOptions := data.FanSpeedOptions{}
	err := fanSpeedOptions.Set("1=256")
	test.Assert(t, strings.Contains(err.Error(), "fan control needs to be in format"), "fan speed error expected")
}

func TestFanSpeedSingleGood(t *testing.T) {
	fanSpeedOptions := data.FanSpeedOptions{}
	err := fanSpeedOptions.Set("1=20")
	test.Assert(t, err == nil, "unsuccessfully set single fan speed")
}

func TestFanSpeedMultipleGood(t *testing.T) {
	fanSpeedOptions := data.FanSpeedOptions{}
	err := fanSpeedOptions.Set("1=20,5=100")
	test.Assert(t, err == nil, "unsuccessfully set multiple fan speed")
}

func TestFanSpeedMultipleOneBadOneGood(t *testing.T) {
	fanSpeedOptions := data.FanSpeedOptions{}
	err := fanSpeedOptions.Set("1=-20,5=100")
	test.Assert(t, strings.Contains(err.Error(), "fan control needs to be in format"), "fan speed error expected")
}

func TestFanSpeedStringMultipleGood(t *testing.T) {
	fanSpeedText := "1=20,5=100"
	fanSpeedOptions := data.FanSpeedOptions{}
	err := fanSpeedOptions.Set(fanSpeedText)
	test.Assert(t, err == nil, "unsuccessfully set multiple fan speed")
	test.Equals(t, fanSpeedOptions.String(), fanSpeedText)
}

func TestFanSpeedStringSingleGood(t *testing.T) {
	fanSpeedText := "1=20"
	fanSpeedOptions := data.FanSpeedOptions{}
	err := fanSpeedOptions.Set(fanSpeedText)
	test.Assert(t, err == nil, "unsuccessfully set single fan speed")
	test.Equals(t, fanSpeedOptions.String(), fanSpeedText)
}
