BINARY_NAME  = fsex
PWD          = $(shell pwd)
BUILDER_IMG   = "golang:1.20.2"
UID = $(shell id -u)
GID = $(shell id -g)

HASH := $(shell git rev-parse --short HEAD)
COMMIT_DATE := $(shell git show -s --format=%ci ${HASH})
BUILD_DATE := $(shell date '+%Y%m%d%H%M%S')
VERSION := ${HASH}
CMD = nerdctl

TARGET_BIN = target/$(BINARY_NAME)


CONTAINER_IMAGE_NAME  = $(BINARY_NAME):$(HASH)

build:
	 $(CMD) run --rm -it -v $(PWD):/workdir -w /workdir -v $(PWD)/cache/go-build:/root/.cache/go-build -v $(PWD)/cache/mod:/go/pkg/mod  $(BUILDER_IMG) sh -c \
  	'cd src &&go build -buildvcs=false -ldflags="-s -w -X main.buildVersion=${VERSION} -X main.buildDate=${BUILD_DATE}" -o ../target/$(BINARY_NAME)  .'

test:
	 $(CMD) run --env-file=.env --rm -it -v $(PWD):/workdir -w /workdir  $(BUILDER_IMG) sh -c \
  	'cd src && go test -v'

clean:
	 $(CMD) run --rm -it -v $(PWD):/workdir -w /workdir  $(BUILDER_IMG) sh -c \
  	'cd src && go clean && rm -rf ../target'

shell:
	 $(CMD) run --rm -it -v $(PWD):/workdir -w /workdir  $(BUILDER_IMG) bash 	



.PHONY: build clean
