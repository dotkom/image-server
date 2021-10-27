package api

func (api *API) setupRoutes() {
	api.echo.POST("/upload", api.upload)
	api.echo.GET("/images/:key/download", api.download)
	api.echo.GET("/images/:key", api.info)
}
