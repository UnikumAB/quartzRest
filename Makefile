default:        test

test:   golangci-lint
	go test -v -race ./...

fmt:
	gofmt -w .

golangci-lint:
ifeq (, $(shell which golangci-lint))
        curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.24.0
endif
	golangci-lint run --fix ./...

build: golangci-lint test
	go build ./...
	GOOS=darwin go build -o bin/quartzRestServer-darwin-x86_64 ./cmd/quartzRestServer
	GOOS=linux go build -o bin/quartzRestServer-linux-x86_64 ./cmd/quartzRestServer

mod:
	go mod tidy

all: fmt mod test

.PHONY: imports test fmt mod docker all default
