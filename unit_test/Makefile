MODULE = $(shell go list -m)
VERSION ?= $(shell git describe --tags --always --dirty --match=v* 2> /dev/null || echo "1.0.0")
LDFLAGS := -ldflags "-X main.Version=${VERSION}"
CMD_NAME = "unit_test"
DOCKER_PACKAGE_NAME="docker.io/mehmetcantas/unit_test"
DOCKER_TAG_VERSION="123123"

.PHONY: default
default: help

.PHONY: build
build:  ## build the source code
	CGO_ENABLED=0 go build ${LDFLAGS} -a -o $(CMD_NAME) $(MODULE)/

.PHONY: test
test:  ## run all the tests
	go test -json -covermode=atomic -coverpkg=../unit_test/...  ./...

.PHONY: docker-build
docker-build: ## build the docker image
	docker build -f Dockerfile -t $(DOCKER_PACKAGE_NAME):$(DOCKER_TAG_VERSION) .

.PHONY: version
version: ## display the version of the cli
	@echo $(VERSION)