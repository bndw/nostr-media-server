package main

import (
	"context"
	"io"

	"github.com/bndw/nostr-media-server/storage"
)

type storageProvider interface {
	Save(context.Context, io.Reader, storage.Options) (string, error)
	Get(ctx context.Context, sum, name string) (io.Reader, error)
}
