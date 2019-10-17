package storage

import (
	"encoding/json"
	"github.com/dgrijalva/jwt-go"
	"github.com/nomkhonwaan/myblog/pkg/auth"
	"github.com/sirupsen/logrus"
	"net/http"
)

func Handler(u Uploader) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		if r.Context().Value(auth.UserProperty) == nil {
			responseError(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}
		authorID := r.Context().Value(auth.UserProperty).(*jwt.Token).Claims.(jwt.MapClaims)["sub"]

		file, header, err := r.FormFile("file")
		if err != nil {
			responseError(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer file.Close()

		path := authorID.(string) + "/" + header.Filename
		logrus.Infof("uploading file %s with size %d to the storage server...", path, header.Size)

		_, err = u.Upload(path, file)
		if err != nil {
			responseError(w, err.Error(), http.StatusInternalServerError)
			return
		}

		v, _ := json.Marshal(map[string]string{"message": "OK"})
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
