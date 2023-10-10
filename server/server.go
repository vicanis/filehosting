package server

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

func Start() error {
	mux := mux.NewRouter()

	mux.PathPrefix("/").Methods(http.MethodGet).HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			parsed, err := url.Parse(r.RequestURI)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				fmt.Fprintf(w, "request parse failed: %s", err)
				return
			}

			if !strings.HasSuffix(parsed.Path, "/") {
				http.FileServer(http.Dir("./static")).ServeHTTP(w, r)
				return
			}

			if parsed.Query().Has("pack") {
				packHandler("./static").ServeHTTP(w, r)
				return
			}

			listingHandler("./static").ServeHTTP(w, r)
		},
	)

	srv := http.Server{
		Addr:         ":9000",
		Handler:      mux,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	log.Printf("start server at %s", srv.Addr)

	return srv.ListenAndServe()
}
