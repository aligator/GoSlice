package modifier

import (
	"GoSlice/clip"
	"GoSlice/data"
	"GoSlice/handler"
	"errors"
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

// OverlapPerimeters extracts the attribute "overlapPerimeters" from the layer.
// If it has the wrong type, a error is returned.
// If it doesn't exist, (nil, nil) is returned.
// If it exists, the perimeters are returned as [part][insetParts]data.LayerPart.
func OverlapPerimeters(layer data.PartitionedLayer) ([][]data.LayerPart, error) {
	if attr, ok := layer.Attributes()["overlapPerimeters"]; ok {
		overlappingPerimeters, ok := attr.([][]data.LayerPart)
		if !ok {
			return nil, errors.New("the attribute overlapPerimeters has the wrong datatype")
		}

		return overlappingPerimeters, nil
	}

	return nil, nil
}

// Perimeters extracts the attribute "perimeters" from the layer.
// If it has the wrong type, a error is returned.
// If it doesn't exist, (nil, nil) is returned.
// If it exists, the perimeters are returned as [part][insetNr][insetParts]data.LayerPart.
func Perimeters(layer data.PartitionedLayer) ([][][]data.LayerPart, error) {
	if attr, ok := layer.Attributes()["perimeters"]; ok {
		perimeters, ok := attr.([][][]data.LayerPart)
		if !ok {
			return nil, errors.New("the attribute perimeters has the wrong datatype")
		}

		return perimeters, nil
	}

	return nil, nil
}

func (m perimeterModifier) Init(model data.OptimizedModel) {}

func (m perimeterModifier) Modify(layers []data.PartitionedLayer) error {
	for layerNr := range layers {
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
					return err
				}
				overlapPerimeter[partNr] = append(overlapPerimeter[partNr], maxOverlapBorder...)
			}
		}

		newLayer := newExtendedLayer(layers[layerNr])
		newLayer.attributes["perimeters"] = insetParts
		newLayer.attributes["overlapPerimeters"] = overlapPerimeter
		layers[layerNr] = newLayer
	}

	return nil
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
