package modifier

import (
	"GoSlice/clip"
	"GoSlice/data"
	"GoSlice/handler"
)

type perimeterModifier struct {
	options *data.Options
}

// NewPerimeterModifier creates a modifier which calculates all perimeters
//
// The perimeters are saved as attribute in the LayerPart.
func NewPerimeterModifier(options *data.Options) handler.LayerModifier {
	return &perimeterModifier{
		options: options,
	}
}

func (m perimeterModifier) Init(model data.OptimizedModel) {}

func (m perimeterModifier) Modify(layerNr int, layers []data.PartitionedLayer) ([]data.PartitionedLayer, error) {
	// generate perimeters
	c := clip.NewClipper()
	insetParts := c.InsetLayer(layers[layerNr].LayerParts(), m.options.Printer.ExtrusionWidth, m.options.Print.InsetCount)

	newLayer := newExtendedLayer(layers[layerNr])
	newLayer.attributes["perimeters"] = insetParts
	layers[layerNr] = newLayer

	return layers, nil
}
