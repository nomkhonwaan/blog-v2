package facebook

import (
	"github.com/nomkhonwaan/myblog/pkg/blog"
	"github.com/nomkhonwaan/myblog/pkg/storage"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"regexp"
	"strings"
	"text/template"
	"time"
)

var (
	// DefaultTimeZone uses to format date-time in the specific time zone, default is Asia/Bangkok which is GMT + 7
	DefaultTimeZone, _ = time.LoadLocation("Asia/Bangkok")
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

	// A File repository
	File() storage.FileRepository
}

type service struct {
	postRepo blog.PostRepository
	fileRepo storage.FileRepository
}

func (s service) Post() blog.PostRepository {
	return s.postRepo
}

func (s service) File() storage.FileRepository {
	return s.fileRepo
}

// CrawlerMiddleware is a Facebook specific middleware
// for rendering server-side HTML which contains only Facebook's extra meta tags but empty content
type CrawlerMiddleware struct {
	url      string
	service  Service
	template *template.Template
}

// NewCrawlerMiddleware returns a Facebook's crawler specific middleware instance
func NewCrawlerMiddleware(url string, openGraphTemplate string, postRepo blog.PostRepository, fileRepo storage.FileRepository) (CrawlerMiddleware, error) {
	t, err := template.New("facebook-opengraph-template").Parse(openGraphTemplate)
	if err != nil {
		return CrawlerMiddleware{}, err
	}

	return CrawlerMiddleware{
		url: url,
		service: service{
			postRepo: postRepo,
			fileRepo: fileRepo,
		},
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
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	if p.Status != blog.Published {
		http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
		return
	}

	featuredImage := mw.url + "/assets/images/303589.webp"
	if !p.FeaturedImage.ID.IsZero() {
		file, _ := mw.service.File().FindByID(r.Context(), p.FeaturedImage.ID)
		if file.Slug != "" {
			featuredImage = mw.url + "/api/v2/storage/" + file.Slug
		}
	}

	data := struct {
		URL           string
		Type          string
		Title         string
		Description   string
		FeaturedImage string
	}{
		URL:           mw.url + "/" + p.PublishedAt.In(DefaultTimeZone).Format("2006/1/2") + "/" + p.Slug,
		Type:          "article",
		Title:         p.Title,
		Description:   strings.Split(p.Markdown, "\n")[0],
		FeaturedImage: featuredImage,
	}

	_ = mw.template.Execute(w, data)
}
