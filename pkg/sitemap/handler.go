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
				defer body.Close()
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
