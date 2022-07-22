.PHONY: lint
## lint: lint the source code
lint:
	golangci-lint run -c .golangci.yaml

.PHONY: test
## test: run all tests
test:
	go test -race -cover -v ./...

.PHONY: cover
## cover: run all tests
cover:
	go test -race -coverprofile=coverage.out -v ./...

.PHONY: gen
## gen: run go generate
gen:
	@go install github.com/golang/mock/gomock
	go generate -v ./...

.PHONY: help
## help: prints this help message
help:
	@echo "Usage: \n"
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':'
