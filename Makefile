.PHONY: test quality quality-fix fmt build run

test:
	go test ./...

quality:
	golangci-lint run

quality-fix:
	golangci-lint run --fix

fmt:
	gofumpt -w .

build:
	go build ./...

run:
	go run ./cmd/server/main.go
