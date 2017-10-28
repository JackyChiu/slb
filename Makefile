all: bin build

bin:
	mkdir -p bin

clean:
	rm -rf bin

build: bin/tlb

bin/tlb: bin
	go build -o bin/tlb ./cmd/tlb

run: bin/tlb
	bin/tlb

.PHONY: all
