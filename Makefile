it: build install

clean:
	rm -rf bin || true

build:
	go build -o bin/docker-stackx .

install:
	install -Dm755 bin/docker-stackx ${HOME}/.docker/cli-plugins/docker-stackx
