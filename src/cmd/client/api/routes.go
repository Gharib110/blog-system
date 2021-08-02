package api

import (
	"github.com/go-chi/chi"
	"net/http"
)

// Routes make http handler for HTTP server
func Routes() http.Handler {
	mux := chi.NewRouter()

	mux.Use(conf.enableCORS)
	mux.Get("/create-blog", conf.StatusHandler)

	return mux
}
