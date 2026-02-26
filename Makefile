.PHONY: run build install test fmt vet lint clean

# Default target
run:
	go run .

build:
	go build -o p .

install:
	go install .

test:
	go test ./...

fmt:
	go fmt ./...

vet:
	go vet ./...

lint:
	golangci-lint run

clean:
	rm -f p
