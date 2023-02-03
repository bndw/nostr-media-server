REPO ?= bndw/nostr-media-server
GITSHA=$(shell git rev-parse --short HEAD)
TAG_COMMIT=$(REPO):$(GITSHA)
TAG_LATEST=$(REPO):latest

.PHONY: build
build:
	go build -o ./bin/nostr-media-server .

.PHONY: docker
docker:
	@docker build -t $(TAG_LATEST) .

# run runs the server with a config file
.PHONY: run
run: build
	./bin/nostr-media-server -config config.yml

run-docker:
	@docker run --rm -p 8080:80 $(TAG_LATEST)

.PHONY: publish
publish:
	docker push $(TAG_LATEST)
	@docker tag $(TAG_LATEST) $(TAG_COMMIT)
	docker push $(TAG_COMMIT)

.PHONY: test
test:
	go test -v ./...
