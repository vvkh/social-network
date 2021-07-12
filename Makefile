install: deps tools

deps:
	go mod tidy
	go mod vendor

.PHONY:tools
tools:
	go generate ./tools

test:
	go test ./...

lint:
	golangci-lint run ./...

generate:
	 PATH=$$PATH:./bin go generate ./internal/...

up:
	mkdir -p ./bin
	go build -o ./bin/site ./cmd/site
	./bin/site