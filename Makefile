all: build install

build:
	go build

install: edisco-test
	go install

.PHONY: clean

clean:
	rm edisco-test cmd/ingestEmail/input/*.eml cmd/ingestEmail/jsonl/eml.jsonl