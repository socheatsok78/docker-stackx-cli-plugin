it: build install

clean:
	rm -rf bin || true

build:
	go build -o bin/docker-stackx .

install:
	install -Dm755 bin/docker-stackx ${HOME}/.docker/cli-plugins/docker-stackx

cross-binaries: clean
	for os in darwin linux; do \
		for arch in amd64 arm64; do \
			GOOS=$$os GOARCH=$$arch go build -o bin/docker-stackx-$$os-$$arch .; \
		done; \
	done
