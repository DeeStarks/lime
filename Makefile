.DEFAULT_GOAL := setup

setup:
	@make install && make build

install:
	@go get ./...

build:
	@go build -o ./lime lime.go
