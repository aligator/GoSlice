[![CircleCI](https://circleci.com/gh/aligator/GoSlice.svg?style=svg)](https://circleci.com/gh/aligator/GoSlice)

<img width="200" alt="sliced Gopher logo" src="https://raw.githubusercontent.com/aligator/GoSlice/master/logo.png">

# GoSlice

This is a very experimental slicer for 3d printing.

The initial work of GoSlice is based on the first CuraEngine commits.
As I had no clue where to start, I chose to port the initial Cura commit to Go.
The code of this early Cura version already provides a very simple and working slicer and the code of it is easy to read.
https://github.com/Ultimaker/CuraEngine/tree/80dc349e2014eaa9450086c007118e10bda0b534

Most of the work after "first gcode result" is done from scratch.

## Run
Minimal usage:
```
go run . --file /path/to/stl/file.stl
```
If you get an error while reading the file, it may help to import it in Prusa Slicer and export the stl from there. (e.g. 3DBenchy has this problem)
I didn't look into this error yet, but it may be a problem of the stl lib I use.


To get help for all possible flags take a look at /data/option.go or just run:
```
go run . --help
```

## Contribution
You are welcome to help.  
[Just look for open issues](https://github.com/aligator/GoSlice/issues) and pick one, create new issues or create new pull requests.

## Credits
* CuraEngine for the great first commit, which was a very good starting point for research.
* https://www.thingiverse.com/thing:3413597 for the great Gopher model used as logo. (Original Gopher designed by [Renee French CC BY 3.0](http://reneefrench.blogspot.com/))
* Go for the great language.
* All libs GoSlice uses. (just take a look at go.mod)
