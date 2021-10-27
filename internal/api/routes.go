package api

func (api *API) setupRoutes() {
	api.router.Use(jsonReponsMiddleware)
	api.router.HandleFunc("/upload", api.upload).Methods("POST")
	api.router.HandleFunc("/images/{key}/download", api.download)
	api.router.HandleFunc("/images/{key}", api.info)
}
