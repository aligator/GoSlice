// This file provides a renderer for the brim lines generated by the brim modifier.
package renderer

import (
	"GoSlice/data"
	"GoSlice/gcode"
	"GoSlice/modifier"
)

// Brim just draws the brim lines generated by the brim modifier.
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

	// Use type SKIRT as Cura also does it the same. This is for support of the gcode viewer in Cura.
	b.AddComment("TYPE:SKIRT")

	err = nil
	brim.ForEach(func(part data.LayerPart, _, _, _ int) bool {
		err = b.AddPolygon(nil, part.Outline(), z, false)
		return err != nil
	})

	return err
}
