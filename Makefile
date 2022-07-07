.PHONY: lint test

lint:
	golangci-lint run -c .golangci.yaml

test:
	go test ./...
