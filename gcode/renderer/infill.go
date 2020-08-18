// This file provides a renderer for filling parts.

package renderer

import (
	"GoSlice/clip"
	"GoSlice/data"
	"GoSlice/gcode"
	"GoSlice/modifier"
)

// Infill is a renderer which can fill parts which are defined by a layer part attribute of a specific name.
// The attribute has to be of type []data.LayerPart.
type Infill struct {
	// PatternSetup is called once on init and sets a specific pattern this infill renderer should use.
	// Min and max define the dimension of the model (in X and Y direction)
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

func (i *Infill) Render(b *gcode.Builder, layerNr int, maxLayer int, layer data.PartitionedLayer, z data.Micrometer, options *data.Options) error {
	if i.pattern == nil {
		return nil
	}

	infillParts, err := modifier.PartsAttribute(layer, i.AttrName)
	if err != nil {
		return err
	}
	if infillParts == nil {
		return nil
	}

	for _, part := range infillParts {
		for _, c := range i.Comments {
			b.AddComment(c)
		}

		for _, path := range i.pattern.Fill(layerNr, part) {
			err := b.AddPolygon(layer, path, z, true)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
