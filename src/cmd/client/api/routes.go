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

	mux.Post("/insert-blog", conf.InsertBlogHandler)
	mux.Post("/insert-author", conf.InsertAuthorHandler)

	mux.Get("/get-blog/{id}", conf.GetBlogHandler)
	mux.Get("/get-author/{id}", conf.GetAuthorByIDHandler)

	return mux
}
