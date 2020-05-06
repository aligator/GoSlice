package renderer

import (
	"GoSlice/clip"
	"GoSlice/data"
	"GoSlice/gcode/builder"
)

type Infill struct {
	PatternSetup func(min data.MicroPoint, max data.MicroPoint) clip.Pattern
	AttrName     string
	Comments     []string

	pattern clip.Pattern
}

func (i *Infill) Init(model data.OptimizedModel) {
	i.pattern = i.PatternSetup(model.Min().PointXY(), model.Max().PointXY())
}

func (i *Infill) Render(builder builder.Builder, layerNr int, layers []data.PartitionedLayer, z data.Micrometer, options *data.Options) {
	if i.pattern == nil {
		return
	}

	bottom, ok := layers[layerNr].Attributes()[i.AttrName].([]data.LayerPart)
	if !ok {
		return
	}

	for _, part := range bottom {
		for _, c := range i.Comments {
			builder.AddComment(c)
		}

		for _, path := range i.pattern.Fill(layerNr, part) {
			builder.AddPolygon(path, z)
		}
	}
}
