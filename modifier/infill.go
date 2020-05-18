package modifier

import (
	"GoSlice/clip"
	"GoSlice/data"
	"GoSlice/handler"
	"errors"
	"fmt"
	"strconv"
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

func (m infillModifier) Modify(layerNr int, layers []data.PartitionedLayer) ([]data.PartitionedLayer, error) {
	overlappingPerimeters, err := OverlapPerimeters(layers[layerNr])
	if err != nil || overlappingPerimeters == nil {
		return layers, err
	}

	perimeters, err := Perimeters(layers[layerNr])
	if err != nil || perimeters == nil {
		return layers, err
	}

	var bottomInfill []data.LayerPart
	var topInfill []data.LayerPart

	// calculate the bottom/top parts for each inner perimeter part
	for partNr, part := range perimeters {
		// for the last (most inner) inset of each part
		for insetPartNr, insetPart := range part[len(part)-1] {
			fmt.Println("layerNr " + strconv.Itoa(layerNr) + " partNr " + strconv.Itoa(partNr) + " insertPart " + strconv.Itoa(insetPartNr))

			var bottomInfillParts, topInfillParts []data.LayerPart
			var err error

			// 1. Calculate the area which needs full infill for top and bottom layerS
			if layerNr == 0 {
				// Just fill the bottom layer.
				bottomInfillParts = []data.LayerPart{insetPart}
			} else if layerNr == len(layers)-1 {
				// Just fill the top layer.
				topInfillParts = []data.LayerPart{insetPart}
			} else {
				// Subtract the below / above layer to get the parts which need infill.
				bottomInfillParts, err = partDifference(insetPart, layers[layerNr-1])
				if err != nil {
					return nil, err
				}

				topInfillParts, err = partDifference(insetPart, layers[layerNr+1])
				if err != nil {
					return nil, err
				}
			}

			// 2. Exset the area which needs infill to generate the internal overlap of top and bottom layer.
			var internalOverlappingBottomParts, internalOverlappingTopParts []data.LayerPart
			for _, bottomPart := range bottomInfillParts {
				overlappingParts, err := calculateOverlapPerimeter(bottomPart, m.options.Print.InfillOverlapPercent+m.options.Print.AdditionalInternalInfillOverlapPercent, m.options.Printer.ExtrusionWidth)
				if err != nil {
					return nil, err
				}

				internalOverlappingBottomParts = append(internalOverlappingBottomParts, overlappingParts...)
			}

			for _, topPart := range topInfillParts {
				overlappingParts, err := calculateOverlapPerimeter(topPart, m.options.Print.InfillOverlapPercent+m.options.Print.AdditionalInternalInfillOverlapPercent, m.options.Printer.ExtrusionWidth)
				if err != nil {
					return nil, err
				}

				internalOverlappingTopParts = append(internalOverlappingTopParts, overlappingParts...)
			}

			// 3. Clip the resulting areas by the overlappingPerimeters.
			c := clip.NewClipper()
			if internalOverlappingBottomParts != nil {
				clippedParts, ok := c.Intersection(internalOverlappingBottomParts, overlappingPerimeters[partNr])
				if !ok {
					return nil, errors.New("error while intersecting infill areas by the overlapping border")
				}
				bottomInfill = append(bottomInfill, clippedParts...)
			}

			if internalOverlappingTopParts != nil {
				clippedParts, ok := c.Intersection(internalOverlappingTopParts, overlappingPerimeters[partNr])
				if !ok {
					return nil, errors.New("error while intersecting infill areas by the overlapping border")
				}

				topInfill = append(topInfill, clippedParts...)
			}
		}
	}

	newLayer := newExtendedLayer(layers[layerNr])
	if len(bottomInfill) > 0 {
		newLayer.attributes["bottom"] = bottomInfill
	}
	if len(topInfill) > 0 {
		newLayer.attributes["top"] = topInfill
	}

	layers[layerNr] = newLayer

	return layers, nil
}
