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
	./bin/gendata --names benchmarks/data/names.txt --output benchmarks/requests/register.txt --test-name register
	./bin/gendata --names benchmarks/data/names.txt --output benchmarks/requests/search.txt --test-name search

BENCH_N_CONN=50
BENCH_DURATION=5m
BENCH_TIMEOUT=1s
bench-register:
	wrk --latency --timeout $(BENCH_TIMEOUT) -d $(BENCH_DURATION) -t 3 -c $(BENCH_N_CONN) -s benchmarks/register.lua http://localhost

bench-search:
	wrk --latency --timeout $(BENCH_TIMEOUT) -d $(BENCH_DURATION) -t 3 -c $(BENCH_N_CONN) -s benchmarks/search.lua http://localhost
