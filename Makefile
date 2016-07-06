default: build

.PHONY: build

build: fmt vet
	go build -o ./build/security-scan

install:
	go install github.com/onetwopunch/security-scan

fmt:
	go fmt

vet:
	go vet
