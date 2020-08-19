package modifier

import (
	"GoSlice/clip"
	"GoSlice/data"
	"GoSlice/handler"
	"fmt"
)

type brimModifier struct {
	options *data.Options
}

func (m brimModifier) Init(model data.OptimizedModel) {}

// NewBrimModifier generates the brim lines.
func NewBrimModifier(options *data.Options) handler.LayerModifier {
	return &brimModifier{
		options: options,
	}
}

// Brim extracts the attribute "brim" from the layer.
// If it has the wrong type, a error is returned.
// If it doesn't exist, (nil, nil) is returned.
// If it exists, the infill is returned.
func Brim(layer data.PartitionedLayer) (clip.OffsetResult, error) {
	if attr, ok := layer.Attributes()["brim"]; ok {
		parts, ok := attr.(clip.OffsetResult)
		if !ok {
			return nil, fmt.Errorf("the attribute 'brim' has the wrong datatype")
		}

		return parts, nil
	}

	return nil, nil
}

// BrimOuterDimension extracts the attribute "outerBrim" from the layer.
// This attribute describes the exact outer dimension of the brim.
// Can be clipped from other parts to avoid overlapping with the brim.
//
// If it has the wrong type, a error is returned.
// If it doesn't exist, (nil, nil) is returned.
// If it exists, the infill is returned.
func BrimOuterDimension(layer data.PartitionedLayer) ([]data.LayerPart, error) {
	if attr, ok := layer.Attributes()["outerBrim"]; ok {
		parts, ok := attr.([]data.LayerPart)
		if !ok {
			return nil, fmt.Errorf("the attribute 'outerbrim' has the wrong datatype")
		}

		return parts, nil
	}

	return nil, nil
}

func (m brimModifier) Modify(layers []data.PartitionedLayer) error {
	if m.options.Print.BrimSkirt.BrimCount == 0 {
		return nil
	}

	layer := layers[0]

	// Get the perimeters to base the brim on them.
	perimeters, err := Perimeters(layer)
	if err != nil {
		return err
	}
	if perimeters == nil {
		return nil
	}

	// Extract the outer perimeters of all perimeters.
	var allOuterPerimeters []data.LayerPart

	for _, part := range perimeters {
		for _, wall := range part {
			if len(wall) > 0 {
				// wall[0] is the outer perimeter
				allOuterPerimeters = append(allOuterPerimeters, wall[0])
			}
		}
	}

	cl := clip.NewClipper()

	// Get the top level polys e.g. the polygons which are not inside another.
	topLevelPerimeters, _ := cl.TopLevelPolygons(allOuterPerimeters)
	allOuterPerimeters = nil
	for _, p := range topLevelPerimeters {
		allOuterPerimeters = append(allOuterPerimeters, data.NewBasicLayerPart(p, nil))
	}

	if allOuterPerimeters == nil {
		// No need to go further and prevent fail of union.
		return nil
	}

	// Generate the brim.
	brim := cl.InsetLayer(allOuterPerimeters, -m.options.Printer.ExtrusionWidth, m.options.Print.BrimSkirt.BrimCount, m.options.Printer.ExtrusionWidth)

	// Now we need to generate the outer bounds of the brim (e.g. outer brim line + half line width)
	// That is needed for the support, to remove the support at the places where the brim is.
	var outerBrimLines []data.LayerPart
	// For this we first get only the most outer brim lines.
	for _, part := range brim {
		if len(part) == 0 {
			continue
		}
		for _, insetPart := range part[len(part)-1] {
			outerBrimLines = append(outerBrimLines, insetPart)
		}
	}

	// Then the outer brim lines are exset so that the result matches the exact dimension taking into account the extrusion width.
	outerBrim := cl.InsetLayer(outerBrimLines, -m.options.Printer.ExtrusionWidth, 1, m.options.Printer.ExtrusionWidth/2).ToOneDimension()

	newLayer := newExtendedLayer(layers[0])
	if len(brim) > 0 {
		newLayer.attributes["brim"] = append(brim)
	}

	if len(outerBrim) > 0 {
		newLayer.attributes["outerBrim"] = outerBrim
	}

	layers[0] = newLayer

	return nil
}
