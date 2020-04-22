package modify

import (
	"GoSlice/data"
	"GoSlice/handle"
)

type infillModifier struct {
	options *data.Options
}

// NewInfillModifier calculates the areas which need infill and passes them as "bottom" attribute to the layer.
func NewInfillModifier(options *data.Options) handle.LayerModifier {
	return &infillModifier{
		options: options,
	}
}

func (m infillModifier) Modify(layerNr int, layers []data.PartitionedLayer) ([]data.PartitionedLayer, error) {
	// for the first layer set everything to bottom
	if layerNr == 0 {
		perimeters, ok := layers[layerNr].Attributes()["perimeters"].([3][]data.Paths)
		if !ok {
			return layers, nil
		}

		innerPaths := perimeters[2]
		if len(innerPaths) == 0 {
			innerPaths = perimeters[1]
		}
		if len(innerPaths) == 0 {
			innerPaths = perimeters[0]
		}

		newLayer := newTypedLayer(layers[layerNr])

		if len(innerPaths) > 0 {
			newLayer.attributes["bottom"] = innerPaths
		}

		layers[layerNr] = newLayer
	}

	return layers, nil
}
