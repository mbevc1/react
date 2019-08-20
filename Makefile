.PHONY: help
.DEFAULT: help
ifndef VERBOSE
.SILENT:
endif

NO_COLOR=\033[0m
GREEN=\033[32;01m
YELLOW=\033[33;01m
RED=\033[31;01m

VER?=dev
GHASH:=$(shell git rev-parse --short HEAD)
GO:=            go
GO_BUILD:=      go build -mod vendor -ldflags "-s -w -X main.GitCommit=${GHASH}"
GO_VENDOR:=     go mod vendor
BIN:=           react

help:: ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[33m%-20s\033[0m %s\n", $$1, $$2}'

$(BIN): vendor ## Produce binary
	GO111MODULE=on $(GO_BUILD)

# We always want vendor to run
.PHONY: vendor
vendor: **/*.go ## Build vendor deps
	GO111MODULE=on $(GO_VENDOR)

clean: ## Clean artefacts
	rm -rf $(BIN) $(BIN)_* $(BIN).exe

clean-vendor: ## Clean vendor folder
	rm -rf vendor

build: clean $(BIN) ## Build binary
	upx $(BIN)

run:
	go run .

dev: clean ## Dev test target
	go build -ldflags "-s -w -X main.Version=${VER}" -o $(BIN)
	upx $(BIN)

test: vendor ## Run tests
	go test -v ./...

fmt: **/*.go ## Formt Golang code
	go fmt ./...

lint:
	golint $(BIN)

vet:
	go vet --all $(BIN) ./...

$(BIN)_linux_amd64: vendor **/*.go
	GOOS=linux GOARCH=amd64 CGO_ENABLED=1 go build -o $@ *.go
	upx $@

$(BIN)_linux_alpine: vendor **/*.go
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o $@ *.go
	upx $@

$(BIN)_darwin_amd64: vendor **/*.go
	GOOS=darwin go build -o $@ *.go
	upx $@

$(BIN)_windows_amd64.exe: vendor **/*.go
	GOOS=windows GOARCH=amd64 go build -o $@ *.go
	upx $@

pack: $(BIN)_linux_amd64 $(BIN)_darwin_amd64 $(BIN)_windows_amd64.exe
	zip $(BIN)_linux_amd64.zip $(BIN)_linux_amd64
	zip $(BIN)_darwin_amd64.zip $(BIN)_darwin_amd64
	zip $(BIN)_windows_amd64.zip $(BIN)_windows_amd64.exe

fmtcheck: vendor **/*.go ## Check formatting
	@gofmt_files=$$(gofmt -l `find . -name '*.go' | grep -v vendor`); \
	    if [ -n "$${gofmt_files}" ]; then \
	    	echo 'gofmt needs running on the following files:'; \
	    	echo "$${gofmt_files}"; \
	    	echo "You can use the command: \`make fmt\` to reformat code."; \
	    	exit 1; \
	    fi; \
	    exit 0
