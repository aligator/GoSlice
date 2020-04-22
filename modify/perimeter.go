package modify

import (
	"GoSlice/clip"
	"GoSlice/data"
	"GoSlice/handle"
)

type perimeterModifier struct {
	options *data.Options
}

// NewPerimeterModifier creates a modifier which calculates all perimeters
//
// The perimeters are saved as attribute in the LayerPart.
func NewPerimeterModifier(options *data.Options) handle.LayerModifier {
	return &perimeterModifier{
		options: options,
	}
}

func (m perimeterModifier) Modify(layerNr int, layers []data.PartitionedLayer) ([]data.PartitionedLayer, error) {
	// generate perimeters
	c := clip.NewClipper()
	insetParts := c.InsetLayer(layers[layerNr].LayerParts(), m.options.Printer.ExtrusionWidth, m.options.Print.InsetCount)

	var innerPerimeters []data.Paths
	var outerPerimeters []data.Paths
	var middlePerimeters []data.Paths

	// iterate over all generated perimeters
	for _, part := range insetParts {
		for _, wall := range part {
			for insetNum, wallInset := range wall {
				if insetNum == 0 {
					outerPerimeters = append(outerPerimeters, wallInset)
				} else if insetNum > 0 && insetNum < len(wall)-1 {
					middlePerimeters = append(middlePerimeters, wallInset)
				} else {
					innerPerimeters = append(innerPerimeters, wallInset)
				}
			}
		}
	}

	newLayer := newTypedLayer(layers[layerNr])

	newLayer.attributes["perimeters"] = [3][]data.Paths{
		outerPerimeters,
		middlePerimeters,
		innerPerimeters,
	}

	layers[layerNr] = newLayer

	return layers, nil
}
