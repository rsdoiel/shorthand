
# Assignments and Expansions

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

*shorthand* replaces the LABEL with the value assigned to it whereever it is encountered in the text being passed.
Commonlly this might be curly brackes, dollar signs or even at signs.  Doesn't really matter but it needs to be unique
and cannot be in the pattern of space, colon, string, colon and space.  An assignment statement is not written to stdout output.

operator                    | meaning                                  | example
----------------------------|------------------------------------------|---------------------------------------------------------------------
 :label:                    | Assign String                            | {{name}} :label: Freda
----------------------------|------------------------------------------|---------------------------------------------------------------------
 :import-text:              | Assign the contents of a file            | {{content}} :import-text: myfile.txt
----------------------------|------------------------------------------|---------------------------------------------------------------------
 :import-shorthand:         | Get assignments from a file              | _ :import-shorthand: myfile.shorthand
----------------------------|------------------------------------------|---------------------------------------------------------------------
 :expand:                   | Assign an expansion                      | $reportTitle$ :expand: Report: @title for @date
----------------------------|------------------------------------------|---------------------------------------------------------------------
 :expand-expansion:         | Assign expanded expansion                | {{reportHeading}} :expand-expansion: @reportTitle
----------------------------|------------------------------------------|---------------------------------------------------------------------
 :import-expansion:         | Include Expansion                        | @nav@ :import-expansion: mynav.html
----------------------------|------------------------------------------|---------------------------------------------------------------------
 :bash:                     | Assign Shell output                      | {{date}} :bash: date +%Y-%m-%d
----------------------------|------------------------------------------|---------------------------------------------------------------------
 :expand-and-bash:          | Assign Expand then gete Shell output     | {{entry}} :expand-and-bash: cat header.txt @filename footer.txt
----------------------------|------------------------------------------|---------------------------------------------------------------------
 :markdown:                 | Assign Markdown processed text           | {div} :markdown: # My h1 for a Div
----------------------------|------------------------------------------|---------------------------------------------------------------------
 :expand-markdown:          | Assign Expanded Markdown                 | {{div}} :expand-markdown: Greetings **@name**
----------------------------|------------------------------------------|---------------------------------------------------------------------
 :import-markdown:          | Include Markdown processed text          | $nav$ :import-markdown: mynav.md
----------------------------|------------------------------------------|---------------------------------------------------------------------
 :import-expanded-markdown: | Include Expanded Markdown processed text | {nav} :import-expanded-markdown: mynav.md
----------------------------|------------------------------------------|---------------------------------------------------------------------
 :export-expansion:         | Output Assigned Expansion                | {{content}} :export-expansion: content.txt
----------------------------|------------------------------------------|---------------------------------------------------------------------
 :export-all-expansions:    | Output all assigned Expansions           | _ :export-all-expansions: contents.txt
----------------------------|------------------------------------------|---------------------------------------------------------------------
 :export-label:             | Output Assignment                        | {{content}} :export-label: content.shorthand
----------------------------|------------------------------------------|---------------------------------------------------------------------
 :export-all-labels:        | Output all Assignments                   | _ :export-all-labels: contents.shorthand
----------------------------|------------------------------------------|---------------------------------------------------------------------
 :exit:                     | Exit the shorthand repl                  | :exit:
----------------------------|------------------------------------------|---------------------------------------------------------------------



Notes: Using an underscore as a LABEL means the label will be ignored. There are no guarantees of order when writing values or assignment statements to a file.

The spaces surrounding " :label: ", " :import-text: ", " :bash: ", " :expand: ", " :export-expansion: ", etc. are required.


## Example

In this example a file containing the text of pre-amble is assigned to the label @PREAMBLE, the time 3:30 is assigned to the label @NOW.  
```text
    {{PREAMBLE}} :import-text: /home/me/preamble.text
    {{NOW}} :label: 3:30

    At {{NOW}} I will be reading the {{PREAMBLE}} until everyone falls asleep.
```

If the file preamble.txt contained the phrase "Hello World" (including the quotes but without any carriage return or line feed) the output after processing the shorthand would look like -

```text

    At 3:30 I will be reading the "Hello World" until everyone falls asleep.
```

Notice the lines containing the assignments are not included in the output and that no carriage returns or line feeds are added the the substituted labels.
+ Assign shorthand expansions to a LABEL
    + LABEL :expand: SHORTHAND_TO_BE_EXPANDED
    + @content@ :expand: @report_name@ @report_date@
        + this would concatenate report name and date

### Processing Markdown pages

_shorthand_ also provides a markdown processor. It uses the [blackfriday](https://github.com/russross/blackfriday) markdown library. This is both a convience and also allows you to treat markdown with shorthand assignments as a template that renders HTML or HTML with shorthand ready for expansion. It is a poorman's text rendering engine.

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
    @doctype :bash: echo "<!DOCTYPE html>"
    @headBlock :label: <head><title>@pageTitle</title>
    @pageTemplate :import-markdown: post-template.md
    @dateString :bash: date
    @blogTitle :label:  My Blog
    @pageTitle :label A Post
    @contentBlock :import-markdown: a-post.md
    @output :expand-expansion: @doctype<html>@headBlock<body>@pageTemplate</body></html>
    @output :export-expansion: post.html
```


