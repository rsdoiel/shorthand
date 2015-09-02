
# Shorthand prototype

## Basic ideas

+ Shorthand is a semi-language (lacks loops/recusion, control statements, and lambdas) based on assigning a value to a lable
    + it supports two types of expressions
        + assignments (attaching a value or expansion to a label)
        + expansion (what the resolved label represents)
+ A assignment statement is a triple (three cells of a tuple)
    + label (the receiver of the results)
    + op (the assignment operation to be performed)
    + value (the thing being operated on)
+ In memory it is represented as a five-tuple
    + label
    + op
    + value
    + expansion
    + line number
+ An additional map contains an index in to the 5-tuple pointing at the most recent assignment associated with that label
+ An expansion is a string with zero or more labels expanded
+ By defaut shorthand works on standard in and standard out
+ Operators
    + internal inputs
        + assign a string to a label
        + assign an expansion result to a label
        + assign an expansion or another expansion
        + assign a markdown processed string to a label
        + assign a expansion markdown processed string to a label
    + external inputs
        + assign the contents of a file to a label
        + assign the expanded contents fo a file to a label
        + assign the output of a Bash command to a label
        + expands the command sent to Bash, assign the output to a label
        + assign a markdown processed file to a label
        + assign an expansion markdown processed file to a label 
    + special form
        + read a file for assignments nothing is assigned to the label
    + outputs
        + write to a file the expansion for a label
        + write to a file all the expansions for all labels (order not guaranteed)
        + write to a file an assignment state for a label
        + write to a file all assignment statements (order by parse sequence)
+ two factors prevent shorthand from being a "language"
    + it does not support a control statement directly (bash could be used as a the control statement)
    + it does not support creation of new operators (though you can synthesize this be shelling out to bash)
+ it is not "functional" in the sense that labels are mutable
+ there is only one type - a string
+ there is no comments though you could treat the output to stnout as comments
+ there is no interactive though you could synthesize this from Bash
+ shorthand could be expanded to use another interpretive environment but Bash is very convientent
+ A "compiler" could be built generating a go program that then could be compiled and run



