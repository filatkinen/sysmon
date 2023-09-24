BINSERVICE := "./build/sysmon"
BINCLIENT := "./build/client"

GIT_HASH := $(shell git log --format="%h" -n 1)
LDFLAGS := -X main.release="develop" -X main.buildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S) -X main.gitHash=$(GIT_HASH)

build:
	go build -v -o $(BINSERVICE) -ldflags "$(LDFLAGS)" ./cmd/service
	go build -v -o $(BINCLIENT)  ./cmd/client

run: build
	$(BINSERVICE) -config ./configs/service.yaml

run-sudo: build
	sudo $(BINSERVICE) -config ./configs/service.yaml


run-client: build
	$(BINCLIENT)

build-img:
	docker-compose -f  deployments/docker-compose.yaml build

run-img: build-img
	docker-compose -f deployments/docker-compose.yaml up -d

down-img:
	docker-compose -f deployments/docker-compose.yaml down \
		 --rmi local \
		--volumes \
		--remove-orphans \
		--timeout 60; \



test:
	go test -race -v --count=2 ./...

test-integration:
	go test -race -v ./... -tags integration

#generate-proto:
#    protoc -I  ./internal/grpc/ ./internal/grpc/sysmon.proto   --go_out=./internal/grpc/
#    protoc -I  ./internal/grpc/ ./internal/grpc/sysmon.proto   --go-grpc_out=require_unimplemented_servers=false:./internal/grpc/
#	apt install  protobuf-compiler
#    go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28
#    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2


install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.54.2

lint: install-lint-deps
	golangci-lint run ./...


.PHONY: build test test-integration lint
