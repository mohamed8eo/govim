BINARY = govim
FILE = test.txt
ARG ?= $(FILE)

.PHONY: all build run test fmt vet clean

all: clean build

build:
	@echo "Building $(BINARY)..."
	go build -o $(BINARY) ./cmd/govim

run: build
	@echo "Running $(BINARY)..."
	./$(BINARY) $(ARG)

test: build
	@echo "Testing $(BINARY) with $(FILE)..."
	./$(BINARY) $(FILE)

fmt:
	go fmt ./...

vet:
	go vet ./...

clean:
	@echo "Cleaning up..."
	@rm -f $(BINARY)
