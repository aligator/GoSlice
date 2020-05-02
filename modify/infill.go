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
	options         *data.Options
	pattern         clip.Pattern
	internalPattern clip.Pattern
}

func (m *infillModifier) Init(model data.OptimizedModel) {
	m.pattern = clip.NewLinearPattern(model.Min().PointXY(), model.Max().PointXY(), m.options.Printer.ExtrusionWidth)

	// TODO: the calculation of the percentage is currently very basic and may not be correct.
	// It needs improvement.

	if m.options.Print.InfillPercent != 0 {
		mm10 := data.Millimeter(10).ToMicrometer()
		linesPer10mmFor100Percent := mm10 / m.options.Printer.ExtrusionWidth
		linesPerArea10x10ForInfillPercent := float64(linesPer10mmFor100Percent) * float64(m.options.Print.InfillPercent) / 100.0

		lineWidth := data.Micrometer(float64(mm10) / linesPerArea10x10ForInfillPercent)

		m.internalPattern = clip.NewLinearPattern(model.Min().PointXY(), model.Max().PointXY(), lineWidth)
	}
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
	var internalInfill []data.Paths

	// calculate the bottom parts for each inner perimeter part
	for partNr, part := range perimeters {
		// for the last (most inner) inset of each part
		for insertPart, insetPart := range part[len(part)-1] {
			fmt.Println("layerNr " + strconv.Itoa(layerNr) + " partNr " + strconv.Itoa(partNr) + " insertPart " + strconv.Itoa(insertPart))

			infill, bottomInfillParts, err := m.genTopBottomInfill(insetPart, layerNr-1, layers)
			if err != nil {
				return nil, err
			}
			for _, paths := range infill {
				bottomInfill = append(bottomInfill, paths)

			}

			infill, topInfillParts, err := m.genTopBottomInfill(insetPart, layerNr+1, layers)
			if err != nil {
				return nil, err
			}
			for _, paths := range infill {
				bottomInfill = append(bottomInfill, paths)
			}

			// if no infill, just ignore the generation
			if m.internalPattern == nil {
				continue
			}
			// add the parts from top and bottom as ignore-parts, to avoid intersections of top/bottom infill and the internal infill
			// TODO: this extra-clipping seems to have a bad performance
			infill, _, err = m.genInternalInfill(insetPart, layerNr, layers, append(topInfillParts, bottomInfillParts...)...)
			if err != nil {
				return nil, err
			}
			for _, paths := range infill {
				internalInfill = append(internalInfill, paths)
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

// genInfill returns the infill for the top or bottom parts.
// It calculates the difference of the layer with layerNr and the given part.
// Then it fills the result by using the given pattern.
func (m *infillModifier) genTopBottomInfill(part data.LayerPart, layerNr int, layers []data.PartitionedLayer) (result []data.Paths, resultParts []data.LayerPart, err error) {
	// for the first or last layer infill everything
	if layerNr == -1 || layerNr == len(layers) {
		infill, resultPart := m.pattern.Fill(layerNr, part, nil, m.options.Printer.ExtrusionWidth, m.options.Print.InfillOverlapPercent, internalInfillOverlap)

		result = append(result, infill)
		resultParts = append(resultParts, resultPart)
		return result, resultParts, nil
	}

	// For the other layers detect the infill parts by calculating the difference between the current most inner perimeter and the given layer.
	// (based on the given layerNr)
	var toClip []data.LayerPart

	for _, otherPart := range layers[layerNr].LayerParts() {
		toClip = append(toClip, otherPart)
	}

	c := clip.NewClipper()

	toInfill, ok := c.Difference(part, toClip)
	if !ok {
		return nil, nil, errors.New("error while calculating difference for detecting bottom/top parts")
	}

	for _, fill := range toInfill {
		infill, resultPart := m.pattern.Fill(layerNr, fill, part, m.options.Printer.ExtrusionWidth, m.options.Print.InfillOverlapPercent, internalInfillOverlap)
		result = append(result, infill)
		resultParts = append(resultParts, resultPart)
	}

	return result, resultParts, nil
}

func (m *infillModifier) genInternalInfill(part data.LayerPart, layerNr int, layers []data.PartitionedLayer, partsToIgnore ...data.LayerPart) (result []data.Paths, resultParts []data.LayerPart, err error) {
	// for the first or last layer do nothing
	if layerNr == 0 || layerNr == len(layers)-1 {
		return []data.Paths{}, []data.LayerPart{}, nil
	}

	// For the other layers detect the infill parts by calculating the intersection between the current most inner perimeter and the layers above and below.
	var toClip []data.LayerPart

	for _, otherPart := range layers[layerNr-1].LayerParts() {
		toClip = append(toClip, otherPart)
	}

	c := clip.NewClipper()

	toFillWithoutIgnoredParts, ok := c.Difference(part, partsToIgnore)
	if !ok {
		return nil, nil, errors.New("error while clipping the ignored parts from the parts to infill")
	}

	for _, part := range toFillWithoutIgnoredParts {
		toInfill, ok := c.Intersection(part, toClip)
		if !ok {
			return nil, nil, errors.New("error while calculating intersection for detecting internal parts")
		}

		for _, fill := range toInfill {
			infill, resultPart := m.internalPattern.Fill(layerNr, fill, part, m.options.Printer.ExtrusionWidth, m.options.Print.InfillOverlapPercent, 0)
			result = append(result, infill)
			resultParts = append(resultParts, resultPart)
		}
	}

	return result, resultParts, nil
}
