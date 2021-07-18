install: deps tools

deps:
	go mod tidy
	go mod vendor

.PHONY:tools
tools:
	go generate ./tools

test:
	SKIP_DB_TEST=1 go test ./...

test-full:
	go test ./...

lint:
	golangci-lint run ./...

generate:
	 PATH=$$PATH:./bin go generate ./internal/...

up:
	mkdir -p ./bin
	go build -o ./bin/site ./cmd/site
	./bin/site

db:
	docker-compose down
	docker-compose up -d db

migrate:
	docker-compose up migrate

env:
	make db
	sleep 30
	make migrate

