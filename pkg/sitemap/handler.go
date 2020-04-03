package sitemap

import (
	"bytes"
	"context"
	"encoding/xml"
	"github.com/nomkhonwaan/myblog/pkg/blog"
	"github.com/nomkhonwaan/myblog/pkg/storage"
	"github.com/nomkhonwaan/myblog/pkg/timeutil"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
	"net/url"
	"time"
)

const (
	cacheFilePath = "sitemap.xml"
)

// ServeSiteMapHandlerFunc provides sitemap.xml file for search engine robot
func ServeSiteMapHandlerFunc(cache storage.Cache, genURLsFunc ...func() ([]URL, error)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/xml")

		if cache.Exists("sitemap.xml") {
			body, err := cache.Retrieve(cacheFilePath)
			if err == nil {
				_, _ = io.Copy(w, body)
				return
			}
			logrus.Errorf("unable to retrieve sitemap.xml: %s", err)
		}

		urlSet, err := generateURLSet(genURLsFunc...)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		data, _ := xml.Marshal(urlSet)
		data = append([]byte(`<?xml version="1.0" encoding="UTF-8"?>`), data...)

		err = cache.Store(bytes.NewReader(data), cacheFilePath)
		if err != nil {
			logrus.Errorf("unable to store sitemap.xml: %s", err)
		}

		_, _ = w.Write(data)
	}
}

func generateURLSet(genURLsFunc ...func() ([]URL, error)) (URLSet, error) {
	urlSet := URLSet{URLs: make([]URL, 0)}
	for _, f := range genURLsFunc {
		urls, err := f()
		if err != nil {
			return urlSet, err
		}
		urlSet.URLs = append(urlSet.URLs, urls...)
	}
	return urlSet, nil
}

// GenerateFixedURLs generates all fixed URLs (which is not required any data source connection)
func GenerateFixedURLs(baseURL string) func() ([]URL, error) {
	return func() ([]URL, error) {
		return []URL{
			{
				Location: baseURL,
				Priority: 1,
			},
		}, nil
	}
}

// GeneratePostURLs generates all Post URLs
func GeneratePostURLs(baseURL string, repository blog.PostRepository) func() ([]URL, error) {
	return func() ([]URL, error) {
		posts, err := repository.FindAll(context.Background(), blog.NewPostQueryBuilder().
			WithStatus(blog.StatusPublished).WithLimit(9999).Build())
		if err != nil {
			return nil, err
		}
		urls := make([]URL, len(posts))
		for i, p := range posts {
			location, _ := url.Parse(baseURL + "/" + p.PublishedAt.In(timeutil.TimeZoneAsiaBangkok).Format("2006/1/2") + "/" + p.Slug)
			lastModify := p.PublishedAt
			if !p.UpdatedAt.IsZero() {
				lastModify = p.UpdatedAt
			}
			urls[i] = URL{
				Location:   location.String(),
				LastModify: lastModify.Format(time.RFC3339),
				Priority:   0.8,
			}
		}
		return urls, nil
	}
}

// GenerateCategoryURLs generates all Category URLs
func GenerateCategoryURLs(baseURL string, repository blog.CategoryRepository) func() ([]URL, error) {
	return func() ([]URL, error) {
		cats, err := repository.FindAll(context.Background())
		if err != nil {
			return nil, err
		}

		urls := make([]URL, len(cats))
		for i, c := range cats {
			location, _ := url.Parse(baseURL + "/category/" + c.Slug)
			urls[i] = URL{
				Location: location.String(),
				Priority: 0.5,
			}
		}
		return urls, err
	}
}

// GenerateTagURLs generates all Tag URLs
func GenerateTagURLs(baseURL string, repository blog.TagRepository) func() ([]URL, error) {
	return func() ([]URL, error) {
		tags, err := repository.FindAll(context.Background())
		if err != nil {
			return nil, err
		}

		urls := make([]URL, len(tags))
		for i, t := range tags {
			location, _ := url.Parse(baseURL + "/tag/" + t.Slug)
			urls[i] = URL{
				Location: location.String(),
				Priority: 0.5,
			}
		}
		return urls, err
	}
}

//const (
//	// CacheFilePath refers to path of the sitemap.xml to be saved in the cache storage
//	CacheFilePath = "sitemap.xml"
//)
//
//// Service helps co-working between data-layer and control-layer
//type Service interface {
//	// Create a new sitemap and return byte array of the XML data
//	Generate() ([]byte, error)
//	// A Cache service
//	Cache() storage.Cache
//	// A category repository
//	Category() blog.CategoryRepository
//	// A post repository
//	Post() blog.PostRepository
//	// A tag repository
//	Tag() blog.TagRepository
//}
//
//type service struct {
//	blog.Service
//
//	baseURL      string
//	cacheService storage.Cache
//}
//
//func (s service) Generate() ([]byte, error) {
//	urlSet := NewURLSet()
//	urlSet.URLs = append(urlSet.URLs, URL{
//		Location: s.baseURL,
//		Priority: "1",
//	})
//
//	// TODO: the limit number should be able to configure rather than fixed value
//	posts, err := s.Post().FindAll(context.Background(), blog.NewPostQueryBuilder().
//		WithStatus(blog.StatusPublished).
//		WithLimit(9999).
//		Build(),
//	)
//	if err != nil {
//		return nil, err
//	}
//
//	for i, p := range posts {
//		// update the base URLNode last modify with the latest published posts publishing time
//		if i == 0 {
//			urlSet.URLs[0].LastModify = p.PublishedAt.Format(time.RFC3339)
//		}
//
//		lastModify := p.PublishedAt
//		if !p.UpdatedAt.IsZero() {
//			lastModify = p.UpdatedAt
//		}
//
//		location, _ := url.Parse(s.baseURL + "/" + p.PublishedAt.In(timeutil.TimeZoneAsiaBangkok).Format("2006/1/2") + "/" + p.Slug)
//		urlSet.URLs = append(urlSet.URLs, URL{
//			Location:   location.String(),
//			LastModify: lastModify.Format(time.RFC3339),
//			Priority:   "0.8",
//		})
//	}
//
//	categories, err := s.Category().FindAll(context.Background())
//	if err != nil {
//		return nil, err
//	}
//	for _, cat := range categories {
//		urlSet.URLs = append(urlSet.URLs, URL{
//			Location: s.baseURL + "/category/" + cat.Slug,
//			Priority: "0.5",
//		})
//	}
//
//	tags, err := s.Tag().FindAll(context.Background())
//	if err != nil {
//		return nil, err
//	}
//	for _, tag := range tags {
//		urlSet.URLs = append(urlSet.URLs, URL{
//			Location: s.baseURL + "/tag/" + tag.Slug,
//			Priority: "0.5",
//		})
//	}
//
//	return xml.Marshal(urlSet)
//}
//
//func (s service) Cache() storage.Cache {
//	return s.cacheService
//}
//
//// Handler provides site-map handlers
//type Handler struct {
//	service Service
//}
//
//// NewHandler returns a new handler instance
//func NewHandler(baseURL string, cache storage.Cache, blogService blog.Service) Handler {
//	return Handler{
//		service: service{
//			Service:      blogService,
//			baseURL:      baseURL,
//			cacheService: cache,
//		},
//	}
//}
//
//// Register does registering sitemap routes under the prefix "/sitemap.xml"
//func (h Handler) Register(r *mux.Router) {
//	r.Path("").HandlerFunc(h.serve).Methods(http.MethodGet)
//}
//
//func (h Handler) serve(w http.ResponseWriter, r *http.Request) {
//	w.Header().Set("Content-Type", "text/xml")
//
//	if h.service.Cache().Exists(CacheFilePath) {
//		body, err := h.service.Cache().Retrieve(CacheFilePath)
//		if err == nil {
//			length, _ := io.Copy(w, body)
//			w.Header().Set("Content-Length", fmt.Sprintf("%d", length))
//			return
//		}
//		logrus.Errorf("unable to retrieve file from sitemap.xml: %s", err)
//	}
//
//	data, err := h.service.Generate()
//	if err != nil {
//		http.Error(w, err.Error(), http.StatusInternalServerError)
//		return
//	}
//	data = append([]byte(`<?xml version="1.0" encoding="UTF-8"?>`), data...)
//
//	err = h.service.Cache().Store(bytes.NewReader(data), CacheFilePath)
//	if err != nil {
//		logrus.Errorf("unable to store file on sitemap.xml: %s", err)
//	}
//
//	_, _ = w.Write(data)
//}
