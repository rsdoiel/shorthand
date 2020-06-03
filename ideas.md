---
title: "Shorthand - ideas and background"
---
    
# Shorthand - ideas and background

# Shorthand prototype

## Basic ideas

+ Shorthand is a incomplete-language (lacks loops, recusion, control statements, and lambdas) 
+ Shorthand is based on the general idea of assigning a value to a label
+ Shorthand only has two types of statements - assignments and expansions
    + Assignments exist on a single line in the form of label, operator and value
        + assignments (attaching a value or expansion to a label)
    + If the line is not an assignment it is an expansion written to stdout
        + expansion (is a string with any labels in string resolved to their previously assigned value)
+ A assignment statement is a triple (three cells of a tuple)
    + op (the assignment operation to be performed)
    + label (the receiver of the results)
    + value (the thing being operated on)
+ In memory it is represented as a five-tuple
    + op
    + label
    + value
    + expansion
    + line number
+ An additional map contains an index into the 5-tuple pointing at the most recent assignment associated with that label
    + this allows labels to mutate but the prior version of the label still exists in memory and can be written out to a file
+ An expansion is a string with zero or more labels expanded
+ By defaut shorthand works on standard in and standard out
    + but the VM in shorthand package contains an Eval for embedding in other projects
+ Operators (functions whos results are assigned to labels)
    + begin and end with colons and cannot contain spaces (e.g. ":set:", :"import-text:")
    + internal inputs
        + assign a string to a label
        + assign an expansion result to a label
        + assign an expanded result to another expansion
        + assign a markdown processed string to a label
        + assign a expanded markdown processed string to a label
    + external inputs
        + assign the contents of a file to a label
        + assign the expanded contents of a file to a label
        + assign the output of a Bash command to a label
        + expand the command sent to Bash, assign the output to a label
        + assign a markdown processed file to a label
        + assign an expanded markdown processed file to a label 
    + special form
        + read a file for assignments nothing is assigned to the label
    + outputs
        + write to a file the expansion for a label
        + write to a file all the expansions for all labels (order not guaranteed)
        + write to a file an assignment statement for a label
        + write to a file all assignment statements (order by parse sequence)
+ two factors prevent shorthand from being a full "language"
    + it does not support a control statements (you might be able to synthesize this via Bash)
    + it does not support creation of new operators (though expansion of labels passed to Bash is similar)
+ it is not "functional" in the sense that labels are mutable (can take a new value over time)
    + but the prior version to the mutation still exists in the queue of assignments and can be played back by dumping to a file.
+ there is only one type - a string
+ there are no built in data structures
+ there is no comments though you could treat the output to standard out as comments
+ shorthand could be expanded to use another interpretive environment but Bash is very convientent
+ text is UTF-8
+ operators begin and end with a colon
+ assignments are a single line (terminated with a \n)


## someday, maybe

+ Write a formal definition of Shorthand to encourage other implementations
+ Allow additional operators to be defined in other easily go embedded languages
    + JavaScript
    + Lua
    + PHP
    + Scheme (LispEx, GLisp)
+ Organize the registration of operators in the code for easy expansion
    + symbol (e.g. :=: )
    + Name
    + Function pointer
+ lamda operator
    + if the "label" if the lamba is in the form of a an operator it receives the value expression and then can be used as an Operator itself
+ Additional data types
    + stack of strings
    + a queue of strings
    + a map of strings
    + numbers in the lisp sense of numbers (e.g. integers, ratios, decimal)
+ Additional built-in operators
    + math
        + add
        + substract (alias to add with second number sign flip)
        + multiply
            + invert sign of number (positive/negative)
        + divide
        + modulo
        + inc by N (positive or negative)
+ io
    + support for URI/URL for file operation
    + support CURL like operations
+ Add some sort of random text or Lorem Ipsum function for testing and filling in templates
    + [Random Text](https://github.com/AhmedZaleh/margopher?utm_source=golangweekly&utm_medium=email) with MarGopher
    + Lorem Ipsum in Golang: [Golorem](https://github.com/drhodes/golorem), other version are probably out there too

