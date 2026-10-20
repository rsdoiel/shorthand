---
title: shorthand syntax
---

shorthand syntax
================

*shorthand* is a simple label expander. A line is either a definition, which
assigns a value to a label, or text. The text is written with every defined
label replaced by its value, and the definition lines are not written.

A line is split on spaces into an operator, a label and the rest, the source:

~~~text
:set: {{name}} Freda
~~~

A label is any text that does not contain an operator. Curly brackets, dollar
signs and at signs are common; what matters is that it does not appear in your
text by accident. The spaces around the operator and the label are required.
Using an underscore as the label means the label is ignored, which is how an
operator is used only for what it does, as with `:import-shorthand:`.

operator | meaning | example
---|---|---
`:set:` | assign a string | `:set: {{name}} Freda`
`:import-text:` | assign the contents of a file | `:import-text: {{content}} myfile.txt`
`:import-shorthand:` | read definitions from a file | `:import-shorthand: _ myfile.shorthand`
`:expand:` | assign an expansion | `:expand: {{reportTitle}} Report: @title for @date`
`:expand-expansion:` | assign an expansion of an expansion | `:expand-expansion: {{reportHeading}} @reportTitle`
`:import:` | include a file, processing its shorthand | `:import: {{nav}} mynav.shorthand`
`:bash:` | assign the output of a Bash command | `:bash: {{date}} date +%Y-%m-%d`
`:expand-and-bash:` | expand, then assign the output of a Bash command | `:expand-and-bash: {{entry}} cat header.txt @filename footer.txt`
`:export:` | write a label's value to a file | `:export: {{content}} content.txt`
`:export-all:` | write every label's value to a file | `:export-all: _ contents.txt`
`:export-shorthand:` | write a label's definition to a file | `:export-shorthand: {{content}} content.shorthand`
`:export-all-shorthand:` | write every definition to a file | `:export-all-shorthand: _ contents.shorthand`

`:export-all:` and `:export-all-shorthand:` write in no particular order.

# How labels are replaced

Labels are replaced in one pass over the text, and the answer never depends on
the order they were defined in.

- Where two labels match at the same place the longer one wins. With `@a` and
  `@ab` defined, `@ab` is the value of `@ab`, and `@abc` is the value of `@ab`
  followed by `c`.
- A value is not expanded again. If `@b` was set to `@a`, then `@b` expands to
  the text `@a`. Use `:expand:` to expand the labels in a source when it is
  assigned, or `:expand-expansion:` to expand them twice, when you want a label
  in a value expanded too.
- The labels `_` and an empty label are never replaced, so using `_` for an
  operator that only does something, like `:import-shorthand:`, leaves the
  underscores in the rest of your text alone.

`:exit:` ends the input, and so does `:quit:`. A line that contains either one
ends it, wherever the word is in the line. When *shorthand* is run with no file
at a terminal, `:help:` lists the operators.

A definition applies to the lines after it, and to the files after it on the
same command line. Defining a label again replaces its value from then on.

# Example

In this example the contents of a file are assigned to `{{PREAMBLE}}` and the
time 3:30 to `{{NOW}}`:

~~~text
:import-text: {{PREAMBLE}} /home/me/preamble.txt
:set: {{NOW}} 3:30
At {{NOW}} I will be reading the {{PREAMBLE}} until everyone falls asleep.
~~~

If *preamble.txt* held `"Hello World"`, with no newline after it, the output is

~~~text
At 3:30 I will be reading the "Hello World" until everyone falls asleep.
~~~

The lines with the definitions are not in the output, and no line ends are
added to the substituted labels.

# A timesheet

*shorthand* is used to put the present time into a timesheet for the
[stngo](https://github.com/rsdoiel/stngo) tools. The timesheet starts with

~~~text
:bash: @now echo -n $(date +%I:%M)
~~~

and an entry that is still going on is written `11:45 - @now; Acme; Widgets`.
`shorthand Time_Sheet.txt | stnparse` then reads it as ending at the present
time.

# Calling it from Go

The package is `github.com/rsdoiel/shorthand`. `New` returns a
`VirtualMachine`, `Eval` and `Apply` expand text, `RegisterOp` adds an
operator, and `Out` and `Err` say where `Run` writes.

~~~go
vm := shorthand.New()
out, err := vm.Apply([]byte(":set: @who Freda\nHello @who\n"))
~~~
