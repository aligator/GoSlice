package modify

import (
	"GoSlice/data"
)

type typedLayer struct {
	data.PartitionedLayer
	typ        string
	attributes map[string]interface{}
}

// returns a new simple PartitionedLayer which just contains several LayerParts.
func newTypedLayer(layer data.PartitionedLayer, typ ...string) typedLayer {
	attributes := layer.Attributes()
	if attributes == nil {
		attributes = map[string]interface{}{}
	}

	newType := ""
	if len(typ) > 0 {
		newType = typ[0]
	}

	return typedLayer{
		PartitionedLayer: layer,
		attributes:       attributes,
		typ:              newType,
	}
}

func (l typedLayer) Type() string {
	if l.typ == "" {
		return l.PartitionedLayer.Type()
	}
	return l.typ
}

func (l typedLayer) Attributes() map[string]interface{} {
	return l.attributes
}

type typedLayerPart struct {
	data.LayerPart
	typ        string
	attributes map[string]interface{}
}

func newTypedLayerPart(layerPart data.LayerPart, typ ...string) typedLayerPart {
	attributes := layerPart.Attributes()
	if attributes == nil {
		attributes = map[string]interface{}{}
	}

	newType := ""
	if len(typ) > 0 {
		newType = typ[0]
	}

	return typedLayerPart{
		LayerPart:  layerPart,
		attributes: attributes,
		typ:        newType,
	}
}

func (l typedLayerPart) Type() string {
	if l.typ == "" {
		return l.LayerPart.Type()
	}
	return l.typ
}

func (l typedLayerPart) Attributes() map[string]interface{} {
	return l.attributes
}
