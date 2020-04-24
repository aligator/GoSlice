# GoSlice

This is a very experimental slicer for 3d printing.

The initial work of GoSlice is based on the first CuraEngine commits.
As I had no clue where to start, I chose to port the initial Cura commit to Go.
The code of this early Cura version already provides a very simple and working slicer and the code of it is easy to read.
https://github.com/Ultimaker/CuraEngine/tree/80dc349e2014eaa9450086c007118e10bda0b534

Most of the work after "first gcode result" is done from scratch.

## Run
go run /path/to/stl/file.stl

## ToDo
* ~~read stl~~ (initially done by using external lib github.com/hschendel/stl)
* ~~implement optimisation as in first Cura Commit~~
* ~~first gcode result~~ YAY!!
* ~~refactor and Go-ify the code~~ (done, for now...)
* ~~perimeters, with configurable outer perimeter speed~~
* ~~bottom layer~~ (only one layer and not perfect but it works)
* top layer
* simple infill
* options as parameters / config file (using cobra / viper)
* add function / interface / struct documentations
* add tests
* lots of other things...