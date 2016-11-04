BINARY := wacc_34

GOGLIDE := $(GOPATH)/bin/glide
GOLINTER := $(GOPATH)/bin/gometalinter.v1

SRC := $(shell find . -name '*.go' -not -path '*/vendor/*')

all: $(BINARY)

$(BINARY): $(SRC) vendor
	go build

vendor: $(GOGLIDE) glide.lock
	$(GOGLIDE) install

format:
	go fmt

lint: $(GOLINTER)
	$(GOLINTER) --exclude=vendor

install: $(BINARY)
	go install

clean:
	go clean

$(GOGLIDE):
	go get -u github.com/Masterminds/glide

$(GOLINTER):
	go get -u gopkg.in/alecthomas/gometalinter.v1
	$(GOLINTER) --install

.PHONY: all clean lint format
