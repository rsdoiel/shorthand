
# Bugs/Ideas

+ make a decision about labels being immutable or not.
    + pros is it simplifies handling significantly, also makes playback easier (makes things functional)
    + cons if you're typing in the repl you're going to want to replace a label's content periodically as you explore and make mistakes
+ add line number ranges when including include files.
+ the glygh :{<: does not seem to work for reading in a file, expanding content and assigning the results to a label (or at least not as expected)
    + do we need anything more than :=<: when you have support for :{{: to double expand content (e.g. when reading in a shorthand file as template)

