package storage

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
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

// Service helps co-working between data-layer and control-layer
type Service interface {
	// Provide storage functions
	Storage
	// Provide image resizing functions
	image.Resizer
	// A Cache service
	Cache() Cache
	// A File repository
	File() FileRepository
}

type service struct {
	Storage
	image.Resizer

	cache    Cache
	fileRepo FileRepository
}

func (s service) Cache() Cache {
	return s.cache
}

func (s service) File() FileRepository {
	return s.fileRepo
}

// Handler provides storage handlers
type Handler struct {
	service Service
}

// NewHandler returns a new handler instance
func NewHandler(cache Cache, storage Storage, fileRepo FileRepository, resizer image.Resizer) Handler {
	return Handler{
		service: service{
			Storage:  storage,
			Resizer:  resizer,
			cache:    cache,
			fileRepo: fileRepo,
		},
	}
}

// Register does registering storage routes under the prefix "/api/v2.1/storage"
func (h Handler) Register(r *mux.Router) {
	r.Path("/{slug}").HandlerFunc(h.download).Methods(http.MethodGet)
	r.Path("/delete/{slug}").HandlerFunc(h.delete).Methods(http.MethodDelete)
	r.Path("/upload").HandlerFunc(h.upload).Methods(http.MethodPost)
}

func (h Handler) delete(w http.ResponseWriter, r *http.Request) {
	authorizedID := auth.GetAuthorizedUserID(r.Context())
	if authorizedID == nil {
		h.responseError(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	var (
		vars = mux.Vars(r)
		slug = Slug(vars["slug"])
	)

	file, err := h.service.File().FindByID(r.Context(), slug.MustGetID())
	if err != nil {
		h.responseError(w, err.Error(), http.StatusNotFound)
		return
	}

	logrus.Infof("deleting file %s from the storage server...", file.Path)
	err = h.service.Delete(r.Context(), file.Path)
	if err != nil {
		h.responseError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = h.service.File().Delete(r.Context(), slug.MustGetID())
	if err != nil {
		h.responseError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	return
}

func (h Handler) download(w http.ResponseWriter, r *http.Request) {
	var (
		vars          = mux.Vars(r)
		slug          = Slug(vars["slug"])
		width, height = getWidthAndHeightFromQuery(r.URL.Query())

		body        io.Reader
		resizedPath string
	)

	file, err := h.service.File().FindByID(r.Context(), slug.MustGetID())
	if err != nil {
		h.responseError(w, err.Error(), http.StatusNotFound)
		return
	}

	path := file.Path
	mimeType := mime.TypeByExtension(filepath.Ext(path))

	if (mimeType == "image/jpeg" || mimeType == "image/png") && (width > 0 || height > 0) {
		resizedPath = fmt.Sprintf("%s-%d-%d%s", path[0:len(path)-len(filepath.Ext(path))], width, height, filepath.Ext(path))

		if h.service.Cache().Exist(resizedPath) {
			body, err = h.service.Cache().Retrieve(resizedPath)
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
		body, err = h.downloadOriginalFile(r.Context(), path)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
	}

	if resizedPath != "" {
		var buf bytes.Buffer
		rdr := io.TeeReader(body, &buf)
		body, err = h.service.Resize(rdr, width, height)
		if err != nil {
			logrus.Errorf("unable to resize image: %s", err)
		}
		if body == nil {
			body = &buf
		}

		rdr = io.TeeReader(body, &buf)
		if err = h.service.Cache().Store(rdr, resizedPath); err != nil {
			logrus.Errorf("unable to store file on %s: %s", path, err)
		}
		body = &buf
	}

	length, _ := io.Copy(w, body)

	w.Header().Set("Content-Type", mimeType)
	w.Header().Set("Content-Length", fmt.Sprintf("%d", length))
}

func (h Handler) downloadOriginalFile(ctx context.Context, path string) (body io.Reader, err error) {
	if h.service.Cache().Exist(path) {
		body, err = h.service.Cache().Retrieve(path)
		if err != nil {
			logrus.Errorf("unable to retrieve file from %s: %s", path, err)
		} else {
			return body, nil
		}
	}

	body, err = h.service.Download(ctx, path)
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	rdr := io.TeeReader(body, &buf)
	if err = h.service.Cache().Store(rdr, path); err != nil {
		logrus.Errorf("unable to store file on %s: %s", path, err)
	}
	body = &buf

	return body, nil
}

func (h Handler) upload(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	authorizedID := auth.GetAuthorizedUserID(r.Context())
	if authorizedID == nil {
		h.responseError(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	f, header, err := r.FormFile("file")
	if err != nil {
		h.responseError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer f.Close()

	id := primitive.NewObjectID()
	fileName := header.Filename
	ext := filepath.Ext(fileName)
	slug := fmt.Sprintf("%s-%s%s", slugify.Make(fileName[0:len(fileName)-len(ext)]), id.Hex(), ext)
	path := authorizedID.(string) + string(filepath.Separator) + slug

	logrus.Infof("uploading file %s with size %d to the storage server...", path, header.Size)
	err = h.service.Upload(r.Context(), f, path)
	if err != nil {
		h.responseError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	file, err := h.service.File().Create(r.Context(), File{
		ID:       id,
		Path:     path,
		FileName: fileName,
		Slug:     slug,
	})
	if err != nil {
		h.responseError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	val, _ := json.Marshal(file)
	_, _ = w.Write(val)
}

func (h Handler) responseError(w http.ResponseWriter, message string, code int) {
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

func getWidthAndHeightFromQuery(values url.Values) (int, int) {
	w, _ := strconv.Atoi(values.Get("width"))
	h, _ := strconv.Atoi(values.Get("height"))
	return w, h
}

func (h Handler) Delete(w http.ResponseWriter, r *http.Request) {

}

func (h Handler) Download(w http.ResponseWriter, r *http.Request) {

}

func (h Handler) Upload(w http.ResponseWriter, r *http.Request) {

}
