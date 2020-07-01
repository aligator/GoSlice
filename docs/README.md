GoSlice
=======
GoSlice is a program used to convert 3D models into GCode for 3D printers.

# Initial design decisions
As I initially ported the first Cura commit to go, some of the design joice's are also just the same as there.
* using Clipper:  
The Clipper lib is not very Go-idiomatic because it is initially ported from Delphy to C++ to Go.
It was the only lib I found which really can do anything I need.  
Also the C++ version of it is used in Cura.
* use clipper only in one package:  
As I am not sure if there will be another, better polygon clipping library and to seperate the lib usage from the GoSlice code
I decided to leave all references to it in the clip package.  
However this has one downside: I have to convert the polygons to the clipper representation each time. For now the performance is ok I would say.
(keeping in mind that I have no concurrency, yet)
* use int for internal calculations:  
This is mainly also because of the use of Clipper. It uses ints internally.
Also the Cura uses it and it was more easy to port it.
I think it may also avoid rounding errors and it may be a bit faster, but I have not verified.
* modularity
I made interfaces for the different steps. So is very easy to add new ones or replace one.
Also they can be mocked easily.
The slicer.go in the root combines them all in one struct and executes them.  
I think this was a good decision as it is clean and easy to extend.
* functional approach:  
I often use the functional approach and avoid pointers and pointer receivers. 
My very first version was using some more pointers (more or less the same as in the Cura commit) 
but then I changed it and was surprised that it ran even faster than the c++ version (with exactly the same features, back then)

# How it works
(This describes only the functionality of GoSlice. There may be other ways to slice 3D models.)

For in depth documentation also see the comments directly in the code (and/or using the `go doc` tool.)

Basically the slicing process can be split into several steps:
1. [load STL file](load_stl.md)
1. [optimize model](optimize_model.md)
1. [slice model](slice_model.md)
1. [modify slices](modify_slices.md)
1. [generate GCode](generate_gcode.md)

