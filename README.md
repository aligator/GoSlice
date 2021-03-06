[![CircleCI](https://circleci.com/gh/aligator/GoSlice.svg?style=svg)](https://circleci.com/gh/aligator/GoSlice)

<img width="200" alt="sliced Gopher logo" src="https://raw.githubusercontent.com/aligator/GoSlice/master/logo.png">

# GoSlice

This is a very experimental slicer for 3d printing. It is currently in a very early stage but it can already slice models:

Supported features:
* perimeters
* simple linear infill
* rotated infill
* top / bottom layer
* simple temperature control
* simple speed control
* simple retraction on crossing perimeters
* several options to customize slicing output
* simple support generation
* brim and skirt

Example:  
<img width="200" alt="sliced Gopher logo" src="https://raw.githubusercontent.com/aligator/GoSlice/master/docs/GoSlice-print.png">

## Try it out - for users
Download latest release matching your platform from here:
https://github.com/aligator/GoSlice/releases

Unpack the executable and run it in the commandline.  
linux / mac:  
```
./goslice /path/to/stl/file.stl
```

windows:  
```
goslice.exe /path/to/stl/file.stl` 
```

If you need the usage of all possible flags, run it with the `--help` flag:
```
./goslice --help
```

## Try it out - for developers
Just running GoSlice:
```
go run ./cmd/goslice /path/to/stl/file.stl
```
To get help for all possible flags take a look at /data/option.go or just run:
```
go run ./cmd/goslice --help
```

Building GoSlice:
Ideally you should have make installed:
```
make
```
The resulting binary will be in the `.target` folder.

If you do not have make, you can still run the build command manually, but it is not recommended:
```
go build -ldflags "-X=main.Version=$(git describe --tags) -X=main.Build=$(git rev-parse --short HEAD)" -o .target ./cmd/goslice
```
## How does it work
[see here](docs/README.md)

## Contribution
You are welcome to help.  
[Just look for open issues](https://github.com/aligator/GoSlice/issues) and pick one, create new issues or create new pull requests.

For debugging of the GCode I suggest you to use Cura to open the resulting GCode.
Cura can open it without any problem and I try to add the markings into the GCode which Cura understands (e.g. mark what is infill, perimeter, etc.).

## Credits
* CuraEngine for the great first commit, which was a very good starting point for research.
* https://www.thingiverse.com/thing:3413597 for the great Gopher model used as logo. (Original Gopher designed by [Renee French CC BY 3.0](http://reneefrench.blogspot.com/))
* Go for the great language.
* All libs GoSlice uses. (just take a look at go.mod)
