package main

import (
	_ "gocloud.dev/docstore/memdocstore"
)

import (
	"context"
	"github.com/sfomuseum/go-url-shortener/app/server"
	"log"
)

func main() {

	ctx := context.Background()
	logger := log.Default()

	err := server.Run(ctx, logger)

	if err != nil {
		logger.Fatalf("Failed to run server, %w", err)
	}
}
