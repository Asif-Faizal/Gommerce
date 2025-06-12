build:
	@go build -o bin/gommerce cmd/main.go

test:
	@go test -v ./...

run:
	@go run cmd/main.go