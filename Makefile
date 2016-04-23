#
# Shorthand a text label expander.
#
# @author R. S. Doiel, <rsdoiel@gmail.com>
# copyright (c) 2015 all rights reserved.
# Released under the BSD 2-Clause license
# See: http://opensource.org/licenses/BSD-2-Clause
#
build: cmds/shorthand/shorthand.go shorthand.go ops.go
	go build -o bin/shorthand cmds/shorthand/shorthand.go

lint:
	gofmt -w shorthand.go && golint shorthand.go
	gofmt -w shorthand_test.go && golint shorthand_test.go
	gofmt -w cmds/shorthand/shorthand.go && golint cmds/shorthand/shorthand.go

test:
	go test

clean:
	if [ -f bin/shorthand ]; then rm bin/shorthand; fi

install:
	go install cmds/shorthand/shorthand.go

uninstall:
	if [ -f $(GOBIN)/shorthand ]; then /bin/rm $(GOBIN)/shorthand; fi

doc: build.shorthand nav.md shorthand.md TODO.md ideas.md README.md
	shorthand build.shorthand

