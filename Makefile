.PHONY: test quality fmt build run

test:
	go test ./...

quality:
	golangci-lint run

fmt:
	gofumpt -w .

build:
	go build ./...

serve:
	go run ./cmd/server/main.go
