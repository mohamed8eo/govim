BINARY = govim
FILE = test.txt
ARG ?= $(FILE)

.PHONY: all build run test fmt vet clean

all: clean build test

build:
	@echo "Building $(BINARY)..."
	go build -o $(BINARY) ./cmd/govim

run: build
	@echo "Running $(BINARY)..."
	./$(BINARY) $(ARG)

test:
	@echo "Running unit tests..."
	go test -v ./...

fmt:
	go fmt ./...

vet:
	go vet ./...

clean:
	@echo "Cleaning up..."
	@rm -f $(BINARY)
