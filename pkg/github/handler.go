package github

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/nomkhonwaan/myblog/pkg/storage"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
	"net/url"
	"strings"
)

// Service helps co-working between data-layer and control-layer
type Service interface {
	// A cache service
	Cache() storage.Cache

	// Request to GitHub server for the Gist content
	Retrieve(src string) (*http.Response, error)
}

type service struct {
	cache     storage.Cache
	transport http.RoundTripper
}

func (s service) Cache() storage.Cache {
	return s.cache
}

func (s service) Retrieve(src string) (*http.Response, error) {
	u, _ := url.Parse(src)

	u.Host = "gist.github.com"
	u.Path = strings.Replace(u.Path, ".js", ".json", 1)

	req, _ := http.NewRequest(http.MethodGet, u.String(), nil)
	return s.transport.RoundTrip(req)
}

type Handler struct {
	service Service
}

func NewHandler(cache storage.Cache, transport http.RoundTripper) Handler {
	return Handler{
		service: service{
			cache:     cache,
			transport: transport,
		},
	}
}

// Register does registering GitHub routes under the prefix "/api/v2.1/github"
func (h Handler) Register(r *mux.Router) {
	r.Path("/gist").HandlerFunc(h.serveGist).Methods(http.MethodGet)
}

func (h Handler) serveGist(w http.ResponseWriter, r *http.Request) {
	src := r.URL.Query().Get("src")
	if src == "" {
		http.Error(w, "`src` is required for downloading Gist content", http.StatusBadRequest)
		return
	}

	path := url.QueryEscape(src) + ".json"

	if h.service.Cache().Exist(path) {
		body, err := h.service.Cache().Retrieve(path)
		if err == nil {
			length, _ := io.Copy(w, body)
			w.Header().Set("Content-Length", fmt.Sprintf("%d", length))
			return
		}

		logrus.Errorf("unable to retrieve Gist file from %s: %s", path, err)
	}

	res, err := h.service.Retrieve(src)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer res.Body.Close()

	rdr, wtr := io.Pipe()
	body := io.TeeReader(res.Body, wtr)

	go func(wtr io.WriteCloser, rdr io.Reader) {
		defer wtr.Close()

		if err = h.service.Cache().Store(rdr, path); err != nil {
			logrus.Errorf("unable to store file on %s: %s", path, err)
		}
	}(wtr, rdr)

	length, _ := io.Copy(w, body)

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Length", fmt.Sprintf("%d", length))
}
