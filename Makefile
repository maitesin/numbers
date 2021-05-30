tools:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

test:
	go test -race -cover -v ./...

lint:
	golangci-lint run
