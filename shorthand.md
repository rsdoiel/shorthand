
# Assignments and Expansions

Shorthand is a simple text label expansion utility. It is based on a simple key value substitution.  It supports this following types of definitions

+ Assign a string to a LABEL
+ Assign the contents of a file to a LABEL
+ Assign the output of a Bash shell expression to a LABEL
+ Assign the output of a shorthand expansion to a LABEL
+ Read a file of shorthand assignments and assign any expansions to the LABEL
+ Output a LABEL value to a file
+ Output all LABEL values to a file
+ Output a LABEL assignment statement to a file
+ Output all assignment statements to a file

*shorthand* replaces the LABEL with the value assigned to it whereever it is encountered in the text being passed. The assignment statement is not output by the preprocessor.


+ Assign a string to a label
    + LABEL := STRING
+ Assign the contents of a file to a label
    + LABEL :=< FILENAME
    + @content@ :=< README.md
+ Assign the output of a shell command to a label
    + LABEL :! SHELL_COMMAND
    + @content@ :! cat README.md | markdown
+ Assign shorthand expansions to a LABEL
    + LABEL :{ SHORTHAND_TO_BE_EXPANDED
    + @content@ :{ @report_name@ @report_date@
        + this would concatenate report name and date
+ Read a file of shorthand assignments and assign any expansions to the label
    + LABEL :{< FILE_CONTAINING_SHORTHAND_DEFINISIONS_AND_EXPANSIONS
        + _ :{< myfile.shorthand
        + Assignments are added to the abbreviation list in memory and any expansions to LABEL
+ Render a LABEL value to a file
    + LABEL :> FILENAME
    + @content@ :> page.txt
+ render all LABEL values to a file
    + LABEL :*> FILENAME
    + _ :=> page.txt
+ render a LABEL assignment statement to a file
    + LABEL :}> FILENAME
    + @content@ :}> mydef.shorthand
+ render all LABEL assignment statements to a file
    + LABEL :*}> FILENAME 
    + _ :*}> mydefs.shorthand
+ rendered markdown to a label
    + LABEL :[ MARKDOWN_EXPRESSION
+ render markdown file to a label
    + LABEL :[< FILENAME

Notes: Using an underscore as a LABEL means the label will be ignored. There are no guarantees of order when writing values or assignment statements to a file.

The spaces surrounding " := ", " :=< ", " :! ", " :{ ", " :> ", " :*> ", " :} ", " :}> ", etc. are required.


## Example

In this example a file containing the text of pre-amble is assigned to the label @PREAMBLE, the time 3:30 is assigned to the label @NOW.  
```text
    @PREAMBLE :=< /home/me/preamble.text
    @NOW := 3:30

    At @NOW I will be reading the @PREAMBLE until everyone falls asleep.
```

If the file preamble.txt contained the phrase "Hello World" (including the quotes but without any carriage return or line feed) the output after processing the shorthand would look like -

```text

    At 3:30 I will be reading the "Hello World" until everyone falls asleep.
```

Notice the lines containing the assignments are not included in the output and that no carriage returns or line feeds are added the the substituted labels.

### Processing Markdown pages

_shorthand_ also provides a markdown processor. It uses the [blackfriday](https://github.com/russross/blackfriday) markdown library. This is both a convience and also allows you to treat markdown with shorthand assignments as a template that renders HTML or HTML with shorthand ready for expansion. It can effectivly be a poorman's rendering engine.

In this example we'll build a HTML page with shorthand labels from markdown text. Then
we will use the render HTML as a template for a blog page entry.

Our markdown file serving as a template will be call "post-template.md". It should contain
the outline of the structure of the page plus some shorthand labels we'll expand later.

```markdown

    # @blogTitle

    ## @pageTitle

    ### @dateString

    @contentBlocks

```

For the purposes of this exercise we'll use _shorthand_ as a repl and just enter the
assignments sequencly.  Also rather than use the output of shorthand directly we'll
build up the content for the page in a label and use shorthand itself to write the final
page out.

The steps we'll follow will be to 

1. Read in our markdown file page.md and turn it into an HTML with embedded shorthand labels
2. Assign some values to the labels
3. Expand the labels in the HTML and assign to a new label
4. Write the new label out to are page call "page.html"

Start the repl with this version of the shorthand command:

```shell
    shorthand -p "? "
```

The _-p_ option tells _shorthand_ to use the value "? " as the prompt. When _shorthand_ starts
it will display "? " to indicate it is ready for an assignment or expansion.

The following assumes you are in the _shorthand_ repl.

Load the mardkown file and transform it into HTML with embedded shorthand labels

```shell
    @doctype :! echo "<!DOCTYPE html>"
    @headBlock := <head><title>@pageTitle</title>
    @pageTemplate :[< post-template.md
    @dateString :! date
    @blogTitle :=  My Blog
    @pageTitle := A Post
    @contentBlock :[< a-post.md
    @output :{{ @doctype<html>@headBlock<body>@pageTemplate</body></html>
    @output :> post.html
```


