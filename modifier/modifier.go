package modifier

import (
	"GoSlice/data"
	"sync"
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

type enumeratedPartitionedLayer struct {
	layer   data.PartitionedLayer
	layerNr int
}

func getModifierInputChannel(layers []data.PartitionedLayer) <-chan enumeratedPartitionedLayer {
	inputChannel := make(chan enumeratedPartitionedLayer, len(layers)/2)
	go func() {
		defer close(inputChannel)

		for layerNr, layer := range layers {
			inputChannel <- enumeratedPartitionedLayer{
				layer:   layer,
				layerNr: layerNr,
			}
		}
	}()

	return inputChannel
}

func merge(outputsChan []<-chan enumeratedPartitionedLayer, errorsChan []<-chan error) (<-chan enumeratedPartitionedLayer, <-chan error) {
	var wg sync.WaitGroup
	var errWg sync.WaitGroup

	merged := make(chan enumeratedPartitionedLayer, len(outputsChan)/2)
	mergedError := make(chan error)

	// increase counter to number of channels `len(outputsChan)`
	// as we will spawn number of goroutines equal to number of channels received to merge
	wg.Add(len(outputsChan))
	errWg.Add(len(errorsChan))

	output := func(layerCh <-chan enumeratedPartitionedLayer) {
		// run until channel closes
		for layer := range layerCh {
			merged <- layer
		}

		wg.Done()
	}

	outputErr := func(errCh <-chan error) {
		// run until channel closes
		for err := range errCh {
			mergedError <- err
		}

		errWg.Done()
	}

	// run above `output` function as groutines, `n` number of times
	// where n is equal to number of channels received as argument the function
	// here we are using `for range` loop on `outputsChan` hence no need to manually tell `n`
	for _, optChan := range outputsChan {
		go output(optChan)
	}
	for _, errChan := range errorsChan {
		go outputErr(errChan)
	}

	// run goroutine to close merged channel once done
	go func() {
		// wait until WaitGroup finishesh
		wg.Wait()
		errWg.Wait()
		close(merged)
		close(mergedError)
	}()

	return merged, mergedError
}

func modifyConcurrently(layers []data.PartitionedLayer, modifier func(layerCh <-chan enumeratedPartitionedLayer, outputCh chan<- enumeratedPartitionedLayer, errCh chan<- error)) error {
	inputCh := getModifierInputChannel(layers)
	results := make([]<-chan enumeratedPartitionedLayer, 0)
	errors := make([]<-chan error, 0)

	for i := 0; i < len(layers); i++ {
		outputCh := make(chan enumeratedPartitionedLayer, len(layers)/2)
		errCh := make(chan error)

		go func() {
			defer close(outputCh)
			defer close(errCh)
			modifier(inputCh, outputCh, errCh)
		}()

		results = append(results, outputCh)
		errors = append(errors, errCh)
	}

	mergedLayerCh, mergedErrCh := merge(results, errors)

	for {
		select {
		case err := <-mergedErrCh:
			return err
		case layer, ok := <-mergedLayerCh:
			if !ok {
				break
			}
			layers[layer.layerNr] = layer.layer
		default:
		}
	}

	return nil
}
