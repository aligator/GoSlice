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

	var modifiedLayerParts []data.LayerPart

	for partNr, part := range insetParts {
		attributes := layers[layerNr].LayerParts()[partNr].Attributes()
		if attributes == nil {
			attributes = map[string]interface{}{}
		}
		attributes["perimeters"] = part

		modifiedLayerParts = append(modifiedLayerParts, typedLayerPart{
			LayerPart:  layers[layerNr].LayerParts()[partNr],
			attributes: attributes,
		})
	}

	layers[layerNr] = data.NewPartitionedLayer(modifiedLayerParts)

	return layers, nil
}
