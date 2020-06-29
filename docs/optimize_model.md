# 2. Optimize model

This step is mostly ported and refactored from the [first CuraEngine commit](https://github.com/Ultimaker/CuraEngine/commit/80dc349e2014eaa9450086c007118e10bda0b534).

It is responsible for
* find neighbour faces for each face
* find faces with open sides (which have no neighbour) which is bad for 3D printing.
* moving the model into the center

The code for this can be found in the folder [optimizer](../optimizer).