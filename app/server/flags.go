package server

import (
	"context"
	"flag"
	"github.com/sfomuseum/go-flags/flagset"
)

var server_uri string

var database_uri string

var authenticator_uri string

func DefaultFlagSet(ctx context.Context) *flag.FlagSet {

	fs := flagset.NewFlagSet("server")

	fs.StringVar(&server_uri, "server-uri", "http://localhost:8080", "")

	fs.StringVar(&database_uri, "database-uri", "mem://urls/Source", "")

	fs.StringVar(&authenticator_uri, "authenticator-uri", "null://", "")

	return fs
}
