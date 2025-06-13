build:
	@go build -o bin/gommerce cmd/main.go

test:
	@go test -v ./...

run:
	@go run cmd/main.go

# Migration commands
migration-create:
	@migrate create -ext sql -dir cmd/migrate/migrations -seq $(filter-out $@,$(MAKECMDGOALS))

migrate-up:
	@go run cmd/migrate/main.go up

migrate-down:
	@go run cmd/migrate/main.go down

# Allow passing arguments to migration-create
%:
	@: