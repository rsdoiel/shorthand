#!/bin/bash
#
shorthand \
    -e "@page_title := Shorthand - a simple label expander" \
    -e "@content :! cat README.md | sed -e 's/shorthand\.md/shorthand.html/g' | wsmarkdown" \
    -e "@copyright :! cat copyright.md | wsmarkdown" \
    < page.shorthand > index.html
shorthand \
    -e "@page_title := Shorthand - how it works" \
    -e "@content :! cat shorthand.md | wsmarkdown" \
    -e "@copyright :! cat copyright.md | wsmarkdown" \
    < page.shorthand > shorthand.html

