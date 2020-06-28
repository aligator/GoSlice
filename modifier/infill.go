package modifier

import (
	"GoSlice/clip"
	"GoSlice/data"
	"GoSlice/handler"
	"errors"
)

type infillModifier struct {
	options *data.Options
}

func (m infillModifier) Init(model data.OptimizedModel) {}

// NewInfillModifier calculates the areas which need infill and passes them as "bottom" attribute to the layer.
func NewInfillModifier(options *data.Options) handler.LayerModifier {
	return &infillModifier{
		options: options,
	}
}

// BottomInfill extracts the attribute "bottom" from the layer.
// If it has the wrong type, a error is returned.
// If it doesn't exist, (nil, nil) is returned.
// If it exists, the infill is returned.
func BottomInfill(layer data.PartitionedLayer) ([]data.LayerPart, error) {
	return InfillParts(layer, "bottom")
}

// TopInfill extracts the attribute "top" from the layer.
// If it has the wrong type, a error is returned.
// If it doesn't exist, (nil, nil) is returned.
// If it exists, the infill is returned.
func TopInfill(layer data.PartitionedLayer) ([]data.LayerPart, error) {
	return InfillParts(layer, "top")
}

// TopInfill extracts the given attribute" from the layer.
// If it has the wrong type, a error is returned.
// If it doesn't exist, (nil, nil) is returned.
// If it exists, the infill is returned.
func InfillParts(layer data.PartitionedLayer, typ string) ([]data.LayerPart, error) {
	if attr, ok := layer.Attributes()[typ]; ok {
		parts, ok := attr.([]data.LayerPart)
		if !ok {
			return nil, errors.New("the attribute " + typ + " has the wrong datatype")
		}

		return parts, nil
	}

	return nil, nil
}

func (m infillModifier) Modify(layerNr int, layers []data.PartitionedLayer) error {
	overlappingPerimeters, err := OverlapPerimeters(layers[layerNr])
	if err != nil || overlappingPerimeters == nil {
		return err
	}

	perimeters, err := Perimeters(layers[layerNr])
	if err != nil || perimeters == nil {
		return err
	}

	var bottomInfill []data.LayerPart
	var topInfill []data.LayerPart

	c := clip.NewClipper()

	// Calculate the bottom/top parts for each inner perimeter part.
	// It also takes into account the configured number of top/bottom layers.
	for partNr, part := range perimeters {
		// for the last (most inner) inset of each part
		for _, insetPart := range part[len(part)-1] {
			var bottomInfillParts, topInfillParts []data.LayerPart
			// 1. Calculate the area which needs full infill for top and bottom layerS

			// TODO: maybe merge these two loops in one function somehow?
			// calculate the difference with the layers bellow.
			for i := 0; i < m.options.Print.NumberBottomLayers; i++ {
				var parts []data.LayerPart
				if layerNr-i == 0 {
					// if it's the first layer, use the whole layer
					parts = []data.LayerPart{insetPart}
				} else if i > layerNr {
					// if we are below layer 0 stop calculation
					break
				} else {
					// else calculate the difference and use it
					parts, err = partDifference(insetPart, layers[layerNr-1-i])
					if err != nil {
						return err
					}
				}

				// union the parts if needed
				if len(bottomInfillParts) == 0 {
					bottomInfillParts = parts
				} else {
					var ok bool
					bottomInfillParts, ok = c.Union(bottomInfillParts, parts)
					if !ok {
						return errors.New("could not union bottom parts")
					}
				}
			}

			// calculate the difference with the layers above
			for i := 0; i < m.options.Print.NumberTopLayers; i++ {
				var parts []data.LayerPart
				if layerNr+i == len(layers)-1 {
					// if it's the last layer, use the whole layer
					parts = []data.LayerPart{insetPart}
				} else if layerNr+1+i >= len(layers) {
					// if we are above the top layer stop calculation
					break
				} else {
					// else calculate the difference and use it
					parts, err = partDifference(insetPart, layers[layerNr+1+i])
					if err != nil {
						return err
					}
				}

				// union the parts if needed
				if len(topInfillParts) == 0 {
					topInfillParts = parts
				} else {
					var ok bool
					topInfillParts, ok = c.Union(topInfillParts, parts)
					if !ok {
						return errors.New("could not union top parts")
					}
				}
			}

			// 2. Exset the area which needs infill to generate the internal overlap of top and bottom layer.
			var internalOverlappingBottomParts, internalOverlappingTopParts []data.LayerPart
			for _, bottomPart := range bottomInfillParts {
				overlappingParts, err := calculateOverlapPerimeter(bottomPart, m.options.Print.InfillOverlapPercent+m.options.Print.AdditionalInternalInfillOverlapPercent, m.options.Printer.ExtrusionWidth)
				if err != nil {
					return err
				}

				internalOverlappingBottomParts = append(internalOverlappingBottomParts, overlappingParts...)
			}

			for _, topPart := range topInfillParts {
				overlappingParts, err := calculateOverlapPerimeter(topPart, m.options.Print.InfillOverlapPercent+m.options.Print.AdditionalInternalInfillOverlapPercent, m.options.Printer.ExtrusionWidth)
				if err != nil {
					return err
				}

				internalOverlappingTopParts = append(internalOverlappingTopParts, overlappingParts...)
			}

			// 3. Clip the resulting areas by the overlappingPerimeters.
			if internalOverlappingBottomParts != nil {
				clippedParts, ok := c.Intersection(internalOverlappingBottomParts, overlappingPerimeters[partNr])
				if !ok {
					return errors.New("error while intersecting infill areas by the overlapping border")
				}

				u, ok := c.Union(bottomInfill, clippedParts)
				if !ok {
					return errors.New("error while calculating the union of new infill with already existing one")
				}
				bottomInfill = u
			}

			if internalOverlappingTopParts != nil {
				clippedParts, ok := c.Intersection(internalOverlappingTopParts, overlappingPerimeters[partNr])
				if !ok {
					return errors.New("error while intersecting infill areas by the overlapping border")
				}
				u, ok := c.Union(topInfill, clippedParts)
				if !ok {
					return errors.New("error while calculating the union of new infill with already existing one")
				}
				topInfill = u
			}
		}
	}

	if len(topInfill) > 0 && len(bottomInfill) > 0 {
		diff, ok := c.Difference(topInfill, bottomInfill)
		if !ok {
			return errors.New("error while calculating the difference of new top infill with the bottom infill to avoid duplicates")
		}
		topInfill = diff
	}

	newLayer := newExtendedLayer(layers[layerNr])
	if len(bottomInfill) > 0 {
		newLayer.attributes["bottom"] = bottomInfill
	}
	if len(topInfill) > 0 {
		newLayer.attributes["top"] = topInfill
	}

	layers[layerNr] = newLayer

	return nil
}
