package slicer

import (
	"GoSlicer/util"
	"time"
)

type Slicer struct {
	Path string
}

type config struct {
	layerThickness        int
	initialLayerThickness int
	filamentDiameter      int
	extrusionWidth        int
	insetCount            int
}

func (s *Slicer) Process() {
	c := config{
		layerThickness:        100,
		initialLayerThickness: 200,
		filamentDiameter:      1500,
		extrusionWidth:        400,
		insetCount:            20,
	}

	matrix := util.NewFMatrix3x3()

	t := time.Now()
	model := loadModel()

}
