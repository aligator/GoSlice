package modify

import (
	"GoSlice/data"
	"GoSlice/handle"
)

type partTypeModifier struct {
	options *data.Options
}

type typedLayerPart struct {
	data.LayerPart
	typ string
}

func (l typedLayerPart) Type() string {
	return l.typ
}

// NewPartTypeModifier checks for each part which type it is. (e.g. bottom, top, hanging, etc.)
func NewPartTypeModifier(options *data.Options) handle.LayerModifier {
	return &partTypeModifier{
		options: options,
	}
}

func (m partTypeModifier) Modify(layerNr int, layers []data.PartitionedLayer) ([]data.PartitionedLayer, error) {
	// for the first layer set everything to bottom
	if layerNr == 0 {
		var layerParts []data.LayerPart

		for _, part := range layers[layerNr].LayerParts() {
			layerParts = append(layerParts, typedLayerPart{
				LayerPart: data.NewUnknownLayerPart(part.Outline(), part.Holes()),
				typ:       "bottom",
			})
		}

		layers[layerNr] = data.NewPartitionedLayer(layerParts)
	}

	return layers, nil
}
