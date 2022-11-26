BUILD_TARGET=cotton

GO_FILES:=$(shell find . -type f -name '*.go' -print)

$(BUILD_TARGET): $(GO_FILES)
	CGO_ENABLED=0 go build -o $@ -v

.PHONY: test
test:
	go test -v -timeout 10s ./...
