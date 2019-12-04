build:
	go build -v -o ./bin/api ./cmd/api

run:
	go run ./cmd/api/main.go

test:
	go test -coverprofile cover.out -coverpkg ./internal/album/... ./internal/artist/... ./... && go tool cover -func cover.out

# Please, note:
#     for testing setting avatars you should create folders:
#         2019_2_Covenant/resources/avatars
#     and add any img or png file in last folder

.PHONY: build run test
