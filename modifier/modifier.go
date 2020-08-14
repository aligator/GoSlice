package modifier

import (
	"GoSlice/data"
	"fmt"
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

// PartsAttribute extracts the given attribute from the layer.
// It supports only []data.LayerPart as type.
// If it has the wrong type, a error is returned.
// If it doesn't exist, (nil, nil) is returned.
// If it exists, the infill is returned.
func PartsAttribute(layer data.PartitionedLayer, typ string) ([]data.LayerPart, error) {
	if attr, ok := layer.Attributes()[typ]; ok {
		parts, ok := attr.([]data.LayerPart)
		if !ok {
			return nil, fmt.Errorf("the attribute %s has the wrong datatype", typ)
		}

		return parts, nil
	}

	return nil, nil
}
