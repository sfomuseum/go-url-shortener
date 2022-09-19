package api

import (
	"github.com/sfomuseum/go-url-shortener/database"
	"io"
	"log"
	"net/http"
	"path/filepath"
)

type ResolveURIHandlerOptions struct {
	Database database.Database
	Logger   *log.Logger
}

func ResolveURIHandler(opts *ResolveURIHandlerOptions) http.Handler {

	fn := func(rsp http.ResponseWriter, req *http.Request) {

		ctx := req.Context()

		path := req.URL.Path
		short_url := filepath.Base(path)

		u, err := opts.Database.GetURIWithShortURL(ctx, short_url)

		if err != nil {

			if err == io.EOF {
				http.Error(rsp, "Not found", http.StatusNotFound)
				return
			}

			opts.Logger.Printf("Failed to resolve '%s', %w", short_url, err)
			http.Error(rsp, "Not found", http.StatusNotFound)
			return
		}

		http.Redirect(rsp, req, u.Source, 303)
		return
	}

	return http.HandlerFunc(fn)
}
