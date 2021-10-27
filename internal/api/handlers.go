package api

import (
	"bytes"
	"fmt"
	"image"
	"io/ioutil"
	"mime"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/chai2010/webp"
	"github.com/disintegration/imaging"
	"github.com/dotkom/image-server/internal/models"
	"github.com/labstack/echo/v4"

	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

type infoParams struct {
	Key string `param:"key"`
}

// List info about single image
func (api *API) info(c echo.Context) error {
	vars := new(infoParams)
	if err := c.Bind(vars); err != nil {
		return c.JSON(http.StatusBadRequest, nil)
	}
	meta, err := api.ms.Get(c.Request().Context(), vars.Key)
	if err != nil {
		return c.JSONBlob(http.StatusNotFound, []byte(`{"msg": "Not found"}`))
	}

	return c.JSON(http.StatusOK, meta)
}

type downloadParams struct {
	Key     string  `param:"key"`
	Width   int     `query:"width"`
	Height  int     `query:"height"`
	Quality float64 `query:"quality"`
}

// Download image with given width, height and quality
func (api *API) download(c echo.Context) error {
	vars := new(downloadParams)
	if err := c.Bind(vars); err != nil {
		return c.JSON(http.StatusBadRequest, nil)
	}

	if vars.Quality == 0 {
		vars.Quality = 100
	}

	// Check if the image is cached and return the cached version if it is
	imageCacheKey := fmt.Sprintf("image:%s:%d:%d:%f", vars.Key, vars.Width, vars.Height, vars.Quality)
	buffer, err := api.cache.GetByteBuffer(imageCacheKey)
	if err == nil {
		c.Stream(http.StatusOK, "image/webp", bytes.NewReader(buffer))
		return err
	}

	sizeCacheKey := fmt.Sprintf("meta:size:%s", vars.Key)

	size, err := api.cache.GetUint64(sizeCacheKey)
	if err != nil {
		meta, err := api.ms.Get(c.Request().Context(), vars.Key)
		if err != nil {
			log.Debug(err)
		}
		size = meta.Size
		api.cache.SetUint64(sizeCacheKey, size)
	}

	// Get image buffer
	buffer, err = api.fs.Get(c.Request().Context(), vars.Key, size)
	if err != nil {
		c.JSON(http.StatusNotFound, nil)
		return err
	}

	// If no query parameters are specified, the image is returned without modifications. And added to the cache
	if vars.Width == 0 && vars.Height == 0 && vars.Quality == 100 {
		api.cache.SetByteBuffer(imageCacheKey, buffer)
		c.Stream(http.StatusOK, "image/webp", bytes.NewReader(buffer))
		return err
	}

	// Decode the image to webp
	img, err := webp.Decode(bytes.NewReader(buffer))
	if err != nil {
		log.Error(err)
		return err
	}

	// Resize
	if vars.Width != 0 || vars.Height != 0 {
		img = imaging.Resize(img, vars.Width, vars.Height, imaging.Lanczos)
	}

	// Encode with the requested quality
	buffer, err = webp.EncodeRGBA(img, float32(vars.Quality))
	if err != nil {
		log.Error(err)
		return err
	}

	// Update cache
	api.cache.SetByteBuffer(imageCacheKey, buffer)

	return c.Stream(http.StatusOK, "image/webp", bytes.NewReader(buffer))

}

type uploadParams struct {
	Name        string `form:"name"`
	Description string `form:"description"`
	Tags        string `form:"tags"`
}

// Upload image
func (api *API) upload(c echo.Context) error {
	vars := new(uploadParams)
	if err := c.Bind(vars); err != nil {
		return c.JSON(http.StatusBadRequest, nil)
	}
	// Parse form. Maxium 32MB
	fileHeader, err := c.FormFile("file")
	if err != nil {
		return err
	}
	file, err := fileHeader.Open()
	if err != nil {
		return err
	}

	defer file.Close()

	// Get mime type and form values
	mimeType := mime.TypeByExtension(filepath.Ext(fileHeader.Filename))

	if vars.Name == "" {
		return c.JSONBlob(http.StatusBadRequest, []byte(`{"msg": "name is required"}`))
	}

	// Read file as byte buffer
	buffer, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}

	// Change type of image to webp
	if mimeType != "image/webp" {
		img, _, err := image.Decode(bytes.NewReader(buffer))
		if err != nil {
			log.Debug(err)
			return err
		}
		buffer, err = webp.EncodeExactLosslessRGBA(img)
		if err != nil {
			log.Error(err)
			return err
		}
		mimeType = "image/webp"
	}

	key, err := uuid.NewRandom()
	if err != nil {
		log.Error(err)
		return err
	}

	// Save the image to persistant storage
	err = api.fs.Save(c.Request().Context(), buffer, key.String())
	if err != nil {
		log.Error(err)
		return err
	}

	// Decode tag csv
	tagSlice := strings.Split(vars.Tags, ",")
	for i := range tagSlice {
		tagSlice[i] = strings.TrimSpace(tagSlice[i])
	}

	// Save image data to db
	meta := models.ImageMeta{
		Key:         key,
		Name:        vars.Name,
		Description: vars.Description,
		Tags:        tagSlice,
		Mime:        mimeType,
		Size:        uint64(len(buffer)),
	}
	err = api.ms.Save(c.Request().Context(), meta)
	if err != nil {
		log.Error(err)
		return err
	}

	return c.JSON(http.StatusOK, meta)
}
