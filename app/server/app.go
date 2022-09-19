package server

import (
	"context"
	"flag"
	"fmt"
	"github.com/aaronland/go-http-server"
	"github.com/sfomuseum/go-flags/flagset"
	"github.com/sfomuseum/go-http-auth"
	"github.com/sfomuseum/go-url-shortener/database"
	"github.com/sfomuseum/go-url-shortener/http/api"
	"log"
	"net/http"
)

func Run(ctx context.Context, logger *log.Logger) error {
	fs := DefaultFlagSet(ctx)
	return RunWithFlagSet(ctx, fs, logger)
}

func RunWithFlagSet(ctx context.Context, fs *flag.FlagSet, logger *log.Logger) error {

	flagset.Parse(fs)

	db, err := database.NewDatabase(ctx, database_uri)

	if err != nil {
		return fmt.Errorf("Failed to open docstore, %w", err)
	}

	err = db.SetLogger(ctx, logger)

	if err != nil {
		return fmt.Errorf("Failed to assign logger to database, %w", err)
	}

	authenticator, err := auth.NewAuthenticator(ctx, authenticator_uri)

	if err != nil {
		return fmt.Errorf("Failed to create authenticator, %w", err)
	}

	mux := http.NewServeMux()

	add_opts := &api.AddURIHandlerOptions{
		Authenticator: authenticator,
		Logger:        logger,
		Database:      db,
	}

	add_handler := api.AddURIHandler(add_opts)
	add_handler = authenticator.WrapHandler(add_handler)

	mux.Handle("/add", add_handler)

	resolve_opts := &api.ResolveURIHandlerOptions{
		Logger:   logger,
		Database: db,
	}

	resolve_handler := api.ResolveURIHandler(resolve_opts)

	mux.Handle("/", resolve_handler)

	s, err := server.NewServer(ctx, server_uri)

	if err != nil {
		return fmt.Errorf("Failed to create new server, %w", err)
	}

	logger.Printf("Listening for requests on %s", s.Address())

	err = s.ListenAndServe(ctx, mux)

	if err != nil {
		return fmt.Errorf("Failed to serve requests, %w", err)
	}

	return nil
}
