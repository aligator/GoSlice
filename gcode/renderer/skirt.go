// This file provides renderers for gcode injected at specific layers.

package renderer

import (
	"GoSlice/clip"
	"GoSlice/data"
	"GoSlice/gcode"
	"GoSlice/modifier"
	"errors"
)

// Skirt generates the skirt lines.
type Skirt struct{}

func (Skirt) Init(model data.OptimizedModel) {}

func (Skirt) Render(b *gcode.Builder, layerNr int, layers []data.PartitionedLayer, z data.Micrometer, options *data.Options) error {
	if options.Print.BrimSkirt.SkirtCount == 0 {
		return nil
	}

	// TODO: add comment used by cura
	//b.AddComment("LAYER:%v", layerNr)
	if layerNr == 0 {
		// Get the perimeters and support to base the hull (line around everything) on them.
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

		// Skirt distance + (1/2 extrusion with of the model side + 1/2 extrusion width of the most inner brim line) + the brim width
		// is the distance between the perimeter (or brim) and skirt.
		distance := options.Print.BrimSkirt.SkirtDistance.ToMicrometer() + (options.Printer.ExtrusionWidth * data.Micrometer(options.Print.BrimSkirt.BrimCount)) + options.Printer.ExtrusionWidth

		// Draw the skirt.
		c := clip.NewClipper()
		// Generate the hull around everything.
		hull, ok := c.Hull(append(support, perimeters.ToOneDimension()...))
		if !ok {
			return errors.New("could not generate hull around all perimeters to create the skirt")
		}

		// Generate all skirt lines by exsetting the hull.
		skirt := c.Inset(data.NewBasicLayerPart(hull, nil), -options.Printer.ExtrusionWidth, options.Print.BrimSkirt.SkirtCount, distance)

		for _, wall := range skirt {
			for _, loopPart := range wall {
				// As we use the hull around the whole object there shouldn't be any collision with the model -> currentLayer is nil
				err := b.AddPolygon(nil, loopPart.Outline(), z, false)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}
