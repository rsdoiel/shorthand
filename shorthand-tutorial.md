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

Put the definitions in a file, _now.shorthand_

```text
    :bash: @now@ date
    :import: @report@ template/report.md
    @report@
```

and run it through _shorthand_

```shell
    shorthand now.shorthand | pandoc -s > reports/report-time.html
```

What this command does is read the definitions and it
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

In the end of the definitions file we also have an `@report@`.
First it is defined as an "import" and then it is rendered. The "import" reads the template, evaluates
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


