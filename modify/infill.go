package modify

import (
	"GoSlice/clip"
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
		perimeters, ok := layers[layerNr].Attributes()["perimeters"].([][][]data.LayerPart)
		if !ok {
			return layers, nil
		}
		// perimeters contains them as [part][insetNr][insetParts]

		c := clip.NewClipper()

		var infill []data.Paths

		for _, part := range perimeters {
			// for the last (most inner) inset of each part
			for _, insetPart := range part[len(part)-1] {
				infill = append(infill, c.Fill(insetPart, m.options.Printer.ExtrusionWidth, m.options.Print.InfillOverlapPercent))
			}
		}

		newLayer := newTypedLayer(layers[layerNr])

		if len(infill) > 0 {
			newLayer.attributes["bottom"] = infill
		}

		layers[layerNr] = newLayer
	}

	return layers, nil
}
