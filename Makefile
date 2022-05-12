.DEFAULT_GOAL := install

install:
	@go get ./...

build:
	@go build -o ./bin/lime ./cmd/main.go

lime:
	@make build && ./bin/lime
