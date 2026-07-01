.PHONY: run db-up db-down test

run:
	go run cmd/server/main.go

db-up:
	docker-compose up -d

db-down:
	docker-compose down

test:
	go test -v ./...
