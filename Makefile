LOCAL_BIN := $(CURDIR)/bin
GOLANGCI_BIN := $(LOCAL_BIN)/golangci-lint
GOFUMPT_BIN := $(LOCAL_BIN)/gofumpt
SWAG_BIN := $(LOCAL_BIN)/swag
GO_TEST=$(LOCAL_BIN)/gotest
GO_TEST_ARGS=-race -v ./...

all: generate lint test

lint:
	$(GOFUMPT_BIN) -l -w . || true
	$(GOLANGCI_BIN) run || true

test:
	$(GO_TEST) $(GO_TEST_ARGS)

generate: bin-deps .generate build


bin-deps: .bin-deps
.bin-deps: export GOBIN := $(LOCAL_BIN)
.bin-deps: .create-bin
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.64.8 && \
	go install github.com/rakyll/gotest@v0.0.6 && \
	go install go.uber.org/mock/mockgen@latest && \
	go install github.com/swaggo/swag/cmd/swag@latest
	mv $(LOCAL_BIN)/mockgen $(LOCAL_BIN)/mockgen_uber && \
	go install github.com/swaggo/swag/cmd/swag@v1.16.4 && \
	go install mvdan.cc/gofumpt@latest


.create-bin:
	rm -rf ./bin
	mkdir -p ./bin

.generate:
	$(info Generating code...)
	$(SWAG_BIN) init -g cmd/todo-list/main.go --parseDependency --parseInternal
	(PATH="$(PATH):$(LOCAL_BIN)" && go generate ./...)
	go mod tidy


build:
	go mod tidy
	go build -o ./bin/todo-list ./cmd/todo-list