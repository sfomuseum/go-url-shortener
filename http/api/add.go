package api

import (
	"github.com/aaronland/go-http-sanitize"
	"github.com/sfomuseum/go-http-auth"
	"github.com/sfomuseum/go-url-shortener/database"
	"github.com/sfomuseum/go-url-shortener/uri"
	"io"
	"log"
	"net/http"
	"net/url"
)

type AddURIHandlerOptions struct {
	Authenticator auth.Authenticator
	Database      database.Database
	Logger        *log.Logger
}

func AddURIHandler(opts *AddURIHandlerOptions) http.Handler {

	fn := func(rsp http.ResponseWriter, req *http.Request) {

		ctx := req.Context()

		_, err := opts.Authenticator.GetAccountForRequest(req)

		if err != nil {
			http.Error(rsp, err.Error(), http.StatusBadRequest)
			return
		}

		source_url, err := sanitize.GetString(req, "url")

		if err != nil {
			http.Error(rsp, err.Error(), http.StatusBadRequest)
			return
		}

		if source_url == "" {
			http.Error(rsp, "Missing ?url= parameter", http.StatusBadRequest)
			return
		}

		_, err = url.Parse(source_url)

		if err != nil {
			http.Error(rsp, err.Error(), http.StatusBadRequest)
			return
		}

		u, err := opts.Database.GetURIWithSourceURL(ctx, source_url)

		if err != nil {

			if err != io.EOF {
				http.Error(rsp, err.Error(), http.StatusInternalServerError)
				return
			}

			var short_url string
			max_tries := 10

			for i := 1; i < max_tries; i++ {

				short_url = uri.GenerateShortURI(i)

				short_url, _ := opts.Database.GetURIWithShortURL(ctx, short_url)

				if short_url == nil {
					break
				}
			}

			if short_url == "" {
				opts.Logger.Printf("Failed to generate new short URL, exceeded max tries (%d)", max_tries)
				http.Error(rsp, "Failed to generate new short URL", http.StatusInternalServerError)
				return
			}

			u = &uri.URI{
				Source: source_url,
				Short:  short_url,
			}

			err := opts.Database.AddURI(ctx, u)

			if err != nil {
				opts.Logger.Printf("Failed to put url, %w", err)
				http.Error(rsp, "Failed to store new short URL", http.StatusInternalServerError)
				return
			}
		}

		rsp.Write([]byte(u.Short))
		return
	}

	return http.HandlerFunc(fn)
}
