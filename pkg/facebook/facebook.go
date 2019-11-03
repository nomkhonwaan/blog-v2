package facebook

import (
	"github.com/nomkhonwaan/myblog/pkg/blog"
	"net/http"
	"regexp"
)

// IsFacebookCrawlerRequest does checking the request user-agent strings.
// The Facebook crawler user-agent strings are listed here:
// https://developers.facebook.com/docs/sharing/webmasters/crawler/#identify
func IsFacebookCrawlerRequest(userAgent string) bool {
	return regexp.MustCompile("facebookexternalhit").MatchString(userAgent)
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
	service Service
}

// NewCrawlerMiddleware returns a Facebook's crawler specific middleware instance
func NewCrawlerMiddleware(postRepo blog.PostRepository) CrawlerMiddleware {
	return CrawlerMiddleware{
		service: service{
			postRepo: postRepo,
		},
	}
}

func (mw CrawlerMiddleware) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// It is not a Facebook's crawler request, do nothing...
		if !IsFacebookCrawlerRequest(r.UserAgent()) {
			next.ServeHTTP(w, r)
			return
		}

		// Only a single page supported for now
		//paths := strings.Split(r.URL.Path, "/")
		mw.serveSinglePage(w, r)
	})
}

func (mw CrawlerMiddleware) serveSinglePage(w http.ResponseWriter, id interface{}) {
}
