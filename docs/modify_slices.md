# 4. Modify slices

This step is heavily customized since I cloned the first cura commit.  
In this step we calculate the different types of lines needed for a 3D Print.
It can be seen as _preprocessing_ of the GCode generation.

It is built using a modular set of structs implementing the Modifier interface.  
It is called _Modifier_ because each of them modifies the layers and adds different attributes to them
which can be used later by the GCode generator.

These modifiers make heavy usage of the clipper lib.

Currently, it generated these types:
* perimeter
* infill
* top / bottom layers