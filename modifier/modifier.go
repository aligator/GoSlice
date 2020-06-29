package modifier

import (
	"GoSlice/data"
	"errors"
)

// extendedLayer is a partitioned layer which supports types
type extendedLayer struct {
	data.PartitionedLayer
	typ        string
	attributes map[string]interface{}
}

// newExtendedLayer returns a new PartitionedLayer
// which supports a type and attributes.
// These attributes can be used to add additional parts
// or any other additional data.
func newExtendedLayer(layer data.PartitionedLayer, typ ...string) extendedLayer {
	attributes := layer.Attributes()
	if attributes == nil {
		attributes = map[string]interface{}{}
	}

	newType := ""
	if len(typ) > 0 {
		newType = typ[0]
	}

	return extendedLayer{
		PartitionedLayer: layer,
		attributes:       attributes,
		typ:              newType,
	}
}

func (l extendedLayer) Attributes() map[string]interface{} {
	return l.attributes
}

// extendedLayerPart is a partitioned layer which supports types
type extendedLayerPart struct {
	data.LayerPart
	typ        string
	attributes map[string]interface{}
}

// newExtendedLayerPart returns a new simple PartitionedLayer which just contains several LayerParts.
func newExtendedLayerPart(layerPart data.LayerPart, typ ...string) extendedLayerPart {
	attributes := layerPart.Attributes()
	if attributes == nil {
		attributes = map[string]interface{}{}
	}

	newType := ""
	if len(typ) > 0 {
		newType = typ[0]
	}

	return extendedLayerPart{
		LayerPart:  layerPart,
		attributes: attributes,
		typ:        newType,
	}
}

func (l extendedLayerPart) Attributes() map[string]interface{} {
	return l.attributes
}

// PartsAttribute extracts the given attribute from the layer.
// It supports only []data.LayerPart as type.
// If it has the wrong type, a error is returned.
// If it doesn't exist, (nil, nil) is returned.
// If it exists, the infill is returned.
func PartsAttribute(layer data.PartitionedLayer, typ string) ([]data.LayerPart, error) {
	if attr, ok := layer.Attributes()[typ]; ok {
		parts, ok := attr.([]data.LayerPart)
		if !ok {
			return nil, errors.New("the attribute " + typ + " has the wrong datatype")
		}

		return parts, nil
	}

	return nil, nil
}
