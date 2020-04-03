package storage

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/nomkhonwaan/myblog/pkg/auth"
	"github.com/nomkhonwaan/myblog/pkg/image"
	slugify "github.com/nomkhonwaan/myblog/pkg/slug"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"io"
	"mime"
	"net/http"
	"net/url"
	"path/filepath"
	"strconv"
)

// DeleteHandlerFunc handles deletion request
func DeleteHandlerFunc(storage Storage, repository FileRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authorizedID := auth.GetAuthorizedUserID(r.Context())
		if authorizedID == nil {
			respondError(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		slug := Slug(chi.URLParam(r, "slug"))
		file, err := repository.FindByID(r.Context(), slug.MustGetID())
		if err != nil {
			respondError(w, err.Error(), http.StatusNotFound)
			return
		}

		logrus.Infof("deleting file %s from the storage server...", file.Path)
		err = storage.Delete(r.Context(), file.Path)
		if err != nil {
			respondError(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = repository.Delete(r.Context(), slug.MustGetID())
		if err != nil {
			respondError(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

// DownloadHandlerFunc handles downloading request
func DownloadHandlerFunc(storage Storage, cache Cache, resizer image.Resizer, repository FileRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			body          io.Reader
			resizedPath   string
			slug          = Slug(chi.URLParam(r, "slug"))
			width, height = getWidthAndHeight(r.URL.Query())
		)

		file, err := repository.FindByID(r.Context(), slug.MustGetID())
		if err != nil {
			respondError(w, err.Error(), http.StatusNotFound)
			return
		}

		path := file.Path
		mimeType := mime.TypeByExtension(filepath.Ext(path))

		if (mimeType == "image/jpeg" || mimeType == "image/png") && (width > 0 || height > 0) {
			resizedPath = fmt.Sprintf("%s-%d-%d%s", path[0:len(path)-len(filepath.Ext(path))], width, height, filepath.Ext(path))

			if cache.Exists(resizedPath) {
				body, err = cache.Retrieve(resizedPath)
				if err != nil {
					logrus.Errorf("unable to retrieve file from %s: %s", path, err)
				} else {
					// resized image already on the cache storage,
					// clear `resizedPath` for preventing resize function
					resizedPath = ""
				}
			}
		}

		if body == nil {
			body, err = downloadOriginalFile(r.Context(), storage, cache, path)
			if err != nil {
				http.Error(w, err.Error(), http.StatusNotFound)
				return
			}
		}

		if resizedPath != "" {
			var buf bytes.Buffer
			rdr := io.TeeReader(body, &buf)
			body, err = resizer.Resize(rdr, width, height)
			if err != nil {
				logrus.Errorf("unable to resize image: %s", err)
			}
			if body == nil {
				body = &buf
			}

			rdr = io.TeeReader(body, &buf)
			if err = cache.Store(rdr, resizedPath); err != nil {
				logrus.Errorf("unable to store file on %s: %s", path, err)
			}
			body = &buf
		}

		length, _ := io.Copy(w, body)

		w.Header().Set("Content-Type", mimeType)
		w.Header().Set("Content-Length", fmt.Sprintf("%d", length))
	}
}

func downloadOriginalFile(ctx context.Context, storage Storage, cache Cache, path string) (body io.Reader, err error) {
	if cache.Exists(path) {
		body, err = cache.Retrieve(path)
		if err != nil {
			logrus.Errorf("unable to retrieve file from %s: %s", path, err)
		} else {
			return body, nil
		}
	}

	originalFileBody, err := storage.Download(ctx, path)
	if err != nil {
		return nil, err
	}
	defer originalFileBody.Close()

	body = originalFileBody

	var buf bytes.Buffer
	rdr := io.TeeReader(body, &buf)
	if err = cache.Store(rdr, path); err != nil {
		logrus.Errorf("unable to store file on %s: %s", path, err)
	}
	body = &buf

	return body, nil
}

// UpdateHandlerFunc handles uploading request
func UploadHandlerFunc(storage Storage, repository FileRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		authorizedID := auth.GetAuthorizedUserID(r.Context())
		if authorizedID == nil {
			respondError(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		f, header, err := r.FormFile("file")
		if err != nil {
			respondError(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer f.Close()

		var (
			id       = primitive.NewObjectID()
			fileName = header.Filename
			ext      = filepath.Ext(fileName)
			slug     = fmt.Sprintf("%s-%s%s", slugify.Make(fileName[0:len(fileName)-len(ext)]), id.Hex(), ext)
			path     = authorizedID.(string) + string(filepath.Separator) + slug
		)
		logrus.Infof("uploading file %s with size %d to the storage server...", path, header.Size)
		err = storage.Upload(r.Context(), f, path)
		if err != nil {
			respondError(w, err.Error(), http.StatusInternalServerError)
			return
		}

		file, err := repository.Create(r.Context(), File{
			ID:       id,
			Path:     path,
			FileName: fileName,
			Slug:     slug,
		})
		if err != nil {
			respondError(w, err.Error(), http.StatusInternalServerError)
			return
		}

		val, _ := json.Marshal(file)
		_, _ = w.Write(val)
	}
}

func respondError(w http.ResponseWriter, message string, code int) {
	var data struct {
		Error struct {
			Code    int    `json:"code"`
			Message string `json:"message"`
		} `json:"error"`
	}

	data.Error.Code = code
	data.Error.Message = message

	val, _ := json.Marshal(data)

	w.WriteHeader(code)
	_, _ = w.Write(val)
}

func getWidthAndHeight(values url.Values) (int, int) {
	w, _ := strconv.Atoi(values.Get("width"))
	h, _ := strconv.Atoi(values.Get("height"))
	return w, h
}
