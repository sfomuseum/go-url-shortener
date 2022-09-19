package database

import (
	"context"
	"fmt"
	"github.com/sfomuseum/go-url-shortener/uri"
	"gocloud.dev/docstore"
	"io"
	"log"
	"time"
)

func init() {
	ctx := context.Background()

	for _, scheme := range docstore.DefaultURLMux().CollectionSchemes() {
		RegisterDatabase(ctx, scheme, NewDocstoreDatabase)
	}
}

type DocstoreDatabase struct {
	Database
	col    *docstore.Collection
	logger *log.Logger
}

func NewDocstoreDatabase(ctx context.Context, uri string) (Database, error) {

	col, err := docstore.OpenCollection(ctx, uri)

	if err != nil {
		return nil, fmt.Errorf("Failed to create collection, %w", err)
	}

	logger := log.New(io.Discard, "", 0)

	db := &DocstoreDatabase{
		col:    col,
		logger: logger,
	}

	return db, nil
}

func (db *DocstoreDatabase) AddURI(ctx context.Context, u *uri.URI) error {

	now := time.Now()
	u.Created = now.Unix()

	return db.col.Put(ctx, u)
}

func (db *DocstoreDatabase) GetURIWithShortURL(ctx context.Context, short_url string) (*uri.URI, error) {

	iter := db.col.Query().Where("Short", "=", short_url).Get(ctx)
	defer iter.Stop()

	u := &uri.URI{}
	err := iter.Next(ctx, u)

	if err != nil {
		return nil, err
	}

	return u, nil
}

func (db *DocstoreDatabase) GetURIWithSourceURL(ctx context.Context, source_url string) (*uri.URI, error) {

	iter := db.col.Query().Where("Source", "=", source_url).Get(ctx)
	defer iter.Stop()

	u := &uri.URI{}
	err := iter.Next(ctx, u)

	if err != nil {
		return nil, err
	}

	return u, nil
}

func (db *DocstoreDatabase) SetLogger(ctx context.Context, logger *log.Logger) error {
	db.logger = logger
	return nil
}
