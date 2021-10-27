package api

import (
	_ "image/png"
	"net/http"

	"github.com/dotkom/image-server/internal/cache"
	"github.com/dotkom/image-server/internal/storage"
	"github.com/gorilla/mux"
)

// API implements http.Handler and is used to handle http requests
type API struct {
	fs     storage.FileStorage
	ms     storage.MetaStorage
	cache  cache.Cache
	router *mux.Router
}

// Creates a new instance of API
func New(fs storage.FileStorage, ms storage.MetaStorage, cache cache.Cache, router *mux.Router) *API {
	api := &API{}
	api.fs = fs
	api.ms = ms
	api.router = router
	api.cache = cache
	api.setupRoutes()
	return api
}

// Handles HTTP requests. Implementation of http.Handler
func (api *API) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	api.router.ServeHTTP(w, r)
}
