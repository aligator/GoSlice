# 3. Slice model

This step is also mostly ported and refactored from the [first CuraEngine commit](https://github.com/Ultimaker/CuraEngine/commit/80dc349e2014eaa9450086c007118e10bda0b534).

It is responsible for the actual slicing of the model into slices at different heights.

To do this for each height the faces are sliced into segments (lines).  
This is done in [slicer.go](../slicer/slicer.go) which calls SliceFace in [segment.go](../slicer/segment.go).

After that the loose segments have to be connected together to several polygons and holes of polygons:
1. connect segments which have the exact same points
1. sometimes there are small spaces between the points. 
If it the space is small enough it is just skipped and the segments are also connected together.
1. also it is possible that some poligons are _nearly_ closed or two half polygons _nearly_ can be connected to one.
Then they are just closed.
1. at the end all not yet closed or too small polygons are just removed.

All this happens in the method makePolygons in [layer.go](../slicer/layer.go).

As last step the LayerParts are generated for each layer out of the polygons.
This happens using the method GenerateLayerParts in [clip.go](../clip/clip.go).  
It basically just calls the clipper lib which groups the polygons together and calculates which one are polygons and which one are holes.
(Note: polygons go counter-clockwise, holes go clockwise.)