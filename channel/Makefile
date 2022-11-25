generate:
	go generate ./...

tidy:
	go mod tidy

lint: tidy
	golangci-lint run

vet: tidy
	go vet ./...

check: lint vet test
	echo "check done"

test: tidy
	go test ./...

coverage: tidy
	go test -json -covermode=count -coverprofile=profile.cov ./... > report.json
