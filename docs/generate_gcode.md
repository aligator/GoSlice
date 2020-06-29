# 5. Generate GCode

This step basically just uses the attributes added by the modifiers to generate the GCode.  
For example, it generates the infill pattern for all parts defined in the attribute "infill".

In many cases each modifier has a counterpart renderer.

The actual GCode is built using a GCode builder.