.PHONY: build

build:
	go build -v -o fock ./cmd/fock-cmd

.PHONY: test
test:
	go test -v -race -timeout 30s ./...

.DEFAULT_GOAL := build
