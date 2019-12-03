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
	"strings"
)

// Slug is a valid URL string composes with file name and ID
type Slug string

// GetID returns an ID from the slug string
func (s Slug) GetID() (interface{}, error) {
	sl := strings.Split(string(s), "-")
	fileName := sl[len(sl)-1]
	return primitive.ObjectIDFromHex(fileName[0 : len(fileName)-len(filepath.Ext(fileName))])
}

// MustGetID always return ID from the slug string
func (s Slug) MustGetID() interface{} {
	if id, err := s.GetID(); err == nil {
		return id
	}
	return primitive.NewObjectID()
}

// Service helps co-working between data-layer and control-layer
type Service interface {
	// Provide file downloading function
	Downloader

	// Provide file uploading function
	Uploader

	// Provide image resizing function
	image.Resizer

	// A Cache service
	Cache() Cache

	// A File repository
	File() FileRepository
}

type service struct {
	Downloader
	Uploader
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
func NewHandler(cache Cache, fileRepo FileRepository, downloader Downloader, uploader Uploader, resizer image.Resizer) Handler {
	return Handler{
		service: service{
			Downloader: downloader,
			Uploader:   uploader,
			Resizer:    resizer,
			cache:      cache,
			fileRepo:   fileRepo,
		},
	}
}

// Register does registering storage routes under the prefix "/api/v2.1/storage"
func (h Handler) Register(r *mux.Router) {
	r.Path("/{slug}").HandlerFunc(h.download).Methods(http.MethodGet)
	r.Path("/upload").HandlerFunc(h.upload).Methods(http.MethodPost)
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
				// a resized image found on cache storage,
				// clear the `resizedPath` value for preventing image resize function
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

		rdr, wtr := io.Pipe()
		body = io.TeeReader(body, wtr)

		go func(wtr io.WriteCloser, rdr io.Reader) {
			defer wtr.Close()

			if err = h.service.Cache().Store(rdr, resizedPath); err != nil {
				logrus.Errorf("unable to store file on %s: %s", path, err)
			}
		}(wtr, rdr)
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
		}
	}

	if body == nil {
		body, err = h.service.Download(ctx, path)
	}
	if err != nil {
		return nil, err
	}

	rdr, wtr := io.Pipe()
	body = io.TeeReader(body, wtr)

	go func(wtr io.WriteCloser, rdr io.Reader) {
		defer wtr.Close()

		if err = h.service.Cache().Store(rdr, path); err != nil {
			logrus.Errorf("unable to store file on %s: %s", path, err)
		}
	}(wtr, rdr)

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
	err = h.service.Upload(r.Context(), path, f)
	if err != nil {
		h.responseError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	file, err := h.service.File().Create(r.Context(), File{
		ID:             id,
		Path:           path,
		FileName:       fileName,
		Slug:           slug,
		OptionalField1: "CustomizedAmazonS3Client",
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
