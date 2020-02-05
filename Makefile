NAME := vanity-keygen
VERSION := $(shell git tag | tail -n 1)
BUILD_STRING := $(shell git log --pretty=format:'%h' -n 1)
VERSION_LONG := $(NAME) version $(VERSION)+$(BUILD_STRING)
BUILD_DATE := $(shell date -u)

SRC_FILES := $(shell find . -type f -name *.go)
MAIN_SRC := cmd/vanity-keygen/main.go

LDFLAGS := "-X \"github.com/brannondorsey/vanity-keygen/pkg/vanitykeygen.VERSION=$(VERSION)\" -X \"github.com/brannondorsey/vanity-keygen/pkg/vanitykeygen.VERSION_LONG=$(VERSION_LONG)\" -X \"github.com/brannondorsey/vanity-keygen/pkg/vanitykeygen.BUILD_DATE=$(BUILD_DATE)\""

.PHONY: default deps clean install snapshot

default: deps build

install: deps build
	cp ./bin/$(NAME) /usr/local/bin/$(NAME)

build: $(SRC_FILES)
	go build -ldflags=$(LDFLAGS) -o bin/$(NAME) $(MAIN_SRC)

build-all: $(SRC_FILES) clean default
	mkdir -p bin/$(NAME)-macos bin/$(NAME)-windows bin/$(NAME)-linux-x64 bin/$(NAME)-linux-arm7 bin/$(NAME)-linux-arm6
	GOOS=darwin GOARCH=amd64 go build -ldflags=$(LDFLAGS) -o bin/$(NAME)-macos/$(NAME) $(MAIN_SRC)
	GOOS=windows GOARCH=amd64 go build -ldflags=$(LDFLAGS) -o bin/$(NAME)-windows/$(NAME) $(MAIN_SRC)
	GOOS=linux GOARCH=amd64 go build -ldflags=$(LDFLAGS) -o bin/$(NAME)-linux-x64/$(NAME) $(MAIN_SRC)
	GOOS=linux GOARCH=arm GOARM=7 go build -ldflags=$(LDFLAGS) -o bin/$(NAME)-linux-arm7/$(NAME) $(MAIN_SRC)
	GOOS=linux GOARCH=arm GOARM=6 go build -ldflags=$(LDFLAGS) -o bin/$(NAME)-linux-arm6/$(NAME) $(MAIN_SRC)
	cd bin/ && \
		tar czf $(NAME)-linux-x64.tar.gz $(NAME)-linux-x64/ && \
		tar czf $(NAME)-linux-arm7.tar.gz $(NAME)-linux-arm7/ && \
		tar czf $(NAME)-linux-arm6.tar.gz $(NAME)-linux-arm6/ && \
		zip -r -9 $(NAME)-macos.zip $(NAME)-macos/ && \
		zip -r -9 $(NAME)-windows.zip $(NAME)-windows/
	rm -rf bin/$(NAME)-macos bin/$(NAME)-windows bin/$(NAME)-linux-x64 bin/$(NAME)-linux-arm7 bin/$(NAME)-linux-arm6

deps:
	go mod vendor

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
