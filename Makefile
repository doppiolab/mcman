.PHONY: build-linux
## build-linux: build the project for linux
build-linux:
	GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o bin/mcman.linux.amd64 main.go
	GOOS=linux GOARCH=arm64 go build -ldflags="-w -s" -o bin/mcman.linux.arm64 main.go

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
	go test -race -coverprofile=coverage.out -covermode=atomic -v ./...

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
