all: bin build

bin:
	mkdir -p bin

clean:
	rm -rf bin

build: clean
	go build ./cmd/...

bin/tlb: bin
	go build -o bin/tlb ./cmd/tlb

bin/example: bin
	go build -o bin/example ./cmd/example

tlb: bin/tlb
	bin/tlb

example: bin/example
	bin/example

test:
	go test ./...

.PHONY: all
