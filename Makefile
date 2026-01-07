.PHONY: build test lint clean example

VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
LDFLAGS := -ldflags "-X main.version=$(VERSION)"

build:
	go build $(LDFLAGS) -o bin/prophet ./cmd/prophet

test:
	go test -race -cover ./...

lint:
	golangci-lint run ./...

clean:
	rm -rf bin/

example: build
	./bin/prophet analyze /Users/cblevins/workspace/platform/gitops/k3s/ai/vllm/

whatif-example: build
	./bin/prophet whatif --scale vllm-deployment=2

fmt:
	go fmt ./...

tidy:
	go mod tidy
