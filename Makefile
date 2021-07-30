install: deps tools
	@ls .env 2> /dev/null || cp .env.sample .env

env:
	make db
	sleep 30
	make migrate

up: build
	./bin/site

up-docker:
	docker-compose up service

check: generate test lint fmt

deps:
	go mod tidy
	go mod vendor

.PHONY:tools
tools:
	go generate ./tools

test:
	SKIP_DB_TEST=1 go test -v -coverprofile="coverage.txt" -covermode=atomic -race -count 1 -timeout 20s ./...

test-full:
	go test -coverprofile="coverage.txt" ./...

lint:
	golangci-lint run ./...

fmt:
	go fmt ./cmd/... ./internal/...
	goimports -ungroup -local github.com/vvkh/social-network -w ./cmd ./internal

generate:
	 PATH=$$PATH:./bin go generate ./internal/...

build:
	mkdir -p ./bin
	go build -o ./bin/site ./cmd/site

db:
	docker-compose down
	docker-compose up -d db

migrate:
	docker-compose up migrate


