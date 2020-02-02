NAME := vanity-keygen
VERSION := $(shell git tag | tail -n 1)
BUILD_STRING := $(shell git log --pretty=format:'%h' -n 1)
VERSION_LONG := $(NAME) version $(VERSION)+$(BUILD_STRING)
BUILD_DATE := $(shell date -u)

SRC_FILES := $(shell find . -type f -name *.go)
SRC_DIRS := ./cmd/vanity-keygen/ ./internal/vanitykeygen/
MAIN_SRC := cmd/vanity-keygen/main.go

LDFLAGS := "-X \"github.com/brannondorsey/vanity-keygen/internal/vanitykeygen.VERSION=$(VERSION)\" -X \"github.com/brannondorsey/vanity-keygen/internal/vanitykeygen.VERSION_LONG=$(VERSION_LONG)\" -X \"github.com/brannondorsey/vanity-keygen/internal/vanitykeygen.BUILD_DATE=$(BUILD_DATE)\""

.PHONY: all deps clean install

default: deps build

install: deps build
	cp ./bin/$(NAME) /usr/local/bin/$(NAME)

build: $(SRC_FILES)
	go build -ldflags=$(LDFLAGS) -o bin/$(NAME) $(MAIN_SRC)

build-all: $(SRC_FILES)
	mkdir -p bin/macos bin/linux-x64 bin/linux-arm64 bin/windows
	GOOS=linux GOARCH=amd64 go build -o bin/linux-x64/$(NAME) $(MAIN_SRC)
	GOOS=linux GOARCH=arm64 go build -o bin/linux-arm64/$(NAME) $(MAIN_SRC)
	GOOS=darwin GOARCH=amd64 go build -o bin/macos/$(NAME) $(MAIN_SRC)
	GOOS=windows GOARCH=amd64 go build -o bin/windows/$(NAME) $(MAIN_SRC)

deps:
	go get $(SRC_DIRS)

clean:
	go clean
	rm -rf bin/*
	touch bin/.gitkeep

ifneq (,$(findstring snapshot,$(VERSION)))
snapshot:
	echo "The latest tagged version is a snapshot. Not tagging."
else
snapshot:
	echo "The latest version tagged is not a snapshot. Tagging!"
	git tag $(VERSION)-snapshot
	git push --tags
endif
