#
# Shorthand a text label expander.
#
# @author R. S. Doiel, <rsdoiel@gmail.com>
# copyright (c) 2015 all rights reserved.
# Released under the BSD 2-Clause license
# See: http://opensource.org/licenses/BSD-2-Clause
#
build:
	go build -o bin/shorthand cmds/shorthand/shorthand.go
	shorthand build.shorthand

lint:
	gofmt -w shorthand.go && golint shorthand.go
	gofmt -w shorthand_test.go && golint shorthand_test.go
	gofmt -w cmds/shorthand/shorthand.go && golint cmds/shorthand/shorthand.go

test:
	go test

clean:
	if [ -d bin ]; then /bin/rm -fR bin; fi
	if [ -d dist ]; then /bin/rm -fR dist; fi
	if [ -f shorthand-binary-release.zip ]; then /bin/rm shorthand-binary-release.zip; fi

install:
	GOBIN=$(HOME)/bin go install cmds/shorthand/shorthand.go

uninstall:
	if [ -f $(GOBIN)/shorthand ]; then /bin/rm $(GOBIN)/shorthand; fi

doc:
	shorthand build.shorthand

release:
	./mk-release.sh

publish:
	./publish.sh
