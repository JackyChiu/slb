all: clean test install

clean:
	rm -rf bin

build:
	go build ./...

install: 
	GOBIN=$(shell pwd)/bin go install ./cmd/...

test:
	go test ./...

.PHONY: all
