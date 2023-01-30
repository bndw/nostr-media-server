CONFIG_PATH=config.yml
API_PATH=http://localhost:9000/upload
MEDIA_PATH=http://localhost:9000
ACCEPTED_MIMETYPES=image/jpg,image/png,image/gif
STORAGE_TYPE=filesystem
STORAGE_CONFIG="media_dir:./files"
NAMES=alice:npub1xxx,bob:npub1yyy

build:
	go build -o ./bin/nostr-media-server .

# run runs the server with a config file
run: build
	./bin/nostr-media-server -config config.yml

# run-env runs the server with config from the environment
run-env: build
	API_PATH=$(API_PATH) \
	MEDIA_PATH=$(MEDIA_PATH) \
	ACCEPTED_MIMETYPES=$(ACCEPTED_MIMETYPES) \
	STORAGE_TYPE=$(STORAGE_TYPE) \
	STORAGE_CONFIG=$(STORAGE_CONFIG) \
	NAMES=$(NAMES) \
	./bin/nostr-media-server
