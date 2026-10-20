[![Project Status: Active – The project has reached a stable, usable state and is being actively developed.](https://www.repostatus.org/badges/latest/active.svg)](https://www.repostatus.org/#active)


# shorthand

A simple label or macro expander

Example use cases:

+ label or abbreviation expansion in Markdown files
+ build html templates from markdown files
+ compose pages from multiple markdown files

The supported command line options can be listed using the _-help_
option.

```shell
    shorthand -help
```
Source code can be found at [github.com/rsdoiel/shorthand](https://github.com/rsdoiel/shorthand)

The project website is [rsdoiel.github.io/shorthand](http://rsdoiel.github.io/shorthand)


## Tutorial

### Timestamp in Markdown file

If the content of the markdown file _testdata/report.md_ was

```markdown

    Report Date: @now

    # Topic: The current local time.

    This report highlights the current local time of rendering this document

    The current local time is @now

```

Put the definitions in a file of their own, _now.shorthand_

```text
    :bash: @now date
```

and give it to _shorthand_ before the document, then render the result
with _pandoc_

```shell
    shorthand now.shorthand testdata/report.md | pandoc -s > testdata/report.html
```

What this command does is read the definitions, then the document, and
replace all occurrences of "@now" in the document with the
output from the Unix program _date_. 

The output (before piping to Pandoc) would look something like

```html
    Report Date: Sat Aug 29 11:25:48 PDT 2015

    # Topic: The current local time.

    This report highlights the current local time of rendering this document

    The current local time is Sat Aug 29 11:25:48 PDT 2015
```

Notice that both "@now" are replace with the same date information.

### embedding shorthand definitions

You could also embed the shorthand definitions straight in the
markdown itself. with something like

```markdown
    :bash: @now date

    Report Date: @now

    # Topic: The current local time.

    This report highlights the current local time of rendering this document

    The current local time is @now

```

That makes the command line a little shorter

```shell
    shorthand testdata/report.md | pandoc -s > testdata/report.html
```


## Installation

_shorthand_ can be installed with the *go install* command.

```
    go install github.com/rsdoiel/shorthand/cmd/shorthand@latest
```

## Documentation

- [shorthand command](shorthand.1.html)
- [shorthand syntax](shorthand-syntax.html)
- [shorthand tutorial](shorthand-tutorial.html)
