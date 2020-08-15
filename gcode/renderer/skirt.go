// This file provides renderers for gcode injected at specific layers.

package renderer

import (
	"GoSlice/clip"
	"GoSlice/data"
	"GoSlice/gcode"
	"GoSlice/modifier"
)

// Skirt generates the skirt lines.
type Skirt struct{}

func (Skirt) Init(model data.OptimizedModel) {}

func (Skirt) Render(b *gcode.Builder, layerNr int, layers []data.PartitionedLayer, z data.Micrometer, options *data.Options) error {
	if options.Print.BrimSkirt.SkirtCount == 0 {
		return nil
	}

	b.AddComment("LAYER:%v", layerNr)
	if layerNr == 0 {
		perimeters, err := modifier.Perimeters(layers[layerNr])
		if err != nil {
			return err
		}
		if perimeters == nil {
			return nil
		}

		// skirt distance + 1/2 extrusion with of the model side + 1/2 extrusion width of the most inner brim line is the distance between the perimeter (or brim) and skirt.
		distance := options.Print.BrimSkirt.SkirtDistance.ToMicrometer() + (options.Printer.ExtrusionWidth)

		// draw skirt
		c := clip.NewClipper()
		skirtMostInner := c.Inset(data.NewBasicLayerPart(c.Hull(perimeters.ToOneDimension()), nil), -distance, 1)

		if len(skirtMostInner) == 0 || len(skirtMostInner[0]) == 0 {
			return nil
		}

		// there should only be one line, so just use it
		if options.Print.BrimSkirt.SkirtCount <= 1 {
			return nil
		}
		// and generate the other skirt lines from it
		skirt := c.Inset(skirtMostInner[0][0], -options.Printer.ExtrusionWidth, options.Print.BrimSkirt.SkirtCount)

		// as we have only one hull there should be always only one "wall"...
		for _, wall := range skirt {
			for _, loopPart := range wall {
				// as we use the hull around the whole object there shouldn't be any collision with the model -> currentLayer is nil
				b.AddPolygon(nil, loopPart.Outline(), z, false)
			}
		}
	}

	return nil
}
