// This file provides a renderer for filling parts.

package renderer

import (
	"GoSlice/clip"
	"GoSlice/data"
	"GoSlice/gcode/builder"
)

// Infill is a renderer which can fill parts which are defined by a layer part attribute of a specific name.
// The attribute has to be of type []data.LayerPart.
type Infill struct {
	// PatternSetup is called once on init and sets a specific pattern this infill renderer should use.
	PatternSetup func(min data.MicroPoint, max data.MicroPoint) clip.Pattern

	// AttrName is the name of the attribute containing the []data.LayerPart's to fill.
	AttrName string

	// Comments is a list of comments to be added before each infill.
	Comments []string

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
