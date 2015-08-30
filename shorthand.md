
# shorthand

Shorthand is a simple text label expansion utility. It is based on a simple key value substitution.  It supports this following types of definitions

+ assigning a string to a LABEL
+ assigning the contents of a file to a LABEL
+ assigning the output of a Bash shell expression to a LABEL
+ assigning the output of a shorthand phrase or definition to a LABEL
+ Output a LABEL value to a file
+ Output all LABEL values to a file
+ Output a LABEL assignment statement to a file
+ Output all assignment statements to a file

*shorthand* replaces the LABEL with the value assigned to it whereever it is encountered in the text being passed. The assignment statement is not output by the preprocessor.


+ Assign a string to a label
    + LABEL := STRING
+ Assign the contents of a file to a label
    + LABEL :< FILENAME
    + @content@ :< README.md
+ Assign the output of a shell command to a label
    + LABEL :! SHELL_COMMAND
    + @content@ :! cat README.md | markdown
+ Assign shorthand expansions to a LABEL
    + LABEL :{ SHORTHAND_TO_BE_EXPANDED
    + @content@ :{ @report_name@ @report_date@
        + this would concatenate report name and date
+ Render a LABEL value to a file
    + LABEL :> FILENAME
    + @content@ :> page.txt
+ render all LABEL values to a file
    + IGNORE_LABEL :=> FILENAME
    + _ :=> page.txt
        + By convention IGNORED_LABEL is an underscore
        + There is no guaranteed order to the values written out
+ render a LABEL assignment statement to a file
    + LABEL :} FILENAME
    + @content@ :} mydef.shorthand
+ render all LABEL assignment statements to a file
    + IGNORED_LABEL :=} FILENAME 
    + _ :=} mydefs.shorthand
        + By convention IGNORED_LABEL is an underscore
        + There is no guaranteed order to the values written out


The spaces surrounding " := ", " :< ", " :! ", " :{ ", " :> ", " :=> ", " :} " and " :=} " are required.

## Example


In this example a file containing the text of pre-amble is assigned to the
label @PREAMBLE, the time 3:30 is assigned to the label @NOW.

```text
    @PREAMBLE :< /home/me/preamble.text
    @NOW := 3:30

    At @NOW I will be reading the @PREAMBLE until everyone falls asleep.
```

If the file preamble.txt contained the phrase "Hello World" (including
the quotes but without any carriage return or line feed) the output after
processing the shorthand would look like -

```text

    At 3:30 I will be reading the "Hello World" until everyone falls asleep.
```

Notice the lines containing the assignments are not included in the output and that no carriage returns or line feeds are added the the substituted labels.


