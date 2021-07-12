test:
	go test ./...

deps:
	go mod tidy
	go mod vendor

lint:
	golangci-lint run ./...

.PHONY:tools
tools:
	go generate ./tools

generate:
	PATH=$$PATH:./bin go generate ./internal/...

up:
	mkdir -p ./bin
	go build -o ./bin/site ./cmd/site
	./bin/site