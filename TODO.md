
# Next

+ need a RegisterOp function that takes a string and a function and stores the result in a function table
+ need an EvalString function that takes a function table, symbol table and input string and either writes a string to stdout, make a new assignment or emits an error message with line number


# Bugs/Ideas

+ make a decision about labels being immutable or not.
    + pros is it simplifies handling significantly, also makes playback easier (makes things functional)
    + cons if you're typing in the repl you're going to want to replace a label's content periodically as you explore and make mistakes
+ add line number ranges when including include files.
+ the glygh :{<: does not seem to work for reading in a file, expanding content and assigning the results to a label (or at least not as expected)
    + do we need anything more than :=<: when you have support for :{{: to double expand content (e.g. when reading in a shorthand file as template)

