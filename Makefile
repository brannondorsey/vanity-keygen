FILES := $(shell find . -type f -name *.go)

.PHONY: all deps clean install

all: install

build: $(FILES)
	go build -o bin/vanity-keygen cmd/vanity-keygen/main.go

deps:
	go get cmd

clean:
	go clean
	rm -rf bin/*
	touch bin/.gitkeep

install: deps build
