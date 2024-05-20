
USAGE: shorthand [OPTIONS] [SHORTHAND_FILES]

SYNOPSIS

shorthand a simple label expander

DESCRIPTION

shorthand is a command line utility to expand labels 
based on their assigned definitions. The render output is the 
transformed text and without the shorthand definitions themselves. 
shorthand reads from standard input and writes to standard output.

OPTIONS

    -examples            display examples
    -generate-markdown   output documentation in Markdown
    -h, -help            display help
    -i, -input           input filename
    -l, -license         display license
    -n, -no-prompt       Turn off the prompt for interactive processing
    -o, -output          output filename
    -p, -prompt          Output a prompt for interactive processing
    -quiet               suppress error messages
    -v, -version         diplsay version


EXAMPLES

ASSIGNMENTS AND EXPANSIONS

Shorthand is a simple label expansion utility. It is based on a simple key value substitution.  It supports this following types of definitions

+ Assign a string to a LABEL
+ Assign the contents of a file to a LABEL
+ Assign the output of a Bash shell expression to a LABEL
+ Assign the output of a shorthand expansion to a LABEL
+ Read a file of shorthand assignments and assign any expansions to the LABEL
+ Output a LABEL value to a file
+ Output all LABEL values to a file
+ Output a LABEL assignment statement to a file
+ Output all assignment statements to a file

*shorthand* replaces the LABEL with the value assigned to it whereever it is encountered in the text being passed. The assignment statement is 
not written to stdout output.

operator                    | meaning                                  | example
----------------------------|------------------------------------------|---------------------------------------------------------------------
 :set:                      | Assign String                            | :set: {{name}} Freda
----------------------------|------------------------------------------|---------------------------------------------------------------------
 :import-text:              | Assign the contents of a file            | :import-text: {{content}} myfile.txt
----------------------------|------------------------------------------|---------------------------------------------------------------------
 :import-shorthand:         | Get assignments from a file              | :import-shorthand: _ myfile.shorthand
----------------------------|------------------------------------------|---------------------------------------------------------------------
 :expand:                   | Assign an expansion                      | :expand: {{reportTitle}} Report: @title for @date
----------------------------|------------------------------------------|---------------------------------------------------------------------
 :expand-expansion:         | Assign expanded expansion                | :expand-expansion: {{reportHeading}} @reportTitle
----------------------------|------------------------------------------|---------------------------------------------------------------------
 :import:                   | Include a file, procesisng the shorthand | :import: {{nav}} mynav.shorthand
----------------------------|------------------------------------------|---------------------------------------------------------------------
 :bash:                     | Assign Shell output                      | :bash: {{date}} date +%Y-%m-%%d
----------------------------|------------------------------------------|---------------------------------------------------------------------
 :expand-and-bash:          | Assign Expand then gete Shell output     | :expand-and-bash: {{entry}} cat header.txt @filename footer.txt
----------------------------|------------------------------------------|---------------------------------------------------------------------
 :export:                   | Output a label's value to a file         | :export: {{content}} content.txt
----------------------------|------------------------------------------|---------------------------------------------------------------------
 :export-all:               | Output all assigned Expansions           | :export-all: _ contents.txt
----------------------------|------------------------------------------|---------------------------------------------------------------------
 :export-shorthand:             | Output Assignment                        | :export-shorthand: {{content}} content.shorthand
----------------------------|------------------------------------------|---------------------------------------------------------------------
 :export-all-shorthand:        | Output all shorthand assignments      | :export-all-shorthand: _ contents.shorthand
----------------------------|------------------------------------------|---------------------------------------------------------------------
 :exit:                     | Exit the shorthand repl                  | :exit:
----------------------------|------------------------------------------|---------------------------------------------------------------------



Notes: Using an underscore as a LABEL means the label will be ignored. There are no guarantees of order when writing values or assignment 
statements to a file.

The spaces following surrounding ":set:", ":import-text:", ":bash:", ":expand:", ":export:", etc. are required.


EXAMPLE

In this example a file containing the text of pre-amble is assigned to the label @PREAMBLE, the time 3:30 is assigned to the label {{NOW}}.

    :import-text: {{PREAMBLE}} /home/me/preamble.text
	:set: {{NOW}} 3:30

    At {{NOW}} I will be reading the {{PREAMBLE}} until everyone falls asleep.


If the file preamble.txt contained the phrase "Hello World" (including the quotes but without any carriage return or line feed) the output 
after processing the shorthand would look like - 

    At 3:30 I will be reading the "Hello World" until everyone falls asleep.

Notice the lines containing the assignments are not included in the output and that no carriage returns or line feeds are added the the 
substituted labels.

+ Assign shorthand expansions to a LABEL
    + :expand: LABEL SHORTHAND_TO_BE_EXPANDED
	+ :expand: @content@ @report_name@ @report_date@
        + this would concatenate report name and date


PROCESSING MARKDOWN PAGES

_shorthand_ is a label expander or light weight macro expande. It
can be combined with programs like *pandoc* as a pre-processor.

In this example we'll build a HTML page with shorthand labels from 
a couple markdown documents. Then we will use the render HTML as a 
template for a blog page entry.

Our markdown file serving as a template will be call 
"post-template.md".  It should contain the outline of the 
structure of the page plus some shorthand labels we'll expand 
later.


    # @blogTitle

    ## @pageTitle

    ### @dateString

    @contentBlocks


For the purposes of this exercise we'll use _shorthand_ as a repl 
and just enter the assignments sequencly.  Also rather than use 
the output of shorthand directly we'll build up the content for the 
page in a label and use shorthand itself to write the final page out.

The steps we'll follow will be to 

1. Read in our markdown file page.md and turn it into an HTML with embedded shorthand labels
2. Assign some values to the labels
3. Expand the labels in the HTML and assign to a new label
4. Write the new label out to are page call "page.html"

Start the repl with this version of the shorthand command:

    shorthand -p "? "

The _-p_ option tells _shorthand_ to use the value "? " as the prompt. When _shorthand_ starts it will display "? " to indicate it is 
ready for an assignment or expansion.

The following assumes you are in the _shorthand_ repl.

Load the mardkown file and transform it into HTML with embedded shorthand labels

	:import: @pageTemplate post-template.md
	:bash: @dateString date
	:set: @blogTitle My Blog
	:set: @pageTitle A Post
	:import: @contentBlock a-post.md
	:export: @output post.md
	:exit:

Then we can run the following command to pipe it through Pandoc.

    pandoc -s -t html post.md > post.html

shorthand 0.2.2
