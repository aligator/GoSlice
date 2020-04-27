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
	var bottomInfill []data.Paths

	min, max := layers[layerNr].Bounds()
	pattern := c.LinearPattern(min, max, m.options.Printer.ExtrusionWidth)

	// calculate the bottom parts for each inner perimeter part
	for partNr, part := range perimeters {
		// for the last (most inner) inset of each part
		for insertPart, insetPart := range part[len(part)-1] {
			fmt.Println("layerNr " + strconv.Itoa(layerNr) + " partNr " + strconv.Itoa(partNr) + " insertPart " + strconv.Itoa(insertPart))
			if layerNr == 0 {
				// for the first layer bottomInfill everything
				bottomInfill = append(bottomInfill, c.Fill(insetPart, nil, m.options.Printer.ExtrusionWidth, pattern, m.options.Print.InfillOverlapPercent))
				continue
			}

			// For the other layers detect the bottom parts by calculating the difference between the layer below
			// and the current most inner perimeter.
			var toRemove []data.LayerPart

			if len(layers[layerNr-1].LayerParts())-1 < partNr {
				continue
			}

			// exset the below parts to fix generating very small infills for oblique walls (e.g. the sides of the 3DBenchy model)
			below := layers[layerNr-1].LayerParts()[partNr]
			belowBigger := c.Inset(below, -m.options.Printer.ExtrusionWidth, 1)

			for _, parts := range belowBigger {
				for _, layerPart := range parts {
					toRemove = append(toRemove, layerPart)
				}
			}

			fmt.Println("calculate difference")
			toInfill, ok := c.Difference(insetPart, toRemove)
			if !ok {
				return nil, errors.New("error while calculating difference with previous layer for detecting bottom parts")
			}

			for _, fill := range toInfill {
				bottomInfill = append(bottomInfill, c.Fill(fill, insetPart, m.options.Printer.ExtrusionWidth, pattern, m.options.Print.InfillOverlapPercent))
			}
		}
	}

	newLayer := newTypedLayer(layers[layerNr])
	if len(bottomInfill) > 0 {
		newLayer.attributes["bottom"] = bottomInfill
	}

	layers[layerNr] = newLayer

	return layers, nil
}
