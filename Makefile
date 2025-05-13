.PHONY: test fmt build run

test:
	go test ./...

fmt:
	gofumpt -w .

build:
	go build ./...

serve:
	go run ./cmd/server/main.go
