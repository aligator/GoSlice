package modifier

import (
	"GoSlice/clip"
	"GoSlice/data"
	"GoSlice/handler"
	"errors"
)

type internalInfillModifier struct {
	options *data.Options
}

func (m internalInfillModifier) Init(model data.OptimizedModel) {}

// NewInfillModifier calculates the areas which need infill and passes them as "bottom" attribute to the layer.
func NewInternalInfillModifier(options *data.Options) handler.LayerModifier {
	return &internalInfillModifier{
		options: options,
	}
}

func (m internalInfillModifier) Modify(layers []data.PartitionedLayer) error {
	for layerNr := range layers {
		overlappingPerimeters, err := OverlapPerimeters(layers[layerNr])
		if err != nil || overlappingPerimeters == nil {
			return err
		}

		bottomInfill, err := BottomInfill(layers[layerNr])
		if err != nil {
			return err
		}

		topInfill, err := TopInfill(layers[layerNr])
		if err != nil {
			return err
		}

		var internalInfill []data.LayerPart

		c := clip.NewClipper()

		// calculate the bottom parts for each inner perimeter part
		for _, overlappingPart := range overlappingPerimeters {
			// Calculate the difference between the overlappingPerimeters and the final top/bottom infills
			// to get the internal infill areas.

			// if no infill, just ignore the generation
			if m.options.Print.InfillPercent == 0 {
				continue
			}

			// calculating the difference would fail if both are nil so just ignore this
			if overlappingPart == nil && bottomInfill == nil && topInfill == nil {
				continue
			}

			parts, ok := c.Difference(overlappingPart, append(bottomInfill, topInfill...))
			if !ok {
				return errors.New("error while calculating the difference between the max overlap border and the bottom infill")
			}

			internalInfill = append(internalInfill, parts...)
		}

		newLayer := newExtendedLayer(layers[layerNr])
		if len(internalInfill) > 0 {
			newLayer.attributes["infill"] = internalInfill
		}
	}

	return nil
}

func partDifference(part data.LayerPart, layerToRemove data.PartitionedLayer) ([]data.LayerPart, error) {
	var toClip []data.LayerPart

	for _, otherPart := range layerToRemove.LayerParts() {
		toClip = append(toClip, otherPart)
	}

	c := clip.NewClipper()

	diff, ok := c.Difference([]data.LayerPart{part}, toClip)
	if !ok {
		return nil, errors.New("error while calculating difference of a part and a layer")
	}

	return diff, nil
}
