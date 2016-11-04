BINARY := wacc_34

GOGLIDE := $(GOPATH)/bin/glide
GOLINTER := $(GOPATH)/bin/gometalinter.v1
GOPEG := $(GOPATH)/bin/peg

SRC := $(shell find . -name '*.go' -not -path '*/vendor/*')
GRM := $(shell find . -name '*.peg' -not -path '*/vendor/*')
SRC += $(patsubst %.peg,%.peg.go,$(GRM))

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

$(patsubst %.peg,%.peg.go,$(GRM)): $(GRM) $(GOPEG)
	$(GOPEG) $(patsubst %.go,%,$@)

$(GOGLIDE):
	go get -u github.com/Masterminds/glide

$(GOLINTER):
	go get -u gopkg.in/alecthomas/gometalinter.v1
	$(GOLINTER) --install

$(GOPEG):
	go get -u github.com/pointlander/peg

.PHONY: all clean lint format
