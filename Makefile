run: build
	bin/challenge-02/main.exe

build:
	go build -o bin/challenge-02/main.exe ./cmd/challenge-02/main.go

test:
	go test -v -cover -short ./...

sqlc:
	sqlc generate -f ./internal/postgresql/sqlc.yaml

.PHONY: run build test