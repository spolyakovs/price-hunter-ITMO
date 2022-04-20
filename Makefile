CONFIG_PATH?=configs/local.toml

.PHONY: build
build:
	go build -v ./cmd/apiserver

.PHONY: test
test:
	make build
	go test -v -race -timeout 30s ./...

.PHONY: start
start:
	make build
	./apiserver -config-path=$(CONFIG_PATH)

.DEFAULT_GOAL := start
