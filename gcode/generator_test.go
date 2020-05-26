package gcode_test

import (
	"GoSlice/data"
	"GoSlice/gcode"
	"GoSlice/util/test"
	"testing"
)

type counter struct {
	c map[string]int
}

func newCounter() *counter {
	return &counter{
		c: map[string]int{},
	}
}

type fakeRenderer struct {
	t testing.TB
	c *counter
}

func (f *fakeRenderer) Init(model data.OptimizedModel) {
	f.c.c["init"]++
}

func (f *fakeRenderer) Render(b *gcode.Builder, layerNr int, layers []data.PartitionedLayer, z data.Micrometer, options *data.Options) {
	f.c.c["render"]++
	test.Assert(f.t, len(layers) > layerNr, "the number of layers should be more than the current layer number")
	b.AddCommand("number %v", layerNr)
}

type fakePartitionedLayer struct {
	t testing.TB
}

func (f *fakePartitionedLayer) LayerParts() []data.LayerPart {
	panic("implement me")
}

func (f *fakePartitionedLayer) Attributes() map[string]interface{} {
	panic("implement me")
}

func (f *fakePartitionedLayer) Bounds() (data.MicroPoint, data.MicroPoint) {
	panic("implement me")
}

func TestGCodeGenerator(t *testing.T) {
	rendererCounter := newCounter()

	layers := []data.PartitionedLayer{
		&fakePartitionedLayer{t: t},
		&fakePartitionedLayer{t: t},
		&fakePartitionedLayer{t: t},
	}

	generator := gcode.NewGenerator(&data.Options{}, gcode.WithRenderer(&fakeRenderer{t: t, c: rendererCounter}))
	generator.Init(nil)
	result := generator.Generate(layers)

	test.Assert(t, rendererCounter.c["init"] == 1, "init should have been called only one time")
	test.Assert(t, rendererCounter.c["render"] == len(layers), "render should have been called %v times (one for each layer)", len(layers))
	test.Equals(t, "number 0\n"+
		"number 1\n"+
		"number 2\n", result)
}
