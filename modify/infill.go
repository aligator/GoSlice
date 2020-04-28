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

// NewInfillModifier calculates the areas which need infill and passes them as "bottom" attribute to the layer.
func NewInfillModifier(options *data.Options) handle.LayerModifier {
	return &infillModifier{
		options: options,
	}
}

// internalInfillOverlap is a magic number needed to compensate the extra inset done for each part which is needed for oblique walls.
const internalInfillOverlap = 200

func (m infillModifier) Modify(layerNr int, layers []data.PartitionedLayer) ([]data.PartitionedLayer, error) {
	perimeters, ok := layers[layerNr].Attributes()["perimeters"].([][][]data.LayerPart)
	if !ok {
		return layers, nil
	}
	// perimeters contains them as [part][insetNr][insetParts]

	c := clip.NewClipper()
	var bottomInfill []data.Paths
	var topInfill []data.Paths

	min, max := layers[layerNr].Bounds()
	pattern := c.LinearPattern(min, max, m.options.Printer.ExtrusionWidth)

	// TODO remove code duplication of top and bottom layer generation
	// calculate the bottom parts for each inner perimeter part
	for partNr, part := range perimeters {
		// for the last (most inner) inset of each part
		for insertPart, insetPart := range part[len(part)-1] {
			fmt.Println("layerNr " + strconv.Itoa(layerNr) + " partNr " + strconv.Itoa(partNr) + " insertPart " + strconv.Itoa(insertPart))
			if layerNr == 0 {
				// for the first layer bottomInfill everything
				bottomInfill = append(bottomInfill, c.Fill(insetPart, nil, m.options.Printer.ExtrusionWidth, pattern, m.options.Print.InfillOverlapPercent, internalInfillOverlap))
				continue
			} else if layerNr == len(layers)-1 {
				// for the last layer topInfill everything
				topInfill = append(topInfill, c.Fill(insetPart, nil, m.options.Printer.ExtrusionWidth, pattern, m.options.Print.InfillOverlapPercent, internalInfillOverlap))
				continue
			}

			// For the other layers detect the bottom parts by calculating the difference between the current most inner perimeter and the layer below.
			// Also detect the top parts by calculating the difference between the current current most inner perimeter and the layer above
			var toClipBelow []data.LayerPart
			var toClipAbove []data.LayerPart

			for _, belowPart := range layers[layerNr-1].LayerParts() {
				toClipBelow = append(toClipBelow, belowPart)
			}

			for _, abovePart := range layers[layerNr+1].LayerParts() {
				toClipAbove = append(toClipAbove, abovePart)
			}

			fmt.Println("calculate difference with layer below")
			toInfillBottom, ok := c.Difference(insetPart, toClipBelow)
			if !ok {
				return nil, errors.New("error while calculating difference with previous layer for detecting bottom parts")
			}

			fmt.Println("calculate difference with layer above")
			toInfillTop, ok := c.Difference(insetPart, toClipAbove)
			if !ok {
				return nil, errors.New("error while calculating difference with next layer for detecting top parts")
			}

			for _, fill := range toInfillBottom {
				bottomInfill = append(bottomInfill, c.Fill(fill, insetPart, m.options.Printer.ExtrusionWidth, pattern, m.options.Print.InfillOverlapPercent, internalInfillOverlap))
			}
			for _, fill := range toInfillTop {
				topInfill = append(topInfill, c.Fill(fill, insetPart, m.options.Printer.ExtrusionWidth, pattern, m.options.Print.InfillOverlapPercent, internalInfillOverlap))
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

	layers[layerNr] = newLayer

	return layers, nil
}
