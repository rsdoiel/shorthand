% shorthand-tutorial(1) tutorial
% R. S. Doiel
% August 4, 2022

# NAME

shorthand - tutorial

# DESCRIPTION

Shorthand is a simple label or macro expander. This is a brief
tutorial an using it.

### Timestamp in Markdown file

If the content of the markdown file _template/report.md_ was

```markdown

    Report Date: @now@

    # Topic: The current local time.

    This report highlights the current local time of rendering this document

    The current local time is @now@

```

From the command line you can do something like this

```shell
    shorthand -e ':bash: @now@ date' \
        -e ":import: @report@ template/report.md" \
        -e "@report@" \
        -e ':exit:' | pandoc -s > reports/report-time.html
```

What this command does is launch the _shorthand_ interpreter and it
replaces all occurrences of "@now@" in the markdown document with the
output from the Unix program _date_. 

The output (before piping to Pandoc) would look something like

```html
    Report Date: Sat Aug 29 11:25:48 PDT 2015

    # Topic: The current local time.

    This report highlights the current local time of rendering this document

    The current local time is Sat Aug 29 11:25:48 PDT 2015
```

Notice that both "@now@" are replace with the same date information.

In the end of the command line example we also have an `@report@`.
First it is defined as an "import" and then it is rendered before
exiting the interpreter. The "import" reads the template, evaluates
the template rendering the result at the second `@report@`. This result
is what is passed via a pipe to Pandoc to generate the HTML page
_reports/report-time.html_

### embedding shorthand definitions

We can make that command line easier by embedding the definitions in our
template.

You could also embed the shorthand definitions command straight in the
markdown itself. with something like

```markdown
    :bash: @now@ date

    Report Date: @now@

    # Topic: The current local time.

    This report highlights the current local time of rendering this document

    The current local time is @now

```

That simpler command is much shorter. This is the more typical
usage.

```shell
    shorthand testdata/report.md | pandoc -s > testdata/report.html
```


## Installation

_shorthand_ can be installed with the *go get* command.

```
    go get github.com/rsdoiel/shorthand/...
```

# ALSO SEE

- [shorthand](shorthand.html)
- [shorthand-syntax](shorthand-syntax.html)
- Website: [https://rsdoiel.github.io/shorthand](https://rsdoiel.github.io/shorthand)


