package modifier

import (
	"GoSlice/clip"
	"GoSlice/data"
	"GoSlice/handler"
	"errors"
	"fmt"
	"math"
)

type supportModifier struct {
	options *data.Options
}

func (m supportModifier) Init(model data.OptimizedModel) {}

// NewSupportModifier calculates the areas which need support.
func NewSupportModifier(options *data.Options) handler.LayerModifier {
	return &supportModifier{
		options: options,
	}
}

func (m supportModifier) Modify(layerNr int, layers []data.PartitionedLayer) error {
	if !m.options.Print.SupportEnabled || layerNr == 0 {
		return nil
	}

	// ### = a layer
	//
	// ############
	// ############
	// ### ___d____
	// ### |     /
	// ### |    /
	// ### h   /
	// ### |  /
	// ### |θ/
	// ### |/
	//
	// d = h * tan θ
	// https://tams.informatik.uni-hamburg.de/publications/2018/MSc_Daniel_Ahlers.pdf
	// 4.1.5  Support Generation
	//
	// "To get the actual areas where the support is later generated,
	//  the previous layer is offset by the calculated d and then subtracted from the current layer.
	//  All areas that remain have a higher angle than the threshold and need to be supported."

	// calculate distance (d):
	distance := float64(m.options.Print.LayerThickness) * math.Tan(data.ToRadians(float64(m.options.Print.SupportThresholdAngle)))

	// offset previous layer by d
	cl := clip.NewClipper()
	var offsetLayer []data.LayerPart

	fmt.Println(data.Micrometer(-math.Round(distance)))
	for _, part := range cl.InsetLayer(layers[layerNr-1].LayerParts(), data.Micrometer(-math.Round(distance)), 1) {
		for _, wall := range part {
			offsetLayer = append(offsetLayer, wall...)
		}
	}

	// subtract from current layer
	support, ok := cl.Difference(layers[layerNr].LayerParts(), offsetLayer)
	if !ok {
		return errors.New("could not calculate the support parts")
	}

	newLayer := newExtendedLayer(layers[layerNr])
	if len(support) > 0 {
		fmt.Println("len", len(support), "layer", layerNr)
		newLayer.attributes["support"] = support
	}
	layers[layerNr] = newLayer

	return nil
}
