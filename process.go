package main

import (
	"GoSlicer/model"
	"GoSlicer/slicer"
	"GoSlicer/util"
	"fmt"
	"time"
)

type Process struct {
	Path string
}

func (s *Process) Process() error {
	c := util.Config{
		LayerThickness:        200,
		InitialLayerThickness: 200,
		FilamentDiameter:      1500,
		ExtrusionWidth:        400,
		InsetCount:            2,
	}

	t := time.Now()
	m, err := model.LoadSTL(s.Path)

	if err != nil {
		return err
	}

	fmt.Println("load from disk time:", time.Now().Sub(t))
	t = time.Now()

	om := model.OptimizeModel(m, 30, util.NewMicroVec3(102500, 102500, 0))
	fmt.Println("#Face count: Model:", len(m.Faces()), "optimized:", len(om.Faces()), "->", float32(len(om.Faces()))/float32(len(m.Faces()))*100, "%")
	fmt.Println("#Vertex count: Model:", len(m.Faces())*3, "optimized:", len(om.Points()), "->", float32(len(om.Points()))/float32(len(m.Faces())*3)*100, "%")
	fmt.Println("optimization time:", time.Now().Sub(t))
	t = time.Now()

	om.SaveDebugSTL("debug.stl")

	fmt.Println("save stl time:", time.Now().Sub(t))
	t = time.Now()

	fmt.Println("Slicing model")
	// Why is initialLayerThickness / 2 ??
	slicer := slicer.NewSlicer(om, c.InitialLayerThickness/2, c.LayerThickness)
	slicer.DumpSegments("output.html")
	fmt.Println("slicing time:", time.Now().Sub(t))
	t = time.Now()
	fmt.Println("Generating layer parts")
	slicer.GenerateLayerParts()
	fmt.Println("layerParts generation time:", time.Now().Sub(t))
	t = time.Now()

	//slicer.DumpLayerParts("layerParts.html")
	t = time.Now()
	slicer.GenerateGCode("output.gcode", c)
	fmt.Println("gcode generation time:", time.Now().Sub(t))

	return nil
}
