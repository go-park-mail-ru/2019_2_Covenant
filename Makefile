build:
	go build -v -o ./bin/api ./cmd/api

run:
	go run ./cmd/api/main.go

test:
	go test -v -cover -race -timeout 30s ./...
Please, note:
    for testing setting avatars you should create folders:
        2019_2_Covenant/resources/avatars
    and add any img or png file in last folder

.PHONY: build run test
