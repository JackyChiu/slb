all: bin build

bin:
	mkdir -p bin

clean:
	rm -rf bin

build: clean bin/tlb

bin/tlb: bin
	go build -o bin/tlb ./cmd/tlb

run: build
	bin/tlb

test:
	go test -v -race ./...

.PHONY: all
