package api

import (
	"github.com/go-chi/chi"
	"net/http"
)

func routes() http.Handler {
	mux := chi.NewRouter()

	mux.Use(conf.enableCORS)
	mux.Get("/create-blog", conf.StatusHandler)

	return mux
}
