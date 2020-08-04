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
	return modifyConcurrently(layers, func(layerCh <-chan enumeratedPartitionedLayer, outputCh chan<- enumeratedPartitionedLayer, errCh chan<- error) {
		for layer := range layerCh {
			overlappingPerimeters, err := OverlapPerimeters(layer.layer)
			if err != nil || overlappingPerimeters == nil {
				errCh <- err
				return
			}

			bottomInfill, err := BottomInfill(layer.layer)
			if err != nil {
				errCh <- err
				return
			}

			topInfill, err := TopInfill(layer.layer)
			if err != nil {
				errCh <- err
				return
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

				if parts, ok := c.Difference(overlappingPart, append(bottomInfill, topInfill...)); !ok {
					errCh <- errors.New("error while calculating the difference between the max overlap border and the bottom infill")
					return
				} else {
					internalInfill = append(internalInfill, parts...)
				}
			}

			newLayer := newExtendedLayer(layer.layer)
			if len(internalInfill) > 0 {
				newLayer.attributes["infill"] = internalInfill
			}
			outputCh <- enumeratedPartitionedLayer{
				layer:   newLayer,
				layerNr: layer.layerNr,
			}
		}
	})
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
