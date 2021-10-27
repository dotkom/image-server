package api

func (api *API) setupRoutes() {
	api.router.Use(jsonReponsMiddleware)
	api.router.HandleFunc("/upload", api.upload)
	api.router.HandleFunc("/images/{key}/download", api.download)
	api.router.HandleFunc("/images/{key}", api.info)
	api.router.HandleFunc("/images", api.list)
}
