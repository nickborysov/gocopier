.PHONY: lint test

all:
	make lint
	make tests

lint:
	golangci-lint run --config .golangci.yml

tests:
	go test -timeout 30s -cover ./...
