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
	pattern clip.Pattern
}

func (m *infillModifier) Init(model data.OptimizedModel) {
	m.pattern = clip.NewLinearPattern(model.Min().PointXY(), model.Max().PointXY(), m.options.Printer.ExtrusionWidth)
}

// NewInfillModifier calculates the areas which need infill and passes them as "bottom" attribute to the layer.
func NewInfillModifier(options *data.Options) handle.LayerModifier {
	return &infillModifier{
		options: options,
	}
}

// internalInfillOverlap is a magic number needed to compensate the extra inset done for each part which is needed for oblique walls.
const internalInfillOverlap = 200

func (m *infillModifier) Modify(layerNr int, layers []data.PartitionedLayer) ([]data.PartitionedLayer, error) {
	perimeters, ok := layers[layerNr].Attributes()["perimeters"].([][][]data.LayerPart)
	if !ok {
		return layers, nil
	}
	// perimeters contains them as [part][insetNr][insetParts]

	var bottomInfill []data.Paths
	var topInfill []data.Paths

	// calculate the bottom parts for each inner perimeter part
	for partNr, part := range perimeters {
		// for the last (most inner) inset of each part
		for insertPart, insetPart := range part[len(part)-1] {
			fmt.Println("layerNr " + strconv.Itoa(layerNr) + " partNr " + strconv.Itoa(partNr) + " insertPart " + strconv.Itoa(insertPart))

			infill, err := m.genTopBottomInfill(insetPart, layerNr-1, layers, m.pattern)
			if err != nil {
				return nil, err
			}
			for _, paths := range infill {
				bottomInfill = append(bottomInfill, paths)

			}

			infill, err = m.genTopBottomInfill(insetPart, layerNr+1, layers, m.pattern)
			if err != nil {
				return nil, err
			}
			for _, paths := range infill {
				bottomInfill = append(bottomInfill, paths)

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

// genInfill returns the infill for the top or bottom parts.
// It calculates the difference of the layer with layerNr and the given part.
// Then it fills the result by using the given pattern.
func (m *infillModifier) genTopBottomInfill(part data.LayerPart, layerNr int, layers []data.PartitionedLayer, pattern clip.Pattern) (result []data.Paths, err error) {
	// for the first or last layer infill everything
	if layerNr == -1 || layerNr == len(layers) {
		result = append(result, m.pattern.Fill(layerNr, part, nil, m.options.Printer.ExtrusionWidth, m.options.Print.InfillOverlapPercent, internalInfillOverlap))
		return result, nil
	}

	// For the other layers detect the bottom parts by calculating the difference between the current most inner perimeter and the layer below.
	// Also detect the top parts by calculating the difference between the current current most inner perimeter and the layer above
	var toClip []data.LayerPart

	for _, otherPart := range layers[layerNr].LayerParts() {
		toClip = append(toClip, otherPart)
	}

	c := clip.NewClipper()

	toInfill, ok := c.Difference(part, toClip)
	if !ok {
		return nil, errors.New("error while calculating difference for detecting bottom/top parts")
	}

	for _, fill := range toInfill {
		result = append(result, m.pattern.Fill(layerNr, fill, part, m.options.Printer.ExtrusionWidth, m.options.Print.InfillOverlapPercent, internalInfillOverlap))
	}

	return result, nil
}
