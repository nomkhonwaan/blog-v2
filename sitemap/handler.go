package sitemap

import (
	"bytes"
	"context"
	"encoding/xml"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/nomkhonwaan/myblog/pkg/blog"
	"github.com/nomkhonwaan/myblog/pkg/facebook"
	"github.com/nomkhonwaan/myblog/pkg/storage"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
	"time"
)

const (
	sitemapFilePath = "sitemap.xml"
)

type blogService struct{ blog.Service }

// Service helps co-working between data-layer and control-layer
type Service interface {
	// Create a new sitemap and return byte array of the XML data
	Generate() ([]byte, error)

	// A Cache service
	Cache() storage.Cache

	// A category repository
	Category() blog.CategoryRepository

	// A post repository
	Post() blog.PostRepository

	// A tag repository
	Tag() blog.TagRepository
}

type service struct {
	blogService

	baseURL string
	cache   storage.Cache
}

func (s service) Generate() ([]byte, error) {
	urlSet := NewURLSet()
	urlSet.URLs = append(urlSet.URLs, URL{
		Location: s.baseURL,
		Priority: "1",
	})

	// TODO: the limit number should be able to configure rather than fixed value
	posts, err := s.Post().FindAll(context.Background(), blog.NewPostQueryBuilder().
		WithStatus(blog.Published).
		WithLimit(9999).
		Build(),
	)
	if err != nil {
		return nil, err
	}

	for i, p := range posts {
		// update the base URL last modify with the latest published posts publishing time
		if i == 0 {
			urlSet.URLs[0].LastModify = p.PublishedAt.Format(time.RFC3339)
		}

		lastModify := p.PublishedAt
		if !p.UpdatedAt.IsZero() {
			lastModify = p.UpdatedAt
		}

		urlSet.URLs = append(urlSet.URLs, URL{
			Location:   s.baseURL + "/" + p.PublishedAt.In(facebook.DefaultTimeZone).Format("2006/1/2") + "/" + p.Slug,
			LastModify: lastModify.Format(time.RFC3339),
			Priority:   "0.8",
		})
	}

	categories, err := s.Category().FindAll(context.Background())
	if err != nil {
		return nil, err
	}
	for _, cat := range categories {
		urlSet.URLs = append(urlSet.URLs, URL{
			Location: s.baseURL + "/category/" + cat.Slug,
			Priority: "0.5",
		})
	}

	tags, err := s.Tag().FindAll(context.Background())
	if err != nil {
		return nil, err
	}
	for _, tag := range tags {
		urlSet.URLs = append(urlSet.URLs, URL{
			Location: s.baseURL + "/tag/" + tag.Slug,
			Priority: "0.5",
		})
	}

	return xml.Marshal(urlSet)
}

func (s service) Cache() storage.Cache {
	return s.cache
}

// Handler provides site-map handlers
type Handler struct {
	service Service
}

// NewHandler returns a new handler instance
func NewHandler(baseURL string, cache storage.Cache, blogSvc blog.Service) Handler {
	return Handler{
		service: service{
			blogService: blogService{blogSvc},
			baseURL:     baseURL,
			cache:       cache,
		},
	}
}

// Register does registering site-map routes under the prefix "/sitemap.xml"
func (h Handler) Register(r *mux.Router) {
	r.Path("").HandlerFunc(h.serve).Methods(http.MethodGet)
}

func (h Handler) serve(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/xml")

	if h.service.Cache().Exist(sitemapFilePath) {
		body, err := h.service.Cache().Retrieve(sitemapFilePath)
		if err == nil {
			length, _ := io.Copy(w, body)
			w.Header().Set("Content-Length", fmt.Sprintf("%d", length))
			return
		}
		logrus.Errorf("unable to retrieve file from sitemap.xml: %s", err)
	}

	data, err := h.service.Generate()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	data = append([]byte(`<?xml version="1.0" encoding="UTF-8"?>`), data...)

	err = h.service.Cache().Store(bytes.NewReader(data), sitemapFilePath)
	if err != nil {
		logrus.Errorf("unable to store file on sitemap.xml: %s", err)
	}

	_, _ = w.Write(data)
}
