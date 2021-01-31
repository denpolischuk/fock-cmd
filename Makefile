BIN_NAME=fock

.PHONY: build
build:
	go build -v -o ${BIN_NAME} ./cmd/fock-cli

.PHONY: test
test:
	go test -v -race -timeout 30s ./...

.DEFAULT_GOAL := build
