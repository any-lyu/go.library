# goimports
GOFMT_PRIVATE="github.com/any-lyu/go.library"
GOFMT_LOCAL="github.com/any-lyu/go.library"
GOFMT_FILES=$(shell find . -name '*.go' | grep -v \.pb\.go$ | xargs)

.PHONY: all
all: fmt lint build vet modtidy errcheck golangci-lint golang-lint

.PHONY: fmt
fmt:
	@goimports -l -w -private "${GOFMT_PRIVATE}" -local "${GOFMT_LOCAL}" ${GOFMT_FILES}

.PHONY: lint
lint:
	@golint ./...

.PHONY: build
build:
	@go build -mod=mod ./...

.PHONY: vet
vet:
	@go vet -mod=mod ./...

.PHONY: modtidy
modtidy:
	@go mod tidy

.PHONY: errcheck
errcheck:
	@errcheck -blank -ignoregenerated -ignoretests -exclude=.errcheck-excludes.txt ./...

.PHONY: golangci-lint
golangci-lint:
	@golangci-lint run --config=.golangci-lint.yml

.PHONY: golang-lint
golang-lint:
	@golang-lint ./...
