#!/bin/bash
#

# Make all my definitions.
cat > index.shorthand <<EOF
@page_title := Shorthand - a simple label expander
@content :! cat README.md | sed -e 's/shorthand\.md/shorthand.html/g' | wsmarkdown
@copyright :! cat copyright.md | wsmarkdown
EOF
# Cat the definitions to page.shorthand and render the result to index.html
cat index.shorthand page.shorthand | shorthand > index.html


# Add the defs for shorthand.md page.
cat > shorthand.shorthand <<EOF
@page_title := Shorthand - how it works
@content :! cat shorthand.md | wsmarkdown
@copyright :! cat copyright.md | wsmarkdown
EOF
cat shorthand.shorthand page.shorthand | shorthand > shorthand.html 

