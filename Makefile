build:
	go build -v -o ./bin/api ./cmd/api

run:
	go run ./cmd/api/main.go

test:
	go test -v -cover -race -timeout 30s ./...

.PHONY: build run test
