BINARY := wacc_34

SRC := $(shell find . -name '*.go' -not -path '*/vendor/*')

all: $(BINARY)

$(BINARY): $(SRC) glide.lock
	go build

vendor:
	glide install

clean:
	go clean

format:
	go fmt

lint:
	gometalinter.v1 --exclude=vendor

.PHONY: all clean vendor lint format
