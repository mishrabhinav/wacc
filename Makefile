BINARY := wacc_34

SRC := $(shell find . -name '*.go' -not -path '*/vendor/*')

all: $(BINARY)

$(BINARY): $(SRC) glide.lock .vendor
	go build

vendor: .vendor

.vendor:
	glide install
	touch .vendor

clean:
	go clean

format:
	go fmt

lint:
	gometalinter.v1 --exclude=vendor

install: $(BINARY)
	go install

.PHONY: all clean vendor lint format
