#
# Shorthand a text label expander.
#
# @author R. S. Doiel, <rsdoiel@gmail.com>
# copyright (c) 2015 all rights reserved.
# Released under the BSD 2-Clause license
# See: http://opensource.org/licenses/BSD-2-Clause
#
PROJECT = shorthand

VERSION = $(shell grep -m1 'Version = ' $(PROJECT).go | cut -d\"  -f 2)

BRANCH = $(shell git branch | cut -d\  -f 2)

OS = $(shell uname)

EXT =
ifeq ($(OS), Windows)
    EXT = .exe
endif

build:
	go build -o bin/shorthand$(EXT) cmds/shorthand/shorthand.go
	bin/shorthand$(EXT) build.shorthand

lint:
	gofmt -w shorthand.go && golint shorthand.go
	gofmt -w shorthand_test.go && golint shorthand_test.go
	gofmt -w cmds/shorthand/shorthand.go && golint cmds/shorthand/shorthand.go

test:
	go test

status:
	git status

save:
	if [ "$(msg)" != "" ]; then git commit -am "$(msg)"; else git commit -am "Quick Save"; fi
	git push origin $(BRANCH)

clean:
	if [ -d bin ]; then rm -fR bin; fi
	if [ -d dist ]; then rm -fR dist; fi

install:
	GOBIN=$(HOME)/bin go install cmds/shorthand/shorthand.go

uninstall:
	if [ -f $(GOBIN)/shorthand ]; then rm $(GOBIN)/shorthand; fi

doc:
	shorthand build.shorthand

dist/linux-amd64:
	mkdir -p dist/bin
	env GOOS=linux GOARCH=amd64 go build -o dist/bin/shorthand cmds/shorthand/shorthand.go
	cd dist && zip -r $(PROJECT)-$(VERSION)-linux-amd64.zip README.md LICENSE INSTALL.md bin/*
	rm -fR dist/bin

dist/macosx-amd64:
	mkdir -p dist/bin
	env GOOS=darwin	GOARCH=amd64 go build -o dist/bin/shorthand cmds/shorthand/shorthand.go
	cd dist && zip -r $(PROJECT)-$(VERSION)-macosx-amd64.zip README.md LICENSE INSTALL.md bin/*
	rm -fR dist/bin

dist/windows-amd64:
	mkdir -p dist/bin
	env GOOS=windows GOARCH=amd64 go build -o dist/bin/shorthand.exe cmds/shorthand/shorthand.go
	cd dist && zip -r $(PROJECT)-$(VERSION)-windows-amd64.zip README.md LICENSE INSTALL.md bin/*
	rm -fR dist/bin

dist/raspbian-arm7:
	mkdir -p dist/bin
	env GOOS=linux GOARCH=arm GOARM=7 go build -o dist/bin/shorthand cmds/shorthand/shorthand.go
	cd dist && zip -r $(PROJECT)-$(VERSION)-raspbian-arm7.zip README.md LICENSE INSTALL.md bin/*
	rm -fR dist/bin

dist/linux-arm64:
	mkdir -p dist/bin
	env GOOS=linux GOARCH=arm64 go build -o dist/bin/shorthand cmds/shorthand/shorthand.go
	cd dist && zip -r $(PROJECT)-$(VERSION)-linux-arm64.zip README.md LICENSE INSTALL.md bin/*
	rm -fR dist/bin

distrubute_docs:
	mkdir -p dist
	cp -v README.md dist/
	cp -v LICENSE dist/
	cp -v INSTALL.md dist/

release: distrubute_docs dist/linux-amd64 dist/windows-amd64 dist/macosx-amd64 dist/raspbian-arm7 dist/linux-arm64

website:
	go build -o bin/shorthand$(EXT) cmds/shorthand/shorthand.go
	#./mk-website.bash

publish:
	go build -o bin/shorthand$(EXT) cmds/shorthand/shorthand.go
	./publish.bash
