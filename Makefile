COLORIZE_PASS=sed ''/PASS/s//$$(printf "$(GREEN)PASS$(RESET)")/''
COLORIZE_FAIL=sed ''/FAIL/s//$$(printf "$(RED)FAIL$(RESET)")/''
SHELL=/bin/bash
TAG := $(shell git rev-parse HEAD)

.PHONY: build bin-linux-amd64 test

build:
	go build -a -tags "netgo" -installsuffix netgo -ldflags="-s -w -X main.version=$(TAG) -extldflags \"-static\"" -o build/bin/ ./...

bin-linux-amd64:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -tags "netgo" -installsuffix netgo -ldflags="-s -w -X main.version=$(TAG) -extldflags \"-static\"" -o build/bin/ ./...

test:
	gofmt -l .
	go vet -v ./...
	staticcheck ./...
	go test -v ./...  | $(COLORIZE_PASS) | $(COLORIZE_FAIL)
