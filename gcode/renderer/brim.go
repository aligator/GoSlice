// This file provides renderers for gcode injected at specific layers.

package renderer

import (
	"GoSlice/clip"
	"GoSlice/data"
	"GoSlice/gcode"
	"GoSlice/modifier"
	"errors"
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
		// Get the perimeters and support to base the brim on them.
		perimeters, err := modifier.Perimeters(layers[layerNr])
		if err != nil {
			return err
		}

		support, err := modifier.FullSupport(layers[layerNr])
		if err != nil {
			return err
		}
		if support == nil && perimeters == nil {
			return nil
		}

		// Extract the outer perimeters of all perimeters.
		var allOuterPerimeters []data.LayerPart

		for _, part := range perimeters {
			for _, wall := range part {
				if len(wall) > 0 {
					allOuterPerimeters = append(allOuterPerimeters, wall[0])
				}
			}
		}

		cl := clip.NewClipper()

		// Combine the outerPerimeters with the support area.
		objectArea, ok := cl.Union(allOuterPerimeters, support)
		if !ok {
			return errors.New("could not union the outer perimeters with the support")
		}

		// Get the top level polys e.g. the polygons which are not inside another.
		topLevelPerimeters, _ := cl.TopLevelPolygons(objectArea)
		allOuterPerimeters = nil
		for _, p := range topLevelPerimeters {
			allOuterPerimeters = append(allOuterPerimeters, data.NewBasicLayerPart(p, nil))
		}

		if allOuterPerimeters == nil {
			// No need to go further and prevent fail of union.
			return nil
		}

		// Generate the brim.
		brim := cl.InsetLayer(objectArea, -options.Printer.ExtrusionWidth, options.Print.BrimSkirt.BrimCount, options.Printer.ExtrusionWidth)

		for _, part := range brim {
			for _, wall := range part {
				for _, path := range wall {
					// Remove support from the brim at the same location to avoid overlapping of them
					res, ok := cl.Difference([]data.LayerPart{path}, support)

					if !ok {
						return errors.New("could not remove the support from the brim line")
					}

					for _, r := range res {
						err := b.AddPolygon(nil, r.Outline(), z, false)
						if err != nil {
							return err
						}
					}

				}
			}
		}
	}

	return nil
}
