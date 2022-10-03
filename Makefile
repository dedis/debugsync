generate:
	go generate ./...

lint:
	# Coding style static check.
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.49.0
	@go mod tidy
	golangci-lint run

vet:
	@echo "⚠️ Warning: the following only works with go >= 1.14" && \
	go vet ./...

# target to run all the possible checks; it's a good habit to run it before
# pushing code
check: lint vet
	go test ./...
