###############################################
#
# Makefile
#
###############################################

.DEFAULT_GOAL := build

.PHONY: test

VERSION := 1.1.2

ver:
	@sed -i '' 's/^const Version = "[0-9]\{1,3\}.[0-9]\{1,3\}.[0-9]\{1,3\}"/const Version = "${VERSION}"/' cron.go

lint:
	golint .

build:
	go build ./...

demo: build
	go build -o demo cmd/demo.go
	./demo

clean:
	rm -f demo

test: build
	go test -count=1 -v .

github:
	open "https://github.com/mlavergn/gocron"

release:
	zip -r gocron.zip .
	hub release create -m "${VERSION} - Go Cron" -a gocron.zip -t master "v${VERSION}"
	open "https://github.com/mlavergn/gocron/releases"
