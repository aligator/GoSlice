// This file provides renderers for gcode injected at specific layers.

package renderer

import (
	"GoSlice/data"
	"GoSlice/gcode"
	"GoSlice/modifier"
)

// Brim generates the brim lines.
type Brim struct{}

func (Brim) Init(model data.OptimizedModel) {}

func (Brim) Render(b *gcode.Builder, layerNr int, maxLayer int, layer data.PartitionedLayer, z data.Micrometer, options *data.Options) error {
	// Get the brim data.
	brim, err := modifier.Brim(layer)
	if err != nil {
		return err
	}
	if brim == nil {
		return nil
	}

	// Use type SKIRT as Cura also does it. This is for support of the gcode viewer in Cura.
	b.AddComment("TYPE:SKIRT")

	for _, part := range brim {
		for _, wall := range part {
			for _, path := range wall {
				err := b.AddPolygon(nil, path.Outline(), z, false)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}
