%shorthand(1) user manual | version 0.3.0 67d8213
% R. S. Doiel
% 2026-10-20

# NAME

shorthand — a simple label expander

# SYNOPSIS

shorthand [OPTIONS] [SHORTHAND_FILES]

# DESCRIPTION

shorthand expands labels in text using the definitions assigned to them. It
reads each file named on the command line in turn, or standard input when there
is none, and writes the text with the labels replaced. The lines that make the
definitions are not written. It works as a naive macro language, for the jobs a
template language is too much for: abbreviations in Markdown, a timestamp in a
timesheet, pages built from several documents. Run shorthand -examples for the
operators.

Without a file, shorthand is an interactive interpreter when standard input is
a terminal: it prompts for each line, and ":help:" lists the operators. Use
":exit:" to leave. A line that contains ":exit:" ends the input in a file too.

The definitions and the labels live for one run, so the files given together
share them. The first line below defines {{who}}, so the second prints
"Hello Freda":

~~~
:set: {{who}} Freda
Hello {{who}}
~~~

# OPTIONS

-h, -help
: display this help page

-l, -license
: display the license

-v, -version
: display the program name, version and release hash

-examples
: display the operators with examples

-i FILENAME, -input FILENAME
: read FILENAME before any files named on the command line

-o FILENAME, -output FILENAME
: write to FILENAME rather than standard output

-p PROMPT, -prompt PROMPT
: the prompt of the interactive interpreter (default "=> ")

-n, -no-prompt
: turn the prompt off

-quiet
: do not report a file that cannot be opened

# EXIT STATUS

0
: success

2
: bad command line

65
: a line could not be evaluated, such as an unsupported operator or a file
  an operator could not read. It is reported with its line number, the rest of
  the input is still processed, and the exit status is 65

66
: an input file cannot be opened. The other files are still processed

73
: the output file cannot be created

# EXAMPLES

Expand a timesheet that starts with a definition of @now:

~~~shell
shorthand Time_Sheet.txt | stnparse
~~~

Build a page from a template, writing it to a file:

~~~shell
shorthand -o index.html build.shorthand
~~~

Expand text on standard input, with no prompt:

~~~shell
echo "Today is @today" | shorthand -n definitions.shorthand -
~~~

# SEE ALSO

pandoc(1), stnparse(1)

