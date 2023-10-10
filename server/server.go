package server

import (
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

func Start() error {
	mux := mux.NewRouter()

	mux.PathPrefix("/files").Methods(http.MethodGet).HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			if !strings.HasSuffix(r.RequestURI, "/") {
				http.StripPrefix("/files", http.FileServer(http.Dir("./static"))).ServeHTTP(w, r)
				return
			}

			listingHandler{"./static"}.ServeHTTP(w, r)
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
