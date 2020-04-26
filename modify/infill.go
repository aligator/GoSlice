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

	min, max := layers[layerNr].Bounds()
	pattern := c.LinearPattern(min, max, m.options.Printer.ExtrusionWidth)

	// calculate the bottom parts for each inner perimeter part
	for partNr, part := range perimeters {
		// for the last (most inner) inset of each part
		for insertPart, insetPart := range part[len(part)-1] {
			fmt.Println("layerNr " + strconv.Itoa(layerNr) + " partNr " + strconv.Itoa(partNr) + " insertPart " + strconv.Itoa(insertPart))
			if layerNr == 0 {
				// for the first layer infill everything
				infill = append(infill, c.Fill(insetPart, nil, m.options.Printer.ExtrusionWidth, pattern, m.options.Print.InfillOverlapPercent))
				continue
			}

			perimetersBelow, ok := layers[layerNr-1].Attributes()["perimeters"].([][][]data.LayerPart)
			if !ok {
				return nil, errors.New("wrong type for attribute perimeters")
			}

			var toRemove []data.LayerPart

			// remove each part below from the current part
			// for _, partBelow := range perimetersBelow {
			// TODO: verify if it is always enough to check only the one layer with the partNr
			if len(perimetersBelow)-1 >= partNr {

				// use the 2nd last perimeters of the below parts to get some overlap.
				// if the partBelow has only one perimeter, use it
				var perimeter []data.LayerPart
				if len(perimetersBelow[partNr]) == 1 {
					// Todo: 3DBenchy with 1 Perimeter-settings need improvement
					perimeter = perimetersBelow[partNr][0]
				} else {
					perimeter = perimetersBelow[partNr][len(perimetersBelow[partNr])-2]
				}

				for _, insetPartBelow := range perimeter {
					toRemove = append(toRemove, insetPartBelow)
				}
			}

			fmt.Println("calculate difference")
			toInfill, ok := c.Difference(insetPart, toRemove)
			if !ok {
				return nil, errors.New("error while calculating difference with previous layer for detecting bottom parts")
			}

			for _, fill := range toInfill {
				infill = append(infill, c.Fill(fill, insetPart, m.options.Printer.ExtrusionWidth, pattern, m.options.Print.InfillOverlapPercent))
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
