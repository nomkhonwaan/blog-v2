package storage

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/nomkhonwaan/myblog/pkg/auth"
	"github.com/sirupsen/logrus"
	"mime"
	"net/http"
	"path/filepath"
)

// Register uses to registering HTTP handlers for each sub-router of storage API (prefix: /api/v2/storage)
func Register(r *mux.Router, cache Cache, downloader Downloader, uploader Uploader) {
	r.Path("/{authorizedID}/{fileName}").Handler(downloadHandler(cache, downloader)).Methods(http.MethodGet)
	r.Path("/upload").Handler(uploadHandler(uploader)).Methods(http.MethodPost)
}

func downloadHandler(cache Cache, downloader Downloader) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var (
			file File
			err  error
		)

		vars := mux.Vars(r)
		path := vars["authorizedID"] + "/" + vars["fileName"]

		if cache.Exist(path) {
			file, err = cache.Retrieve(path)
			if err != nil {
				logrus.Errorf("unable to retrieve file from %s due to error: %s", path, err)
			}
		}

		if file.Body == nil {
			file, err = downloader.Download(r.Context(), path)
			if err != nil {
				responseError(w, err.Error(), http.StatusNotFound)
				return
			}
			if err := cache.Store(file); err != nil {
				logrus.Errorf("unable to store file on %s due to error: %s", path, err)
			}
		}

		mimeType := mime.TypeByExtension(filepath.Ext(path))

		w.Header().Set("Content-Type", mimeType)
		w.Header().Set("Content-Length", fmt.Sprintf("%d", len(file.Body)))

		_, _ = w.Write(file.Body)
	})
}

func uploadHandler(uploader Uploader) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		authorizedID := auth.GetAuthorizedUserID(r.Context())
		if authorizedID == nil {
			responseError(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		f, header, err := r.FormFile("file")
		if err != nil {
			responseError(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer f.Close()

		path := authorizedID.(string) + "/" + header.Filename
		logrus.Infof("uploading file %s with size %d to the storage server...", path, header.Size)

		file, err := uploader.Upload(r.Context(), path, f)
		if err != nil {
			responseError(w, err.Error(), http.StatusInternalServerError)
			return
		}

		v, _ := json.Marshal(file)
		_, _ = w.Write(v)
	})
}

func responseError(w http.ResponseWriter, message string, code int) {
	v, _ := json.Marshal(map[string]map[string]interface{}{
		"error": {
			"code":    code,
			"message": message,
		},
	})
	w.WriteHeader(code)
	_, _ = w.Write(v)
}
