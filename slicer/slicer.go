package slicer

import (
	"GoSlicer/slicer/model"
	"GoSlicer/util"
	"fmt"
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

func (s *Slicer) Process() error {
	/*c := config{
		layerThickness:        100,
		initialLayerThickness: 200,
		filamentDiameter:      1500,
		extrusionWidth:        400,
		insetCount:            20,
	}*/

	t := time.Now()
	m, err := model.LoadSTL(s.Path)

	if err != nil {
		return err
	}

	fmt.Println("load from disk time: ", time.Now().Sub(t))

	om := model.OptimizeModel(m, 30, util.NewMicroVec3(102500, 102500, 0))

	om.SaveDebugSTL("debug.stl")

	fmt.Println("time needed: ", time.Now().Sub(t))

	return nil
}
