package storage

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/nomkhonwaan/myblog/pkg/auth"
	"github.com/sirupsen/logrus"
	"net/http"
)

// Register allows the register HTTP handlers for each sub-router
func Register(r *mux.Router, uploader Uploader) {
	r.Handle("/upload", uploadFileHandler(uploader)).Methods(http.MethodPost)
}

func uploadFileHandler(u Uploader) http.Handler {
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

		file, err := u.Upload(r.Context(), path, f)
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
