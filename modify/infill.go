package modify

import (
	"GoSlice/clip"
	"GoSlice/data"
	"GoSlice/handle"
	"errors"
	"fmt"
	"strconv"
)

type infillModifier struct {
	options *data.Options
}

func (m infillModifier) Init(model data.OptimizedModel) {}

// NewInfillModifier calculates the areas which need infill and passes them as "bottom" attribute to the layer.
func NewInfillModifier(options *data.Options) handle.LayerModifier {
	return &infillModifier{
		options: options,
	}
}

// internalInfillOverlap is a magic number needed to compensate the extra inset done for each part which is needed for oblique walls.
const internalInfillOverlap = 400

func (m infillModifier) Modify(layerNr int, layers []data.PartitionedLayer) ([]data.PartitionedLayer, error) {
	perimeters, ok := layers[layerNr].Attributes()["perimeters"].([][][]data.LayerPart)
	if !ok {
		return layers, nil
	}
	// perimeters contains them as [part][insetNr][insetParts]

	var bottomInfill []data.LayerPart
	var topInfill []data.LayerPart
	var internalInfill []data.LayerPart

	// calculate the bottom parts for each inner perimeter part
	for partNr, part := range perimeters {
		// for the last (most inner) inset of each part
		for insertPart, insetPart := range part[len(part)-1] {
			fmt.Println("layerNr " + strconv.Itoa(layerNr) + " partNr " + strconv.Itoa(partNr) + " insertPart " + strconv.Itoa(insertPart))

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
				bottomInfillParts, err = m.partDifference(insetPart, layers[layerNr-1])
				if err != nil {
					return nil, err
				}

				topInfillParts, err = m.partDifference(insetPart, layers[layerNr+1])
				if err != nil {
					return nil, err
				}
			}

			// 2. Calculate the overlapping perimeter.
			maxOverlapBorder, err := m.overlapPerimeter(insetPart, m.options.Print.InfillOverlapPercent)
			if err != nil {
				return nil, err
			}

			// 3. Exset the area which needs infill to generate the internal overlap of top and bottom layer.
			var internalOverlappingBottomParts, internalOverlappingTopParts []data.LayerPart
			for _, bottomPart := range bottomInfillParts {
				overlappingParts, err := m.overlapPerimeter(bottomPart, m.options.Print.InfillOverlapPercent+internalInfillOverlap)
				if err != nil {
					return nil, err
				}

				internalOverlappingBottomParts = append(internalOverlappingBottomParts, overlappingParts...)
			}

			for _, topPart := range topInfillParts {
				overlappingParts, err := m.overlapPerimeter(topPart, m.options.Print.InfillOverlapPercent+internalInfillOverlap)
				if err != nil {
					return nil, err
				}

				internalOverlappingTopParts = append(internalOverlappingTopParts, overlappingParts...)
			}

			// 4. Clip the resulting areas by the maxOverlapBorder.
			c := clip.NewClipper()
			for _, part := range internalOverlappingBottomParts {
				clippedParts, ok := c.Intersection([]data.LayerPart{part}, maxOverlapBorder)
				if !ok {
					return nil, errors.New("error while intersecting infill areas by the max overlap border")
				}

				bottomInfill = append(bottomInfill, clippedParts...)
			}
			for _, part := range internalOverlappingTopParts {
				clippedParts, ok := c.Intersection([]data.LayerPart{part}, maxOverlapBorder)
				if !ok {
					return nil, errors.New("error while intersecting infill areas by the max overlap border")
				}

				topInfill = append(topInfill, clippedParts...)
			}

			// 5. Calculate the difference between the maxOverlapBorder and the final top/bottom infills
			//    to get the internal infill areas.

			// if no infill, just ignore the generation
			if m.options.Print.InfillPercent == 0 {
				continue
			}

			// calculating the difference would fail if both are nil so just ignore this
			if maxOverlapBorder == nil && bottomInfill == nil {
				continue
			}

			if parts, ok := c.Difference(maxOverlapBorder, bottomInfill); !ok {
				return nil, errors.New("error while calculating the difference between the max overlap border and the bottom infill")
			} else {
				internalInfill = append(internalInfill, parts...)
			}
		}
	}

	newLayer := newTypedLayer(layers[layerNr])
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

func (m infillModifier) partDifference(part data.LayerPart, layerToRemove data.PartitionedLayer) ([]data.LayerPart, error) {
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

func (m infillModifier) overlapPerimeter(part data.LayerPart, overlapPercent int) ([]data.LayerPart, error) {
	perimeterOverlap := data.Micrometer(float32(m.options.Printer.ExtrusionWidth) * (100.0 - float32(overlapPercent)) / 100.0)

	if perimeterOverlap != 0 {
		c := clip.NewClipper()
		// as we use only one inset, just return index 0
		return c.Inset(part, perimeterOverlap, 1)[0], nil
	} else {
		return []data.LayerPart{part}, nil
	}
}
