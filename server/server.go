package server

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

func Start() error {
	mux := mux.NewRouter()

	mux.PathPrefix("/files").Methods(http.MethodGet).Handler(
		http.StripPrefix("/files", http.FileServer(http.Dir("./static"))),
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
