// This file provides modifiers needed to generate support.
// It contains one supportDetectorModifier and an supportGenerationModifier which is meant to run after the detector,
// so that it can use the information of all layers at once.

package modifier

import (
	"GoSlice/clip"
	"GoSlice/data"
	"GoSlice/handler"
	"errors"
	"fmt"
	"math"
)

type supportDetectorModifier struct {
	options *data.Options
}

func (m supportDetectorModifier) Init(model data.OptimizedModel) {}

// NewSupportDetectorModifier calculates the areas which need support.
// It saves them as the attribute "support" as []data.LayerPart.
// It is meant as a preprocessing modifier.
// Another modifier can use this information to generate the actual support.
func NewSupportDetectorModifier(options *data.Options) handler.LayerModifier {
	return &supportDetectorModifier{
		options: options,
	}
}

func (m supportDetectorModifier) Modify(layers []data.PartitionedLayer) error {
	for layerNr := range layers {
		if !m.options.Print.Support.Enabled || layerNr == len(layers)-1 {
			return nil
		}

		// ### = a layer
		//
		// ############
		// ############
		// ### ___d____  |
		// ### |     /   |
		// ### |    /    |
		// ### h   /     | h = 1 layer height
		// ### |  /      |
		// ### |θ/       |
		// ### |/        |
		//
		// d = h * tan θ
		// https://tams.informatik.uni-hamburg.de/publications/2018/MSc_Daniel_Ahlers.pdf
		// 4.1.5  Support Generation
		//
		// "To get the actual areas where the support is later generated,
		//  the previous layer is offset by the calculated d and then subtracted from the current layer.
		//  All areas that remain have a higher angle than the threshold and need to be supported."

		// calculate distance (d):
		distance := float64(m.options.Print.LayerThickness) * math.Tan(data.ToRadians(float64(m.options.Print.Support.ThresholdAngle)))

		// offset previous layer by d
		cl := clip.NewClipper()
		var offsetLayer []data.LayerPart

		for _, part := range cl.InsetLayer(layers[layerNr].LayerParts(), data.Micrometer(-math.Round(distance)), 1) {
			for _, wall := range part {
				offsetLayer = append(offsetLayer, wall...)
			}
		}

		// subtract from current layer
		support, ok := cl.Difference(layers[layerNr+1].LayerParts(), offsetLayer)
		if !ok {
			return errors.New("could not calculate the support parts")
		}

		// Save the result at the layer below.
		newLayer := newExtendedLayer(layers[layerNr+1])
		if len(support) > 0 {
			newLayer.attributes["support"] = support
		}
		layers[layerNr+1] = newLayer
	}

	return nil
}

type supportGeneratorModifier struct {
	options *data.Options
}

func (m supportGeneratorModifier) Init(model data.OptimizedModel) {}

// NewSupportGeneratorModifier generates the actual areas for the support out of the areas which need support.
// It grows these areas down till the first layer or till it touches the model.
func NewSupportGeneratorModifier(options *data.Options) handler.LayerModifier {
	return &supportGeneratorModifier{
		options: options,
	}
}

func (m supportGeneratorModifier) Modify(layers []data.PartitionedLayer) error {
	var lastSupport []data.LayerPart

	// for each layer starting at the top layer
	for layerNr := len(layers) - 1; layerNr >= 0; layerNr-- {
		if !m.options.Print.Support.Enabled || layerNr == 0 {
			return nil
		}

		// load support for the current layer (or use the result from the last round to avoid loading it again??)
		currentSupport := lastSupport
		if layerNr == len(layers)-1 {
			var err error
			currentSupport, err = PartsAttribute(layers[layerNr], "support")
			if err != nil {
				return err
			}
		}

		// load support needed for the layer below
		belowSupport, err := PartsAttribute(layers[layerNr-1], "support")
		if err != nil {
			return err
		}

		if len(currentSupport) == 0 && len(belowSupport) == 0 {
			continue
		}
		// union them
		cl := clip.NewClipper()
		result, ok := cl.Union(currentSupport, belowSupport)
		if !ok {
			return errors.New(fmt.Sprintf("could not union the supports for layer %d to generate support", layerNr))
		}

		// subtract the (exset) model from the result
		// - exset below
		//exset := cl.InsetLayer(layers[layerNr-1].LayerParts(), data.Millimeter(0.5).ToMicrometer(), -1)
		// - subtract
		actualSupport, ok := cl.Difference(result, layers[layerNr-1].LayerParts())
		if !ok {
			return errors.New(fmt.Sprintf("could not subtract the model from the supports for layer %d", layerNr))
		}

		// save the support as actual support to render for the layer below
		newLayer := newExtendedLayer(layers[layerNr-1])
		if len(actualSupport) > 0 {
			newLayer.attributes["support"] = actualSupport
		} else {
			// remove maybe existing support from the detection modifier
			newLayer.attributes["support"] = []data.LayerPart{}
		}
		layers[layerNr-1] = newLayer
	}
	return nil
}
