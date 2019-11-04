package facebook

import (
	"compress/gzip"
	"github.com/nomkhonwaan/myblog/pkg/blog"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"regexp"
	"strings"
	"text/template"
)

// IsFacebookCrawlerRequest does checking the request user-agent strings.
// The Facebook crawler user-agent strings are listed here: https://developers.facebook.com/docs/sharing/webmasters/crawler/#identify
func IsFacebookCrawlerRequest(userAgent string) bool {
	return regexp.MustCompile("facebookexternalhit").MatchString(userAgent)
}

// IsSingle validates against URL and return its ID if it is a single
func IsSingle(url string) (string, bool) {
	re := regexp.MustCompile(`\d{4}/\d{1,2}/\d{1,2}/.+-(.+)$`)
	if !re.MatchString(url) {
		return "", false
	}

	matches := re.FindStringSubmatch(url)
	return matches[1], true
}

// Service helps co-working between data-layer and control-layer
type Service interface {
	// A Post repository
	Post() blog.PostRepository
}

type service struct {
	postRepo blog.PostRepository
}

func (s service) Post() blog.PostRepository {
	return s.postRepo
}

// CrawlerMiddleware is a Facebook specific middleware
// for rendering server-side HTML which contains only Facebook's extra meta tags but empty content
type CrawlerMiddleware struct {
	service  Service
	template *template.Template
}

// NewCrawlerMiddleware returns a Facebook's crawler specific middleware instance
func NewCrawlerMiddleware(openGraphTemplate string, postRepo blog.PostRepository) (CrawlerMiddleware, error) {
	t, err := template.New("facebook-opengraph-template").Parse(openGraphTemplate)
	if err != nil {
		return CrawlerMiddleware{}, err
	}

	return CrawlerMiddleware{
		service:  service{postRepo: postRepo},
		template: t,
	}, nil
}

func (mw CrawlerMiddleware) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if IsFacebookCrawlerRequest(r.UserAgent()) {
			if id, yes := IsSingle(r.URL.Path); yes {
				if postID, err := primitive.ObjectIDFromHex(id); err == nil {
					mw.serveSingle(w, r, postID)
					return
				}

			}
		}

		next.ServeHTTP(w, r)
	})
}

func (mw CrawlerMiddleware) serveSingle(w http.ResponseWriter, r *http.Request, id interface{}) {
	p, err := mw.service.Post().FindByID(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Do not allow unpublished post to be shared
	if p.Status != blog.Published {
		http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
		return
	}

	// TODO: a post's featured image function is not ready yet
	data := struct {
		URL           string
		Type          string
		Title         string
		Description   string
		FeaturedImage string
	}{
		URL:           "https://beta.nomkhonwaan.com/" + p.PublishedAt.Format("2006/1/2") + "/" + p.Slug,
		Type:          "article",
		Title:         p.Title,
		Description:   strings.Split(p.Markdown, "\n")[0],
		FeaturedImage: "https://beta.nomkhonwaan.com/assets/images/303589.webp",
	}

	w.Header().Set("Content-Type", "text/html")
	w.Header().Set("Content-Encoding", "gzip")

	wtr, _ := gzip.NewWriterLevel(w, gzip.BestCompression)
	defer wtr.Close()

	if err = mw.template.Execute(wtr, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
