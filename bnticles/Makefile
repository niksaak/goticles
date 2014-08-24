.PHONY: all build test benchmark gofmt
all: build test

build:
	go build -gcflags '-s'

test:
	go test -v

benchmark:
	go test -v -bench . -benchmem -benchtime 10s

gofmt:
	go fmt
