# nostr-media-server

Implements the proposed [nostr-media-spec](https://github.com/michaelhall923/nostr-media-spec).

This server handles both uploads and downloads of media content. By default 
the uploaded media is stored on the local filesystem, but the code is written 
with a pluggabe storage interface to allow additional storage providers.

### Status

- [x] Getting host info at https://<hostname>/.well-known/nostr.json
- [x] Uploading media
- [x] Getting media
- [x] Getting media: filename optional
- [x] Getting media: type optional
- [ ] Getting media: potential query params (size, aspect ratio)

### Quickstart

1. In one terminal, start the server on port 9000:
  ```
  make run
  ```

2. In a second terminal, upload a test image:
  ```
  cd curls && ./upload.sh
  ```

See `config.yml` for configuration options.
