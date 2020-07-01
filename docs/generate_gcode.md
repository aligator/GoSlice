# 5. Generate GCode

This step basically just uses the attributes added by the modifiers to generate the GCode.  
For example, it generates the infill pattern for all parts defined in the attribute "infill".

In many cases each modifier has a counterpart renderer.

The actual GCode is built using a GCode builder.

There is one problem I had:
the generated gcode produced many small steps in some places which made the printer slow down and hang many times.  
At first I searched a bit in the (current) cura sources and found a [Simplify-method](../data/layer.go).
I ported it over without any idea what it does.

But it didn't do exactly what I wanted. It somehow smoothes the model a bit and therefore reduces the slice time a bit,
but not enough to avoid the problem.  
So after some research I found that the PrusaSlicer uses the [DouglasPeucker](https://en.wikipedia.org/wiki/Ramer%E2%80%93Douglas%E2%80%93Peucker_algorithm) right before generating the GCode.  
So I [implemented it myself](../data/2d.go) (it's a relatively easy recursive algorithm) and __yes__, now the printer prints everything smoothly.