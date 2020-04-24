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

func (m infillModifier) Modify(layerNr int, layers []data.PartitionedLayer) ([]data.PartitionedLayer, error) {
	perimeters, ok := layers[layerNr].Attributes()["perimeters"].([][][]data.LayerPart)
	if !ok {
		return layers, nil
	}
	// perimeters contains them as [part][insetNr][insetParts]

	c := clip.NewClipper()
	var infill []data.Paths

	// calculate the bottom parts for each inner perimeter part
	for partNr, part := range perimeters {
		// for the last (most inner) inset of each part
		for insertPart, insetPart := range part[len(part)-1] {
			fmt.Println("layerNr " + strconv.Itoa(layerNr) + " partNr " + strconv.Itoa(partNr) + " insertPart " + strconv.Itoa(insertPart))
			if layerNr == 0 {
				// for the first layer infill everything
				infill = append(infill, c.Fill(insetPart, m.options.Printer.ExtrusionWidth, m.options.Print.InfillOverlapPercent))
				continue
			}

			perimetersBelow, ok := layers[layerNr-1].Attributes()["perimeters"].([][][]data.LayerPart)
			if !ok {
				return nil, errors.New("wrong type for attribute perimeters")
			}

			var toRemove []data.LayerPart

			// remove each part below from the current part
			for _, partBelow := range perimetersBelow {
				// for the last (most inner) inset of each part
				for _, insetPartBelow := range partBelow[len(partBelow)-1] {
					toRemove = append(toRemove, insetPartBelow)
				}
			}

			fmt.Println("calculate difference")
			toInfill, ok := c.Difference(insetPart, toRemove)
			if !ok {
				return nil, errors.New("error while calculating difference with previous layer for detecting bottom parts")
			}

			for _, fill := range toInfill {
				infill = append(infill, c.Fill(fill, m.options.Printer.ExtrusionWidth, m.options.Print.InfillOverlapPercent))
			}
		}
	}

	newLayer := newTypedLayer(layers[layerNr])
	if len(infill) > 0 {
		newLayer.attributes["bottom"] = infill
	}

	layers[layerNr] = newLayer

	return layers, nil
}
