# 1. Load STL file

This step is currently done by using an external lib (https://github.com/hschendel/stl).
I wanted to focus on the slicing process itself.
That's why I cannot write much about it.

__But I can tell about some problems you can face with stl files:__

Basically a STL file exists in two flavours:
* ASCII
* Binary

Both of them have a similar structure, but the binary version needs less space.

[Wikipedia](https://en.wikipedia.org/wiki/STL_(file_format))

However it is a bit tricky to check if it is ASCII or Binary.  
The specification defines that the Binary files never start with "solid ".
But sadly there are some files out there which do not respect this rule.
(even the well known [3DBenchy.stl](https://www.thingiverse.com/thing:763622) has this problem.)

[That's why I created a PR for the lib I use](https://github.com/hschendel/stl/pull/3) and use my fixed version till it gets merged.
