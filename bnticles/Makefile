.PHONY: build test gofmt
build:
	go build -gcflags '-s'

test:
	go test -v

benchmark:
	go test -bench . -benchmem -benchtime 10s

gofmt:
	go fmt
