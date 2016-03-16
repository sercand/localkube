mkfile_dir := $(dir $(abspath $(lastword $(MAKEFILE_LIST))))

GO ?= CGO_ENABLED=0 go
GOFMT ?= gofmt -s
GOLINT ?= golint
GODEP ?= godep
DOCKER ?= docker

BUILDER ?=

GOFILES := find . -name '*.go' -not -path "./vendor/*"
GOOS := $(shell go list -f '{{context.GOOS}}')

GOBUILD_LDFLAGS ?=
GOBUILD_FLAGS ?= -a -installsuffix cgo

PKG ?= rsprd.com/localkube
EXEC_PKG := $(PKG)/cmd/localkube
DOCKER_DIR := /go/src/$(PKG)

MNT_DOCKER_SOCK ?= -v "/var/run/docker.sock:/var/run/docker.sock"
MNT_WEAVE_SOCK ?= -v "/var/run/weave.sock:/var/run/weave.sock"
MNT_REPO ?= -v "$(mkfile_dir):$(DOCKER_DIR)"

DOCKER_OPTS ?=
DOCKER_RUN_OPTS ?= $(MNT_DOCKER_SOCK)

# image data
ORG ?= ethernetdan
NAME ?= localkube
TAG ?= latest

DOCKER_IMAGE_NAME = "$(ORG)/$(NAME):$(TAG)"
DOCKER_DEV_IMAGE ?= "golang:1.6"

.PHONY: all
all: deps clean validate build build-image

.PHONY: validate
validate: checkgofmt

.PHONY: build
build: build/localkube-$(GOOS)

.PHONY: docker-build
docker-build: validate
	mkdir build
	$(DOCKER) run -w $(DOCKER_DIR) $(DOCKER_OPTS) $(MNT_REPO) $(DOCKER_DEV_IMAGE) make build

.PHONY: clean
clean:
	rm -rf ./build

build/localkube-$(GOOS):
	$(GO) build -o $@ $(GOBUILD_FLAGS) $(GOBUILD_LDFLAGS) $(EXEC_PKG)

build/localkube-linux:
	GOOS=linux $(GO) build -o $@ $(GOBUILD_FLAGS) $(GOBUILD_LDFLAGS) $(EXEC_PKG)

.PHONY: build-image
build-image: context
	$(DOCKER) build $(DOCKER_OPTS) -t $(DOCKER_IMAGE_NAME) ./build/context

.PHONY: run-image
run-image: build-image
	$(DOCKER) run -it $(DOCKER_OPTS) $(DOCKER_RUN_OPTS) $(DOCKER_IMAGE_NAME)

.PHONY: context
context: build/localkube-linux
	rm -rf ./build/context
	cp -r ./image ./build/context
	cp ./build/localkube-linux ./build/context
	chmod +x ./build/context/localkube-linux


.PHONY: checkgofmt
checkgofmt:
	# get all go files and run go fmt on them
	@files=$$($(GOFILES) | xargs $(GOFMT) -l); if [[ -n "$$files" ]]; then \
		  echo "Error: '$(GOFMT)' needs to be run on:"; \
		  echo "$${files}"; \
		  exit 1; \
		  fi;

.PHONY: deps
deps:
	go get github.com/tools/godep
	$(GODEP) restore -v
