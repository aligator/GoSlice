package modifier

import (
	"GoSlice/clip"
	"GoSlice/data"
	"GoSlice/handler"
)

type perimeterModifier struct {
	options *data.Options
}

// NewPerimeterModifier creates a modifier which calculates all perimeters
//
// The perimeters are saved as attribute in the LayerPart.
func NewPerimeterModifier(options *data.Options) handler.LayerModifier {
	return &perimeterModifier{
		options: options,
	}
}

func (m perimeterModifier) Init(model data.OptimizedModel) {}

func (m perimeterModifier) Modify(layerNr int, layers []data.PartitionedLayer) ([]data.PartitionedLayer, error) {
	// Generate the perimeters.
	c := clip.NewClipper()
	insetParts := c.InsetLayer(layers[layerNr].LayerParts(), m.options.Printer.ExtrusionWidth, m.options.Print.InsetCount)

	// Also generate the overlapping perimeter, which helps with calculating the infill.
	// This is derived from the most inner perimeters and offset by the options.Print.InfillOverlapPercent option.

	var overlapPerimeter [][]data.LayerPart

	for partNr, part := range insetParts {
		if len(overlapPerimeter) >= partNr {
			overlapPerimeter = append(overlapPerimeter, nil)
		}

		// Use only the most inner perimeter.
		for _, insetPart := range part[len(part)-1] {

			maxOverlapBorder, err := calculateOverlapPerimeter(insetPart, m.options.Print.InfillOverlapPercent, m.options.Printer.ExtrusionWidth)
			if err != nil {
				return nil, err
			}
			overlapPerimeter[partNr] = append(overlapPerimeter[partNr], maxOverlapBorder...)
		}
	}

	newLayer := newExtendedLayer(layers[layerNr])
	newLayer.attributes["perimeters"] = insetParts
	newLayer.attributes["overlapPerimeters"] = overlapPerimeter
	layers[layerNr] = newLayer

	return layers, nil
}

// calculateOverlapPerimeter helper function for calculating the overlap-perimeter out of a layer part.
func calculateOverlapPerimeter(part data.LayerPart, overlapPercent int, extrusionWidth data.Micrometer) ([]data.LayerPart, error) {
	perimeterOverlap := data.Micrometer(float32(extrusionWidth) * (100.0 - float32(overlapPercent)) / 100.0)

	if perimeterOverlap != 0 {
		c := clip.NewClipper()
		// As we use only one inset, just return index 0.
		return c.Inset(part, perimeterOverlap, 1)[0], nil
	} else {
		// If no overlap needed, just return the input part.
		return []data.LayerPart{part}, nil
	}
}
