tools:
	go install github.com/matryer/moq@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

generate:
	go generate ./...

test: generate
	go test -race -cover -v ./...

lint: generate
	golangci-lint run

run:
	cd cmd/numbers && go run main.go

send-numbers-for-30-seconds:
	cd cmd/generator && go run main.go | nc localhost 4000 &
	sleep 21
	echo "terminate" | nc localhost 4000
