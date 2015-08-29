#!/bin/bash
#
GO=$(which go)
if [ "$GO" = "" ]; then
    echo "Must install Golang first"
    echo "See http://golang.org for instructions"
    exit 1
fi

# Install dependent libraries
# Add test package for shorthand
go get github.com/rsdoiel/ok
# Markdown library
go get github.com/russross/blackfriday

make
make test
