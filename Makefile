GOOS ?= $(shell go env GOOS)
GOARCH ?= $(shell go env GOARCH)

it: build install

tidy:
	go mod tidy
	go mod vendor
clean:
	rm -rf dist || true

build:
	GOOS=$(GOOS) GOARCH=$(GOARCH) goreleaser build --snapshot --clean --single-target

install:
	install -Dm755 dist/docker-stackx-cli-plugin_$(GOOS)_$(GOARCH)_v8.0/docker-stackx-cli-plugin ${HOME}/.docker/cli-plugins/docker-stackx

cross-binaries: clean
	goreleaser release --snapshot --clean
