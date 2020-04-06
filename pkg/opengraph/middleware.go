package opengraph

import (
	"bytes"
	"github.com/nomkhonwaan/myblog/pkg/blog"
	"github.com/nomkhonwaan/myblog/pkg/graphql"
	"github.com/nomkhonwaan/myblog/pkg/storage"
	"github.com/nomkhonwaan/myblog/pkg/timeutil"
	"html/template"
	"net/http"
	"regexp"
	"strings"
)

var (
	singlePageRegExp = regexp.MustCompile(`\d{4}/\d{1,2}/\d{1,2}/(.+)$`)
)

// ServeStaticSinglePageMiddleware provides a static HTML page for crawler bot on the single page
func ServeStaticSinglePageMiddleware(baseURL string, ogTmpl *template.Template, postRepository blog.PostRepository, fileRepository storage.FileRepository) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if isFacebookCrawler(r.UserAgent()) {
				if slug, yes := isSinglePage(r.URL.Path); yes {
					id := graphql.Slug(slug).MustGetID()
					p, err := postRepository.FindByID(r.Context(), id)
					if err != nil {
						http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
						return
					}

					if p.Status == blog.StatusPublished {
						var f storage.File
						if !p.FeaturedImage.ID.IsZero() {
							f, _ = fileRepository.FindByID(r.Context(), p.FeaturedImage.ID)
						}

						_, _ = w.Write(renderStaticSinglePage(baseURL, ogTmpl, p, f))
						return
					}
				}
			}

			next.ServeHTTP(w, r)
		})
	}
}

func renderStaticSinglePage(baseURL string, tmpl *template.Template, p blog.Post, f storage.File) []byte {
	featuredImage := baseURL + "/assets/images/303589.webp"
	if f.Slug != "" {
		featuredImage = baseURL + "/api/v2.1/storage/" + f.Slug
	}

	buf := bytes.Buffer{}
	_ = tmpl.Execute(&buf, struct {
		URL           string
		Type          string
		Title         string
		Description   string
		FeaturedImage string
	}{
		URL:           baseURL + "/" + p.PublishedAt.In(timeutil.TimeZoneAsiaBangkok).Format("2006/1/2") + "/" + p.Slug,
		Type:          "article",
		Title:         p.Title,
		Description:   strings.Split(p.Markdown, "\n")[0],
		FeaturedImage: featuredImage,
	})

	return buf.Bytes()
}

func isFacebookCrawler(userAgent string) bool {
	return regexp.MustCompile("facebookexternalhit").MatchString(userAgent)
}

func isSinglePage(url string) (string, bool) {
	if !singlePageRegExp.MatchString(url) {
		return "", false
	}
	matches := singlePageRegExp.FindStringSubmatch(url)
	return matches[1], true
}
