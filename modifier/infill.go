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

// internalInfillOverlap is added to the normal overlap to allow the infill to grow into the model.
const internalInfillOverlap = 400

/*
// modifyInfillLayer adds the addedPart to the given layer to the attribute with the given typ.
// It creates a union of them and clips by the most inner perimeter.
// TODO: maybe just add the "maybe to fill areas" in the modifier (also in the main modify method) and clip by the outline in another modifier at the end.
// 		 This would prevent multiple clips by the most inner perimeters.
func (m infillModifier) modifyInfillLayer(layer data.PartitionedLayer, typ string, addedParts []data.LayerPart) (data.PartitionedLayer, error) {
	oldInfill, ok := layer.Attributes()[typ].([]data.LayerPart)
	if !ok {
		// just use the new added part
		oldInfill = addedParts
	}

	c := clip.NewClipper()

	var parts []data.LayerPart

	if len(oldInfill) != 0 && len(addedParts) != 0 {
		parts, ok = c.Union(oldInfill, addedParts)
		if !ok {
			fmt.Println("", oldInfill, addedParts)
			return layer, nil// errors.New("could not combine the old infill parts with the new ones")
		}
	} else if len(oldInfill) == 0 {
		parts = addedParts
	} else {
		parts = oldInfill
	}

	// get the perimeters to clip the part-to-add by the most inner one
	perimeters, ok := layer.Attributes()["perimeters"].([][][]data.LayerPart)
	if !ok {
		return layer, nil
	}

	var toRemove []data.LayerPart

	for _, part := range perimeters {
		// for the last (most inner) inset of each part
		for _, insetPart := range part[len(part)-1] {
			maxOverlapBorder, err := calculateOverlapPerimeter(insetPart, m.options.Print.InfillOverlapPercent, m.options.Printer.ExtrusionWidth)
			if err != nil {
				return nil, err
			}

			toRemove = append(toRemove, maxOverlapBorder...)
		}
	}

	newParts, ok := c.Intersection(parts, toRemove)
	if !ok {
		return nil, errors.New("could not clip the infill parts by the border")
	}

	newLayer := newTypedLayer(layer)
	if len(newParts) > 0 {
		newLayer.attributes[typ] = newParts
	}

	return newLayer, nil
}*/

func (m infillModifier) Modify(layerNr int, layers []data.PartitionedLayer) ([]data.PartitionedLayer, error) {
	overlappingPerimeters, ok := layers[layerNr].Attributes()["overlapPerimeters"].([][]data.LayerPart)
	// overlappingPerimeters contains them as [part][insetParts]
	if !ok {
		return layers, nil
	}
	perimeters, ok := layers[layerNr].Attributes()["perimeters"].([][][]data.LayerPart)
	// perimeters contains them as [part][insetNr][insetParts]
	if !ok {
		return layers, nil
	}

	var bottomInfill []data.LayerPart
	var topInfill []data.LayerPart
	var internalInfill []data.LayerPart

	// calculate the bottom parts for each inner perimeter part
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
				overlappingParts, err := calculateOverlapPerimeter(bottomPart, m.options.Print.InfillOverlapPercent+internalInfillOverlap, m.options.Printer.ExtrusionWidth)
				if err != nil {
					return nil, err
				}

				internalOverlappingBottomParts = append(internalOverlappingBottomParts, overlappingParts...)
			}

			for _, topPart := range topInfillParts {
				overlappingParts, err := calculateOverlapPerimeter(topPart, m.options.Print.InfillOverlapPercent+internalInfillOverlap, m.options.Printer.ExtrusionWidth)
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
			// 4. Calculate the difference between the overlappingPerimeters and the final top/bottom infills
			//    to get the internal infill areas.

			// if no infill, just ignore the generation
			if m.options.Print.InfillPercent == 0 {
				continue
			}

			// calculating the difference would fail if both are nil so just ignore this
			if overlappingPerimeters[partNr] == nil && bottomInfill == nil {
				continue
			}

			if parts, ok := c.Difference(overlappingPerimeters[partNr], bottomInfill); !ok {
				return nil, errors.New("error while calculating the difference between the max overlap border and the bottom infill")
			} else {
				internalInfill = append(internalInfill, parts...)
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
	if len(internalInfill) > 0 {
		newLayer.attributes["infill"] = internalInfill
	}

	layers[layerNr] = newLayer

	return layers, nil
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
