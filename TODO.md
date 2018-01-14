
# Action Items

## Bugs

## Next

## Someday maybe

## Think about ideas

## Completed

+ [x] Need to switch notation from infix to prefix for assignment ops
+ [x] Added EvalSymbol which lets you construct your own Symbol and send it to the VM (like Eval without the parse step)
+ [x] need an EvalString function that takes a function table, symbol table and input string and either writes a string to stdout, make a new assignment or emits an error message with line number
+ [x] need a RegisterOp function that takes a string and a function and stores the result in a function table
+ [x] make a decision about labels being immutable or not.
    + pros is it simplifies handling significantly, also makes playback easier (makes things functional)
    + cons if you're typing in the repl you're going to want to replace a label's content periodically as you explore and make mistakes
+ [x] add line number ranges when including include files.

