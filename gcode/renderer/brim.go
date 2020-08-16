// This file provides renderers for gcode injected at specific layers.

package renderer

import (
	"GoSlice/clip"
	"GoSlice/data"
	"GoSlice/gcode"
	"GoSlice/modifier"
)

// Brim generates the brim lines.
type Brim struct{}

func (Brim) Init(model data.OptimizedModel) {}

func (Brim) Render(b *gcode.Builder, layerNr int, layers []data.PartitionedLayer, z data.Micrometer, options *data.Options) error {
	if options.Print.BrimSkirt.BrimCount == 0 {
		return nil
	}

	// TODO: add comment used by cura
	//b.AddComment("LAYER:%v", layerNr)
	if layerNr == 0 {
		// get the perimeters to base the brim on them
		perimeters, err := modifier.Perimeters(layers[layerNr])
		if err != nil {
			return err
		}
		if perimeters == nil {
			return nil
		}

		var allOuterPerimeters []data.LayerPart

		for _, part := range perimeters {
			for _, wall := range part {
				if len(wall) > 0 {
					allOuterPerimeters = append(allOuterPerimeters, wall[0])
				}
			}
		}

		cl := clip.NewClipper()
		// Get the top level polys e.g. the polygons which are not inside another.
		topLevelPerimeters, _ := cl.TopLevelPolygons(allOuterPerimeters)
		allOuterPerimeters = nil
		for _, p := range topLevelPerimeters {
			allOuterPerimeters = append(allOuterPerimeters, data.NewBasicLayerPart(p, nil))
		}

		brim := cl.InsetLayer(allOuterPerimeters, -options.Printer.ExtrusionWidth, options.Print.BrimSkirt.BrimCount, options.Printer.ExtrusionWidth)

		for _, part := range brim {
			for _, wall := range part {
				for _, path := range wall {
					b.AddPolygon(nil, path.Outline(), z, false)
				}
			}
		}

	}

	return nil
}
