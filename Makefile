generate:
	go generate ./...

lint:
	# Coding style static check.
	@go install honnef.co/go/tools/cmd/staticcheck@latest
	@go mod tidy
	staticcheck `go list ./...`

vet:
	@echo "⚠️ Warning: the following only works with go >= 1.14" && \
	go install ./internal/mcheck && \
	go vet -vettool=`go env GOPATH`/bin/mcheck -commentLen -ifInit ./...

# target to run all the possible checks; it's a good habit to run it before
# pushing code
check: lint vet
	go test ./...
