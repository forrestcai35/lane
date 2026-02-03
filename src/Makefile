.PHONY: build install clean run test

BINARY_NAME=lane
VERSION=0.1.0
LDFLAGS=-ldflags "-X github.com/forrestcai35/lane/cmd.Version=$(VERSION)"

# Build the binary
build:
	go build $(LDFLAGS) -o bin/$(BINARY_NAME) .

# Install to GOPATH/bin
install:
	go install $(LDFLAGS) .

# Clean build artifacts
clean:
	rm -rf bin/

# Run with arguments (usage: make run ARGS="500 --client Apple --desc Test")
run:
	go run . $(ARGS)

# Run tests
test:
	go test -v ./...

# Download dependencies
deps:
	go mod tidy

# Build for multiple platforms
release:
	GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o bin/$(BINARY_NAME)-darwin-amd64 .
	GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) -o bin/$(BINARY_NAME)-darwin-arm64 .
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o bin/$(BINARY_NAME)-linux-amd64 .
	GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o bin/$(BINARY_NAME)-windows-amd64.exe .
