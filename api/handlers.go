package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"image"
	"io/ioutil"
	"mime"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/chai2010/webp"
	"github.com/disintegration/imaging"
	"github.com/dotkom/image-server/models"
	"github.com/gofrs/uuid"
	"github.com/gorilla/mux"
)

func (api *API) list(w http.ResponseWriter, r *http.Request) {
	var images []models.Image
	api.db.Find(&images)
	json.NewEncoder(w).Encode(images)
}

func (api *API) info(w http.ResponseWriter, r *http.Request) {
	key := mux.Vars(r)["key"]

	var image models.Image
	api.db.Limit(1).Find(&image, "id = ?", key)

	json.NewEncoder(w).Encode(image)
}

func (api *API) download(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "image/webp")
	buffer, err := api.cache.Get(fmt.Sprintf("%s/%s", r.URL.Path, r.URL.Query().Encode()))
	if err == nil {
		w.Write(buffer)
		return
	}

	key := mux.Vars(r)["key"]
	width, err := strconv.Atoi(r.URL.Query().Get("width"))
	if err != nil {
		width = 0
	}
	height, err := strconv.Atoi(r.URL.Query().Get("height"))
	if err != nil {
		height = 0
	}
	quality, err := strconv.ParseFloat(r.URL.Query().Get("quality"), 32)
	if err != nil {
		quality = 100
	}

	var image models.Image
	api.db.Limit(1).Find(&image, "id = ?", key)
	buffer, err = api.storage.Get(r.Context(), key)
	if err != nil {
		fmt.Print(err)
		return
	}

	img, err := webp.Decode(bytes.NewReader(buffer))
	if err != nil {
		fmt.Print(err)
		return
	}

	if width != 0 || height != 0 {
		img = imaging.Resize(img, width, height, imaging.Lanczos)
	}

	buffer, err = webp.EncodeRGBA(img, float32(quality))
	if err != nil {
		fmt.Print(err)
		return
	}

	err = api.cache.Set(fmt.Sprintf("%s/%s", r.URL.Path, r.URL.Query().Encode()), buffer)
	if err != nil {
		fmt.Print(err)
	}

	w.Write(buffer)
}

func (api *API) upload(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(32 << 20)
	file, fileInfo, err := r.FormFile("file")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}
	defer file.Close()

	mimeType := mime.TypeByExtension(filepath.Ext(fileInfo.Filename))

	name := r.FormValue("name")
	description := r.FormValue("description")
	tags := r.FormValue("tags")

	if name == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	buffer, err := ioutil.ReadAll(file)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Print(err)
		return
	}

	if mimeType != "image/webp" {
		img, _, err := image.Decode(bytes.NewReader(buffer))
		if err != nil {
			fmt.Print(err)
			return
		}
		buffer, err = webp.EncodeExactLosslessRGBA(img)
		if err != nil {
			fmt.Print(err)
			return
		}
		mimeType = "image/webp"
	}

	key, err := uuid.NewV4()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Print(err)
		return
	}

	err = api.storage.Save(r.Context(), buffer, key.String())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Print(err)
		return
	}

	tagSlice := strings.Split(tags, ",")
	for i := range tagSlice {
		tagSlice[i] = strings.TrimSpace(tagSlice[i])
	}

	var image models.Image
	newImage := &models.Image{
		ID:          key,
		Name:        name,
		Description: description,
		Tags:        tagSlice,
		Mime:        mimeType,
		Size:        uint64(fileInfo.Size),
	}
	api.db.Model(&image).Create(newImage)

	json.NewEncoder(w).Encode(newImage)
}
