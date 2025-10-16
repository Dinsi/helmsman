build:
	go build -o bin/helmsman cmd/server/main.go

dep:
	glide install

docker:
	docker run --rm -e GO111MODULE=off -v $(shell pwd)/bin:/go/src/github.com/andrepinto/helmsman/bin $(shell docker build -f Dockerfile.build --no-cache -q .) go build -o bin/helmsman cmd/server/main.go
	docker build -f Dockerfile.dist -t andrepinto/helmsman:1.1.1 .

.PHONY: build dep docker