package api

import (
	_ "image/png"
	"io"
	"net/http"
	"strings"

	"github.com/dotkom/image-server/internal/cache"
	"github.com/dotkom/image-server/internal/storage"
	"github.com/gorilla/mux"
	"github.com/labstack/echo-contrib/jaegertracing"
	"github.com/labstack/echo-contrib/prometheus"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// API implements http.Handler and is used to handle http requests
type API struct {
	fs     storage.FileStorage
	ms     storage.MetaStorage
	cache  cache.Cache
	echo   *echo.Echo
	jeager io.Closer
}

// Creates a new instance of API
func New(fs storage.FileStorage, ms storage.MetaStorage, cache cache.Cache, router *mux.Router) *API {
	api := &API{}
	api.fs = fs
	api.ms = ms
	api.cache = cache
	api.echo = echo.New()
	api.echo.Pre(middleware.RemoveTrailingSlash())
	prom := prometheus.NewPrometheus("image_server", nil)
	prom.Use(api.echo)
	api.jeager = jaegertracing.New(api.echo, nil)

	api.echo.Use(middleware.GzipWithConfig(middleware.GzipConfig{
		Skipper: func(c echo.Context) bool {
			return strings.Contains(c.Path(), "metrics")
		},
	}))
	api.setupRoutes()

	return api
}

func (api *API) Close() error {
	api.jeager.Close()
	return nil
}

// Handles HTTP requests. Implementation of http.Handler
func (api *API) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	api.echo.ServeHTTP(w, r)
}
