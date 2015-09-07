#
# Shorthand a text label expander.
#
# @author R. S. Doiel, <rsdoiel@gmail.com>
# copyright (c) 2015 all rights reserved.
# Released under the BSD 2-Clause license
# See: http://opensource.org/licenses/BSD-2-Clause
#
build: cmd/shorthand/shorthand.go shorthand.go
	go build -o bin/shorthand cmd/shorthand/shorthand.go

lint:
	gofmt -w shorthand.go && golint shorthand.go
	gofmt -w shorthand_test.go && golint shorthand_test.go
	gofmt -w cmd/shorthand/shorthand.go && golint cmd/shorthand/shorthand.go

test:
	go test

clean:
	if [ -f bin/shorthand ]; then rm bin/shorthand; fi

install:
	go install cmd/shorthand/shorthand.go

uninstall:
	if [ -f $(GOBIN)/shorthand ]; then /bin/rm $(GOBIN)/shorthand; fi

doc: build.shorthand nav.md shorthand.md TODO.md ideas.md README.md
	shorthand build.shorthand

