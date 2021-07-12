test:
	go test ./...

deps:
	go mod tidy
	go mod vendor

lint:
	golangci-lint run ./...

make up:
	mkdir -p ./bin
	go build -o ./bin/site ./cmd/site
	./bin/site