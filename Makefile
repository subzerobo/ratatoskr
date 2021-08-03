SHELL := /bin/bash

GOBASE= $(shell pwd)
GOCMD=go
GORUN=$(GOCMD) run
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get

REV := $(shell git rev-parse HEAD)
CHANGES := $(shell test -n "$$(git status --porcelain)" && echo '+CHANGES' || true)

#  Binary Name
BINARY := yggdrasil
VERSION := 1.0

BUILD_TIME:= $(shell date +"%Y-%m-%d %H:%M:%S")
GIT_COMMIT:= $(shell git rev-parse --short HEAD)

# List of Target OS to build the binaries
PLATFORMS=darwin linux windows
ARCHITECTURES=386 amd64
# LDFLAGS
LDFLAGS := -X 'main.GitCommit=$(GIT_COMMIT)' -X 'main.BuildTime=$(BUILD_TIME)'
GPG_SIGNING_KEY := 0ED8F693E316FBE1

ENTRY_FILE_YGGDRASIL := $(shell ls -1 cmd/yggdrasil/*.go | grep -v _test.go)
ENTRY_FILE_BIFROST := $(shell ls -1 cmd/bifrost/*.go | grep -v _test.go)

.PHONY: \
	help \
	gitready \
	default \
	clean \
	clean-artifacts \
	clean-releases \
	tools \
	test \
	vet \
	errors \
	lint \
	imports \
	fmt \
	env \
	debug-yggdrasil \
	debug-bifrost \
	build \
	build-all \
	doc \
	release \
	package-release \
	sign-release \
	version \
	swagger

all: imports fmt lint vet errors build

help:
	@echo 'Usage: make <OPTIONS> ... <TARGETS>'
	@echo ''
	@echo 'Available targets are:'
	@echo ''
	@echo '    help               Show this help screen.'
	@echo '    gitready           Make repo ready to commit.'
	@echo '    debug-yggdrasil    Runs the Yggdrasil project in debug mode.'
	@echo '    debug-bifrost      Runs the Bifrost project in debug mode.'
	@echo '    clean              Remove binaries, artifacts and releases.'
	@echo '    clean-artifacts    Remove build artifacts only.'
	@echo '    clean-releases     Remove releases only.'
	@echo '    tools              Install tools needed by the project.'
	@echo '    deps               Download and install build time dependencies.'
	@echo '    test               Run unit tests.'
	@echo '    vet                Run go vet.'
	@echo '    errors             Run errcheck.'
	@echo '    lint               Run golint.'
	@echo '    imports            Run goimports.'
	@echo '    fmt                Run go fmt.'
	@echo '    env                Display Go environment.'
	@echo '    build              Build project for current platform.'
	@echo '    build-all          Build project for all supported platforms.'
	@echo '    docker-build       Build project using docker.'
	@echo '    doc                Start Go documentation server on port 8080.'
	@echo '    version            Check the Go version.'
	@echo '    swagger            Generate the Swagger documentation files.'
	@echo ''
	@echo 'Targets run by default are: imports, fmt, lint, vet, errors and build.'
	@echo ''

gitready: imports fmt lint vet

print-%:
	@echo $* = $($*)

clean: clean-artifacts clean-releases
	go clean -i ./...
	rm -vf $(CURDIR)/bin/$(BINARY)

clean-artifacts:
	rm -Rf bin/artifacts

clean-releases:
	rm -Rf bin/releases

tools:
	go get golang.org/x/tools/cmd/goimports
	go get github.com/kisielk/errcheck
	go get github.com/golang/lint/golint
	go get github.com/axw/gocov/gocov
	go get github.com/matm/gocov-html
	go get github.com/tools/godep
	go get github.com/mitchellh/gox

deps:
	go get -v ./...

test:
	go test -v ./...

vet:
	go vet -v ./...

errors:
	errcheck -ignoretests -blank ./...

lint:
	golint ./..

imports:
	goimports -l -w .

fmt:
	go fmt ./...

env:
	@go env

debug-yggdrasil:
	$(GORUN) -ldflags "$(LDFLAGS)" $(ENTRY_FILE_YGGDRASIL) -c configs/yggdrasil/config-local.yaml

debug-bifrost:
	$(GORUN) -ldflags "$(LDFLAGS)" $(ENTRY_FILE_BIFROST) -c configs/bifrost/config-local.yaml

build:
	$(GOBUILD) -ldflags "$(LDFLAGS)" -o bin/$(BINARY) $(ENTRY_FILE)

build-all:
	mkdir -v -p $(CURDIR)/bin/artifacts/$(VERSION)
	$(foreach GOOS, $(PLATFORMS),\
	$(foreach GOARCH, $(ARCHITECTURES), $(shell export GOOS=$(GOOS); export GOARCH=$(GOARCH); go build -v -o bin/artifacts/$(VERSION)/$(GOOS)_$(GOARCH)/$(BINARY) $(ENTRY_FILE))))

release: package-release sign-release

package-release:
	@test -x $(CURDIR)/bin/artifacts/$(VERSION) || exit 1
	mkdir -v -p $(CURDIR)/bin/releases/$(VERSION)
	for release in $$(find $(CURDIR)/bin/artifacts/$(VERSION) -mindepth 1 -maxdepth 1 -type d 2>/dev/null); do \
  		platform=$$(basename $$release); \
  		pushd $$release &>/dev/null; \
  		zip $(CURDIR)/bin/releases/$(VERSION)/$(BINARY)_$${platform}.zip $(BINARY); \
  		popd &>/dev/null; \
  	done

sign-release:
	@test -x $(CURDIR)/bin/releases/$(VERSION) || exit 1
	pushd $(CURDIR)/bin/releases/$(VERSION) &>/dev/null; \
	shasum -a 256 -b $(BINARY)_* > SHA256SUMS; \
	if test -n "$(GPG_SIGNING_KEY)"; then \
	  gpg --default-key $(GPG_SIGNING_KEY) -a \
      	      -o SHA256SUMS.sign -b SHA256SUMS; \
	fi; \
    popd &>/dev/null

doc:
	godoc -index

version:
	@go version

swagger:
	swag init -g cmd/yggdrasil/main.go --output ./docs/yggdrasil --exclude ./cmd/bifrost
	swag init -g cmd/bifrost/main.go --output ./docs/bifrost --exclude ./cmd/yggdrasil