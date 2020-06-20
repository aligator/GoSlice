// This file provides a renderer for perimeters.

package renderer

import (
	"GoSlice/data"
	"GoSlice/gcode"
	"GoSlice/modifier"
)

// Perimeter is a renderer which generates the gcode for the attribute "perimeters".
type Perimeter struct{}

func (p Perimeter) Init(model data.OptimizedModel) {}

func (p Perimeter) Render(b *gcode.Builder, layerNr int, layers []data.PartitionedLayer, z data.Micrometer, options *data.Options) {
	perimeters, err := modifier.Perimeters(layers[layerNr])
	if err != nil {
		panic(err)
	}
	if perimeters == nil {
		return
	}

	for _, part := range perimeters {
		for insetNr := range part {
			// print the outer perimeter as last perimeter
			if insetNr >= len(part)-1 {
				insetNr = 0
			} else {
				insetNr++
			}

			for _, insetParts := range part[insetNr] {
				if insetNr == 0 {
					b.AddComment("TYPE:WALL-OUTER")
					b.SetExtrudeSpeed(options.Print.OuterPerimeterSpeed)
				} else {
					b.AddComment("TYPE:WALL-INNER")
					b.SetExtrudeSpeed(options.Print.LayerSpeed)
				}

				for _, hole := range insetParts.Holes() {
					b.AddPolygon(layers[layerNr], hole, z, false)
				}

				b.AddPolygon(layers[layerNr], insetParts.Outline(), z, false)
			}
		}
	}
}
