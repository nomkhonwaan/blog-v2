package facebook

import (
	"encoding/json"
	"github.com/nomkhonwaan/myblog/pkg/blog"
	"github.com/nomkhonwaan/myblog/pkg/storage"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"regexp"
	"strings"
	"text/template"
	"time"
)

const (
	// GraphAPIEndpoint is an endpoint to the Facebook Graph API
	GraphAPIEndpoint = "https://graph.facebook.com"
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
	fileRepo storage.FileRepository
	postRepo blog.PostRepository
}

func (s service) File() storage.FileRepository {
	return s.fileRepo
}

func (s service) Post() blog.PostRepository {
	return s.postRepo
}

// Client uses to handling with Facebook services
// such as: crawler bot on the single-page application or engagement querying
type Client struct {
	// A base URL to be composed with sharing URL on the open-graph tag
	baseURL string

	// A permanent application access_token for querying engagement object via Facebook Graph API
	appAccessToken string

	// A text template instance which parses the open-graph HTML template already
	openGraphTemplate *template.Template

	service   Service
	transport http.RoundTripper
}

// NewClient returns a new Facebook client instance
func NewClient(baseURL string, appAccessToken string, openGraphTemplate string, fileRepo storage.FileRepository, postRepo blog.PostRepository, transport http.RoundTripper) (Client, error) {
	tmpl, err := template.New("facebook-open-graph-template").Parse(openGraphTemplate)
	if err != nil {
		return Client{}, err
	}

	return Client{
		baseURL:           baseURL,
		appAccessToken:    appAccessToken,
		openGraphTemplate: tmpl,
		service: service{
			fileRepo: fileRepo,
			postRepo: postRepo,
		},
		transport: transport,
	}, nil
}

// CrawlerHandler uses to handling Facebook sharing bot crawler
// for rendering static HTML which will be uses to display on the Facebook feed
func (c Client) CrawlerHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if IsFacebookCrawlerRequest(r.UserAgent()) {
			if id, yes := IsSingle(r.URL.Path); yes {
				if postID, err := primitive.ObjectIDFromHex(id); err == nil {
					c.serveSingle(w, r, postID)
					return
				}

			}
		}

		next.ServeHTTP(w, r)
	})
}

func (c Client) serveSingle(w http.ResponseWriter, r *http.Request, id interface{}) {
	p, err := c.service.Post().FindByID(r.Context(), id)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	if p.Status != blog.Published {
		http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
		return
	}

	featuredImage := c.baseURL + "/assets/images/303589.webp"
	if !p.FeaturedImage.ID.IsZero() {
		file, _ := c.service.File().FindByID(r.Context(), p.FeaturedImage.ID)
		if file.Slug != "" {
			featuredImage = c.baseURL + "/api/v2/storage/" + file.Slug
		}
	}

	data := struct {
		URL           string
		Type          string
		Title         string
		Description   string
		FeaturedImage string
	}{
		URL:           c.baseURL + "/" + p.PublishedAt.In(DefaultTimeZone).Format("2006/1/2") + "/" + p.Slug,
		Type:          "article",
		Title:         p.Title,
		Description:   strings.Split(p.Markdown, "\n")[0],
		FeaturedImage: featuredImage,
	}

	_ = c.openGraphTemplate.Execute(w, data)
}

// GetURL returns a URL shared on a timeline on in a comment
func (c Client) GetURL(id string) (URL, error) {
	req, _ := http.NewRequest(http.MethodGet, GraphAPIEndpoint+"/v5.0/", nil)
	q := req.URL.Query()
	q.Add("id", c.baseURL+id)
	q.Add("access_token", c.appAccessToken)
	q.Add("fields", "engagement")
	req.URL.RawQuery = q.Encode()

	res, err := c.transport.RoundTrip(req)
	if err != nil {
		return URL{}, err
	}
	defer res.Body.Close()

	var body URL
	err = json.NewDecoder(res.Body).Decode(&body)

	return body, err
}
