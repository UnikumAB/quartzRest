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

release-build: golangci-lint test
	go build ./...
	GOOS=darwin go build -o bin/quartzRestServer-darwin-x86_64 ./cmd/quartzRestServer
	GOOS=linux go build -o bin/quartzRestServer-linux-x86_64 ./cmd/quartzRestServer

build: golangci-lint test
	go build ./cmd/quartzRestServer

mod:
	go mod tidy

all: fmt mod test

docker:
	docker build -t quartzrestserver .

.PHONY: imports test fmt mod docker all default release-build docker

release:
ifeq (, $(shell which goreleaser))
        curl -sfL https://install.goreleaser.com/github.com/goreleaser/goreleaser.sh | sh
endif
	goreleaser --snapshot --skip-publish --rm-dist

