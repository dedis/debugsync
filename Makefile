generate:
	make -C channel generate
	make -C sync generate

tidy:
	make -C channel tidy
	make -C sync tidy

lint:
	# Coding style static check.
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.54.0
	make -C channel lint
	make -C sync lint

vet:
	@echo "⚠️ Warning: the following only works with go >= 1.14"
	make -C channel vet
	make -C sync vet

check:
# target to run all the possible checks; it's a good habit to run it before
# pushing code
	make -C channel check
	make -C sync check

test:
	make -C channel test
	make -C sync test

coverage:
	make -C channel coverage
	make -C sync coverage
