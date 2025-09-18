VERSION := $(shell git describe --tags --abbrev=0 2>/dev/null || echo dev)
LDFLAGS := -ldflags "-X 'main.version=$(VERSION)'"

.PHONY: build
build:
	CGO_ENABLED=0 GOOS=linux go build $(LDFLAGS) -o . ./cmd/...

.PHONY: update
update:
	go get -u -t ./...
	go mod tidy
	go mod vendor

.PHONY: docs
docs:
	goa gen github.com/neatflowcv/key-stone/design

.PHONY: lint
lint:
	go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.4.0
	golangci-lint run --fix

.PHONY: test
test:
	go test -race -shuffle=on ./...

.PHONY: cover
cover:
	go test ./... --coverpkg ./... -coverprofile=c.out
	go tool cover -html="c.out"
	rm c.out