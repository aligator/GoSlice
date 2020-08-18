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

func (f *fakeRenderer) Render(b *gcode.Builder, layerNr int, maxLayer int, layer data.PartitionedLayer, z data.Micrometer, options *data.Options) error {
	f.c.c["render"]++
	test.Assert(f.t, maxLayer >= layerNr, "the number of layers should be more or equal than the current layer number")
	b.AddCommand("number %v", layerNr)
	return nil
}

func TestGCodeGenerator(t *testing.T) {
	rendererCounter := newCounter()

	layers := make([]data.PartitionedLayer, 3)

	generator := gcode.NewGenerator(&data.Options{}, gcode.WithRenderer(&fakeRenderer{t: t, c: rendererCounter}))
	generator.Init(nil)
	result, err := generator.Generate(layers)

	test.Ok(t, err)

	test.Assert(t, rendererCounter.c["init"] == 1, "init should have been called only one time")
	test.Assert(t, rendererCounter.c["render"] == len(layers), "render should have been called %v times (one for each layer)", len(layers))
	test.Equals(t, "number 0\n"+
		"number 1\n"+
		"number 2\n", result)
}
