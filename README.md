[![Project Status: Suspended – Initial development has started, but there has not yet been a stable, usable release; work has been stopped for the time being but the author(s) intend on resuming work.](https://www.repostatus.org/badges/latest/suspended.svg)](https://www.repostatus.org/#suspended)


# shorthand

A simple label expander and markdown utility

Example use cases:

+ label or abbreviation expansion in Markdown files
+ build html templates from markdown files
+ compose pages from multiple markdown files

The supported command line options can be listed using the _--help_
options.

```shell
    shorthand --help
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

From the command line you can do something like this

```shell
    shorthand -e ':bash: @now date' \
        -e ":import-text: @report testdata/report.md" \
        -e ":markdown: @html @report" \
        -e "@html" \
        -e ':exit:' > testdata/report.html
```

What this command does is launch the _shorthand_ interpreter and it
replaces all occurrences of "@now" in the markdown document with the
output from the Unix program _date_. 

The output would look something like

```html
    <p>Report Date: Sat Aug 29 11:25:48 PDT 2015</p>

    <h1>Topic: The current local time.</h1>

    <p>This report highlights the current local time of rendering this document</p>

    <p>The current local time is Sat Aug 29 11:25:48 PDT 2015</p>
```

Notice that both "@now" are replace with the same date information.

### embedding shorthand definitions

You could also embed the shorthand definitions command straight in the
markdown itself. with something like

```markdown
    @now :bash: date

    Report Date: @now

    # Topic: The current local time.

    This report highlights the current local time of rendering this document

    The current local time is @now

```

That makes the command line a little shorter

```shell
    shorthand -markdown testdata/report.md > testdata/report.html
```


## Installation

_shorthand_ can be installed with the *go get* command.

```
    go get github.com/rsdoiel/shorthand/...
```


