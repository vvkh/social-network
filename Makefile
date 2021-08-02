install: deps tools
	@ls .env 2> /dev/null || cp .env.sample .env

env:
	make db
	sleep 30
	make migrate

stop-env:
	docker-compose stop

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
	go build -o ./bin/gendata ./cmd/gendata

db:
	docker-compose down
	docker-compose up -d db

migrate:
	docker-compose up migrate

gen-bench: build
	./bin/gendata --names benchmarks/data/names.txt --output benchmarks/data/requests.txt

bench:
	wrk --latency -d 60s -t 6 -c 6 -s benchmarks/register.lua http://localhost
