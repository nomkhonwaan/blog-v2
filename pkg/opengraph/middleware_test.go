package opengraph

import (
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/nomkhonwaan/myblog/pkg/blog"
	mock_blog "github.com/nomkhonwaan/myblog/pkg/blog/mock"
	"github.com/nomkhonwaan/myblog/pkg/mongo"
	"github.com/nomkhonwaan/myblog/pkg/storage"
	mock_storage "github.com/nomkhonwaan/myblog/pkg/storage/mock"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"html/template"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestServeStaticSinglePageMiddleware(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	var (
		postRepository = mock_blog.NewMockPostRepository(ctrl)
		fileRepository = mock_storage.NewMockFileRepository(ctrl)
	)

	tmpl := template.Must(template.New("test-opengraph-template").Parse(`
{{.URL}}
{{.Type}}
{{.Title}}
{{.Description}}
{{.FeaturedImage}}
`))
	middleware := ServeStaticSinglePageMiddleware("http://localhost", tmpl, postRepository, fileRepository)
	next := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte("OK"))
	})

	newFacebookCrawlerBot := func(url string) *http.Request {
		req := httptest.NewRequest(http.MethodGet, url, nil)
		req.Header.Set("User-Agent", "facebookexternalhit/1.1")
		return req
	}

	t.Run("With successful rendering template for Facebook crawler bot", func(t *testing.T) {
		// Given
		w := httptest.NewRecorder()
		now := time.Now()
		id := primitive.NewObjectID()
		p := blog.Post{ID: primitive.NewObjectID(), Title: "Test", Slug: "test-" + id.Hex(), Status: blog.StatusPublished,
			Markdown:    "Lorem ipsum dolor sit amet, consectetur adipiscing elit.\nAenean at ornare ipsum.",
			PublishedAt: now, FeaturedImage: mongo.DBRef{ID: primitive.NewObjectID()}}

		postRepository.EXPECT().FindByID(gomock.Any(), id).Return(p, nil)
		fileRepository.EXPECT().FindByID(gomock.Any(), p.FeaturedImage.ID).Return(storage.File{Slug: "test-featured-image"}, nil)

		expected := `
http://localhost/` + now.Format("2006/1/2") + `/test-` + id.Hex() + `
article
Test
Lorem ipsum dolor sit amet, consectetur adipiscing elit.
http://localhost/api/v2.1/storage/test-featured-image
`

		// When
		middleware(next).ServeHTTP(w, newFacebookCrawlerBot("http://localhost/"+now.Format("2006/1/2")+"/test-"+id.Hex()))

		// Then
		assert.Equal(t, expected, w.Body.String())
	})

	t.Run("With successful rendering template with default featured image", func(t *testing.T) {
		// Given
		w := httptest.NewRecorder()
		now := time.Now()
		id := primitive.NewObjectID()
		p := blog.Post{ID: primitive.NewObjectID(), Title: "Test", Slug: "test-" + id.Hex(), Status: blog.StatusPublished,
			Markdown: "Lorem ipsum dolor sit amet, consectetur adipiscing elit.\nAenean at ornare ipsum.", PublishedAt: now}

		postRepository.EXPECT().FindByID(gomock.Any(), id).Return(p, nil)

		expected := `
http://localhost/` + now.Format("2006/1/2") + `/test-` + id.Hex() + `
article
Test
Lorem ipsum dolor sit amet, consectetur adipiscing elit.
http://localhost/assets/images/303589.webp
`

		// When
		middleware(next).ServeHTTP(w, newFacebookCrawlerBot("http://localhost/"+now.Format("2006/1/2")+"/test-"+id.Hex()))

		// Then
		assert.Equal(t, expected, w.Body.String())
	})

	t.Run("When accessing to non single page", func(t *testing.T) {
		// Given
		w := httptest.NewRecorder()

		// When
		middleware(next).ServeHTTP(w, newFacebookCrawlerBot("http://localhost"))

		// Then
		assert.Equal(t, "OK", w.Body.String())
	})

	t.Run("When unable to find a post", func(t *testing.T) {
		// Given
		w := httptest.NewRecorder()

		postRepository.EXPECT().FindByID(gomock.Any(), gomock.Any()).Return(blog.Post{}, errors.New("test unable to find a post"))

		// When
		middleware(next).ServeHTTP(w, newFacebookCrawlerBot("http://localhost/2006/1/2/test-1"))

		// Then
		assert.Equal(t, "Not Found\n", w.Body.String())
	})
}
