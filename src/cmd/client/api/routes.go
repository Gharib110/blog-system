package api

import (
	"github.com/go-chi/chi"
	"net/http"
)

// Routes make http handlers for HTTP1.X API server
func Routes() http.Handler {
	mux := chi.NewRouter()

	mux.Use(conf.enableCORS)
	mux.Get("/create-blog", conf.StatusHandler)

	mux.Post("/insert-blog", conf.InsertBlogHandler) // Implemented

	mux.Get("/get-blog/{id}", conf.GetBlogHandler) // Implemented
	mux.Get("/get-all-blog/{num}", conf.GetAllBlogsHandler)

	return mux
}
