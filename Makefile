.PHONY: lint test gen

lint:
	golangci-lint run -c .golangci.yaml

test:
	go test -v ./...

gen:
	go generate ./...
