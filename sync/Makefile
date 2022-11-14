generate:
	go generate ./...

tidy:
	go mod tidy

lint: tidy
	golangci-lint run

vet: tidy
	go vet ./...

check: lint vet
	go test ./...
